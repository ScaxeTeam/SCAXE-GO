package player

import (
	"sync"

	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/entity"
	"github.com/scaxe/scaxe-go/pkg/inventory"
	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/protocol"
	"github.com/scaxe/scaxe-go/pkg/raknet"
	"github.com/scaxe/scaxe-go/pkg/world"
)

var DebugItemPickup = false

var _ level.ChunkLoader = (*Player)(nil)

type Player struct {
	*entity.Human
	mu sync.RWMutex

	Session   *raknet.Session
	ClientID  uint64
	IPAddress string
	Port      int
	Protocol  int32

	DisplayName string

	Spawned        bool
	Connected      bool
	LoggedIn       bool
	LoadingChunks  bool
	SpawnReadyTick int64
	Op             bool
	Gamemode       int

	ChunkRadius    int32
	LoadedChunks   map[int64]bool
	chunkLoadQueue []int64
	usedChunks     map[int64]bool
	loadCounter    int
	chunksPerTick  int
	spawnThreshold int
	ChunksToSend   []int64

	LastMoveTime int64
	Ping         int

	Inventory *inventory.PlayerInventory
}

func NewPlayer(session *raknet.Session, ip string, port int) *Player {
	p := &Player{
		Human:          entity.NewHuman(),
		Session:        session,
		ClientID:       0,
		IPAddress:      ip,
		Port:           port,
		Protocol:       0,
		DisplayName:    "",
		Spawned:        false,
		Connected:      true,
		LoggedIn:       false,
		LoadingChunks:  false,
		ChunkRadius:    4,
		LoadedChunks:   make(map[int64]bool),
		usedChunks:     make(map[int64]bool),
		chunksPerTick:  4,
		spawnThreshold: 16,
		ChunksToSend:   nil,
		LastMoveTime:   0,
		Ping:           0,
		Inventory:      inventory.NewPlayerInventory(),
	}

	p.Inventory.OnSlotChange = func(slot int, it item.Item) {

		if !p.Spawned {
			return
		}

		if slot >= 36 {

			armorSlot := slot - 36
			pk := protocol.NewContainerSetSlotPacket(0x78, uint16(armorSlot), it)
			p.SendPacket(pk)

			p.BroadcastArmorChange()
			return
		}

		pk := protocol.NewContainerSetSlotPacket(0, uint16(slot), it)
		p.SendPacket(pk)
		logger.Debug("SyncInventory", "slot", slot, "item", it.ID, "count", it.Count)

		heldIndex := p.Inventory.GetHeldItemIndex()
		if slot == heldIndex {
			mobPk := protocol.NewMobEquipmentPacket()
			mobPk.EntityID = 0
			mobPk.ItemID = int16(it.ID)
			mobPk.ItemCount = int8(it.Count)
			mobPk.ItemMeta = uint16(it.Meta)
			mobPk.Slot = uint8(slot)
			mobPk.SelectedSlot = uint8(heldIndex)

			p.SendPacket(mobPk)
			logger.Debug("SyncHand", "eid", 0, "slot", slot, "item", it.ID)
		}
	}

	p.Inventory.AddItem(item.NewItem(item.IRON_PICKAXE, 0, 1))
	p.Inventory.AddItem(item.NewItem(item.DIAMOND_SWORD, 0, 1))
	p.Inventory.AddItem(item.NewItem(block.PLANKS, 0, 64))

	return p
}

func (p *Player) GetName() string {
	return p.Username
}

func (p *Player) SendMessage(message string) {
	pk := protocol.NewTextPacket()
	pk.TextType = protocol.TextTypeRaw
	pk.Message = message
	p.SendPacket(pk)
}

func (p *Player) SetOp(op bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Op = op
}

func (p *Player) IsOp() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Op
}

func (p *Player) sendInventoryContents() {

	realSize := p.Inventory.GetSize()
	hotbarSize := 9
	totalSlots := realSize + hotbarSize

	items := make([]item.Item, totalSlots)

	for i := 0; i < realSize; i++ {
		items[i] = p.Inventory.GetItem(i)
	}

	for i := realSize; i < totalSlots; i++ {
		items[i] = item.NewItem(0, 0, 0)
	}

	hotbar := make([]int32, hotbarSize)
	for i := 0; i < hotbarSize; i++ {

		index := p.Inventory.GetHotbarSlotIndex(i)
		if index <= -1 {
			hotbar[i] = -1
		} else {

			hotbar[i] = int32(index + 9)
		}
	}

	pk := protocol.NewContainerSetContentPacket(0, items)
	pk.HotbarTypes = hotbar
	p.SendPacket(pk)
}

