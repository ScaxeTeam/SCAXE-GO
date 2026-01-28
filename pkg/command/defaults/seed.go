package defaults

import (
	"strconv"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type SeedCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewSeedCommand(server ServerInterface) *SeedCommand {
	return &SeedCommand{
		BaseCommand: command.BaseCommand{
			Name:        "seed",
			Description: "Shows the world seed",
			Usage:       "/seed",
			Permission:  "pocketmine.command.seed",
		},
		server: server,
	}
}

func (c *SeedCommand) Execute(sender command.CommandSender, args []string) bool {
	sender.SendMessage("§aSeed: " + strconv.FormatInt(c.server.GetSeed(), 10))
	return true
}
