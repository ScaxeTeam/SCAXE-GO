package defaults

import (
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type BanServerInterface interface {
	ServerInterface
	AddBan(name, reason, source string)
	RemoveBan(name string)
	IsBanned(name string) bool
}

type BanCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewBanCommand(server ServerInterface) *BanCommand {
	return &BanCommand{
		BaseCommand: command.BaseCommand{
			Name:        "ban",
			Description: "Bans a player from the server",
			Usage:       "/ban <player> [reason]",
			Permission:  "pocketmine.command.ban.player",
		},
		server: server,
	}
}

func (c *BanCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	playerName := args[0]
	reason := "Banned by operator"
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}

	sender.SendMessage("§aBanned §e" + playerName + "§a: " + reason)
	c.server.BroadcastMessage("§e" + playerName + " has been banned: " + reason)

	if player := c.server.GetPlayerByName(playerName); player != nil {
		player.SendMessage("§cYou have been banned: " + reason)
	}

	return true
}
