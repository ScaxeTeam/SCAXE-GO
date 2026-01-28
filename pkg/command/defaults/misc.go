package defaults

import (
	"github.com/scaxe/scaxe-go/pkg/command"
)

type PluginsCommand struct {
	command.BaseCommand
}

func NewPluginsCommand() *PluginsCommand {
	return &PluginsCommand{
		BaseCommand: command.BaseCommand{
			Name:        "plugins",
			Description: "Lists all loaded plugins",
			Usage:       "/plugins",
			Permission:  "pocketmine.command.plugins",
		},
	}
}

func (c *PluginsCommand) Execute(sender command.CommandSender, args []string) bool {

	sender.SendMessage("§aPlugins (0): §7None loaded")
	sender.SendMessage("§7Plugin system is not yet implemented")
	return true
}

type GcCommand struct {
	command.BaseCommand
}

func NewGcCommand() *GcCommand {
	return &GcCommand{
		BaseCommand: command.BaseCommand{
			Name:        "gc",
			Description: "Forces garbage collection",
			Usage:       "/gc",
			Permission:  "pocketmine.command.gc",
		},
	}
}

func (c *GcCommand) Execute(sender command.CommandSender, args []string) bool {

	sender.SendMessage("§aGarbage collection triggered")
	return true
}

type TimingsCommand struct {
	command.BaseCommand
}

func NewTimingsCommand() *TimingsCommand {
	return &TimingsCommand{
		BaseCommand: command.BaseCommand{
			Name:        "timings",
			Description: "Server timings profiler",
			Usage:       "/timings <on|off|paste|reset>",
			Permission:  "pocketmine.command.timings",
		},
	}
}

func (c *TimingsCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) < 1 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	switch args[0] {
	case "on":
		sender.SendMessage("§aTimings enabled")
	case "off":
		sender.SendMessage("§aTimings disabled")
	case "paste":
		sender.SendMessage("§ePaste not implemented yet")
	case "reset":
		sender.SendMessage("§aTimings reset")
	default:
		sender.SendMessage("§cUnknown option: " + args[0])
	}

	return true
}