func (p *Player) SetGamemode(mode int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Gamemode = mode
}

func (p *Player) GetGamemode() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Gamemode
}

func (p *Player) GetEntityID() int64 {
	return p.GetID()
}

func (p *Player) GetDisplayName() string {
	if p.DisplayName != "" {
		return p.DisplayName
	}
	return p.Username
}

func (p *Player) SetDisplayName(name string) {
	p.DisplayName = name
}

func (p *Player) GetAddress() string {
	return p.IPAddress
}

func (p *Player) GetPort() int {
	return p.Port
}

func (p *Player) GetPing() int {
	return p.Ping
}

func (p *Player) UpdatePing(ping int) {
	p.Ping = ping
}

func (p *Player) IsConnected() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Connected
}

func (p *Player) IsOnline() bool {
	return p.IsConnected() && p.LoggedIn
}

func (p *Player) IsSpawned() bool {
	return p.Spawned
}

func (p *Player) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Connected = false
	p.Spawned = false
	if p.Session != nil {
		p.Session.Close()
	}
}

func (p *Player) SendPacket(pk protocol.DataPacket) {
	if p.Session != nil && p.Connected {
		buf := protocol.NewBinaryStream()
		pk.Encode(buf)
		p.Session.SendPacket(buf.Bytes())
	}
}

func (p *Player) BroadcastArmorChange() {

	levelInterface := p.Human.Level
	if levelInterface == nil {
		return
	}

	pk := protocol.NewMobArmorEquipmentPacket()
	pk.EntityID = p.GetID()
	armor := p.Inventory.GetArmorContents()
	for i, it := range armor {
		pk.Slots[i] = protocol.ArmorItem{
			ID:    int16(it.ID),
			Count: int8(it.Count),
			Meta:  uint16(it.Meta),
		}
	}

	entities := p.Human.Level.GetEntities()
	for _, e := range entities {
		if viewer, ok := e.(*Player); ok {
			if viewer != p && viewer.Spawned {
				viewer.SendPacket(pk)
			}
		}
	}
}

func (p *Player) Tick(currentTick int64) bool {
	if !p.IsConnected() {
		return false
	}

	p.checkNearEntities()

	return true
}

func (p *Player) getViewers() []*Player {

	if p.Human.Level == nil {
		return []*Player{}
	}

	var viewers []*Player

	entities := p.Human.Level.GetEntities()
	for _, e := range entities {

		if player, ok := e.(*Player); ok && player.Spawned {
			viewers = append(viewers, player)
		}
	}
	return viewers
}

func (p *Player) checkNearEntities() {
	if p.Human.Level == nil {
		return
	}
	if p.BoundingBox == nil {
		return
	}

	level := p.Human.Level

	grow := 0.6
	yGrowDown := 1.5
	bb := p.BoundingBox
	searchBB := entity.NewAxisAlignedBB(
		bb.MinX-grow, bb.MinY-yGrowDown, bb.MinZ-grow,
		bb.MaxX+grow, bb.MaxY+grow, bb.MaxZ+grow,
	)

	entities := level.GetNearbyEntities(searchBB, p)

	for _, e := range entities {

		itemEnt, ok := e.(*entity.ItemEntity)
		if !ok {
			continue
		}

		if DebugItemPickup {
			logger.Warn("checkNearEntities: checking item",
				"itemID", itemEnt.GetID(),
				"itemBB", itemEnt.BoundingBox,
				"itemPos", itemEnt.Position,
				"PickupDelay", itemEnt.PickupDelay,
				"Closed", itemEnt.Closed)
		}

		if itemEnt.Closed {
			continue
		}

		if itemEnt.PickupDelay > 0 {
			continue
		}

		it := itemEnt.Item

		if p.Gamemode == 0 && !p.Inventory.CanAddItem(it) {
			if DebugItemPickup {
				logger.Warn("Pickup skipped: inventory full or logic", "player", p.Username, "item", it)
			}
			continue
		}

		if DebugItemPickup {
			logger.Warn("Pickup SUCCESS: Processing pickup", "player", p.Username, "item", it.ID)
		}

		pk := protocol.NewTakeItemEntityPacket()
		pk.Target = p.GetID()
		pk.EntityID = itemEnt.GetID()

		p.SendPacket(pk)

		pkSelf := protocol.NewTakeItemEntityPacket()
		pkSelf.Target = 0
		pkSelf.EntityID = itemEnt.GetID()
		p.SendPacket(pkSelf)

		p.Inventory.AddItem(it)

		itemEnt.Close()

		pkRemove := protocol.NewRemoveEntityPacket()
		pkRemove.EntityID = itemEnt.GetID()
		p.SendPacket(pkRemove)

		p.sendInventoryContents()

		logger.Debug("Picked up item", "player", p.Username, "item", itemEnt.Item.ID)

		level.RemoveEntity(itemEnt)
	}
}

