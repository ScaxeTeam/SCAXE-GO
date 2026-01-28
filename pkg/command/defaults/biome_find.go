package defaults

import (
	"fmt"
	"strconv"

	"github.com/scaxe/scaxe-go/pkg/command"
	"github.com/scaxe/scaxe-go/pkg/level/generator/gorigional/layer"
	"github.com/scaxe/scaxe-go/pkg/level/generator/gorigional/structure"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
)

type LocateCommand struct {
	command.BaseCommand
	server ServerInterface
}

type BiomeInfo struct {
	ID   uint8
	Name string
}

var knownBiomes = []BiomeInfo{
	{0, "Ocean"},
	{1, "Plains"},
	{2, "Desert"},
	{3, "Extreme Hills"},
	{4, "Forest"},
	{5, "Taiga"},
	{6, "Swamp"},
	{7, "River"},
	{12, "Ice Plains"},
	{14, "Mushroom Island"},
	{16, "Beach"},
	{21, "Jungle"},
	{23, "Jungle Edge"},
	{24, "Deep Ocean"},
	{27, "Birch Forest"},
	{29, "Roofed Forest"},
	{30, "Cold Taiga"},
	{32, "Mega Taiga"},
	{35, "Savanna"},
	{36, "Savanna Plateau"},
	{37, "Mesa"},
	{38, "Mesa Plateau F"},
	{39, "Mesa Plateau"},
	{129, "Sunflower Plains"},
	{140, "Ice Plains Spikes"},
}

func NewBiomeFindCommand(server ServerInterface) *LocateCommand {
	return &LocateCommand{
		BaseCommand: command.BaseCommand{
			Name:        "bfind",
			Description: "Find and teleport to specific biomes or structures",
			Usage:       "/bfind - list | /bfind <num> - biome | /bfind village - village",
			Permission:  "scaxe.command.bfind",
		},
		server: server,
	}
}

func (c *LocateCommand) Execute(sender command.CommandSender, args []string) bool {
	seed := c.server.GetSeed()

	if len(args) == 0 {
		sender.SendMessage("§e=== Locate Command ===")
		sender.SendMessage("§7/bfind <num>    - Find biome by number")
		sender.SendMessage("§7/bfind village  - Find nearest village")
		sender.SendMessage("§7/bfind list     - List all biomes")
		return true
	}

	if args[0] == "list" {
		sender.SendMessage("§e=== Biome List ===")
		for i, b := range knownBiomes {
			sender.SendMessage(fmt.Sprintf("§a[%d]§f %s §7(ID %d)", i, b.Name, b.ID))
		}
		return true
	}

	if args[0] == "village" {
		sender.SendMessage("§eSearching for village...")

		found, x, z := findNearestVillage(seed, 0, 0, 5000)
		if found {
			y := 80
			sender.SendMessage(fmt.Sprintf("§aFound village at X=%d Z=%d!", x, z))
			sender.SendMessage(fmt.Sprintf("§7Use: §f/tp %d %d %d", x, y, z))

			if ps, ok := sender.(PlayerTeleportable); ok {
				ps.Teleport(float64(x), float64(y), float64(z))
				sender.SendMessage("§aTeleported!")
			}
		} else {
			sender.SendMessage("§cCould not find village within 5000 blocks")
		}
		return true
	}

	index, err := strconv.Atoi(args[0])
	if err != nil || index < 0 || index >= len(knownBiomes) {
		sender.SendMessage(fmt.Sprintf("§cInvalid. Use: 0-%d, 'village', or 'list'", len(knownBiomes)-1))
		return true
	}

	targetBiome := knownBiomes[index]
	sender.SendMessage(fmt.Sprintf("§eSearching for §a%s§e...", targetBiome.Name))

	found, x, z := findBiomeWithGenLayer(seed, targetBiome.ID, 0, 0, 10000)

	if found {
		y := 100
		sender.SendMessage(fmt.Sprintf("§aFound %s at X=%d Z=%d!", targetBiome.Name, x, z))
		sender.SendMessage(fmt.Sprintf("§7Use: §f/tp %d %d %d", x, y, z))

		if ps, ok := sender.(PlayerTeleportable); ok {
			ps.Teleport(float64(x), float64(y), float64(z))
			sender.SendMessage("§aTeleported!")
		}
	} else {
		sender.SendMessage(fmt.Sprintf("§cCould not find %s within 10000 blocks", targetBiome.Name))
	}

	return true
}

