package server

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/command"
	"github.com/scaxe/scaxe-go/pkg/command/defaults"
	"github.com/scaxe/scaxe-go/pkg/config"
	"github.com/scaxe/scaxe-go/pkg/entity"
	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/level/anvil"
	"github.com/scaxe/scaxe-go/pkg/logger"
	luapkg "github.com/scaxe/scaxe-go/pkg/lua"
	"github.com/scaxe/scaxe-go/pkg/permission"
	"github.com/scaxe/scaxe-go/pkg/player"
	"github.com/scaxe/scaxe-go/pkg/protocol"
	"github.com/scaxe/scaxe-go/pkg/raknet"
)

const (
	TicksPerSecond = 20
	TickDuration   = time.Second / TicksPerSecond
)

type Server struct {
	mu sync.RWMutex

	Config *config.ServerConfig

	RakNet  *raknet.Server
	Address string

	Players       map[string]*player.Player
	PlayersByName map[string]*player.Player

	Level  *level.Level
	Levels map[string]*level.Level

	Running     bool
	CurrentTick int64
	StartTime   time.Time

	tickTimes    [20]time.Duration
	tickTimeIdx  int
	lastTickTime time.Time

	packetBuffers   map[*player.Player][][]byte
	packetBuffersMu sync.Mutex

	stopChan chan struct{}

	CommandMap *command.CommandMap

	OpManager *permission.OpManager

	PluginManager *luapkg.PluginManager
}

func NewServer(cfg *config.ServerConfig) *Server {
	address := fmt.Sprintf("%s:%d", cfg.ServerIP, cfg.ServerPort)

	player.DebugItemPickup = cfg.DebugItemPickup

	s := &Server{
		Config:        cfg,
		Address:       address,
		Players:       make(map[string]*player.Player),
		PlayersByName: make(map[string]*player.Player),
		Levels:        make(map[string]*level.Level),
		Running:       false,
		CurrentTick:   0,
		packetBuffers: make(map[*player.Player][][]byte),
		stopChan:      make(chan struct{}),
	}

	return s
}

func (s *Server) Start() error {
	logger.Server("Starting server", "address", s.Address)

	s.RakNet = raknet.NewServer(s.Address)

	s.CommandMap = command.NewCommandMap()
	s.CommandMap.Register(defaults.NewListCommand(s))
	s.CommandMap.Register(defaults.NewStatusCommand(s))
	s.CommandMap.Register(defaults.NewVersionCommand())
	s.registerConsoleCommands()

	s.OpManager = permission.NewOpManager("ops.json")
	if err := s.OpManager.Load(); err != nil {
		logger.Error("Failed to load ops.json", "error", err)
	}

	motd := fmt.Sprintf("MCPE;%s;%d;%s;%d;%d;%d;%s;Survival",
		s.Config.MOTD,
		60,
		"0.14.2",
		s.GetOnlineCount(),
		s.Config.MaxPlayers,
		s.RakNet.ServerID(),
		s.Config.LevelName,
	)
	s.RakNet.SetPongData([]byte(motd))

	s.RakNet.OnConnect = s.handleConnect
	s.RakNet.OnDisconnect = s.handleDisconnect
	s.RakNet.OnPacket = s.handlePacket

	if err := s.RakNet.Start(); err != nil {
		return err
	}

	s.Running = true
	s.StartTime = time.Now()

	logger.Banner(s.Config.ServerName, "SCAXE-GO v0.1.0", s.Address, s.Config.MaxPlayers)
	logger.Server("Server started successfully", "tps", TicksPerSecond)

	levelPath := "worlds/" + s.Config.LevelName
	provider, err := anvil.NewAnvilProvider(levelPath)
	if err != nil {
		return fmt.Errorf("failed to create level provider: %v", err)
	}

	s.Level = level.NewLevel(s.Config.LevelName, levelPath, provider, s.Config.LevelType)
	s.Levels[s.Config.LevelName] = s.Level

	// Initialize Lua plugin system
	s.PluginManager = luapkg.NewPluginManager(NewServerAPIAdapter(s), "plugins")
	if err := s.PluginManager.LoadAll(); err != nil {
		logger.Warn("Failed to load some plugins", "error", err)
	}
	defaults.SetPluginManager(s.PluginManager)

	go s.tickLoop()

	return nil
}

func (s *Server) StopChan() <-chan struct{} {
	return s.stopChan
}

