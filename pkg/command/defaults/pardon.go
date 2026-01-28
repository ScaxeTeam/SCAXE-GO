package defaults

import (
	"github.com/scaxe/scaxe-go/pkg/command"
)

type PardonCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewPardonCommand(server ServerInterface) *PardonCommand {
	return &PardonCommand{
		BaseCommand: command.BaseCommand{
			Name:        "pardon",
			Description: "Unbans a player from the server",
			Usage:       "/pardon <player>",
			Permission:  "pocketmine.command.unban.player",
		},
		server: server,
	}
}

func (c *PardonCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) != 1 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	playerName := args[0]

	sender.SendMessage("§aUnbanned §e" + playerName)
	c.server.BroadcastMessage("§e" + playerName + " has been unbanned")

	return true
}
