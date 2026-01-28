package defaults

import (
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type SayCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewSayCommand(server ServerInterface) *SayCommand {
	return &SayCommand{
		BaseCommand: command.BaseCommand{
			Name:        "say",
			Description: "Broadcasts a message to all players",
			Usage:       "/say <message>",
			Permission:  "pocketmine.command.say",
		},
		server: server,
	}
}

func (c *SayCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	message := strings.Join(args, " ")
	c.server.BroadcastMessage("§d[Server] " + message)
	return true
}