func (s *Server) Stop() {
	s.mu.Lock()
	if !s.Running {
		s.mu.Unlock()
		return
	}
	s.Running = false
	s.mu.Unlock()

	logger.Server("Stopping server...")

	// Disable all plugins first
	if s.PluginManager != nil {
		s.PluginManager.DisableAll()
	}

	select {
	case <-s.stopChan:

	default:
		close(s.stopChan)
	}

	logger.Debug("Disconnecting all players")
	for _, p := range s.Players {
		p.Kick("Server closed", false)
	}

	logger.Debug("Saving all levels")
	for name, lvl := range s.Levels {
		if lvl != nil {
			lvl.Save()
			logger.Debug("Saved level", "name", name)
		}
	}

	logger.Debug("Stopping network interfaces")
	if s.RakNet != nil {
		s.RakNet.Stop()
	}

	logger.Server("Server stopped")
}

func (s *Server) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Running
}

func (s *Server) GetOnlineCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.PlayersByName)
}

func (s *Server) GetTPS() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var totalTime time.Duration
	for _, t := range s.tickTimes {
		totalTime += t
	}
	avgTickTime := totalTime / 20

	if avgTickTime <= 0 {
		return 20.0
	}

	tps := float64(time.Second) / float64(avgTickTime)
	if tps > 20.0 {
		tps = 20.0
	}
	return tps
}

func (s *Server) GetMSPT() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var totalTime time.Duration
	for _, t := range s.tickTimes {
		totalTime += t
	}
	avgTickTime := totalTime / 20

	return float64(avgTickTime.Microseconds()) / 1000.0
}

func (s *Server) GetPlayer(username string) *player.Player {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.PlayersByName[username]
}

func (s *Server) GetOnlinePlayers() []*player.Player {
	s.mu.RLock()
	defer s.mu.RUnlock()
	players := make([]*player.Player, 0, len(s.PlayersByName))
	for _, p := range s.PlayersByName {
		players = append(players, p)
	}
	return players
}

func (s *Server) BroadcastPacket(pk protocol.DataPacket) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.PlayersByName {
		if p.Spawned {
			s.sendPacketUnsafe(p, pk)
		}
	}
}

func (s *Server) broadcastPacketExcept(pkt protocol.DataPacket, except *player.Player) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.PlayersByName {
		if p.Spawned && p != except {
			s.sendPacketUnsafe(p, pkt)
		}
	}
}

func (s *Server) sendPacketUnsafe(p *player.Player, pkt protocol.DataPacket) {
	stream := protocol.NewBinaryStream()
	pkt.Encode(stream)
	p.Session.SendPacket(stream.Bytes())
}

func (s *Server) updatePlayerListAdd(p *player.Player) {
	pk := protocol.NewPlayerListPacket()
	pk.Type = protocol.PlayerListTypeAdd
	pk.Entries = []protocol.PlayerListEntry{{
		UUID:     p.UUID,
		EntityID: p.GetID(),
		Username: p.Username,
		SkinName: p.SkinName,
		SkinData: p.SkinData,
	}}

	payload, err := protocol.CreateBatch([]protocol.DataPacket{pk})
	if err != nil {
		logger.Error("Failed to batch updatePlayerListAdd", "error", err)
		return
	}

	batchPk := protocol.NewBatchPacket()
	batchPk.Payload = payload
	s.BroadcastPacket(batchPk)
}

func (s *Server) updatePlayerListRemove(p *player.Player) {
	pk := protocol.NewPlayerListPacket()
	pk.Type = protocol.PlayerListTypeRemove
	pk.Entries = []protocol.PlayerListEntry{{
		UUID: p.UUID,
	}}
	s.BroadcastPacket(pk)
}

func (s *Server) sendExistingPlayersTo(newPlayer *player.Player) {

	players := s.GetOnlinePlayers()

	pk := protocol.NewPlayerListPacket()
	pk.Type = protocol.PlayerListTypeAdd
	for _, p := range players {
		if p != newPlayer && p.Spawned {
			pk.Entries = append(pk.Entries, protocol.PlayerListEntry{
				UUID:     p.UUID,
				EntityID: p.GetID(),
				Username: p.Username,
				SkinName: p.SkinName,
				SkinData: p.SkinData,
			})
		}
	}
	if len(pk.Entries) > 0 {
		s.sendPacket(newPlayer, pk)
	}

	for _, p := range players {
		if p != newPlayer && p.Spawned {
			addPk := protocol.NewAddPlayerPacket()
			addPk.UUID = p.UUID
			addPk.Username = p.Username
			addPk.EntityID = p.GetID()
			addPk.X = float32(p.Position.X)
			addPk.Y = float32(p.Position.Y)
			addPk.Z = float32(p.Position.Z)
			addPk.SpeedX = float32(p.Motion.X)
			addPk.SpeedY = float32(p.Motion.Y)
			addPk.SpeedZ = float32(p.Motion.Z)
			addPk.Yaw = float32(p.Yaw)
			addPk.Pitch = float32(p.Pitch)

			s.sendPacket(newPlayer, addPk)
		}
	}
}