func (p *Player) Teleport(x, y, z float64) {
	p.mu.Lock()
	p.Position = entity.NewVector3(x, y, z)
	p.mu.Unlock()

	pk := protocol.NewMovePlayerPacket()
	pk.EntityID = p.GetID()
	pk.X = float32(x)
	pk.Y = float32(y)
	pk.Z = float32(z)
	pk.Yaw = float32(p.Yaw)
	pk.BodyYaw = float32(p.Yaw)
	pk.Pitch = float32(p.Pitch)
	pk.Mode = 1
	pk.OnGround = true
	p.SendPacket(pk)

	p.OrderChunks()
}

func (p *Player) SwitchLevel(lvl interface{}) bool {
	p.mu.Lock()
	currentLevel := p.Human.Level
	p.mu.Unlock()

	if currentLevel == lvl {
		return true
	}

	targetLevel, ok := lvl.(*level.Level)
	if !ok {
		return false
	}
	oldLevel, oldOk := currentLevel.(*level.Level)

	if oldOk && oldLevel != nil {
		oldLevel.RemoveEntity(p)

	}

	p.Human.Level = targetLevel
	targetLevel.AddEntity(p)

	p.mu.Lock()
	p.LoadedChunks = make(map[int64]bool)
	p.chunkLoadQueue = nil
	p.usedChunks = make(map[int64]bool)
	p.mu.Unlock()

	spawn := targetLevel.GetSafeSpawn()

	setSpawn := protocol.NewSetSpawnPositionPacket()
	setSpawn.X = int32(spawn.X)
	setSpawn.Y = int32(spawn.Y)
	setSpawn.Z = int32(spawn.Z)
	p.SendPacket(setSpawn)

	p.Teleport(spawn.X, spawn.Y, spawn.Z)

	p.Tick(0)

	return true
}

func (p *Player) GetLoaderId() int64 {
	return int64(p.ClientID)
}

func (p *Player) OnChunkLoaded(chunk *world.Chunk) {
	if chunk == nil {
		return
	}
	hash := world.ChunkHash(chunk.X, chunk.Z)

	pk := protocol.NewFullChunkDataPacket()
	pk.ChunkX = chunk.X
	pk.ChunkZ = chunk.Z
	pk.Data = chunk.GetPacketBytes()

	p.SendPacket(pk)

	p.usedChunks[hash] = true
	p.loadCounter++
}

func (p *Player) OnChunkUnloaded(chunk *world.Chunk) {

}

func (p *Player) OrderChunks() {
	viewDistance := int32(4)

	chunkX := int32(p.Position.X) >> 4
	chunkZ := int32(p.Position.Z) >> 4

	newQueue := make([]int64, 0)

	for x := -viewDistance; x <= viewDistance; x++ {
		for z := -viewDistance; z <= viewDistance; z++ {

			if x*x+z*z > viewDistance*viewDistance {
				continue
			}

			cx := chunkX + x
			cz := chunkZ + z
			hash := world.ChunkHash(cx, cz)

			if _, exists := p.usedChunks[hash]; !exists {
				p.usedChunks[hash] = false
				newQueue = append(newQueue, hash)
			}
		}
	}
	p.chunkLoadQueue = newQueue
}

func (p *Player) SendNextChunk() {
	if len(p.chunkLoadQueue) == 0 {
		return
	}

	if p.Human.Level == nil {
		return
	}

	lvl, ok := p.Human.Level.(*level.Level)
	if !ok {

		return
	}

	loaded := 0
	for i := 0; i < len(p.chunkLoadQueue); i++ {
		if loaded >= p.chunksPerTick {
			break
		}

		hash := p.chunkLoadQueue[i]
		cx := int32(hash >> 32)
		cz := int32(hash & 0xFFFFFFFF)

		lvl.RequestChunk(cx, cz, p)
		loaded++
	}

	if loaded > 0 {
		if loaded >= len(p.chunkLoadQueue) {
			p.chunkLoadQueue = make([]int64, 0)
		} else {
			p.chunkLoadQueue = p.chunkLoadQueue[loaded:]
		}
	}
}

