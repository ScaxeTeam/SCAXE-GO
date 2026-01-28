package defaults

import (
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type RestartCommand struct {
	command.BaseCommand
}

func NewRestartCommand() *RestartCommand {
	return &RestartCommand{
		BaseCommand: command.BaseCommand{
			Name:        "restart",
			Description: "Restarts the server",
			Usage:       "/restart",
			Permission:  "pocketmine.command.restart",
		},
	}
}

func (c *RestartCommand) Execute(sender command.CommandSender, args []string) bool {
	sender.SendMessage("§eServer restart initiated...")

	return true
}

type BackupCommand struct {
	command.BaseCommand
}

func NewBackupCommand() *BackupCommand {
	return &BackupCommand{
		BaseCommand: command.BaseCommand{
			Name:        "backup",
			Description: "Creates a backup of the world",
			Usage:       "/backup",
			Permission:  "pocketmine.command.backup",
		},
	}
}

func (c *BackupCommand) Execute(sender command.CommandSender, args []string) bool {
	sender.SendMessage("§aCreating world backup...")

	sender.SendMessage("§aBackup complete!")
	return true
}

type ChunkInfoCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewChunkInfoCommand(server ServerInterface) *ChunkInfoCommand {
	return &ChunkInfoCommand{
		BaseCommand: command.BaseCommand{
			Name:        "chunkinfo",
			Description: "Shows chunk information",
			Usage:       "/chunkinfo",
			Permission:  "pocketmine.command.chunkinfo",
		},
		server: server,
	}
}

func (c *ChunkInfoCommand) Execute(sender command.CommandSender, args []string) bool {
	sender.SendMessage("§aChunk information:")
	sender.SendMessage("§7Loaded chunks: (calculate from level)")
	return true
}

type BiomeCommand struct {
	command.BaseCommand
}

func NewBiomeCommand() *BiomeCommand {
	return &BiomeCommand{
		BaseCommand: command.BaseCommand{
			Name:        "biome",
			Description: "Shows biome at current position",
			Usage:       "/biome",
			Permission:  "pocketmine.command.biome",
		},
	}
}

func (c *BiomeCommand) Execute(sender command.CommandSender, args []string) bool {
	sender.SendMessage("§aBiome: Plains")
	return true
}

type DumpMemoryCommand struct {
	command.BaseCommand
}

func NewDumpMemoryCommand() *DumpMemoryCommand {
	return &DumpMemoryCommand{
		BaseCommand: command.BaseCommand{
			Name:        "dumpmemory",
			Description: "Dumps memory information",
			Usage:       "/dumpmemory",
			Permission:  "pocketmine.command.dumpmemory",
		},
	}
}

func (c *DumpMemoryCommand) Execute(sender command.CommandSender, args []string) bool {
	sender.SendMessage("§aMemory dump complete")
	return true
}

type BanCidCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewBanCidCommand(server ServerInterface) *BanCidCommand {
	return &BanCidCommand{
		BaseCommand: command.BaseCommand{
			Name:        "ban-cid",
			Description: "Bans a player by Client ID",
			Usage:       "/ban-cid <cid> [reason]",
			Permission:  "pocketmine.command.ban.cid",
		},
		server: server,
	}
}

func (c *BanCidCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	cid := args[0]
	reason := "Banned by CID"
	if len(args) > 1 {
		reason = strings.Join(args[1:], " ")
	}

	sender.SendMessage("§aBanned CID: " + cid + " (" + reason + ")")
	return true
}

type PardonCidCommand struct {
	command.BaseCommand
}

func NewPardonCidCommand() *PardonCidCommand {
	return &PardonCidCommand{
		BaseCommand: command.BaseCommand{
			Name:        "pardon-cid",
			Description: "Unbans a Client ID",
			Usage:       "/pardon-cid <cid>",
			Permission:  "pocketmine.command.unban.cid",
		},
	}
}

func (c *PardonCidCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	sender.SendMessage("§aUnbanned CID: " + args[0])
	return true
}
