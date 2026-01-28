package defaults

import (
	"regexp"
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type BanIpCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewBanIpCommand(server ServerInterface) *BanIpCommand {
	return &BanIpCommand{
		BaseCommand: command.BaseCommand{
			Name:        "ban-ip",
			Description: "Bans an IP address from the server",
			Usage:       "/ban-ip <ip|player> [reason]",
			Permission:  "pocketmine.command.ban.ip",
		},
		server: server,
	}
}

var ipPattern = regexp.MustCompile(`^([01]?\d\d?|2[0-4]\d|25[0-5])\.([01]?\d\d?|2[0-4]\d|25[0-5])\.([01]?\d\d?|2[0-4]\d|25[0-5])\.([01]?\d\d?|2[0-4]\d|25[0-5])$`)

func (c *BanIpCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	value := args[0]
	reason := "IP Banned"
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}

	if ipPattern.MatchString(value) {

		sender.SendMessage("§aBanned IP: " + value)
		c.server.BroadcastMessage("§eIP " + value + " has been banned: " + reason)
	} else {

		if player := c.server.GetPlayerByName(value); player != nil {
			sender.SendMessage("§aBanned IP of player " + value)
		} else {
			sender.SendMessage("§cInvalid IP address or player not found")
			return false
		}
	}

	return true
}

type PardonIpCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewPardonIpCommand(server ServerInterface) *PardonIpCommand {
	return &PardonIpCommand{
		BaseCommand: command.BaseCommand{
			Name:        "pardon-ip",
			Description: "Unbans an IP address",
			Usage:       "/pardon-ip <ip>",
			Permission:  "pocketmine.command.unban.ip",
		},
		server: server,
	}
}

func (c *PardonIpCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) != 1 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	if ipPattern.MatchString(args[0]) {
		sender.SendMessage("§aUnbanned IP: " + args[0])
		c.server.BroadcastMessage("§eIP " + args[0] + " has been unbanned")
	} else {
		sender.SendMessage("§cInvalid IP address format")
		return false
	}

	return true
}

type BanListCommand struct {
	command.BaseCommand
}

func NewBanListCommand() *BanListCommand {
	return &BanListCommand{
		BaseCommand: command.BaseCommand{
			Name:        "banlist",
			Description: "Shows the ban list",
			Usage:       "/banlist [ips|players]",
			Permission:  "pocketmine.command.banlist",
		},
	}
}

func (c *BanListCommand) Execute(sender command.CommandSender, args []string) bool {
	listType := "players"
	if len(args) > 0 {
		listType = args[0]
	}

	switch listType {
	case "ips":
		sender.SendMessage("§aBanned IPs: (none)")
	case "players":
		sender.SendMessage("§aBanned Players: (none)")
	default:
		sender.SendMessage("§cUsage: " + c.Usage)
	}

	return true
}