func (s *Server) spawnPlayerTo(p *player.Player, viewer *player.Player) {

	pk := protocol.NewPlayerListPacket()
	pk.Type = protocol.PlayerListTypeAdd
	pk.Entries = []protocol.PlayerListEntry{{
		UUID:     p.UUID,
		EntityID: p.GetID(),
		Username: p.Username,
		SkinName: p.SkinName,
		SkinData: p.SkinData,
	}}
	s.sendPacket(viewer, pk)

	addPk := protocol.NewAddPlayerPacket()
	addPk.UUID = p.UUID
	addPk.Username = p.Username
	addPk.EntityID = p.GetID()
	addPk.X = float32(p.Position.X)
	addPk.Y = float32(p.Position.Y)
	addPk.Z = float32(p.Position.Z)
	addPk.SpeedX = float32(p.Motion.X)
	addPk.SpeedY = float32(p.Motion.Y)
	addPk.SpeedZ = float32(p.Motion.Z)
	addPk.Yaw = float32(p.Yaw)
	addPk.Pitch = float32(p.Pitch)

	s.sendPacket(viewer, addPk)
}

func (s *Server) syncInventory(p *player.Player) {

	inventoryItems := make([]item.Item, 0)
	contents := p.Inventory.GetContents()
	maxSlot := p.Inventory.GetSize()
	for i := 0; i < maxSlot; i++ {
		if it, ok := contents[i]; ok {
			inventoryItems = append(inventoryItems, it)
		} else {
			inventoryItems = append(inventoryItems, item.NewItem(0, 0, 0))
		}
	}

	for i := 0; i < 9; i++ {
		inventoryItems = append(inventoryItems, item.NewItem(0, 0, 0))
	}

	containerPk := protocol.NewContainerSetContentPacket(0, inventoryItems)
	containerPk.HotbarTypes = make([]int32, 9)
	for i := 0; i < 9; i++ {
		slotIndex := p.Inventory.GetHotbarSlotIndex(i)
		if slotIndex == -1 {
			containerPk.HotbarTypes[i] = -1
		} else {
			containerPk.HotbarTypes[i] = int32(slotIndex + 9)
		}
	}
	s.sendPacket(p, containerPk)

	armorItems := p.Inventory.GetArmorContents()
	armorPk := protocol.NewContainerSetContentPacket(120, armorItems)
	s.sendPacket(p, armorPk)

}

func (s *Server) tickLoop() {
	ticker := time.NewTicker(TickDuration)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.tick()
		}
	}
}

func (s *Server) NextEntityID() int64 {
	return entity.NextEntityID()
}

func (s *Server) tick() {
	tickStart := time.Now()

	s.mu.Lock()
	s.CurrentTick++
	s.mu.Unlock()

	if s.Level != nil {
		s.Level.Tick()

		for _, e := range s.Level.GetEntities() {
			if e.HasMovementUpdate() {

				pk := protocol.NewMoveEntityPacket()
				pk.EntityID = e.GetID()
				pos := e.GetPosition()
				pk.X = float32(pos.X)
				pk.Y = float32(pos.Y)
				pk.Z = float32(pos.Z)
				pk.Yaw = float32(e.GetYaw())
				pk.HeadYaw = float32(e.GetYaw())
				pk.Pitch = float32(e.GetPitch())

				s.BroadcastPacket(pk)
			}
		}
	}

	for _, p := range s.GetOnlinePlayers() {
		if p.IsSpawned() {
			func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Error("Panic in Player.Tick", "player", p.Username, "error", r)
					}
				}()
				p.Tick(s.CurrentTick)
			}()
		}
	}

	s.mu.Lock()
	s.tickTimes[s.tickTimeIdx] = time.Since(tickStart)
	s.tickTimeIdx = (s.tickTimeIdx + 1) % 20
	s.lastTickTime = tickStart
	currentTick := s.CurrentTick
	s.mu.Unlock()

	// Tick plugin scheduler tasks
	if s.PluginManager != nil {
		s.PluginManager.Tick(currentTick)
	}

	s.flushPackets()
}

func (s *Server) handleConnect(session *raknet.Session) {
	addr := session.Address()
	logger.Server("New connection", "address", addr)

	p := player.NewPlayer(session, addr, 0)

	s.mu.Lock()
	s.Players[addr] = p
	s.mu.Unlock()
}

func (s *Server) handleDisconnect(session *raknet.Session) {
	addr := session.Address()

	s.mu.Lock()
	p, exists := s.Players[addr]
	if !exists {
		s.mu.Unlock()
		return
	}

	delete(s.Players, addr)
	s.mu.Unlock()

	if p.LoggedIn {
		s.handlePlayerQuit(p)
	}

	p.Close()
}

