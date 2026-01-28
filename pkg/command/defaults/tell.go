package defaults

import (
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type TellCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewTellCommand(server ServerInterface) *TellCommand {
	return &TellCommand{
		BaseCommand: command.BaseCommand{
			Name:        "tell",
			Description: "Sends a private message to a player",
			Usage:       "/tell <player> <message>",
			Permission:  "",
		},
		server: server,
	}
}

func (c *TellCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) < 2 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	playerName := args[0]
	message := strings.Join(args[1:], " ")

	target := c.server.GetPlayerByName(playerName)
	if target == nil {
		sender.SendMessage("§cPlayer not found: " + playerName)
		return true
	}

	target.SendMessage("§7[" + sender.GetName() + " -> You] " + message)
	sender.SendMessage("§7[You -> " + target.GetName() + "] " + message)
	return true
}
