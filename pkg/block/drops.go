package block

import (
	"math/rand"

	"github.com/scaxe/scaxe-go/pkg/item"
)

func GetDrops(id uint8, meta uint8, tool item.Item) []item.Item {

	switch id {
	case AIR, WATER, STILL_WATER, LAVA, STILL_LAVA, BEDROCK:
		return []item.Item{}

	case GRASS:
		return []item.Item{item.NewItem(DIRT, 0, 1)}

	case STONE:
		return []item.Item{item.NewItem(COBBLESTONE, 0, 1)}

	case COAL_ORE:
		return []item.Item{item.NewItem(item.COAL, 0, 1)}

	case DIAMOND_ORE:
		return []item.Item{item.NewItem(item.DIAMOND, 0, 1)}

	case IRON_ORE:
		return []item.Item{item.NewItem(IRON_ORE, 0, 1)}

	case GOLD_ORE:
		return []item.Item{item.NewItem(GOLD_ORE, 0, 1)}

	case LAPIS_ORE:
		count := 4 + rand.Intn(5)
		return []item.Item{item.NewItem(item.DYE, 4, count)}

	case REDSTONE_ORE, GLOWING_REDSTONE_ORE:
		count := 4 + rand.Intn(2)
		return []item.Item{item.NewItem(item.REDSTONE_DUST, 0, count)}

	case GRAVEL:
		if rand.Float32() < 0.1 {
			return []item.Item{item.NewItem(item.FLINT, 0, 1)}
		}
		return []item.Item{item.NewItem(GRAVEL, 0, 1)}

	case LEAVES:
		if rand.Float32() < 0.05 {

			return []item.Item{item.NewItem(SAPLING, int(meta&0x03), 1)}
		}
		return []item.Item{}

	case GLASS, GLASS_PANE:
		return []item.Item{}

	case TNT:
		return []item.Item{item.NewItem(TNT, 0, 1)}

	case CLAY_BLOCK:
		return []item.Item{item.NewItem(item.CLAY, 0, 4)}

	case GLOWSTONE_BLOCK:
		count := 2 + rand.Intn(3)
		return []item.Item{item.NewItem(item.GLOWSTONE_DUST, 0, count)}

	}

	return []item.Item{item.NewItem(int(id), int(meta), 1)}
}