func (s *Server) UpdatePong() {
	motd := fmt.Sprintf("MCPE;%s;%d;%s;%d;%d;%d;%s;Survival",
		s.Config.MOTD,
		60,
		"0.14.2",
		s.GetOnlineCount(),
		s.Config.MaxPlayers,
		s.RakNet.ServerID(),
		s.Config.LevelName,
	)
	s.RakNet.SetPongData([]byte(motd))
}

func (s *Server) handlePlayerQuit(p *player.Player) {
	username := p.Username

	s.mu.Lock()
	delete(s.PlayersByName, username)
	s.mu.Unlock()

	s.UpdatePong()

	removePlayerPk := protocol.NewRemovePlayerPacket()
	removePlayerPk.EntityID = p.GetID()

	if uuid, err := uuid.Parse(p.UUID); err == nil {
		removePlayerPk.UUID = uuid
	}
	s.broadcastPacketExcept(removePlayerPk, p)

	s.updatePlayerListRemove(p)

	quitMsg := protocol.NewTextPacket()
	quitMsg.TextType = protocol.TextTypeTranslation
	quitMsg.Message = "multiplayer.player.left"
	quitMsg.Parameters = []string{username}
	s.BroadcastPacket(quitMsg)

	logger.PlayerLeave(username, "disconnected")
}

func (s *Server) handlePacket(session *raknet.Session, data []byte) {
	if len(data) == 0 {
		return
	}

	addr := session.Address()
	s.mu.RLock()
	p, exists := s.Players[addr]
	s.mu.RUnlock()

	if !exists {
		logger.Warn("Packet from unknown session", "address", addr)
		return
	}

	packetID := data[0]
	pkt := protocol.GetPacket(packetID)
	if pkt == nil {
		logger.Debug("Unknown packet", "id", fmt.Sprintf("0x%02x", packetID), "from", addr)
		return
	}

	logger.PacketIn(pkt.Name(), addr, "id", fmt.Sprintf("0x%02x", packetID), "size", len(data))

	stream := protocol.NewBinaryStreamFromBytes(data[1:])
	if err := pkt.Decode(stream); err != nil {
		logger.Error("Failed to decode packet", "packet", pkt.Name(), "error", err)
		return
	}

	switch pk := pkt.(type) {
	case *protocol.LoginPacket:
		s.handleLogin(p, pk)
	case *protocol.TextPacket:
		s.handleText(p, pk)
	case *protocol.RequestChunkRadiusPacket:
		s.handleRequestChunkRadius(p, pk)
	case *protocol.MovePlayerPacket:
		s.handleMovePlayer(p, pk)
	case *protocol.PlayerActionPacket:
		s.handlePlayerAction(p, pk)
	case *protocol.AnimatePacket:
		s.handleAnimate(p, pk)
	case *protocol.UseItemPacket:
		s.handleUseItem(p, pk)
	case *protocol.RemoveBlockPacket:
		s.handleRemoveBlock(p, pk)
	case *protocol.MobEquipmentPacket:
		s.handleMobEquipment(p, pk)
	case *protocol.DropItemPacket:
		s.handleDropItem(p, pk)
	default:
		logger.Debug("Unhandled packet", "packet", pkt.Name())
	}
}

func (s *Server) sendPacket(p *player.Player, pkt protocol.DataPacket) {
	stream := protocol.NewBinaryStream()
	pkt.Encode(stream)

	s.packetBuffersMu.Lock()
	s.packetBuffers[p] = append(s.packetBuffers[p], stream.Bytes())
	s.packetBuffersMu.Unlock()

	logger.PacketOut(pkt.Name(), p.GetAddress(), "id", fmt.Sprintf("0x%02x", pkt.ID()), "buffered", true)
}

func (s *Server) sendPacketImmediate(p *player.Player, pkt protocol.DataPacket) {
	stream := protocol.NewBinaryStream()
	pkt.Encode(stream)
	p.Session.SendPacket(stream.Bytes())

	logger.PacketOut(pkt.Name(), p.GetAddress(), "id", fmt.Sprintf("0x%02x", pkt.ID()))
}

func (s *Server) flushPackets() {
	s.packetBuffersMu.Lock()
	defer s.packetBuffersMu.Unlock()

	for p, packets := range s.packetBuffers {
		if len(packets) == 0 {
			continue
		}

		for _, data := range packets {
			p.Session.SendPacket(data)
		}
	}

	s.packetBuffers = make(map[*player.Player][][]byte)
}

