package defaults

import (
	"strconv"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type TimeCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewTimeCommand(server ServerInterface) *TimeCommand {
	return &TimeCommand{
		BaseCommand: command.BaseCommand{
			Name:        "time",
			Description: "Changes or queries the world time",
			Usage:       "/time <set|add|query> <value>",
			Permission:  "pocketmine.command.time",
		},
		server: server,
	}
}

func (c *TimeCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	switch args[0] {
	case "set":
		if len(args) < 2 {
			sender.SendMessage("§cUsage: /time set <value|day|night|noon|midnight>")
			return false
		}
		var time int
		switch args[1] {
		case "day":
			time = 1000
		case "night":
			time = 13000
		case "noon":
			time = 6000
		case "midnight":
			time = 18000
		default:
			var err error
			time, err = strconv.Atoi(args[1])
			if err != nil {
				sender.SendMessage("§cInvalid time value")
				return true
			}
		}
		c.server.SetTime(time)
		sender.SendMessage("§aSet time to " + strconv.Itoa(time))

	case "add":
		if len(args) < 2 {
			sender.SendMessage("§cUsage: /time add <value>")
			return false
		}
		add, err := strconv.Atoi(args[1])
		if err != nil {
			sender.SendMessage("§cInvalid time value")
			return true
		}
		newTime := c.server.GetTime() + add
		c.server.SetTime(newTime)
		sender.SendMessage("§aAdded " + strconv.Itoa(add) + " to time")

	case "query":
		sender.SendMessage("§aCurrent time: " + strconv.Itoa(c.server.GetTime()))

	default:
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	return true
}
