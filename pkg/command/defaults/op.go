package defaults

import (
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type OpCommand struct {
	command.BaseCommand
	Server Server
}

func NewOpCommand(s Server) *OpCommand {
	return &OpCommand{
		BaseCommand: command.BaseCommand{
			Name:        "op",
			Description: "Grants operator status to a player.",
			Usage:       "/op <player>",
			Permission:  "",
		},
		Server: s,
	}
}

func (c *OpCommand) Execute(sender command.CommandSender, args []string) bool {
	if sender.GetName() != "CONSOLE" {
		sender.SendMessage("§cThis command can only be used from the console.")
		return true
	}

	if len(args) < 1 {
		sender.SendMessage("Usage: /op <player>")
		return true
	}

	targetName := args[0]
	playerFound := false

	for _, p := range c.Server.GetOnlinePlayers() {
		if strings.EqualFold(p.GetName(), targetName) {

			c.Server.AddOp(p.GetName(), int64(p.ClientID))
			p.SetOp(true)
			p.SendMessage("§7You have been slightly Oped")
			sender.SendMessage("Oped " + p.GetName())
			playerFound = true
			break
		}
	}

	if !playerFound {
		sender.SendMessage("Player not found online. Cannot verify CID for strict OP.")
	}

	return true
}