func (s *Server) handleLogin(p *player.Player, pkt *protocol.LoginPacket) {
	logger.PlayerJoin(pkt.Username, p.GetAddress(), int(pkt.Protocol))

	if p.LoggedIn {
		logger.Warn("Ignoring duplicate login packet", "player", pkt.Username)
		return
	}

	if !protocol.IsProtocolSupported(int(pkt.Protocol)) {
		logger.Warn("Unsupported protocol", "player", pkt.Username, "protocol", pkt.Protocol)
	}

	logger.Debug("Login Check", "online", s.GetOnlineCount(), "max", s.Config.MaxPlayers)
	if s.GetOnlineCount() >= s.Config.MaxPlayers {
		p.Kick("disconnectionScreen.serverFull", false)
		return
	}

	p.HandleLogin(pkt.Username, pkt.ClientUUID, pkt.SkinID, pkt.SkinData, pkt.Protocol)
	p.ClientID = uint64(pkt.ClientID)

	if s.OpManager.IsOp(pkt.Username, pkt.ClientID) {
		p.SetOp(true)
		logger.Info("Operator logged in", "player", pkt.Username, "cid", pkt.ClientID)
	} else {
		p.SetOp(false)
	}

	s.mu.Lock()
	s.PlayersByName[pkt.Username] = p
	s.mu.Unlock()

	s.UpdatePong()

	playStatus := protocol.NewPlayStatusPacket()
	playStatus.Status = protocol.PlayStatusLoginSuccess
	s.sendPacket(p, playStatus)

	var batchPackets []protocol.DataPacket

	spawn := s.Level.GetSpawnLocation()
	spawnX, spawnY, spawnZ := int32(spawn.X), int32(spawn.Y), int32(spawn.Z)

	startGame := protocol.NewStartGamePacket()
	startGame.Seed = int32(s.Level.GetSeed())
	startGame.Dimension = 0

	genID := int32(1)
	if s.Level.Generator != nil {
		name := s.Level.Generator.GetName()
		if name == "flat" {
			genID = 2
		} else if name == "old" {
			genID = 0
		}
	}
	startGame.Generator = genID
	startGame.Gamemode = int32(s.Config.Gamemode)
	startGame.EntityID = p.GetID()
	startGame.SpawnX = spawnX
	startGame.SpawnY = spawnY
	startGame.SpawnZ = spawnZ
	startGame.X = float32(spawnX)
	startGame.Y = float32(spawnY)
	startGame.Z = float32(spawnZ)
	startGame.LevelID = "d29ybGQ="
	batchPackets = append(batchPackets, startGame)

	advSettings := protocol.NewAdventureSettingsPacket()

	advSettings.Flags = 0
	if s.Config.Gamemode == 1 || s.Config.AllowFlight {
		advSettings.Flags |= 0x80
	}

	advSettings.UserPermission = 2
	advSettings.GlobalPermission = 2

	batchPackets = append(batchPackets, advSettings)

	setTime := protocol.NewSetTimePacket()
	setTime.Time = 0
	setTime.Started = true
	batchPackets = append(batchPackets, setTime)

	setSpawn := protocol.NewSetSpawnPositionPacket()
	setSpawn.X = spawnX
	setSpawn.Y = spawnY
	setSpawn.Z = spawnZ
	batchPackets = append(batchPackets, setSpawn)

	setDiff := protocol.NewSetDifficultyPacket()
	setDiff.Difficulty = 1
	batchPackets = append(batchPackets, setDiff)

	setHealth := protocol.NewSetHealthPacket()
	setHealth.Health = 20
	batchPackets = append(batchPackets, setHealth)

	spawnChunkX := spawnX >> 4
	spawnChunkZ := spawnZ >> 4

	radius := int32(2)
	for cx := spawnChunkX - radius; cx <= spawnChunkX+radius; cx++ {
		for cz := spawnChunkZ - radius; cz <= spawnChunkZ+radius; cz++ {
			chunk := s.Level.GetChunk(cx, cz, true)

			fullChunk := protocol.NewFullChunkDataPacket()
			fullChunk.ChunkX = cx
			fullChunk.ChunkZ = cz
			fullChunk.Order = protocol.ChunkOrderLayered
			fullChunk.Data = chunk.ToPacketBytes()
			batchPackets = append(batchPackets, fullChunk)
			p.MarkChunkLoaded(cx, cz)
		}
	}

	logger.Server("Sending game data", "player", pkt.Username, "packets", len(batchPackets))

	batchPayload, err := protocol.CreateBatch(batchPackets)
	if err != nil {
		logger.Error("Failed to create batch", "error", err)
		return
	}

	batchPkt := protocol.NewBatchPacket()
	batchPkt.Payload = batchPayload
	s.sendPacket(p, batchPkt)

	if s.Config.Gamemode != 3 {
		creativeItems := item.GetCreativeItems()
		s.sendPacket(p, protocol.NewContainerSetContentPacket(121, creativeItems))
	} else {
		s.sendPacket(p, protocol.NewContainerSetContentPacket(121, nil))
	}

	adventurePk := protocol.NewAdventureSettingsPacket()
	flags := int32(0)
	if p.GetGamemode() == 1 {
		flags |= 0x80
	}
	adventurePk.Flags = flags
	adventurePk.UserPermission = 2
	adventurePk.GlobalPermission = 2
	s.sendPacket(p, adventurePk)

	playStatusSpawn := protocol.NewPlayStatusPacket()
	playStatusSpawn.Status = protocol.PlayStatusPlayerSpawn
	s.sendPacket(p, playStatusSpawn)

	welcome := protocol.NewTextPacket()
	welcome.TextType = protocol.TextTypeRaw
	welcome.Message = fmt.Sprintf("§aWelcome to %s!", s.Config.ServerName)
	s.sendPacket(p, welcome)

	p.Spawned = true
	p.Human.Level = s.Level
	p.SetPosition(entity.NewVector3(float64(spawnX), float64(spawnY), float64(spawnZ)))

	s.sendExistingPlayersTo(p)

	for _, other := range s.GetOnlinePlayers() {
		if other != p && other.Spawned {
			s.spawnPlayerTo(p, other)
		}
	}

	logger.Player("Login complete", "player", pkt.Username, "online", s.GetOnlineCount())

	joinMsg := protocol.NewTextPacket()
	joinMsg.TextType = protocol.TextTypeTranslation
	joinMsg.Message = "multiplayer.player.joined"
	joinMsg.Parameters = []string{pkt.Username}
	s.BroadcastPacket(joinMsg)
}

