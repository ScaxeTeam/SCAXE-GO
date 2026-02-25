package entity

// primed_tnt.go — 点燃的TNT实体
// 对应 PHP: entity/PrimedTNT.php (~114行)
//
// 核心功能:
//   - 引信倒计时（默认80 tick = 4秒）
//   - 引信到 0 时 Kill + Explode
//   - Explode: 触发 ExplosionPrimeEvent → 创建 Explosion(force=4)
//   - SpawnTo: 发送 AddEntityPacket (NetworkID=65)
//   - NBT: 保存/加载 Fuse
//   - 不参与实体碰撞

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
)

// ============ 常量 ============

const (
	PrimedTNTNetworkID = 65

	// DefaultFuse 默认引信时间（tick）
	DefaultFuse = 80

	// ShortFuse 短引信（爆炸连锁时的引信 10-30 tick）
	ShortFuseMin = 10
	ShortFuseMax = 30

	// DefaultExplosionForce TNT 默认爆炸力度
	DefaultExplosionForce = 4.0
)

// ============ PrimedTNT 实体 ============

// PrimedTNT 已点燃的 TNT 实体
type PrimedTNT struct {
	*Entity

	// Fuse 引信剩余 tick
	Fuse int

	// ExplosionForce 爆炸力度（默认4.0）
	ExplosionForce float64

	// BlockBreaking 爆炸是否破坏方块
	BlockBreaking bool
}

// NewPrimedTNT 创建点燃的 TNT 实体
// 参数:
//   - fuse: 引信时间（tick），通常为 80
func NewPrimedTNT(fuse int) *PrimedTNT {
	t := &PrimedTNT{
		Entity:         NewEntity(),
		Fuse:           fuse,
		ExplosionForce: DefaultExplosionForce,
		BlockBreaking:  true,
	}

	t.Entity.NetworkID = PrimedTNTNetworkID
	t.Entity.Width = 0.98
	t.Entity.Height = 0.98
	t.Entity.Gravity = 0.04
	t.Entity.Drag = 0.02
	t.Entity.Health = 1
	t.Entity.MaxHealth = 1
	t.Entity.CanCollide = false

	return t
}

// ============ Tick ============

// PrimedTNTTickResult TNT tick 结果
type PrimedTNTTickResult struct {
	HasUpdate     bool    // 需要同步
	ShouldExplode bool    // 引信到 0，应该爆炸
	Force         float64 // 爆炸力度
	BlockBreaking bool    // 是否破坏方块
	ExplodeX      float64 // 爆炸中心位置
	ExplodeY      float64
	ExplodeZ      float64
}

// TickTNT TNT 的每 tick 逻辑
// 对应 PHP PrimedTNT::entityBaseTick()
func (t *PrimedTNT) TickTNT() PrimedTNTTickResult {
	result := PrimedTNTTickResult{}

	t.Entity.TicksLived++
	t.Fuse--

	// 设置引信数据属性（客户端显示闪烁效果）
	result.HasUpdate = true

	if t.Fuse <= 0 {
		result.ShouldExplode = true
		result.Force = t.ExplosionForce
		result.BlockBreaking = t.BlockBreaking
		result.ExplodeX = t.Entity.Position.X
		result.ExplodeY = t.Entity.Position.Y + t.Entity.Height/2
		result.ExplodeZ = t.Entity.Position.Z
	}

	return result
}

// ============ NBT ============

// SavePrimedTNTNBT 保存 PrimedTNT 的 NBT 数据
// 对应 PHP PrimedTNT::saveNBT()
func (t *PrimedTNT) SavePrimedTNTNBT() {
	t.Entity.SaveNBT()
	t.Entity.NamedTag.Set(nbt.NewByteTag("Fuse", int8(t.Fuse)))
}

// LoadFuseFromNBT 从 NBT 加载引信数据
// 对应 PHP PrimedTNT::initEntity()
func (t *PrimedTNT) LoadFuseFromNBT() {
	if t.Entity.NamedTag != nil {
		fuse := t.Entity.NamedTag.GetByte("Fuse")
		if fuse > 0 {
			t.Fuse = int(fuse)
		}
	}
}

// ============ 辅助 ============

// GetFuse 获取剩余引信 tick
func (t *PrimedTNT) GetFuse() int {
	return t.Fuse
}

// SetFuse 设置引信 tick
func (t *PrimedTNT) SetFuse(fuse int) {
	t.Fuse = fuse
}

// SetExplosionForce 设置爆炸力度（由 ExplosionPrimeEvent 修改）
func (t *PrimedTNT) SetExplosionForce(force float64) {
	t.ExplosionForce = force
}

// SetBlockBreaking 设置是否破坏方块（由 ExplosionPrimeEvent 修改）
func (t *PrimedTNT) SetBlockBreaking(breaking bool) {
	t.BlockBreaking = breaking
}
