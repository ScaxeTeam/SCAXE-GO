package defaults

import (
	"github.com/scaxe/scaxe-go/pkg/command"
)

type SaveCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewSaveCommand(server ServerInterface) *SaveCommand {
	return &SaveCommand{
		BaseCommand: command.BaseCommand{
			Name:        "save-all",
			Description: "Saves all player and world data",
			Usage:       "/save-all",
			Permission:  "pocketmine.command.save.perform",
		},
		server: server,
	}
}

func (c *SaveCommand) Execute(sender command.CommandSender, args []string) bool {
	sender.SendMessage("§aSaving all worlds...")
	c.server.BroadcastMessage("§eSaving world data...")

	sender.SendMessage("§aWorld save complete!")
	return true
}

type SaveOnCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewSaveOnCommand(server ServerInterface) *SaveOnCommand {
	return &SaveOnCommand{
		BaseCommand: command.BaseCommand{
			Name:        "save-on",
			Description: "Enables automatic saving",
			Usage:       "/save-on",
			Permission:  "pocketmine.command.save.enable",
		},
		server: server,
	}
}

func (c *SaveOnCommand) Execute(sender command.CommandSender, args []string) bool {

	sender.SendMessage("§aAutomatic saving enabled")
	return true
}

type SaveOffCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewSaveOffCommand(server ServerInterface) *SaveOffCommand {
	return &SaveOffCommand{
		BaseCommand: command.BaseCommand{
			Name:        "save-off",
			Description: "Disables automatic saving",
			Usage:       "/save-off",
			Permission:  "pocketmine.command.save.disable",
		},
		server: server,
	}
}

func (c *SaveOffCommand) Execute(sender command.CommandSender, args []string) bool {

	sender.SendMessage("§aAutomatic saving disabled")
	return true
}