func (s *Server) handleText(p *player.Player, pkt *protocol.TextPacket) {

	if pkt.TextType != protocol.TextTypeChat {
		return
	}

	if pkt.Message == "" {
		return
	}

	if pkt.Message[0] == '/' {
		cmdLine := pkt.Message[1:]
		if s.CommandMap.Dispatch(p, cmdLine) {
			return
		}
		p.SendMessage("Unknown command. Try /help type commands.")
		return
	}

	logger.Player("Chat", "player", p.Username, "message", pkt.Message)

	broadcast := protocol.NewTextPacket()
	broadcast.TextType = protocol.TextTypeChat
	broadcast.SourceName = p.Username
	broadcast.Message = pkt.Message

	s.mu.RLock()
	for _, other := range s.PlayersByName {
		s.sendPacket(other, broadcast)
	}
	s.mu.RUnlock()
}

func (s *Server) handleRequestChunkRadius(p *player.Player, pkt *protocol.RequestChunkRadiusPacket) {
	radius := pkt.Radius
	if radius < 4 {
		radius = 4
	}
	if radius > 16 {
		radius = 16
	}

	p.SetChunkRadius(radius)

	response := protocol.NewChunkRadiusUpdatedPacket()
	response.Radius = radius
	s.sendPacket(p, response)

	logger.Debug("Chunk radius updated", "player", p.Username, "radius", radius)

	s.checkChunks(p)

	s.syncInventory(p)
}

func (s *Server) checkChunks(p *player.Player) {
	radius := p.GetChunkRadius()
	cx := int32(p.Position.X) >> 4
	cz := int32(p.Position.Z) >> 4

	for x := cx - radius; x <= cx+radius; x++ {
		for z := cz - radius; z <= cz+radius; z++ {
			if !p.IsChunkLoaded(x, z) {
				chunk := s.Level.GetChunk(x, z, true)
				if chunk == nil {
					continue
				}

				fullChunk := protocol.NewFullChunkDataPacket()
				fullChunk.ChunkX = x
				fullChunk.ChunkZ = z
				fullChunk.Order = protocol.ChunkOrderLayered
				fullChunk.Data = chunk.ToPacketBytes()
				s.sendPacket(p, fullChunk)

				p.MarkChunkLoaded(x, z)
			}
		}
	}

}

func (s *Server) handleMovePlayer(p *player.Player, pkt *protocol.MovePlayerPacket) {

	oldCX := int32(p.Position.X) >> 4
	oldCZ := int32(p.Position.Z) >> 4

	p.HandleMove(
		float64(pkt.X), float64(pkt.Y), float64(pkt.Z),
		pkt.Yaw, pkt.BodyYaw, pkt.Pitch,
		pkt.OnGround,
	)

	newCX := int32(p.Position.X) >> 4
	newCZ := int32(p.Position.Z) >> 4

	if oldCX != newCX || oldCZ != newCZ {
		s.checkChunks(p)
	}

	broadcastPkt := protocol.NewMovePlayerPacket()
	broadcastPkt.EntityID = p.GetID()
	broadcastPkt.X = pkt.X
	broadcastPkt.Y = pkt.Y
	broadcastPkt.Z = pkt.Z
	broadcastPkt.Yaw = pkt.Yaw
	broadcastPkt.BodyYaw = pkt.BodyYaw
	broadcastPkt.Pitch = pkt.Pitch
	broadcastPkt.Mode = pkt.Mode
	broadcastPkt.OnGround = pkt.OnGround
	s.broadcastPacketExcept(broadcastPkt, p)
}

