package entity

// passive_mobs.go — 更多被动生物（Pig/Chicken/Sheep）
// 对应 PHP: entity/Pig.php, entity/Chicken.php, entity/Sheep.php
//
// 使用 passive.go 中的 Animal 基类创建具体生物。

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/nbt"
)

// ============ 网络ID ============

const (
	PigNetworkID   = 12
	SheepNetworkID = 13
)

// ============ 物品ID（掉落用） ============

const (
	ItemRawPorkchop    = 319
	ItemCookedPorkchop = 320
	ItemRawChicken     = 365
	ItemCookedChicken  = 366
	ItemFeather        = 288
	ItemCarrot         = 391
	// ItemWheat = 296 (已在 passive.go 定义)
	BlockWool = 35
)

// ================================================================
//                            Pig 猪
// ================================================================

// NewPig 创建猪实体
// NetworkID=12, HP=10, 胡萝卜喂食, 掉落猪排
func NewPig() *Animal {
	pig := NewAnimal(PigNetworkID, "Pig", 10, 0.9, 0.9)
	pig.FeedFoodID = ItemCarrot
	pig.DropExpMin = 1
	pig.DropExpMax = 3
	return pig
}

// PigDrops 猪的掉落物计算
// 参数:
//   - isOnFire: 着火→掉熟猪排
//   - count: 掉落数量 (1-2)
func PigDrops(isOnFire bool, count int) (int, int, int) {
	if isOnFire {
		return ItemCookedPorkchop, 0, count
	}
	return ItemRawPorkchop, 0, count
}

// ================================================================
//                         Chicken 鸡
// ================================================================

// Chicken 鸡实体（扩展 Animal，增加下蛋计时器）
type Chicken struct {
	*Animal

	// EggTimer 下蛋倒计时（tick）
	EggTimer int
}

// NewChicken 创建鸡实体
// NetworkID=10, HP=4, 种子喂食, 掉落鸡肉/羽毛, 定时下蛋, 无坠落伤害
func NewChicken() *Chicken {
	c := &Chicken{
		Animal: NewAnimal(ChickenNetworkID, "Chicken", 4, 0.6, 0.7),
	}
	c.Animal.FeedFoodID = ItemWheat // 种子 (简化为小麦)
	c.Animal.DropExpMin = 1
	c.Animal.DropExpMax = 3
	c.Animal.Entity.Gravity = 0.04
	c.ResetEggTimer()
	return c
}

// ResetEggTimer 重置下蛋计时（2.5-5分钟随机）
// 对应 PHP Chicken::resetEggTimer()
func (c *Chicken) ResetEggTimer() {
	c.EggTimer = 3000 + rand.Intn(3001) // 3000-6000 tick
}

// ChickenTickResult 鸡 tick 结果
type ChickenTickResult struct {
	AnimalTickResult
	ShouldLayEgg bool // 是否应该下蛋
}

// TickChicken 鸡的逻辑 tick
// 对应 PHP Chicken::entityBaseTick()
func (c *Chicken) TickChicken() ChickenTickResult {
	base := c.Animal.TickAnimal()
	result := ChickenTickResult{AnimalTickResult: base}

	c.EggTimer--
	if c.EggTimer <= 0 {
		result.ShouldLayEgg = true
		c.ResetEggTimer()
	}

	return result
}

// IsImmuneToFallDamage 鸡免疫坠落伤害
// 对应 PHP Chicken::attack() 中的 CAUSE_FALL 取消
func (c *Chicken) IsImmuneToFallDamage() bool {
	return true
}

// ChickenDrops 鸡的掉落物计算
// 参数:
//   - isOnFire: 着火→熟鸡肉
//   - rand01: 0=鸡肉, 1=羽毛
//   - count: 数量 (1-2)
func ChickenDrops(isOnFire bool, rand01 int, count int) (int, int, int) {
	if rand01 == 0 {
		if isOnFire {
			return ItemCookedChicken, 0, count
		}
		return ItemRawChicken, 0, count
	}
	return ItemFeather, 0, count
}

// ================================================================
//                          Sheep 羊
// ================================================================

// Sheep 羊实体（扩展 Animal，增加颜色和剪毛状态）
type Sheep struct {
	*Animal

	// Color 羊毛颜色 (0-15, 对应 Wool meta)
	Color int

	// Sheared 是否已被剪毛
	Sheared bool
}

