package server

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/entity"
	"github.com/scaxe/scaxe-go/pkg/entity/ai"
	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/player"
)

type levelAccessAdapter struct {
	level  *level.Level
	server *Server
}

type blockInfoAdapter struct {
	id byte
}

func (b *blockInfoAdapter) IsSolid() bool {
	return block.Registry.IsSolid(b.id)
}

func (b *blockInfoAdapter) IsAir() bool {
	return b.id == block.AIR
}

func (a *levelAccessAdapter) GetBlock(x, y, z int) ai.BlockInfo {
	id := a.level.GetBlockId(int32(x), int32(y), int32(z))
	return &blockInfoAdapter{id: id}
}

func (a *levelAccessAdapter) GetNearestPlayer(x, y, z float64, maxDistance float64) ai.PlayerEntity {
	var nearest ai.PlayerEntity
	nearestDist := maxDistance * maxDistance

	a.server.mu.RLock()
	defer a.server.mu.RUnlock()

	for _, p := range a.server.PlayersByName {
		if !p.Spawned || p.Position == nil {
			continue
		}
		dx := p.Position.X - x
		dy := p.Position.Y - y
		dz := p.Position.Z - z
		dist := dx*dx + dy*dy + dz*dz
		if dist < nearestDist {
			nearestDist = dist
			nearest = &playerEntityAdapter{player: p}
		}
	}
	return nearest
}

func (a *levelAccessAdapter) GetEntities() []ai.MobEntity {

	return nil
}

type playerEntityAdapter struct {
	player *player.Player
}

func (p *playerEntityAdapter) GetPosition() (x, y, z float64) {
	return p.player.Position.X, p.player.Position.Y, p.player.Position.Z
}

func (p *playerEntityAdapter) IsAlive() bool {
	return p.player.Health > 0
}

func (p *playerEntityAdapter) IsConnected() bool {
	return p.player.Connected
}

func (p *playerEntityAdapter) IsSurvival() bool {
	return p.player.Gamemode == 0
}

func (s *Server) initMobAI(e entity.IEntity) {

	type mobSetter interface {
		SetLevelAccess(ai.LevelAccess)
	}
	if ms, ok := e.(mobSetter); ok {
		ms.SetLevelAccess(&levelAccessAdapter{level: s.Level, server: s})
	}
}
