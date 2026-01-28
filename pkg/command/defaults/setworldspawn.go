package defaults

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/command"
	"github.com/scaxe/scaxe-go/pkg/level"
)

type SetWorldSpawnCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewSetWorldSpawnCommand(server ServerInterface) *SetWorldSpawnCommand {
	return &SetWorldSpawnCommand{
		BaseCommand: command.BaseCommand{
			Name:        "setworldspawn",
			Description: "Sets the world spawn point",
			Usage:       "/setworldspawn [x y z]",
			Permission:  "pocketmine.command.setworldspawn",
		},
		server: server,
	}
}

func (c *SetWorldSpawnCommand) Execute(sender command.CommandSender, args []string) bool {
	var x, y, z float64
	var ok bool

	if len(args) == 0 {
		if posSender, isPos := sender.(Positional); isPos {
			pos := posSender.GetPosition()
			x, y, z = pos.X, pos.Y, pos.Z
		} else {
			sender.SendMessage("§cPlease provide coordinates.")
			return true
		}
	} else if len(args) >= 3 {
		x, y, z, ok = parseLocation(sender, args, 0)
		if !ok {
			sender.SendMessage("§cInvalid coordinates")
			return true
		}
	} else {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
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

	sender.SendMessage(fmt.Sprintf("§aSet world spawn point to %.1f, %.1f, %.1f", x, y, z))
	return true
}
