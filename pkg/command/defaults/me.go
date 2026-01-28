package defaults

import (
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type MeCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewMeCommand(server ServerInterface) *MeCommand {
	return &MeCommand{
		BaseCommand: command.BaseCommand{
			Name:        "me",
			Description: "Broadcasts an action message",
			Usage:       "/me <action>",
			Permission:  "",
		},
		server: server,
	}
}

func (c *MeCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	message := strings.Join(args, " ")
	c.server.BroadcastMessage("§5* " + sender.GetName() + " " + message)
	return true
}
