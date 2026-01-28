package defaults

import (
	"strconv"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type XpCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewXpCommand(server ServerInterface) *XpCommand {
	return &XpCommand{
		BaseCommand: command.BaseCommand{
			Name:        "xp",
			Description: "Adds or removes player experience",
			Usage:       "/xp <amount[L]> <player>",
			Permission:  "pocketmine.command.xp",
		},
		server: server,
	}
}

func (c *XpCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) < 2 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	amountArg := args[0]
	playerName := args[1]

	target := c.server.GetPlayerByName(playerName)
	if target == nil {
		sender.SendMessage("§cPlayer not found: " + playerName)
		return true
	}

	isLevels := false
	if len(amountArg) > 0 && (amountArg[len(amountArg)-1] == 'L' || amountArg[len(amountArg)-1] == 'l') {
		isLevels = true
		amountArg = amountArg[:len(amountArg)-1]
	}

	amount, err := strconv.Atoi(amountArg)
	if err != nil {
		sender.SendMessage("§cInvalid amount: " + args[0])
		return true
	}

	if isLevels {
		sender.SendMessage("§aGave " + strconv.Itoa(amount) + " levels to " + target.GetName())
	} else {
		sender.SendMessage("§aGave " + strconv.Itoa(amount) + " experience to " + target.GetName())
	}

	return true
}
