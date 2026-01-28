package server

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/command"
	"github.com/scaxe/scaxe-go/pkg/command/defaults"
	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/level/anvil"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

func (s *Server) AddOp(name string, cid int64) {
	s.OpManager.AddOp(name, cid)
}

func (s *Server) RemoveOp(name string) {
	s.OpManager.RemoveOp(name)
}

func (s *Server) IsOp(name string) bool {
	return s.OpManager.IsOpByName(name)
}

func (s *Server) HandleConsoleCommand(cmdLine string) {
	if cmdLine == "" {
		return
	}
	s.CommandMap.Dispatch(&command.ConsoleCommandSender{}, cmdLine)
}

func (s *Server) registerConsoleCommands() {

	s.CommandMap.Register(defaults.NewOpCommand(s))
	s.CommandMap.Register(defaults.NewDeopCommand(s))
	s.CommandMap.Register(defaults.NewStopCommand(s))
	s.CommandMap.Register(defaults.NewGamemodeCommand(s))
	s.CommandMap.Register(defaults.NewKickCommand(s))
	s.CommandMap.Register(defaults.NewKillCommand(s))
	s.CommandMap.Register(defaults.NewTeleportCommand(s))
	s.CommandMap.Register(defaults.NewGiveCommand(s))
	s.CommandMap.Register(defaults.NewTimeCommand(s))
	s.CommandMap.Register(defaults.NewSayCommand(s))
	s.CommandMap.Register(defaults.NewTellCommand(s))
	s.CommandMap.Register(defaults.NewHelpCommand())
	s.CommandMap.Register(defaults.NewDifficultyCommand(s))
	s.CommandMap.Register(defaults.NewMeCommand(s))
	s.CommandMap.Register(defaults.NewSeedCommand(s))
	s.CommandMap.Register(defaults.NewTpsCommand(s))
	s.CommandMap.Register(defaults.NewPingCommand())

	s.CommandMap.Register(defaults.NewBanCommand(s))
	s.CommandMap.Register(defaults.NewPardonCommand(s))
	s.CommandMap.Register(defaults.NewBanIpCommand(s))
	s.CommandMap.Register(defaults.NewPardonIpCommand(s))
	s.CommandMap.Register(defaults.NewBanListCommand())
	s.CommandMap.Register(defaults.NewWhitelistCommand(s))
	s.CommandMap.Register(defaults.NewDefaultGamemodeCommand(s))

	s.CommandMap.Register(defaults.NewEffectCommand(s))
	s.CommandMap.Register(defaults.NewEnchantCommand(s))
	s.CommandMap.Register(defaults.NewXpCommand(s))
	s.CommandMap.Register(defaults.NewWeatherCommand(s))
	s.CommandMap.Register(defaults.NewSpawnpointCommand(s))
	s.CommandMap.Register(defaults.NewSetWorldSpawnCommand(s))
	s.CommandMap.Register(defaults.NewSaveCommand(s))
	s.CommandMap.Register(defaults.NewSaveOnCommand(s))
	s.CommandMap.Register(defaults.NewSaveOffCommand(s))

	s.CommandMap.Register(defaults.NewMWCommand(s))
	s.CommandMap.Register(defaults.NewSetBlockCommand(s))
	s.CommandMap.Register(defaults.NewFillCommand(s))
	s.CommandMap.Register(defaults.NewSummonCommand(s))
	s.CommandMap.Register(defaults.NewParticleCommand(s))
	s.CommandMap.Register(defaults.NewBiomeFindCommand(s))
	s.CommandMap.Register(defaults.NewVillageLocateCommand(s))

	s.CommandMap.Register(defaults.NewPluginsCommand())
	s.CommandMap.Register(defaults.NewGcCommand())
	s.CommandMap.Register(defaults.NewTimingsCommand())
	s.CommandMap.Register(defaults.NewRestartCommand())
	s.CommandMap.Register(defaults.NewBackupCommand())
	s.CommandMap.Register(defaults.NewChunkInfoCommand(s))
	s.CommandMap.Register(defaults.NewBiomeCommand())
	s.CommandMap.Register(defaults.NewDumpMemoryCommand())

	s.CommandMap.Register(defaults.NewBanCidCommand(s))
	s.CommandMap.Register(defaults.NewPardonCidCommand())
	s.CommandMap.Register(defaults.NewBanCidByNameCommand(s))
	s.CommandMap.Register(defaults.NewBanIpByNameCommand(s))

	s.CommandMap.Register(defaults.NewMakePluginCommand())
	s.CommandMap.Register(defaults.NewExtractPluginCommand())
	s.CommandMap.Register(defaults.NewExtractPharCommand())
	s.CommandMap.Register(defaults.NewGeneratePluginCommand())
	s.CommandMap.Register(defaults.NewLvdatCommand())

	s.CommandMap.RegisterAlias("gm", "gamemode")
	s.CommandMap.RegisterAlias("tp", "teleport")
	s.CommandMap.RegisterAlias("w", "tell")
	s.CommandMap.RegisterAlias("msg", "tell")
	s.CommandMap.RegisterAlias("unban", "pardon")
	s.CommandMap.RegisterAlias("unban-ip", "pardon-ip")
	s.CommandMap.RegisterAlias("?", "help")
}

