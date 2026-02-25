package entity

// minecart.go — 矿车系列实体
// 对应 PHP: entity/MinecartBase.php, entity/Minecart.php,
//           entity/MinecartChest.php, entity/MinecartHopper.php, entity/MinecartTNT.php
//
// MinecartBase 提供:
//   - 铁轨检测和轨道移动（直线/弯道/上下坡）
//   - 展示方块（Display Block）系统
//   - 移动速度控制
//
// 4种矿车变种:
//   - Minecart (ID 84) — 普通矿车，可乘坐
//   - MinecartChest (ID 98) — 运输矿车，带27格背包
//   - MinecartHopper (ID 96) — 漏斗矿车，自动收集/传输物品
//   - MinecartTNT (ID 97) — TNT矿车，碰撞或激活时爆炸

// ============ 矿车类型常量 ============

const (
	MinecartTypeNormal = 1
	MinecartTypeChest  = 2
	MinecartTypeHopper = 3
	MinecartTypeTNT    = 4
)

// ============ 矿车状态 ============

const (
	MinecartStateInitial = 0
	MinecartStateOnRail  = 1
	MinecartStateOffRail = 2
)

// ============ 方向常量 ============

const (
	DirectionNorth = 2
	DirectionSouth = 3
	DirectionWest  = 4
	DirectionEast  = 5
)

// ============ MinecartBase 基类 ============

// MinecartBase 矿车基类
type MinecartBase struct {
	*Entity

	// State 矿车状态
	State int

	// MoveSpeed 移动速度 (0 - 0.5)
	MoveSpeed float64

	// Direction 当前方向
	Direction int

	// DisplayBlockID 展示方块ID
	DisplayBlockID int

	// DisplayBlockMeta 展示方块 meta
	DisplayBlockMeta int

	// DisplayOffset 展示方块偏移
	DisplayOffset int

	// HasDisplayBlock 是否有展示方块
	HasDisplayBlock bool

	// MinecartType 矿车类型
	MinecartType int

	// CartName 矿车名称
	CartName string

	// DropItemID 破坏时掉落的物品ID
	DropItemID int
}

// NewMinecartBase 创建矿车基类
func NewMinecartBase(networkID int, name string, cartType int, dropItemID int) *MinecartBase {
	m := &MinecartBase{
		Entity:        NewEntity(),
		State:         MinecartStateInitial,
		MoveSpeed:     0.5,
		Direction:     -1,
		DisplayOffset: 6,
		MinecartType:  cartType,
		CartName:      name,
		DropItemID:    dropItemID,
	}

	m.Entity.NetworkID = networkID
	m.Entity.Width = 0.98
	m.Entity.Height = 0.7
	m.Entity.MaxHealth = 1
	m.Entity.Health = 1
	m.Entity.Gravity = 0.5
	m.Entity.Drag = 0.1

	return m
}

// ============ 展示方块 ============

// SetDisplayBlock 设置展示方块
// 对应 PHP MinecartBase::setDisplayBlock()
func (m *MinecartBase) SetDisplayBlock(blockID, blockMeta int) {
	m.DisplayBlockID = blockID
	m.DisplayBlockMeta = blockMeta
}

// GetDisplayBlock 获取展示方块
func (m *MinecartBase) GetDisplayBlock() (blockID, blockMeta int) {
	return m.DisplayBlockID, m.DisplayBlockMeta
}

// SetHasDisplay 设置是否有展示方块
func (m *MinecartBase) SetHasDisplay(has bool) {
	m.HasDisplayBlock = has
}

// SetDisplayOffset 设置展示方块偏移
func (m *MinecartBase) SetDisplayOffset(offset int) {
	m.DisplayOffset = offset
}

// GetDisplayOffset 获取展示方块偏移
func (m *MinecartBase) GetDisplayOffset() int {
	return m.DisplayOffset
}

// ============ 移动 ============

// SetMoveSpeed 设置移动速度
func (m *MinecartBase) SetMoveSpeed(speed float64) {
	if speed < 0 {
		speed = 0
	}
	if speed > 0.5 {
		speed = 0.5
	}
	m.MoveSpeed = speed
}

// GetMoveSpeed 获取移动速度
func (m *MinecartBase) GetMoveSpeed() float64 {
	return m.MoveSpeed
}

// IsOnRail 是否在铁轨上
func (m *MinecartBase) IsOnRail() bool {
	return m.State == MinecartStateOnRail
}

// GetName 获取矿车名称
func (m *MinecartBase) GetName() string {
	return m.CartName
}

// GetType 获取矿车类型
func (m *MinecartBase) GetType() int {
	return m.MinecartType
}

