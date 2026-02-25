package entity

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

const (
	FallingBlockNetworkID = 66
	AnvilBlockID = 145
	AnvilMaxFallDamage = 40.0
)
type FallingBlock struct {
	*Entity
	BlockID int
	BlockMeta int
}
func NewFallingBlock(blockID, blockMeta int) *FallingBlock {
	fb := &FallingBlock{
		Entity:    NewEntity(),
		BlockID:   blockID,
		BlockMeta: blockMeta,
	}

	fb.Entity.NetworkID = FallingBlockNetworkID
	fb.Entity.Width = 0.98
	fb.Entity.Height = 0.98
	fb.Entity.Gravity = 0.04
	fb.Entity.Drag = 0.02
	fb.Entity.Health = 1
	fb.Entity.MaxHealth = 1
	fb.Entity.CanCollide = false

	return fb
}
type FallingBlockTickResult struct {
	HasUpdate    bool
	Landed       bool
	ShouldPlace  bool
	ShouldDrop   bool
	PlaceX       int
	PlaceY       int
	PlaceZ       int
	PlaceBlockID int
	PlaceMeta    int
	IsAnvil      bool
	AnvilDamage  float64
}
func (fb *FallingBlock) TickFallingBlock(landingBlockID int, landingBlockSolid bool, landingBlockLiquid bool) FallingBlockTickResult {
	result := FallingBlockTickResult{}

	fb.Entity.TicksLived++
	result.HasUpdate = true

	if fb.Entity.OnGround {
		result.Landed = true
		result.PlaceX = int(math.Round(fb.Entity.Position.X - 0.5))
		result.PlaceY = int(math.Round(fb.Entity.Position.Y))
		result.PlaceZ = int(math.Round(fb.Entity.Position.Z - 0.5))
		result.PlaceBlockID = fb.BlockID
		result.PlaceMeta = fb.BlockMeta
		if landingBlockID > 0 && !landingBlockSolid && !landingBlockLiquid {
			result.ShouldDrop = true
		} else {
			result.ShouldPlace = true
			if fb.BlockID == AnvilBlockID {
				result.IsAnvil = true
				damage := (fb.Entity.FallDistance - 1) * 2
				if damage > AnvilMaxFallDamage {
					damage = AnvilMaxFallDamage
				}
				if damage > 0 {
					result.AnvilDamage = damage
				}
			}
		}
	}

	return result
}
func (fb *FallingBlock) SaveFallingBlockNBT() {
	fb.Entity.SaveNBT()
	fb.Entity.NamedTag.Set(nbt.NewIntTag("TileID", int32(fb.BlockID)))
	fb.Entity.NamedTag.Set(nbt.NewByteTag("Data", int8(fb.BlockMeta)))
}
func (fb *FallingBlock) LoadFromNBT() bool {
	if fb.Entity.NamedTag == nil {
		return false
	}
	tileID := fb.Entity.NamedTag.GetInt("TileID")
	if tileID > 0 {
		fb.BlockID = int(tileID)
	} else {
		tile := fb.Entity.NamedTag.GetByte("Tile")
		if tile > 0 {
			fb.BlockID = int(tile)
		}
	}

	data := fb.Entity.NamedTag.GetByte("Data")
	fb.BlockMeta = int(data)
	return fb.BlockID > 0
}
func (fb *FallingBlock) GetBlockID() int {
	return fb.BlockID
}
func (fb *FallingBlock) GetBlockMeta() int {
	return fb.BlockMeta
}
func (fb *FallingBlock) GetBlockInfo() int32 {
	return int32(fb.BlockID) | (int32(fb.BlockMeta) << 8)
}
func (fb *FallingBlock) IsAnvil() bool {
	return fb.BlockID == AnvilBlockID
}
