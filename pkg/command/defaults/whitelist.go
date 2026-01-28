package defaults

import (
	"github.com/scaxe/scaxe-go/pkg/command"
)

type WhitelistCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewWhitelistCommand(server ServerInterface) *WhitelistCommand {
	return &WhitelistCommand{
		BaseCommand: command.BaseCommand{
			Name:        "whitelist",
			Description: "Manages the server whitelist",
			Usage:       "/whitelist <on|off|add|remove|list|reload> [player]",
			Permission:  "pocketmine.command.whitelist",
		},
		server: server,
	}
}

func (c *WhitelistCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) < 1 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	switch args[0] {
	case "on":

		sender.SendMessage("§aWhitelist enabled")
	case "off":

		sender.SendMessage("§aWhitelist disabled")
	case "add":
		if len(args) < 2 {
			sender.SendMessage("§cUsage: /whitelist add <player>")
			return false
		}

		sender.SendMessage("§aAdded " + args[1] + " to whitelist")
	case "remove":
		if len(args) < 2 {
			sender.SendMessage("§cUsage: /whitelist remove <player>")
			return false
		}

		sender.SendMessage("§aRemoved " + args[1] + " from whitelist")
	case "list":

		sender.SendMessage("§aWhitelist: (none)")
	case "reload":

		sender.SendMessage("§aWhitelist reloaded")
	default:
		sender.SendMessage("§cUnknown whitelist command: " + args[0])
		return false
	}

	return true
}
