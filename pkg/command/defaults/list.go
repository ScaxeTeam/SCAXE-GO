package defaults

import (
	"fmt"
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type ListCommand struct {
	command.BaseCommand
	Server Server
}

func NewListCommand(s Server) *ListCommand {
	return &ListCommand{
		BaseCommand: command.BaseCommand{
			Name:        "list",
			Description: "Lists all online players",
			Usage:       "/list",
			Permission:  "pocketmine.command.list",
		},
		Server: s,
	}
}

func (c *ListCommand) Execute(sender command.CommandSender, args []string) bool {
	online := []string{}
	count := 0

	for _, p := range c.Server.GetOnlinePlayers() {

		if p.IsOnline() {
			online = append(online, p.GetDisplayName())
			count++
		}
	}

	sender.SendMessage(fmt.Sprintf("There are %d/%d players online:", count, c.Server.GetMaxPlayers()))
	sender.SendMessage(strings.Join(online, ", "))
	return true
}