type PlayerTeleportable interface {
	command.CommandSender
	Teleport(x, y, z float64)
}

func findBiomeWithGenLayer(seed int64, targetID uint8, startX, startZ, maxDist int) (bool, int, int) {

	layers := layer.InitializeAll(seed)
	if len(layers) == 0 {
		return false, 0, 0
	}

	genLayer := layers[0]

	step := 32
	for dist := 0; dist < maxDist; dist += step {

		for offset := -dist; offset <= dist; offset += step {
			positions := [][2]int{
				{startX + offset, startZ - dist},
				{startX + offset, startZ + dist},
				{startX - dist, startZ + offset},
				{startX + dist, startZ + offset},
			}

			for _, pos := range positions {
				x, z := pos[0], pos[1]

				ids := genLayer.GetInts(x, z, 1, 1)
				if len(ids) > 0 && uint8(ids[0]) == targetID {
					return true, x, z
				}
			}
		}
	}

	return false, 0, 0
}

func findNearestVillage(seed int64, startX, startZ, maxDist int) (bool, int, int) {

	distance := 32

	startChunkX := startX >> 4
	startChunkZ := startZ >> 4
	maxChunkDist := maxDist >> 4

	rnd := rand.NewRandom(seed)

	for chunkDist := 0; chunkDist < maxChunkDist; chunkDist += 8 {
		for offsetX := -chunkDist; offsetX <= chunkDist; offsetX += 8 {
			for offsetZ := -chunkDist; offsetZ <= chunkDist; offsetZ += 8 {

				if offsetX != -chunkDist && offsetX != chunkDist &&
					offsetZ != -chunkDist && offsetZ != chunkDist {
					continue
				}

				chunkX := startChunkX + offsetX
				chunkZ := startChunkZ + offsetZ

				if canVillageSpawnAt(seed, chunkX, chunkZ, distance, rnd) {

					return true, chunkX*16 + 8, chunkZ*16 + 8
				}
			}
		}
	}

	return false, 0, 0
}

func canVillageSpawnAt(seed int64, chunkX, chunkZ, distance int, rnd *rand.Random) bool {
	i := chunkX
	j := chunkZ

	if chunkX < 0 {
		chunkX -= distance - 1
	}
	if chunkZ < 0 {
		chunkZ -= distance - 1
	}

	k := chunkX / distance
	l := chunkZ / distance

	villageSeed := int64(k)*341873128712 + int64(l)*132897987541 + seed + 10387312
	rnd.SetSeed(villageSeed)

	k *= distance
	l *= distance

	k += rnd.NextBoundedInt(distance - 8)
	l += rnd.NextBoundedInt(distance - 8)

	return i == k && j == l
}

type VillageLocateCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewVillageLocateCommand(server ServerInterface) *VillageLocateCommand {
	return &VillageLocateCommand{
		BaseCommand: command.BaseCommand{
			Name:        "locate",
			Description: "Locate structures like villages",
			Usage:       "/locate village",
			Permission:  "scaxe.command.locate",
		},
		server: server,
	}
}

