package entity

// passive.go — 被动生物模板（Animal 基类 + Cow 示例）
// 对应 PHP: entity/Animal.php (~82行), entity/Cow.php (~101行)
//
// Animal 基类:
//   - 幼崽系统: IsBaby + Age（-24000 tick 成长到成年）
//   - 繁殖系统: InLove 状态
//   - 掉落经验值
//
// Cow 示例:
//   - NetworkID = 11, MaxHealth = 8
//   - 喂食: 小麦 (ID=296)
//   - 掉落: 生牛肉/皮革（着火→熟牛肉）
//   - AI 行为: 恐慌/闲逛/寻食/看玩家/随机看

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
)

// ============ 成长常量 ============

const (
	BabyGrowAge = -24000 // 幼崽成长tick（负数=幼崽，每tick+1，到0成年）
)

// ============ Animal 基类 ============

// Animal 被动生物基类
type Animal struct {
	*Entity

	// IsBabyFlag 是否为幼崽
	IsBabyFlag bool

	// InLove 是否处于繁殖状态
	InLove bool

	// AnimalAge 年龄计数（负数=幼崽，每tick递增，到0成年）
	AnimalAge int

	// DropExpMin/Max 掉落经验范围
	DropExpMin int
	DropExpMax int

	// FeedFoodID 可喂食的食物物品ID（0=不可喂食）
	FeedFoodID int

	// MobName 生物名称
	MobName string
}

// NewAnimal 创建被动生物
func NewAnimal(networkID int, name string, maxHealth int, width, height float64) *Animal {
	a := &Animal{
		Entity:    NewEntity(),
		MobName:   name,
		AnimalAge: 0,
	}

	a.Entity.NetworkID = networkID
	a.Entity.Width = width
	a.Entity.Height = height
	a.Entity.MaxHealth = maxHealth
	a.Entity.Health = maxHealth
	a.Entity.Gravity = 0.04
	a.Entity.Drag = 0.02
	a.Entity.StepHeight = 0.6

	return a
}

// ============ 幼崽系统 ============

// SetBaby 设置幼崽状态
// 对应 PHP Animal::setBaby()
func (a *Animal) SetBaby(baby bool) {
	a.IsBabyFlag = baby
	if baby {
		a.AnimalAge = BabyGrowAge
	}
}

// IsBaby 是否为幼崽
func (a *Animal) IsBaby() bool {
	return a.IsBabyFlag
}

// ============ 繁殖系统 ============

// SetInLove 设置繁殖状态
func (a *Animal) SetInLove(inLove bool) {
	a.InLove = inLove
}

// IsInLove 是否在繁殖状态
func (a *Animal) IsInLove() bool {
	return a.InLove
}

// ============ Tick ============

// AnimalTickResult 被动生物 tick 结果
type AnimalTickResult struct {
	HasUpdate bool
	GrewUp    bool // 本 tick 是否成年
}

// TickAnimal 被动生物 tick 逻辑
// 对应 PHP Animal::entityBaseTick()
func (a *Animal) TickAnimal() AnimalTickResult {
	result := AnimalTickResult{}

	a.Entity.TicksLived++
	a.AnimalAge++

	// 幼崽成长检查
	if a.IsBabyFlag && a.AnimalAge >= 0 {
		a.SetBaby(false)
		result.GrewUp = true
		result.HasUpdate = true
	}

	return result
}

// ============ NBT ============

// SaveAnimalNBT 保存动物 NBT
func (a *Animal) SaveAnimalNBT() {
	a.Entity.SaveNBT()
	baby := int8(0)
	if a.IsBabyFlag {
		baby = 1
	}
	a.Entity.NamedTag.Set(nbt.NewByteTag("IsBaby", baby))
	a.Entity.NamedTag.Set(nbt.NewShortTag("Age", int16(a.AnimalAge)))
}

// LoadAnimalFromNBT 从 NBT 加载动物数据
func (a *Animal) LoadAnimalFromNBT() {
	if a.Entity.NamedTag == nil {
		return
	}
	if a.Entity.NamedTag.GetByte("IsBaby") == 1 {
		a.SetBaby(true)
	}
	age := a.Entity.NamedTag.GetShort("Age")
	if age != 0 {
		a.AnimalAge = int(age)
	}
}

// ============ 辅助 ============

// GetName 获取生物名称
func (a *Animal) GetName() string {
	return a.MobName
}

// GetFeedFoodID 获取可喂食的物品ID
func (a *Animal) GetFeedFoodID() int {
	return a.FeedFoodID
}

// CanBeFedWith 判断物品是否可以喂食此动物
func (a *Animal) CanBeFedWith(itemID int) bool {
	return a.FeedFoodID > 0 && itemID == a.FeedFoodID
}

// ============================================================
//                       具体生物: Cow
// ============================================================

// CowNetworkID 牛的网络ID
const CowNetworkID = 11

// ItemWheat 小麦物品ID（牛的食物）
const ItemWheat = 296

// NewCow 创建牛实体
// 对应 PHP Cow.php
func NewCow() *Animal {
	cow := NewAnimal(CowNetworkID, "Cow", 8, 0.9, 1.3)
	cow.FeedFoodID = ItemWheat
	cow.DropExpMin = 1
	cow.DropExpMax = 3
	return cow
}

// CowDrops 牛的掉落物计算
// 对应 PHP Cow::getDrops()
//
// 参数:
//   - isOnFire: 牛是否着火（着火则掉熟牛肉）
//   - rand01: 随机数 0 或 1（控制掉落牛肉还是皮革）
//   - count: 掉落数量 (1-2)
//
// 返回: (itemID, meta, count)
func CowDrops(isOnFire bool, rand01 int, count int) (int, int, int) {
	const (
		RawBeef    = 363
		CookedBeef = 364
		Leather    = 334
	)

	if rand01 == 0 {
		if isOnFire {
			return CookedBeef, 0, count
		}
		return RawBeef, 0, count
	}
	return Leather, 0, count
}