func (s *Server) handlePlayerAction(p *player.Player, pkt *protocol.PlayerActionPacket) {
	p.HandleAction(pkt.Action)

	if pkt.Action == 2 {
		s.breakBlock(p, pkt.X, pkt.Y, pkt.Z)
	}
}

func (s *Server) breakBlock(p *player.Player, x, y, z int32) {

	chunk := s.Level.GetChunk(int32(x>>4), int32(z>>4), false)
	if chunk == nil {
		return
	}
	bid := chunk.GetBlockId(int(x&0xf), int(y), int(z&0xf))
	meta := chunk.GetBlockData(int(x&0xf), int(y), int(z&0xf))

	if bid == 0 {
		return
	}

	tool := item.NewItem(0, 0, 0)
	drops := block.GetDrops(uint8(bid), uint8(meta), tool)

	chunk.SetBlock(int(x&0xf), int(y), int(z&0xf), 0, 0)

	upk := protocol.NewUpdateBlockPacket(x, int32(y), z, 0, 0)

	s.BroadcastPacket(upk)

	levPk := protocol.NewLevelEventPacket()
	levPk.EventID = 2001
	levPk.X = float32(x) + 0.5
	levPk.Y = float32(y) + 0.5
	levPk.Z = float32(z) + 0.5
	levPk.Data = int32(bid) | (int32(meta) << 12)
	s.BroadcastPacket(levPk)

	if p.Gamemode == 0 {
		for _, drop := range drops {
			if drop.Count > 0 {
				s.dropItem(float32(x)+0.5, float32(y)+0.5, float32(z)+0.5, drop)
			}
		}
	}
}

func (s *Server) dropItem(x, y, z float32, it item.Item) {

	mx := float32(0.0)
	my := float32(0.2)
	mz := float32(0.0)

	s.dropItemWithMotion(x, y, z, it, mx, my, mz, 10)
}

func (s *Server) dropItemWithMotion(x, y, z float32, it item.Item, mx, my, mz float32, delay int) {

	itemEnt := entity.NewItemEntity(it)
	itemEnt.Position = entity.NewVector3(float64(x), float64(y), float64(z))
	itemEnt.Motion = entity.NewVector3(float64(mx), float64(my), float64(mz))
	itemEnt.PickupDelay = delay
	itemEnt.Level = s.Level

	itemEnt.Entity.SetPosition(itemEnt.Position)

	s.Level.AddEntity(itemEnt)

	pk := protocol.NewAddItemEntityPacket()
	pk.EntityID = itemEnt.GetID()
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.SpeedX = mx
	pk.SpeedY = my
	pk.SpeedZ = mz
	pk.Item = it
	s.BroadcastPacket(pk)

	dataPk := protocol.NewSetEntityDataPacket()
	dataPk.EntityID = itemEnt.GetID()
	dataPk.Metadata = itemEnt.Metadata.Encode()
	s.BroadcastPacket(dataPk)
}

