package defaults

import (
	"fmt"
	"strconv"

	"github.com/scaxe/scaxe-go/pkg/command"
	"github.com/scaxe/scaxe-go/pkg/level"
)

type SetBlockCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewSetBlockCommand(server ServerInterface) *SetBlockCommand {
	return &SetBlockCommand{
		BaseCommand: command.BaseCommand{
			Name:        "setblock",
			Description: "Sets a block at a position",
			Usage:       "/setblock <x> <y> <z> <block> [data]",
			Permission:  "pocketmine.command.setblock",
		},
		server: server,
	}
}

func (c *SetBlockCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) < 4 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	x, y, z, ok := parseLocation(sender, args, 0)
	if !ok {
		sender.SendMessage("§cInvalid coordinates")
		return true
	}

	if y < 0 || y > 128 {
		sender.SendMessage("§cUnsuccessful: Y coordinate out of bounds (0-128)")
		return true
	}

	blockID, blockMeta, ok := ParseBlockArg(args[3])
	if !ok {
		sender.SendMessage("§cInvalid block ID")
		return true
	}

	if len(args) >= 5 {
		if m, err := strconv.Atoi(args[4]); err == nil {
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

	if targetLevel.SetBlock(int32(x), int32(y), int32(z), byte(blockID), byte(blockMeta), true) {
		sender.SendMessage(fmt.Sprintf("§aBlock placed at %.0f, %.0f, %.0f", x, y, z))
		return true
	}

	sender.SendMessage("§cUnsuccessful: Could not place block")
	return true
}

type LevelAware interface {
	GetLevel() interface{}
}
