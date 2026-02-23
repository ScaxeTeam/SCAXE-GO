package entity

const (
	MinecartTypeNormal = 1
	MinecartTypeChest  = 2
	MinecartTypeHopper = 3
	MinecartTypeTNT    = 4
)

const (
	MinecartStateInitial = 0
	MinecartStateOnRail  = 1
	MinecartStateOffRail = 2
)

const (
	DirectionNorth = 2
	DirectionSouth = 3
	DirectionWest  = 4
	DirectionEast  = 5
)

type MinecartBase struct {
	*Entity

	State int

	MoveSpeed float64

	Direction int

	DisplayBlockID int

	DisplayBlockMeta int

	DisplayOffset int

	HasDisplayBlock bool

	MinecartType int

	CartName string

	DropItemID int
}

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

func (m *MinecartBase) SetDisplayBlock(blockID, blockMeta int) {
	m.DisplayBlockID = blockID
	m.DisplayBlockMeta = blockMeta
}

func (m *MinecartBase) GetDisplayBlock() (blockID, blockMeta int) {
	return m.DisplayBlockID, m.DisplayBlockMeta
}

func (m *MinecartBase) SetHasDisplay(has bool) {
	m.HasDisplayBlock = has
}

func (m *MinecartBase) SetDisplayOffset(offset int) {
	m.DisplayOffset = offset
}

func (m *MinecartBase) GetDisplayOffset() int {
	return m.DisplayOffset
}

func (m *MinecartBase) SetMoveSpeed(speed float64) {
	if speed < 0 {
		speed = 0
	}
	if speed > 0.5 {
		speed = 0.5
	}
	m.MoveSpeed = speed
}

func (m *MinecartBase) GetMoveSpeed() float64 {
	return m.MoveSpeed
}

func (m *MinecartBase) IsOnRail() bool {
	return m.State == MinecartStateOnRail
}

func (m *MinecartBase) GetName() string {
	return m.CartName
}

func (m *MinecartBase) GetType() int {
	return m.MinecartType
}

func (m *MinecartBase) GetDropItemID() int {
	return m.DropItemID
}

const MinecartNetworkID = 84

func NewMinecart() *MinecartBase {
	const MinecartItemID = 328
	return NewMinecartBase(MinecartNetworkID, "Minecart", MinecartTypeNormal, MinecartItemID)
}

const MinecartChestNetworkID = 98

type MinecartChest struct {
	*MinecartBase
}

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

const MinecartHopperNetworkID = 96

type MinecartHopper struct {
	*MinecartBase

	Cooldown int
}

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

func (h *MinecartHopper) ResetCooldown() {
	h.Cooldown = 1
}

func (h *MinecartHopper) HasCooldown() bool {
	return h.Cooldown > 0
}

func (h *MinecartHopper) TickCooldown() {
	if h.Cooldown > 0 {
		h.Cooldown--
	}
}

const MinecartTNTNetworkID = 97

type MinecartTNT struct {
	*MinecartBase

	Primed bool

	FuseTicks int
}

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

func (t *MinecartTNT) Prime() {
	t.Primed = true
	t.FuseTicks = 80
}

func (t *MinecartTNT) IsPrimed() bool {
	return t.Primed
}

type TNTTickResult struct {
	ShouldExplode bool
}

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

func (t *MinecartTNT) GetExplosionPower() float64 {
	return 4.0
}