func (p *Player) DoFirstSpawn() {
	p.Spawned = true

	pk := protocol.NewPlayStatusPacket()
	pk.Status = protocol.PlayStatusPlayerSpawn
	p.SendPacket(pk)

	pkMove := protocol.NewMovePlayerPacket()
	pkMove.EntityID = p.GetID()
	pkMove.X = float32(p.Position.X)
	pkMove.Y = float32(p.Position.Y) + 1.62
	pkMove.Z = float32(p.Position.Z)
	pkMove.Yaw = float32(p.Yaw)
	pkMove.Pitch = float32(p.Pitch)
	pkMove.BodyYaw = float32(p.Yaw)
	pkMove.Mode = 0
	pkMove.OnGround = true
	p.SendPacket(pkMove)

}

func (p *Player) Kick(message string, hideScreen bool) {
	if !p.IsConnected() {
		return
	}

	pk := protocol.NewDisconnectPacket()
	pk.Message = message
	pk.HideDisconnectionScreen = hideScreen

	p.SendPacket(pk)

	p.Close()
}

func (p *Player) HandleLogin(username, uuid, skinName string, skinData []byte, protocol int32) {
	p.Username = username
	p.UUID = uuid
	p.SkinName = skinName
	p.SkinData = string(skinData)
	p.Protocol = protocol
	p.DisplayName = username
	p.LoggedIn = true
}

func (p *Player) HandleMove(x, y, z float64, yaw, bodyYaw, pitch float32, onGround bool) {

	p.mu.Lock()
	p.Human.HandleMove(x, y, z, yaw, bodyYaw, pitch, onGround)
	p.mu.Unlock()
}

func (p *Player) HandleAction(action int32) {
	switch action {
	case ActionJump:

	case ActionStartSprint:
		p.SetSprinting(true)
	case ActionStopSprint:
		p.SetSprinting(false)
	case ActionStartSneak:
		p.SetSneaking(true)
	case ActionStopSneak:
		p.SetSneaking(false)
	case ActionRespawn:
		p.handleRespawn()
	}
}

func (p *Player) handleRespawn() {
	if p.GetHealth() <= 0 {
		p.SetHealth(p.GetMaxHealth())
		p.SetFood(20)
		p.SetSaturation(20)
	}
}

func (p *Player) SetChunkRadius(radius int32) {
	if radius < 2 {
		radius = 2
	}
	if radius > 16 {
		radius = 16
	}
	p.ChunkRadius = radius
}

func (p *Player) GetChunkRadius() int32 {
	return p.ChunkRadius
}

func (p *Player) GetSpawnThreshold() int {
	return p.spawnThreshold
}

func ChunkHash(x, z int32) int64 {
	return int64(x)<<32 | int64(uint32(z))
}

func (p *Player) IsChunkLoaded(x, z int32) bool {
	hash := ChunkHash(x, z)
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.LoadedChunks[hash]
}

func (p *Player) MarkChunkLoaded(x, z int32) {
	hash := ChunkHash(x, z)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.LoadedChunks[hash] = true
}

func (p *Player) UnloadChunk(x, z int32) {
	hash := ChunkHash(x, z)
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.LoadedChunks, hash)
}

func (p *Player) QueueChunks(hashes []int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.chunkLoadQueue = append(p.chunkLoadQueue, hashes...)
}

func (p *Player) GetLoadedChunkList() []int64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]int64, 0, len(p.LoadedChunks))
	for hash := range p.LoadedChunks {
		result = append(result, hash)
	}
	return result
}

func (p *Player) GetLoadedChunkCount() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return len(p.LoadedChunks)
}

const (
	ActionStartBreak   int32 = 0
	ActionAbortBreak   int32 = 1
	ActionStopBreak    int32 = 2
	ActionReleaseItem  int32 = 5
	ActionStopSleeping int32 = 6
	ActionRespawn      int32 = 7
	ActionJump         int32 = 8
	ActionStartSprint  int32 = 9
	ActionStopSprint   int32 = 10
	ActionStartSneak   int32 = 11
	ActionStopSneak    int32 = 12
)
