package defaults

import (
	"github.com/scaxe/scaxe-go/pkg/command"
)

type KillCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewKillCommand(server ServerInterface) *KillCommand {
	return &KillCommand{
		BaseCommand: command.BaseCommand{
			Name:        "kill",
			Description: "Kills a player",
			Usage:       "/kill [player]",
			Permission:  "pocketmine.command.kill",
		},
		server: server,
	}
}

func (c *KillCommand) Execute(sender command.CommandSender, args []string) bool {
	var targetName string

	if len(args) > 0 {
		targetName = args[0]
	} else {
		targetName = sender.GetName()
	}

	target := c.server.GetPlayerByName(targetName)
	if target == nil {
		sender.SendMessage("§cPlayer not found: " + targetName)
		return true
	}

	target.SendMessage("§cYou have been killed")
	sender.SendMessage("§aKilled " + target.GetName())
	return true
}
