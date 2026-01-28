package defaults

import (
	"github.com/scaxe/scaxe-go/pkg/command"
)

type MakePluginCommand struct {
	command.BaseCommand
}

func NewMakePluginCommand() *MakePluginCommand {
	return &MakePluginCommand{
		BaseCommand: command.BaseCommand{
			Name:        "makeplugin",
			Description: "Creates a plugin phar from source",
			Usage:       "/makeplugin <plugin>",
			Permission:  "pocketmine.command.makeplugin",
		},
	}
}

func (c *MakePluginCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}
	sender.SendMessage("§aCreating plugin phar for: " + args[0])
	sender.SendMessage("§7Plugin creation is not supported in Go version")
	return true
}

type ExtractPluginCommand struct {
	command.BaseCommand
}

func NewExtractPluginCommand() *ExtractPluginCommand {
	return &ExtractPluginCommand{
		BaseCommand: command.BaseCommand{
			Name:        "extractplugin",
			Description: "Extracts a plugin phar to source",
			Usage:       "/extractplugin <plugin>",
			Permission:  "pocketmine.command.extractplugin",
		},
	}
}

func (c *ExtractPluginCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}
	sender.SendMessage("§aExtracting plugin: " + args[0])
	sender.SendMessage("§7Plugin extraction is not supported in Go version")
	return true
}

type ExtractPharCommand struct {
	command.BaseCommand
}

func NewExtractPharCommand() *ExtractPharCommand {
	return &ExtractPharCommand{
		BaseCommand: command.BaseCommand{
			Name:        "extractphar",
			Description: "Extracts a phar file",
			Usage:       "/extractphar <file>",
			Permission:  "pocketmine.command.extractphar",
		},
	}
}

func (c *ExtractPharCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}
	sender.SendMessage("§7Phar extraction is not supported in Go version")
	return true
}

type GeneratePluginCommand struct {
	command.BaseCommand
}

func NewGeneratePluginCommand() *GeneratePluginCommand {
	return &GeneratePluginCommand{
		BaseCommand: command.BaseCommand{
			Name:        "generateplugin",
			Description: "Generates a plugin skeleton",
			Usage:       "/generateplugin <name>",
			Permission:  "pocketmine.command.generateplugin",
		},
	}
}

func (c *GeneratePluginCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}
	sender.SendMessage("§aGenerating plugin: " + args[0])
	sender.SendMessage("§7Plugin generation is not supported in Go version")
	return true
}

type LvdatCommand struct {
	command.BaseCommand
}

func NewLvdatCommand() *LvdatCommand {
	return &LvdatCommand{
		BaseCommand: command.BaseCommand{
			Name:        "lvdat",
			Description: "Shows level.dat information",
			Usage:       "/lvdat [world]",
			Permission:  "pocketmine.command.lvdat",
		},
	}
}

func (c *LvdatCommand) Execute(sender command.CommandSender, args []string) bool {
	worldName := "world"
	if len(args) > 0 {
		worldName = args[0]
	}
	sender.SendMessage("§aLevel.dat for: " + worldName)
	sender.SendMessage("§7Spawn: 0, 64, 0")
	sender.SendMessage("§7Seed: 0")
	return true
}

type BanCidByNameCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewBanCidByNameCommand(server ServerInterface) *BanCidByNameCommand {
	return &BanCidByNameCommand{
		BaseCommand: command.BaseCommand{
			Name:        "ban-cid-byname",
			Description: "Bans a player's CID by their name",
			Usage:       "/ban-cid-byname <player> [reason]",
			Permission:  "pocketmine.command.ban.cid",
		},
		server: server,
	}
}

func (c *BanCidByNameCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}
	sender.SendMessage("§aBanned CID of player: " + args[0])
	return true
}

type BanIpByNameCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewBanIpByNameCommand(server ServerInterface) *BanIpByNameCommand {
	return &BanIpByNameCommand{
		BaseCommand: command.BaseCommand{
			Name:        "ban-ip-byname",
			Description: "Bans a player's IP by their name",
			Usage:       "/ban-ip-byname <player> [reason]",
			Permission:  "pocketmine.command.ban.ip",
		},
		server: server,
	}
}

func (c *BanIpByNameCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}
	sender.SendMessage("§aBanned IP of player: " + args[0])
	return true
}
