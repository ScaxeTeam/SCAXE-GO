package defaults

import (
	"github.com/scaxe/scaxe-go/pkg/command"
)

type DefaultGamemodeCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewDefaultGamemodeCommand(server ServerInterface) *DefaultGamemodeCommand {
	return &DefaultGamemodeCommand{
		BaseCommand: command.BaseCommand{
			Name:        "defaultgamemode",
			Description: "Sets the default gamemode for new players",
			Usage:       "/defaultgamemode <mode>",
			Permission:  "pocketmine.command.defaultgamemode",
		},
		server: server,
	}
}

func (c *DefaultGamemodeCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	gamemode := -1
	switch args[0] {
	case "0", "survival", "s":
		gamemode = 0
	case "1", "creative", "c":
		gamemode = 1
	case "2", "adventure", "a":
		gamemode = 2
	case "3", "spectator", "sp":
		gamemode = 3
	}

	if gamemode == -1 {
		sender.SendMessage("§cUnknown gamemode: " + args[0])
		return true
	}

	gamemodeNames := []string{"Survival", "Creative", "Adventure", "Spectator"}
	sender.SendMessage("§aDefault gamemode set to " + gamemodeNames[gamemode])

	return true
}
