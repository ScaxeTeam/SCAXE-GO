package block
type OreBlock struct {
	SolidBase
	DropItemID   int
	DropItemMeta int
	DropMin      int
	DropMax      int
	MinTier      int
	HasFortune   bool
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
func (b *OreBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe || toolTier < b.MinTier {
		return nil
	}
	return []Drop{{ID: b.DropItemID, Meta: b.DropItemMeta, Count: b.DropMin}}
}

const (
	ItemCoal         = 263
	ItemDiamond      = 264
	ItemDye          = 351
	ItemRedstoneDust = 331
	ItemEmerald      = 388
	ItemNetherQuartz = 406
)
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
func (b *GlowingRedstoneOreBlock) GetDrops(toolType, toolTier int) []Drop {
	if toolType != ToolTypePickaxe || toolTier < TierIron {
		return nil
	}
	return []Drop{{ID: ItemRedstoneDust, Meta: 0, Count: b.DropMin}}
}
func IsRedstoneOre(blockID uint8) bool {
	return blockID == REDSTONE_ORE || blockID == GLOWING_REDSTONE_ORE
}

func init() {
	Registry.Register(newOre(COAL_ORE, "Coal Ore", ItemCoal, 0, 1, 1, TierWooden, true))
	Registry.Register(newOre(IRON_ORE, "Iron Ore", int(IRON_ORE), 0, 1, 1, TierStone, false))
	Registry.Register(newOre(GOLD_ORE, "Gold Ore", int(GOLD_ORE), 0, 1, 1, TierIron, false))
	Registry.Register(newOre(DIAMOND_ORE, "Diamond Ore", ItemDiamond, 0, 1, 1, TierDiamond, true))
	Registry.Register(newOre(LAPIS_ORE, "Lapis Lazuli Ore", ItemDye, 4, 4, 8, TierStone, true))
	Registry.Register(newOre(REDSTONE_ORE, "Redstone Ore", ItemRedstoneDust, 0, 4, 5, TierIron, true))
	Registry.Register(NewGlowingRedstoneOreBlock())
	Registry.Register(newOre(EMERALD_ORE, "Emerald Ore", ItemEmerald, 0, 1, 1, TierIron, true))
	Registry.Register(newOre(NETHER_QUARTZ_ORE, "Nether Quartz Ore", ItemNetherQuartz, 0, 1, 1, TierWooden, true))
}
