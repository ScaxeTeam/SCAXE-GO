package defaults

import (
	"github.com/scaxe/scaxe-go/pkg/command"
)

var difficultyNames = map[string]int{
	"peaceful": 0, "p": 0, "0": 0,
	"easy": 1, "e": 1, "1": 1,
	"normal": 2, "n": 2, "2": 2,
	"hard": 3, "h": 3, "3": 3,
}

var difficultyLabels = []string{"Peaceful", "Easy", "Normal", "Hard"}

type DifficultyCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewDifficultyCommand(server ServerInterface) *DifficultyCommand {
	return &DifficultyCommand{
		BaseCommand: command.BaseCommand{
			Name:        "difficulty",
			Description: "Sets the server difficulty",
			Usage:       "/difficulty <0-3|peaceful|easy|normal|hard>",
			Permission:  "pocketmine.command.difficulty",
		},
		server: server,
	}
}

func (c *DifficultyCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§aCurrent difficulty: " + difficultyLabels[c.server.GetDifficulty()])
		return true
	}

	diff, ok := difficultyNames[args[0]]
	if !ok {
		sender.SendMessage("§cInvalid difficulty. Use 0-3 or peaceful/easy/normal/hard")
		return true
	}

	c.server.SetDifficulty(diff)
	c.server.BroadcastMessage("§aDifficulty set to " + difficultyLabels[diff])
	return true
}
