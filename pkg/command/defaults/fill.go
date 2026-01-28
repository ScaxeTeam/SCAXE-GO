package defaults

import (
	"fmt"
	"math"
	"strconv"

	"github.com/scaxe/scaxe-go/pkg/command"
	"github.com/scaxe/scaxe-go/pkg/level"
)

type FillCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewFillCommand(server ServerInterface) *FillCommand {
	return &FillCommand{
		BaseCommand: command.BaseCommand{
			Name:        "fill",
			Description: "Fills a region with blocks",
			Usage:       "/fill <x1> <y1> <z1> <x2> <y2> <z2> <block> [data]",
			Permission:  "pocketmine.command.fill",
		},
		server: server,
	}
}

func (c *FillCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) < 7 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	x1, y1, z1, ok1 := parseLocation(sender, args, 0)
	x2, y2, z2, ok2 := parseLocation(sender, args, 3)

	if !ok1 || !ok2 {
		sender.SendMessage("§cInvalid coordinates")
		return true
	}

	blockID, blockMeta, ok := ParseBlockArg(args[6])
	if !ok {
		sender.SendMessage("§cInvalid block ID")
		return true
	}

	if len(args) >= 8 {
		if m, err := strconv.Atoi(args[7]); err == nil {
			blockMeta = m
		}
	}

	var targetLevel *level.Level
	if la, ok := sender.(LevelAware); ok {
		if lvl, ok := la.GetLevel().(*level.Level); ok {
			targetLevel = lvl
		}
	}
	if targetLevel == nil {
		if defaultLvl, ok := c.server.GetLevelManager().GetDefaultLevel().(*level.Level); ok {
			targetLevel = defaultLvl
		}
	}
	if targetLevel == nil {
		sender.SendMessage("§cInternal Error: No target level found.")
		return true
	}

	minX, maxX := minMax(x1, x2)
	minY, maxY := minMax(y1, y2)
	minZ, maxZ := minMax(z1, z2)

	volume := (maxX - minX + 1) * (maxY - minY + 1) * (maxZ - minZ + 1)
	if volume > 32768 {
		sender.SendMessage(fmt.Sprintf("§cToo many blocks in the specified area (%.0f > 32768)", volume))
		return true
	}

	sender.SendMessage(fmt.Sprintf("§aFilling %.0f blocks...", volume))

	count := 0
	for x := minX; x <= maxX; x++ {
		for y := minY; y <= maxY; y++ {
			for z := minZ; z <= maxZ; z++ {
				if targetLevel.SetBlock(int32(x), int32(y), int32(z), byte(blockID), byte(blockMeta), false) {
					count++
				}
			}
		}
	}

	sender.SendMessage(fmt.Sprintf("§aSuccessfully filled %d blocks", count))
	return true
}

func minMax(a, b float64) (float64, float64) {
	if a < b {
		return math.Floor(a), math.Floor(b)
	}
	return math.Floor(b), math.Floor(a)
}