// GetDropItemID 获取掉落物品ID
func (m *MinecartBase) GetDropItemID() int {
	return m.DropItemID
}

// ============================================================
//               Minecart 普通矿车 (ID 84)
// ============================================================

const MinecartNetworkID = 84

// NewMinecart 创建普通矿车
// 对应 PHP Minecart.php
func NewMinecart() *MinecartBase {
	const MinecartItemID = 328
	return NewMinecartBase(MinecartNetworkID, "Minecart", MinecartTypeNormal, MinecartItemID)
}

// ============================================================
//           MinecartChest 运输矿车 (ID 98)
// ============================================================

const MinecartChestNetworkID = 98

// MinecartChest 运输矿车（带 27 格背包）
type MinecartChest struct {
	*MinecartBase
}

// NewMinecartChest 创建运输矿车
// 对应 PHP MinecartChest.php
// 展示方块: Chest (54)
func NewMinecartChest() *MinecartChest {
	const (
		ChestMinecartItemID = 342
		ChestBlockID        = 54
	)
	base := NewMinecartBase(MinecartChestNetworkID, "Minecart with Chest", MinecartTypeChest, ChestMinecartItemID)
	base.SetDisplayBlock(ChestBlockID, 0)
	base.SetHasDisplay(true)

	return &MinecartChest{MinecartBase: base}
}

// ============================================================
//           MinecartHopper 漏斗矿车 (ID 96)
// ============================================================

const MinecartHopperNetworkID = 96

// MinecartHopper 漏斗矿车（带 5 格背包，自动收集物品）
type MinecartHopper struct {
	*MinecartBase

	// Cooldown 物品传输冷却
	Cooldown int
}

// NewMinecartHopper 创建漏斗矿车
// 对应 PHP MinecartHopper.php
// 展示方块: Hopper (154)
func NewMinecartHopper() *MinecartHopper {
	const (
		HopperMinecartItemID = 408
		HopperBlockID        = 154
	)
	base := NewMinecartBase(MinecartHopperNetworkID, "Minecart with Hopper", MinecartTypeHopper, HopperMinecartItemID)
	base.SetDisplayBlock(HopperBlockID, 0)
	base.SetDisplayOffset(1)
	base.SetHasDisplay(true)

	return &MinecartHopper{MinecartBase: base}
}

// ResetCooldown 重置冷却
func (h *MinecartHopper) ResetCooldown() {
	h.Cooldown = 1
}

// HasCooldown 是否在冷却中
func (h *MinecartHopper) HasCooldown() bool {
	return h.Cooldown > 0
}

// TickCooldown 冷却 tick
func (h *MinecartHopper) TickCooldown() {
	if h.Cooldown > 0 {
		h.Cooldown--
	}
}

// ============================================================
//            MinecartTNT TNT矿车 (ID 97)
// ============================================================

const MinecartTNTNetworkID = 97

// MinecartTNT TNT矿车
type MinecartTNT struct {
	*MinecartBase

	// Primed 是否已点燃
	Primed bool

	// FuseTicks 引信倒计时
	FuseTicks int
}

// NewMinecartTNT 创建TNT矿车
// 对应 PHP MinecartTNT.php
// 展示方块: TNT (46)
func NewMinecartTNT() *MinecartTNT {
	const (
		TNTMinecartItemID = 407
		TNTBlockID        = 46
	)
	base := NewMinecartBase(MinecartTNTNetworkID, "Minecart with TNT", MinecartTypeTNT, TNTMinecartItemID)
	base.SetDisplayBlock(TNTBlockID, 0)
	base.SetHasDisplay(true)

	return &MinecartTNT{MinecartBase: base}
}

// Prime 点燃TNT矿车
func (t *MinecartTNT) Prime() {
	t.Primed = true
	t.FuseTicks = 80 // 4秒
}

// IsPrimed 是否已点燃
func (t *MinecartTNT) IsPrimed() bool {
	return t.Primed
}

// TNTTickResult TNT矿车 tick 结果
type TNTTickResult struct {
	ShouldExplode bool
}

// TickTNT TNT矿车引信 tick
func (t *MinecartTNT) TickTNT() TNTTickResult {
	if !t.Primed {
		return TNTTickResult{}
	}

	t.FuseTicks--
	if t.FuseTicks <= 0 {
		return TNTTickResult{ShouldExplode: true}
	}
	return TNTTickResult{}
}

// GetExplosionPower TNT矿车爆炸威力
func (t *MinecartTNT) GetExplosionPower() float64 {
	return 4.0
}
