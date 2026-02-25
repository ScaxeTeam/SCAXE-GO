package block

// ores.go — 矿石方块（7种模板化）
// 对应 PHP: CoalOre, IronOre, GoldOre, DiamondOre, LapisOre, RedstoneOre, EmeraldOre
//          + GlowingRedstoneOre, NetherQuartzOre
//
// 矿石共同特点:
//   - 继承 Solid，硬度 3，需镐挖
//   - 需要最低镐等级才掉落
//   - 精准采集掉落自身（由 Level 层判断附魔）
//   - 时运增加掉落数（由 Level 层判断附魔）
//
// 此文件为纯数据模板，掉落逻辑（精准采集/时运）由 Level 层处理。

// OreBlock 矿石方块
type OreBlock struct {
	SolidBase
	DropItemID   int  // 正常掉落的物品ID
	DropItemMeta int  // 掉落物品的 meta
	DropMin      int  // 最少掉落数
	DropMax      int  // 最多掉落数（不含时运）
	MinTier      int  // 最低镐等级 (TierWooden=1, TierStone=2, TierIron=3, TierDiamond=4)
	HasFortune   bool // 是否受时运影响
}

func newOre(blockID uint8, name string, dropItemID, dropMeta, dropMin, dropMax, minTier int, hasFortune bool) *OreBlock {
	return &OreBlock{
		SolidBase: SolidBase{
			BlockID:       blockID,
			BlockName:     name,
			BlockHardness: 3,
			BlockToolType: ToolTypePickaxe,
			BlockToolTier: minTier,
		},
		DropItemID:   dropItemID,
		DropItemMeta: dropMeta,
		DropMin:      dropMin,
		DropMax:      dropMax,
		MinTier:      minTier,
		HasFortune:   hasFortune,
	}
}

// GetDrops 矿石掉落（不含精准采集/时运，由 Level 层处理）
func (b *OreBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe || toolTier < b.MinTier {
		return nil
	}
	// 默认掉落最少数量，时运由 Level 层乘算
	return []Drop{{ID: b.DropItemID, Meta: b.DropItemMeta, Count: b.DropMin}}
}

// ---------- 物品ID常量 ----------

const (
	ItemCoal         = 263
	ItemDiamond      = 264
	ItemDye          = 351 // Lapis = Dye:4
	ItemRedstoneDust = 331
	ItemEmerald      = 388
	ItemNetherQuartz = 406
)

// ---------- 发光红石矿特殊逻辑 ----------

// GlowingRedstoneOreBlock 发光红石矿（碰触后变成的方块）
type GlowingRedstoneOreBlock struct {
	OreBlock
}

func NewGlowingRedstoneOreBlock() *GlowingRedstoneOreBlock {
	return &GlowingRedstoneOreBlock{
		OreBlock: OreBlock{
			SolidBase: SolidBase{
				BlockID:         GLOWING_REDSTONE_ORE,
				BlockName:       "Glowing Redstone Ore",
				BlockHardness:   3,
				BlockLightLevel: 9,
				BlockToolType:   ToolTypePickaxe,
				BlockToolTier:   TierIron,
			},
			DropItemID: ItemRedstoneDust,
			DropMin:    4,
			DropMax:    5,
			MinTier:    TierIron,
			HasFortune: true,
		},
	}
}

// GetDrops 发光红石矿掉落红石粉（与普通红石矿相同）
func (b *GlowingRedstoneOreBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe || toolTier < TierIron {
		return nil
	}
	return []Drop{{ID: ItemRedstoneDust, Meta: 0, Count: b.DropMin}}
}

// IsRedstoneOre 判断是否为红石矿（包含发光版本）
func IsRedstoneOre(blockID uint8) bool {
	return blockID == REDSTONE_ORE || blockID == GLOWING_REDSTONE_ORE
}

// ---------- 注册 ----------

func init() {
	// 需要最低等级:
	//   木镐(1): 煤矿、下界石英
	//   石镐(2): 铁矿、青金石
	//   铁镐(3): 金矿、红石矿、绿宝石矿
	//   钻镐(4): 钻石矿
	Registry.Register(newOre(COAL_ORE, "Coal Ore", ItemCoal, 0, 1, 1, TierWooden, true))
	Registry.Register(newOre(IRON_ORE, "Iron Ore", int(IRON_ORE), 0, 1, 1, TierStone, false)) // 精准=掉自身
	Registry.Register(newOre(GOLD_ORE, "Gold Ore", int(GOLD_ORE), 0, 1, 1, TierIron, false))  // 精准=掉自身
	Registry.Register(newOre(DIAMOND_ORE, "Diamond Ore", ItemDiamond, 0, 1, 1, TierDiamond, true))
	Registry.Register(newOre(LAPIS_ORE, "Lapis Lazuli Ore", ItemDye, 4, 4, 8, TierStone, true))
	Registry.Register(newOre(REDSTONE_ORE, "Redstone Ore", ItemRedstoneDust, 0, 4, 5, TierIron, true))
	Registry.Register(NewGlowingRedstoneOreBlock())
	Registry.Register(newOre(EMERALD_ORE, "Emerald Ore", ItemEmerald, 0, 1, 1, TierIron, true))
	Registry.Register(newOre(NETHER_QUARTZ_ORE, "Nether Quartz Ore", ItemNetherQuartz, 0, 1, 1, TierWooden, true))
}