func (s *Server) handleDropItem(p *player.Player, pkt *protocol.DropItemPacket) {
	if !p.Spawned || !p.IsAlive() {
		return
	}

	if pkt.Item.ID == 0 {
		return
	}

	droppedItem := pkt.Item

	if !p.Inventory.Contains(droppedItem) && p.Gamemode == 0 {
		logger.Debug("DropItem failed parity check", "player", p.Username, "item", droppedItem.ID)
		s.syncInventory(p)
		return
	}

	if p.Gamemode == 0 {

		if !p.Inventory.Contains(droppedItem) {
			logger.Debug("DropItem failed parity check: item not in inventory", "player", p.Username, "item", droppedItem.ID)
			s.syncInventory(p)
			return
		}

		leftovers := p.Inventory.RemoveItem(droppedItem)
		if len(leftovers) > 0 {

			logger.Warn("Failed to remove dropped item", "player", p.Username)
			s.syncInventory(p)
			return
		}

		held := p.Inventory.GetItemInHand()
		equipPk := protocol.NewMobEquipmentPacket()
		equipPk.EntityID = p.GetEntityID()
		equipPk.ItemID = int16(held.ID)
		equipPk.ItemCount = int8(held.Count)
		equipPk.ItemMeta = uint16(held.Meta)
		equipPk.Slot = byte(p.Inventory.GetHeldItemIndex())
		equipPk.SelectedSlot = byte(p.Inventory.GetHeldItemIndex())
		s.sendPacket(p, equipPk)
		s.broadcastPacketExcept(equipPk, p)

	}

	yaw := float64(p.Yaw)
	pitch := float64(p.Pitch)

	x := -math.Sin(yaw/180*math.Pi) * math.Cos(pitch/180*math.Pi)
	y := -math.Sin(pitch / 180 * math.Pi)
	z := math.Cos(yaw/180*math.Pi) * math.Cos(pitch/180*math.Pi)

	len := math.Sqrt(x*x + y*y + z*z)
	if len > 0 {
		x /= len
		y /= len
		z /= len
	}

	force := 0.3
	motionX := float32(x * force)
	motionY := float32(y * force)
	motionZ := float32(z * force)

	dropX := float32(p.Position.X)
	dropY := float32(p.Position.Y + 1.3)
	dropZ := float32(p.Position.Z)

	s.dropItemWithMotion(dropX, dropY, dropZ, droppedItem, motionX, motionY, motionZ, 40)
	logger.Debug("DropItem spawned", "player", p.Username, "x", dropX, "y", dropY, "z", dropZ)
}

func (s *Server) handleRemoveBlock(p *player.Player, pkt *protocol.RemoveBlockPacket) {

	s.breakBlock(p, pkt.X, int32(pkt.Y), pkt.Z)
}

func (s *Server) handleAnimate(p *player.Player, pkt *protocol.AnimatePacket) {

	broadcastPkt := protocol.NewAnimatePacket()
	broadcastPkt.Action = pkt.Action
	broadcastPkt.EntityID = p.GetID()
	broadcastPkt.Float = pkt.Float
	s.broadcastPacketExcept(broadcastPkt, p)
}

func (s *Server) handleUseItem(p *player.Player, pkt *protocol.UseItemPacket) {

	if pkt.Face <= 5 {

		tx, ty, tz := pkt.X, pkt.Y, pkt.Z
		switch pkt.Face {
		case 0:
			ty--
		case 1:
			ty++
		case 2:
			tz--
		case 3:
			tz++
		case 4:
			tx--
		case 5:
			tx++
		}

		held := p.Inventory.GetItemInHand()

		if held.ID > 0 && held.ID < 256 {

			s.Level.SetBlock(tx, ty, tz, byte(held.ID), byte(held.Meta), false)

			logger.Player("Placed block", "player", p.Username, "block", held.ID, "x", tx, "y", ty, "z", tz)

			updatePk := protocol.NewUpdateBlockPacket(tx, ty, tz, uint8(held.ID), uint8(held.Meta))
			s.BroadcastPacket(updatePk)

			if p.GetGamemode() == 0 {
				held.Count--
				if held.Count <= 0 {
					held = item.NewItem(0, 0, 0)
				}
				p.Inventory.SetItemInHand(held)

				equipPk := protocol.NewMobEquipmentPacket()
				equipPk.EntityID = p.GetEntityID()
				equipPk.ItemID = int16(held.ID)
				equipPk.ItemCount = int8(held.Count)
				equipPk.ItemMeta = uint16(held.Meta)
				equipPk.Slot = byte(p.Inventory.GetHeldItemIndex())
				equipPk.SelectedSlot = byte(p.Inventory.GetHeldItemIndex())
				s.sendPacket(p, equipPk)
			}
		} else {
			logger.Player("Used item on block", "player", p.Username, "item", held.ID, "x", pkt.X, "y", pkt.Y, "z", pkt.Z)
		}
	} else if pkt.Face == 0xff {

		logger.Player("Used item in air", "player", p.Username, "item", pkt.Item.ID)
	}
}

func (s *Server) handleMobEquipment(p *player.Player, pkt *protocol.MobEquipmentPacket) {

	p.Inventory.SetHeldItemIndex(int(pkt.SelectedSlot))

	broadcastPkt := protocol.NewMobEquipmentPacket()
	broadcastPkt.EntityID = p.GetID()
	broadcastPkt.ItemID = pkt.ItemID
	broadcastPkt.ItemCount = pkt.ItemCount
	broadcastPkt.ItemMeta = pkt.ItemMeta
	broadcastPkt.Slot = pkt.Slot
	broadcastPkt.SelectedSlot = pkt.SelectedSlot
	s.broadcastPacketExcept(broadcastPkt, p)

	logger.Debug("Equipment changed", "player", p.Username, "slot", pkt.SelectedSlot)
}
