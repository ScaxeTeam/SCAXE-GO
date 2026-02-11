package server

import (
	"github.com/scaxe/scaxe-go/pkg/command"
	luapkg "github.com/scaxe/scaxe-go/pkg/lua"
	"github.com/scaxe/scaxe-go/pkg/player"
)

type ServerAPIAdapter struct {
	server *Server
}

func NewServerAPIAdapter(s *Server) *ServerAPIAdapter {
	return &ServerAPIAdapter{server: s}
}

func (a *ServerAPIAdapter) BroadcastMessage(message string) {
	a.server.BroadcastMessage(message)
}

func (a *ServerAPIAdapter) GetOnlineCount() int {
	return a.server.GetOnlineCount()
}

func (a *ServerAPIAdapter) GetMaxPlayers() int {
	return a.server.Config.MaxPlayers
}

func (a *ServerAPIAdapter) GetTPS() float64 {
	return a.server.GetTPS()
}

func (a *ServerAPIAdapter) GetServerName() string {
	return a.server.Config.ServerName
}

func (a *ServerAPIAdapter) GetPlayer(username string) luapkg.PlayerAPI {
	p := a.server.GetPlayer(username)
	if p == nil {
		return nil
	}
	return &PlayerAPIAdapter{player: p, server: a.server}
}

func (a *ServerAPIAdapter) GetOnlinePlayers() []luapkg.PlayerAPI {
	players := a.server.GetOnlinePlayers()
	result := make([]luapkg.PlayerAPI, len(players))
	for i, p := range players {
		result[i] = &PlayerAPIAdapter{player: p, server: a.server}
	}
	return result
}

func (a *ServerAPIAdapter) KickPlayer(username string, reason string) {
	p := a.server.GetPlayer(username)
	if p != nil {
		p.Kick(reason, false)
	}
}

func (a *ServerAPIAdapter) GetLevel() luapkg.LevelAPI {
	if a.server.Level == nil {
		return nil
	}
	return &LevelAPIAdapter{server: a.server}
}

func (a *ServerAPIAdapter) RegisterCommand(cmd command.Command) {
	a.server.CommandMap.Register(cmd)
}

func (a *ServerAPIAdapter) UnregisterCommand(name string) {
	a.server.CommandMap.Unregister(name)
}

func (a *ServerAPIAdapter) Stop() {
	a.server.Stop()
}

func (a *ServerAPIAdapter) GetCurrentTick() int64 {
	a.server.mu.RLock()
	defer a.server.mu.RUnlock()
	return a.server.CurrentTick
}

type PlayerAPIAdapter struct {
	player *player.Player
	server *Server
}

func (p *PlayerAPIAdapter) GetName() string {
	return p.player.Username
}

func (p *PlayerAPIAdapter) GetPosition() (x, y, z float64) {
	pos := p.player.GetPosition()
	if pos == nil {
		return 0, 0, 0
	}
	return pos.X, pos.Y, pos.Z
}

func (p *PlayerAPIAdapter) SetPosition(x, y, z float64) {
	p.player.Teleport(x, y, z)
}

func (p *PlayerAPIAdapter) SendMessage(msg string) {
	p.player.SendMessage(msg)
}

func (p *PlayerAPIAdapter) GetGamemode() int {
	return p.player.GetGamemode()
}

func (p *PlayerAPIAdapter) SetGamemode(mode int) {
	p.server.SetPlayerGamemode(p.player.Username, mode)
}

func (p *PlayerAPIAdapter) IsOp() bool {
	return p.player.IsOp()
}

func (p *PlayerAPIAdapter) Kick(reason string) {
	p.player.Kick(reason, false)
}

func (p *PlayerAPIAdapter) GetHealth() int {
	return p.player.GetHealth()
}

func (p *PlayerAPIAdapter) SetHealth(health int) {
	p.player.SetHealth(health)
}

func (p *PlayerAPIAdapter) GetEntityID() int64 {
	return p.player.GetEntityID()
}

type LevelAPIAdapter struct {
	server *Server
}

func (l *LevelAPIAdapter) GetBlock(x, y, z int32) (id, meta uint8) {
	if l.server.Level == nil {
		return 0, 0
	}
	bs := l.server.Level.GetBlock(x, y, z)
	return bs.ID, bs.Meta
}

func (l *LevelAPIAdapter) SetBlock(x, y, z int32, id, meta uint8) {
	if l.server.Level == nil {
		return
	}
	l.server.Level.SetBlock(x, y, z, id, meta, true)
}

func (l *LevelAPIAdapter) GetTime() int64 {
	if l.server.Level == nil {
		return 0
	}
	return l.server.Level.GetTime()
}

func (l *LevelAPIAdapter) SetTime(t int64) {
	if l.server.Level == nil {
		return
	}
	l.server.Level.SetTime(t)
}

func (l *LevelAPIAdapter) GetSeed() int64 {
	if l.server.Level == nil {
		return 0
	}
	return l.server.Level.GetSeed()
}

func (l *LevelAPIAdapter) GetSpawnLocation() (x, y, z float64) {
	if l.server.Level == nil {
		return 0, 0, 0
	}
	spawn := l.server.Level.GetSpawnLocation()
	if spawn == nil {
		return 0, 64, 0
	}
	return spawn.X, spawn.Y, spawn.Z
}