// SheepColorWeights 羊颜色权重表（白色最常见）
var SheepColorWeights = []struct {
	Color  int
	Weight int
}{
	{0, 20}, // White
	{1, 5},  // Orange
	{2, 5},  // Magenta
	{5, 5},  // Lime
	{3, 5},  // Light Blue
	{4, 5},  // Yellow
	{6, 5},  // Pink
	{7, 10}, // Gray
	{8, 5},  // Light Gray
	{9, 5},  // Cyan
	{10, 5}, // Purple
	{11, 5}, // Blue
	{12, 5}, // Brown
	{13, 5}, // Green
	{14, 5}, // Red
	{15, 5}, // Black
}

// GetRandomSheepColor 获取随机羊毛颜色（按权重）
// 对应 PHP Sheep::getRandomColor()
func GetRandomSheepColor() int {
	totalWeight := 0
	for _, cw := range SheepColorWeights {
		totalWeight += cw.Weight
	}

	r := rand.Intn(totalWeight)
	cumulative := 0
	for _, cw := range SheepColorWeights {
		cumulative += cw.Weight
		if r < cumulative {
			return cw.Color
		}
	}
	return 0 // 默认白色
}

// NewSheep 创建羊实体
// NetworkID=13, HP=8, 小麦喂食, 随机颜色
func NewSheep() *Sheep {
	s := &Sheep{
		Animal:  NewAnimal(SheepNetworkID, "Sheep", 8, 0.9, 1.3),
		Color:   GetRandomSheepColor(),
		Sheared: false,
	}
	s.Animal.FeedFoodID = ItemWheat
	s.Animal.DropExpMin = 1
	s.Animal.DropExpMax = 3
	return s
}

// NewSheepWithColor 创建指定颜色的羊
func NewSheepWithColor(color int) *Sheep {
	s := &Sheep{
		Animal:  NewAnimal(SheepNetworkID, "Sheep", 8, 0.9, 1.3),
		Color:   color,
		Sheared: false,
	}
	s.Animal.FeedFoodID = ItemWheat
	return s
}

// SetColor 设置羊毛颜色
func (s *Sheep) SetColor(color int) {
	s.Color = color
}

// GetColor 获取羊毛颜色
func (s *Sheep) GetColor() int {
	return s.Color
}

// SetSheared 设置剪毛状态
// 对应 PHP Sheep::setSheared()
func (s *Sheep) SetSheared(sheared bool) {
	s.Sheared = sheared
}

// IsSheared 是否已剪毛
func (s *Sheep) IsSheared() bool {
	return s.Sheared
}

// RegrowWool 羊吃草后恢复毛（对应 eatGrassBehavior）
func (s *Sheep) RegrowWool() {
	s.Sheared = false
}

// SheepDrops 羊的掉落物（羊毛，颜色=meta）
// 如果已剪毛则不掉
func (s *Sheep) SheepDrops() (int, int, int) {
	if s.Sheared {
		return 0, 0, 0
	}
	return BlockWool, s.Color, 1
}

// ShearDrops 剪毛掉落（1-3 个同色羊毛）
func (s *Sheep) ShearDrops(count int) (int, int, int) {
	s.Sheared = true
	return BlockWool, s.Color, count
}

// SaveSheepNBT 保存羊 NBT
func (s *Sheep) SaveSheepNBT() {
	s.Animal.SaveAnimalNBT()
	s.Animal.Entity.NamedTag.Set(nbt.NewByteTag("Color", int8(s.Color)))
	sheared := int8(0)
	if s.Sheared {
		sheared = 1
	}
	s.Animal.Entity.NamedTag.Set(nbt.NewByteTag("Sheared", sheared))
}

// LoadSheepFromNBT 从 NBT 加载羊数据
func (s *Sheep) LoadSheepFromNBT() {
	s.Animal.LoadAnimalFromNBT()
	if s.Animal.Entity.NamedTag != nil {
		s.Color = int(s.Animal.Entity.NamedTag.GetByte("Color"))
		s.Sheared = s.Animal.Entity.NamedTag.GetByte("Sheared") == 1
	}
}

// ================================================================
//                      工厂注册表
// ================================================================

// MobFactory 被动生物工厂方法表
var MobFactory = map[int]func() *Animal{
	CowNetworkID: NewCow,
	PigNetworkID: NewPig,
}

// CreatePassiveMob 根据网络ID创建被动生物
// 返回 nil 表示未知生物
func CreatePassiveMob(networkID int) *Animal {
	if factory, ok := MobFactory[networkID]; ok {
		return factory()
	}
	return nil
}