func (c *VillageLocateCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: /locate village")
		return true
	}

	seed := c.server.GetSeed()

	switch args[0] {
	case "village":
		sender.SendMessage("§eLocating nearest village...")
		found, x, z := findNearestVillage(seed, 0, 0, 10000)
		if found {
			sender.SendMessage(fmt.Sprintf("§aVillage found at §f%d, %d", x, z))
			sender.SendMessage(fmt.Sprintf("§7Use: /tp %d 80 %d", x, z))

			if ps, ok := sender.(PlayerTeleportable); ok {
				ps.Teleport(float64(x), 80, float64(z))
				sender.SendMessage("§aTeleported!")
			}
		} else {
			sender.SendMessage("§cNo village found within 10000 blocks")
		}

	case "temple":
		sender.SendMessage("§eLocating nearest temple...")
		found, x, z := findNearestTemple(seed, 0, 0, 10000)
		if found {
			sender.SendMessage(fmt.Sprintf("§aTemple found at §f%d, %d", x, z))
			sender.SendMessage(fmt.Sprintf("§7Use: /tp %d 80 %d", x, z))

			if ps, ok := sender.(PlayerTeleportable); ok {
				ps.Teleport(float64(x), 80, float64(z))
				sender.SendMessage("§aTeleported!")
			}
		} else {
			sender.SendMessage("§cNo temple found within 10000 blocks")
		}

	case "stronghold":
		sender.SendMessage("§eLocating stronghold...")
		found, x, z := findNearestStronghold(seed)
		if found {
			sender.SendMessage(fmt.Sprintf("§aStronghold found at §f%d, %d", x, z))
			sender.SendMessage(fmt.Sprintf("§7Use: /tp %d 40 %d", x, z))
		} else {
			sender.SendMessage("§cNo stronghold found")
		}

	default:
		sender.SendMessage("§cUnknown structure. Available: village, temple, stronghold")
	}

	return true
}

func findNearestTemple(seed int64, startX, startZ, maxDist int) (bool, int, int) {

	distance := 32
	minDistance := 8

	startChunkX := startX >> 4
	startChunkZ := startZ >> 4
	maxChunkDist := maxDist >> 4

	rnd := rand.NewRandom(seed)

	for chunkDist := 0; chunkDist < maxChunkDist; chunkDist += 8 {
		for offsetX := -chunkDist; offsetX <= chunkDist; offsetX += 8 {
			for offsetZ := -chunkDist; offsetZ <= chunkDist; offsetZ += 8 {
				if offsetX != -chunkDist && offsetX != chunkDist &&
					offsetZ != -chunkDist && offsetZ != chunkDist {
					continue
				}

				chunkX := startChunkX + offsetX
				chunkZ := startChunkZ + offsetZ

				if canTempleSpawnAt(seed, chunkX, chunkZ, distance, minDistance, rnd) {
					return true, chunkX*16 + 8, chunkZ*16 + 8
				}
			}
		}
	}

	return false, 0, 0
}

func canTempleSpawnAt(seed int64, chunkX, chunkZ, distance, minDistance int, rnd *rand.Random) bool {
	i := chunkX
	j := chunkZ

	if chunkX < 0 {
		chunkX -= distance - 1
	}
	if chunkZ < 0 {
		chunkZ -= distance - 1
	}

	k := chunkX / distance
	l := chunkZ / distance

	templeSeed := int64(k)*341873128712 + int64(l)*132897987541 + seed + 14357617
	rnd.SetSeed(templeSeed)

	k *= distance
	l *= distance

	k += rnd.NextBoundedInt(distance - minDistance)
	l += rnd.NextBoundedInt(distance - minDistance)

	return i == k && j == l
}

func findNearestStronghold(seed int64) (bool, int, int) {

	sh := structure.NewMapGenStronghold(seed)
	positions := sh.GetStructureCoords()
	if len(positions) > 0 {

		return true, int(positions[0].X) * 16, int(positions[0].Z) * 16
	}
	return false, 0, 0
}

func RegisterLocateCommands(server ServerInterface, cmdMap *command.CommandMap) {
	cmdMap.Register(NewBiomeFindCommand(server))
	cmdMap.Register(NewVillageLocateCommand(server))
}
