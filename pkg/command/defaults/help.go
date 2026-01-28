package defaults

import (
	"github.com/scaxe/scaxe-go/pkg/command"
)

type HelpCommand struct {
	command.BaseCommand
}

func NewHelpCommand() *HelpCommand {
	return &HelpCommand{
		BaseCommand: command.BaseCommand{
			Name:        "help",
			Description: "Shows help for commands",
			Usage:       "/help [command]",
			Permission:  "",
		},
	}
}

func (c *HelpCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) > 0 {
		sender.SendMessage("§eHelp for /" + args[0])
		return true
	}

	sender.SendMessage("§e--- Available Commands ---")
	sender.SendMessage("§7/help - Shows this help message")
	sender.SendMessage("§7/version - Shows server version")
	sender.SendMessage("§7/list - Lists online players")
	sender.SendMessage("§7/stop - Stops the server")
	sender.SendMessage("§7/gamemode - Changes game mode")
	sender.SendMessage("§7/kick - Kicks a player")
	sender.SendMessage("§7/kill - Kills a player")
	sender.SendMessage("§7/tp - Teleports a player")
	sender.SendMessage("§7/time - Changes the time")
	sender.SendMessage("§7/say - Broadcasts a message")
	sender.SendMessage("§7/tell - Sends a private message")
	sender.SendMessage("§7/me - Broadcasts an action")
	sender.SendMessage("§7/op - Gives operator status")
	sender.SendMessage("§7/deop - Removes operator status")
	sender.SendMessage("§7/difficulty - Sets difficulty")
	sender.SendMessage("§7/seed - Shows world seed")
	sender.SendMessage("§7/tps - Shows server TPS")
	sender.SendMessage("§7/ping - Shows your ping")
	sender.SendMessage("§e--------------------------")
	return true
}
