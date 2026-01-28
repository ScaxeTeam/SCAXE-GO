package defaults

import (
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type DeopCommand struct {
	command.BaseCommand
	Server Server
}

func NewDeopCommand(s Server) *DeopCommand {
	return &DeopCommand{
		BaseCommand: command.BaseCommand{
			Name:        "deop",
			Description: "Revokes operator status from a player.",
			Usage:       "/deop <player>",
			Permission:  "",
		},
		Server: s,
	}
}

func (c *DeopCommand) Execute(sender command.CommandSender, args []string) bool {
	if sender.GetName() != "CONSOLE" {
		sender.SendMessage("§cThis command can only be used from the console.")
		return true
	}

	if len(args) < 1 {
		sender.SendMessage("Usage: /deop <player>")
		return true
	}

	targetName := args[0]
	c.Server.RemoveOp(targetName)

	for _, p := range c.Server.GetOnlinePlayers() {
		if strings.EqualFold(p.GetName(), targetName) {
			p.SetOp(false)
			p.SendMessage("§7You are no longer op!")
			break
		}
	}

	sender.SendMessage("De-oped " + targetName)
	return true
}
