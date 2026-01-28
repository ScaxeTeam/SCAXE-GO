package objects

type OreType struct {
	Material   int
	Meta       byte
	MinY       int
	MaxY       int
	Size       int
	Count      int
	TargetType int
}

const (
	BlockStone       = 1
	BlockDirt        = 3
	BlockSand        = 12
	BlockGravel      = 13
	BlockGoldOre     = 14
	BlockIronOre     = 15
	BlockCoalOre     = 16
	BlockLapisOre    = 21
	BlockDiamondOre  = 56
	BlockRedstoneOre = 73
	BlockEmeraldOre  = 129
	BlockAndesite    = 1
	BlockDiorite     = 1
	BlockGranite     = 1
)

var (
	OreCoal = OreType{
		Material:   BlockCoalOre,
		Meta:       0,
		MinY:       0,
		MaxY:       128,
		Size:       17,
		Count:      20,
		TargetType: BlockStone,
	}

	OreIron = OreType{
		Material:   BlockIronOre,
		Meta:       0,
		MinY:       0,
		MaxY:       64,
		Size:       9,
		Count:      20,
		TargetType: BlockStone,
	}

	OreGold = OreType{
		Material:   BlockGoldOre,
		Meta:       0,
		MinY:       0,
		MaxY:       32,
		Size:       9,
		Count:      2,
		TargetType: BlockStone,
	}

	OreDiamond = OreType{
		Material:   BlockDiamondOre,
		Meta:       0,
		MinY:       0,
		MaxY:       16,
		Size:       8,
		Count:      1,
		TargetType: BlockStone,
	}

	OreRedstone = OreType{
		Material:   BlockRedstoneOre,
		Meta:       0,
		MinY:       0,
		MaxY:       16,
		Size:       8,
		Count:      8,
		TargetType: BlockStone,
	}

	OreLapis = OreType{
		Material:   BlockLapisOre,
		Meta:       0,
		MinY:       0,
		MaxY:       32,
		Size:       7,
		Count:      1,
		TargetType: BlockStone,
	}

	OreEmerald = OreType{
		Material:   BlockEmeraldOre,
		Meta:       0,
		MinY:       4,
		MaxY:       32,
		Size:       1,
		Count:      1,
		TargetType: BlockStone,
	}

	OreDirt = OreType{
		Material:   BlockDirt,
		Meta:       0,
		MinY:       0,
		MaxY:       128,
		Size:       33,
		Count:      10,
		TargetType: BlockStone,
	}

	OreGravel = OreType{
		Material:   BlockGravel,
		Meta:       0,
		MinY:       0,
		MaxY:       128,
		Size:       33,
		Count:      8,
		TargetType: BlockStone,
	}

	OreGranite = OreType{
		Material:   BlockStone,
		Meta:       1,
		MinY:       0,
		MaxY:       80,
		Size:       33,
		Count:      10,
		TargetType: BlockStone,
	}

	OreDiorite = OreType{
		Material:   BlockStone,
		Meta:       3,
		MinY:       0,
		MaxY:       80,
		Size:       33,
		Count:      10,
		TargetType: BlockStone,
	}

	OreAndesite = OreType{
		Material:   BlockStone,
		Meta:       5,
		MinY:       0,
		MaxY:       80,
		Size:       33,
		Count:      10,
		TargetType: BlockStone,
	}
)

func OverworldOres() []OreType {
	return []OreType{
		OreDirt,
		OreGravel,
		OreGranite,
		OreDiorite,
		OreAndesite,
		OreCoal,
		OreIron,
		OreGold,
		OreRedstone,
		OreDiamond,
		OreLapis,
	}
}
