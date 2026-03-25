package defaults

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type PingCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewPingCommand(server ServerInterface) *PingCommand {
	return &PingCommand{
		BaseCommand: command.BaseCommand{
			Name:        "ping",
			Description: "Shows your ping latency",
			Usage:       "/ping",
			Permission:  "",
		},
		server: server,
	}
}

func (c *PingCommand) Execute(sender command.CommandSender, args []string) bool {
	ps, ok := sender.(command.PlayerSender)
	if !ok {
		sender.SendMessage("§cThis command can only be used by players")
		return true
	}

	for _, p := range c.server.GetOnlinePlayers() {
		if p.GetEntityID() == ps.GetEntityID() {
			sender.SendMessage(fmt.Sprintf("§aPong! %dms", p.Ping))
			return true
		}
	}

	sender.SendMessage("§aPong!")
	return true
}
