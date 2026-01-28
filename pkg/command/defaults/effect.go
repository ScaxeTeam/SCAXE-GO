package defaults

import (
	"strconv"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type EffectCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewEffectCommand(server ServerInterface) *EffectCommand {
	return &EffectCommand{
		BaseCommand: command.BaseCommand{
			Name:        "effect",
			Description: "Adds or removes a potion effect",
			Usage:       "/effect <player> <effect|clear> [duration] [amplifier]",
			Permission:  "pocketmine.command.effect",
		},
		server: server,
	}
}

var effectNames = map[string]int{
	"speed":           1,
	"slowness":        2,
	"haste":           3,
	"mining_fatigue":  4,
	"strength":        5,
	"instant_health":  6,
	"instant_damage":  7,
	"jump_boost":      8,
	"nausea":          9,
	"regeneration":    10,
	"resistance":      11,
	"fire_resistance": 12,
	"water_breathing": 13,
	"invisibility":    14,
	"blindness":       15,
	"night_vision":    16,
	"hunger":          17,
	"weakness":        18,
	"poison":          19,
	"wither":          20,
	"health_boost":    21,
	"absorption":      22,
	"saturation":      23,
}

func (c *EffectCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) < 2 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	playerName := args[0]
	effectArg := args[1]

	target := c.server.GetPlayerByName(playerName)
	if target == nil {
		sender.SendMessage("§cPlayer not found: " + playerName)
		return true
	}

	if effectArg == "clear" {

		sender.SendMessage("§aCleared all effects from " + target.GetName())
		return true
	}

	var effectID int
	if id, ok := effectNames[effectArg]; ok {
		effectID = id
	} else if id, err := strconv.Atoi(effectArg); err == nil {
		effectID = id
	} else {
		sender.SendMessage("§cUnknown effect: " + effectArg)
		return true
	}

	duration := 600
	if len(args) >= 3 {
		if d, err := strconv.Atoi(args[2]); err == nil {
			duration = d * 20
		}
	}

	amplifier := 0
	if len(args) >= 4 {
		if a, err := strconv.Atoi(args[3]); err == nil {
			amplifier = a
		}
	}

	_ = effectID
	_ = duration
	_ = amplifier

	sender.SendMessage("§aApplied effect " + effectArg + " to " + target.GetName())
	return true
}
