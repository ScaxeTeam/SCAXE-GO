package defaults

import (
	"fmt"
	"strconv"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type TeleportCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewTeleportCommand(server ServerInterface) *TeleportCommand {
	return &TeleportCommand{
		BaseCommand: command.BaseCommand{
			Name:        "tp",
			Description: "Teleports a player",
			Usage:       "/tp <x> <y> <z> or /tp <player> or /tp <player> <x> <y> <z>",
			Permission:  "pocketmine.command.teleport",
		},
		server: server,
	}
}

func (c *TeleportCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	switch len(args) {
	case 1:
		target := c.server.GetPlayerByName(args[0])
		if target == nil {
			sender.SendMessage("§cPlayer not found: " + args[0])
			return true
		}
		sender.SendMessage("§aTeleported to " + target.GetName())

	case 3:
		x, err1 := strconv.ParseFloat(args[0], 64)
		y, err2 := strconv.ParseFloat(args[1], 64)
		z, err3 := strconv.ParseFloat(args[2], 64)
		if err1 != nil || err2 != nil || err3 != nil {
			sender.SendMessage("§cInvalid coordinates")
			return true
		}
		sender.SendMessage(fmt.Sprintf("§aTeleported to %.2f, %.2f, %.2f", x, y, z))

	case 4:
		target := c.server.GetPlayerByName(args[0])
		if target == nil {
			sender.SendMessage("§cPlayer not found: " + args[0])
			return true
		}
		x, err1 := strconv.ParseFloat(args[1], 64)
		y, err2 := strconv.ParseFloat(args[2], 64)
		z, err3 := strconv.ParseFloat(args[3], 64)
		if err1 != nil || err2 != nil || err3 != nil {
			sender.SendMessage("§cInvalid coordinates")
			return true
		}
		sender.SendMessage(fmt.Sprintf("§aTeleported %s to %.2f, %.2f, %.2f", target.GetName(), x, y, z))

	default:
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	return true
}
