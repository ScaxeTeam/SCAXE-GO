package defaults

import (
	"github.com/scaxe/scaxe-go/pkg/command"
)

type PingCommand struct {
	command.BaseCommand
}

func NewPingCommand() *PingCommand {
	return &PingCommand{
		BaseCommand: command.BaseCommand{
			Name:        "ping",
			Description: "Shows your ping latency",
			Usage:       "/ping",
			Permission:  "",
		},
	}
}

func (c *PingCommand) Execute(sender command.CommandSender, args []string) bool {
	sender.SendMessage("§aPong! Your latency will be shown here once implemented.")
	return true
}