func (s *Server) GetPlayerByName(name string) command.PlayerSender {
	p := s.GetPlayer(name)
	if p == nil {
		return nil
	}
	return p
}

func (s *Server) SetPlayerGamemode(playerName string, gamemode int) {
	p := s.GetPlayer(playerName)
	if p == nil {
		return
	}

	p.SetGamemode(gamemode)

	pk := protocol.NewSetPlayerGameTypePacket()
	pk.Gamemode = int32(gamemode & 0x01)
	s.sendPacket(p, pk)

	advSettings := protocol.NewAdventureSettingsPacket()
	flags := int32(0)
	if gamemode == 2 {
		flags |= 0x01
	}
	if gamemode == 1 || gamemode == 3 || s.Config.AllowFlight {
		flags |= 0x80
	}
	if gamemode == 3 {
		flags |= 0x100
	}
	advSettings.Flags = flags
	advSettings.UserPermission = 2
	advSettings.GlobalPermission = 2
	s.sendPacket(p, advSettings)

	if gamemode == 1 {
		creativeItems := item.GetCreativeItems()
		containerPk := protocol.NewContainerSetContentPacket(121, creativeItems)
		s.sendPacket(p, containerPk)
	} else if gamemode != 3 {
		containerPk := protocol.NewContainerSetContentPacket(121, nil)
		s.sendPacket(p, containerPk)
	}

	entityData := protocol.NewSetEntityDataPacket()
	entityData.EntityID = 0

	s.sendPacket(p, entityData)

	logger.Player("Gamemode changed", "player", playerName, "mode", gamemode)
}

type ServerLevelManager struct {
	server *Server
}

func (s *Server) GetLevelManager() defaults.LevelManager {
	return &ServerLevelManager{server: s}
}

func (m *ServerLevelManager) GetLevelNames() []string {
	m.server.mu.RLock()
	defer m.server.mu.RUnlock()
	names := make([]string, 0, len(m.server.Levels))
	for name := range m.server.Levels {
		names = append(names, name)
	}
	return names
}

func (m *ServerLevelManager) GetLevel(name string) interface{} {
	m.server.mu.RLock()
	defer m.server.mu.RUnlock()
	return m.server.Levels[name]
}

func (m *ServerLevelManager) GetDefaultLevel() interface{} {
	return m.server.Level
}

func (m *ServerLevelManager) LoadLevel(name string) (interface{}, error) {
	m.server.mu.Lock()
	defer m.server.mu.Unlock()

	if _, exists := m.server.Levels[name]; exists {
		return m.server.Levels[name], nil
	}

	levelPath := "worlds/" + name
	provider, err := anvil.NewAnvilProvider(levelPath)
	if err != nil {
		return nil, err
	}

	lvl := level.NewLevel(name, levelPath, provider, "normal")
	m.server.Levels[name] = lvl
	logger.Info("Loaded level", "name", name)
	return lvl, nil
}

func (m *ServerLevelManager) GenerateLevel(name string, generatorName string, seed int64) (interface{}, error) {
	m.server.mu.Lock()
	defer m.server.mu.Unlock()

	if _, exists := m.server.Levels[name]; exists {
		return nil, fmt.Errorf("level '%s' already exists", name)
	}

	levelPath := "worlds/" + name
	provider, err := anvil.NewAnvilProvider(levelPath)
	if err != nil {
		return nil, err
	}

	lvl := level.NewLevel(name, levelPath, provider, generatorName)
	lvl.Seed = seed

	m.server.Levels[name] = lvl
	logger.Info("Generated level", "name", name, "generator", generatorName, "seed", seed)
	return lvl, nil
}

func (m *ServerLevelManager) UnloadLevel(name string) bool {
	m.server.mu.Lock()
	defer m.server.mu.Unlock()

	lvl, exists := m.server.Levels[name]
	if !exists {
		return false
	}

	if lvl == m.server.Level {
		return false
	}

	lvl.Close()
	delete(m.server.Levels, name)
	logger.Info("Unloaded level", "name", name)
	return true
}
