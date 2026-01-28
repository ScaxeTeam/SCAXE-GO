package defaults

import (
	"github.com/scaxe/scaxe-go/internal/version"
	"github.com/scaxe/scaxe-go/pkg/command"
)

type VersionCommand struct {
	command.BaseCommand
}

func NewVersionCommand() *VersionCommand {
	return &VersionCommand{
		BaseCommand: command.BaseCommand{
			Name:        "ver",
			Description: "Gets the version of this server including any plugins in use",
			Usage:       "/ver",
			Permission:  "pocketmine.command.version",
		},
	}
}

func (c *VersionCommand) Execute(sender command.CommandSender, args []string) bool {
	sender.SendMessage("This server is running " + version.Full())
	return true
}
