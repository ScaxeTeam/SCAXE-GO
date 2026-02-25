package entity

// falling_block.go — 下落方块实体（沙子/砂砾/铁砧等）
// 对应 PHP: entity/FallingSand.php (~167行)
//
// 核心功能:
//   - 携带方块ID和meta（TileID/Data from NBT）
//   - 受重力下落
//   - 落地时根据目标位置决定：放置方块 or 掉落物品
//   - 铁砧落地对附近实体造成坠落伤害 (damage = (height-1)*2, max 40)
//   - 不参与/不受除虚空外的伤害

import (
	"math"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

// ============ 常量 ============

const (
	FallingBlockNetworkID = 66

	// 铁砧方块ID（落地伤害特殊处理）
	AnvilBlockID = 145

	// 铁砧最大坠落伤害
	AnvilMaxFallDamage = 40.0
)

// ============ FallingBlock 实体 ============

// FallingBlock 下落方块实体
type FallingBlock struct {
	*Entity

	// BlockID 携带的方块ID
	BlockID int

	// BlockMeta 携带的方块 meta/damage
	BlockMeta int
}

// NewFallingBlock 创建下落方块实体
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

// ============ Tick ============

// FallingBlockTickResult 下落方块 tick 结果
type FallingBlockTickResult struct {
	HasUpdate    bool // 需要同步
	Landed       bool // 已落地
	ShouldPlace  bool // 应该放置方块
	ShouldDrop   bool // 应该掉落为物品
	PlaceX       int  // 放置位置
	PlaceY       int
	PlaceZ       int
	PlaceBlockID int     // 放置的方块ID
	PlaceMeta    int     // 放置的方块meta
	IsAnvil      bool    // 是否为铁砧（需要播放音效+伤害实体）
	AnvilDamage  float64 // 铁砧坠落伤害
}

// TickFallingBlock 下落方块的逻辑 tick
// 对应 PHP FallingSand::entityBaseTick()
//
// 参数:
//   - landingBlockID: 落地位置的方块ID
//   - landingBlockSolid: 落地位置方块是否实心
//   - landingBlockLiquid: 落地位置方块是否液体
func (fb *FallingBlock) TickFallingBlock(landingBlockID int, landingBlockSolid bool, landingBlockLiquid bool) FallingBlockTickResult {
	result := FallingBlockTickResult{}

	fb.Entity.TicksLived++
	result.HasUpdate = true

	if fb.Entity.OnGround {
		result.Landed = true

		// 计算放置位置
		result.PlaceX = int(math.Round(fb.Entity.Position.X - 0.5))
		result.PlaceY = int(math.Round(fb.Entity.Position.Y))
		result.PlaceZ = int(math.Round(fb.Entity.Position.Z - 0.5))
		result.PlaceBlockID = fb.BlockID
		result.PlaceMeta = fb.BlockMeta

		// 落地位置有非实心、非液体方块 → 掉落为物品
		if landingBlockID > 0 && !landingBlockSolid && !landingBlockLiquid {
			result.ShouldDrop = true
		} else {
			// 可以放置方块
			result.ShouldPlace = true

			// 铁砧特殊处理
			if fb.BlockID == AnvilBlockID {
				result.IsAnvil = true
				// 伤害 = (fallDistance - 1) * 2, max 40
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

// ============ NBT ============

// SaveFallingBlockNBT 保存下落方块 NBT
// 对应 PHP FallingSand::saveNBT()
func (fb *FallingBlock) SaveFallingBlockNBT() {
	fb.Entity.SaveNBT()
	fb.Entity.NamedTag.Set(nbt.NewIntTag("TileID", int32(fb.BlockID)))
	fb.Entity.NamedTag.Set(nbt.NewByteTag("Data", int8(fb.BlockMeta)))
}

// LoadFromNBT 从 NBT 加载方块数据
// 对应 PHP FallingSand::initEntity()
func (fb *FallingBlock) LoadFromNBT() bool {
	if fb.Entity.NamedTag == nil {
		return false
	}

	// 优先读取 TileID，兼容旧版 Tile
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

	// 无效方块ID → 关闭
	return fb.BlockID > 0
}

// ============ 辅助 ============

// GetBlockID 获取携带的方块ID
func (fb *FallingBlock) GetBlockID() int {
	return fb.BlockID
}

// GetBlockMeta 获取携带的方块 meta
func (fb *FallingBlock) GetBlockMeta() int {
	return fb.BlockMeta
}

// GetBlockInfo 获取客户端显示用的方块信息 (blockID | (meta << 8))
func (fb *FallingBlock) GetBlockInfo() int32 {
	return int32(fb.BlockID) | (int32(fb.BlockMeta) << 8)
}

// IsAnvil 是否为铁砧
func (fb *FallingBlock) IsAnvil() bool {
	return fb.BlockID == AnvilBlockID
}
