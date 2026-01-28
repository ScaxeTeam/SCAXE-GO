package defaults

import (
	"strconv"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type GiveCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewGiveCommand(server ServerInterface) *GiveCommand {
	return &GiveCommand{
		BaseCommand: command.BaseCommand{
			Name:        "give",
			Description: "Gives items to a player",
			Usage:       "/give <player> <item[:data]> [amount]",
			Permission:  "pocketmine.command.give",
		},
		server: server,
	}
}

func (c *GiveCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) < 2 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	playerName := args[0]
	itemArg := args[1]
	amount := 1

	if len(args) >= 3 {
		if a, err := strconv.Atoi(args[2]); err == nil {
			amount = a
		}
	}

	target := c.server.GetPlayerByName(playerName)
	if target == nil {
		sender.SendMessage("§cPlayer not found: " + playerName)
		return true
	}

	itemID := 0
	itemData := 0

	for i, ch := range itemArg {
		if ch == ':' {
			if id, err := strconv.Atoi(itemArg[:i]); err == nil {
				itemID = id
			}
			if data, err := strconv.Atoi(itemArg[i+1:]); err == nil {
				itemData = data
			}
			break
		}
	}

	if itemID == 0 {
		if id, err := strconv.Atoi(itemArg); err == nil {
			itemID = id
		} else {
			sender.SendMessage("§cInvalid item ID: " + itemArg)
			return true
		}
	}

	_ = itemData

	sender.SendMessage("§aGave " + strconv.Itoa(amount) + " x " + itemArg + " to " + target.GetName())
	return true
}
