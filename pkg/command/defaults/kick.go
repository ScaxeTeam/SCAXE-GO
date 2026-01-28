package defaults

import (
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type KickCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewKickCommand(server ServerInterface) *KickCommand {
	return &KickCommand{
		BaseCommand: command.BaseCommand{
			Name:        "kick",
			Description: "Kicks a player from the server",
			Usage:       "/kick <player> [reason]",
			Permission:  "pocketmine.command.kick",
		},
		server: server,
	}
}

func (c *KickCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	playerName := args[0]
	reason := "Kicked by operator"
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}

	player := c.server.GetPlayerByName(playerName)
	if player == nil {
		sender.SendMessage("§cPlayer not found: " + playerName)
		return true
	}

	player.SendMessage("§cYou have been kicked: " + reason)
	c.server.BroadcastMessage("§e" + player.GetName() + " has been kicked: " + reason)
	return true
}
