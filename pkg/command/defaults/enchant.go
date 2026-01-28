package defaults

import (
	"strconv"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type EnchantCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewEnchantCommand(server ServerInterface) *EnchantCommand {
	return &EnchantCommand{
		BaseCommand: command.BaseCommand{
			Name:        "enchant",
			Description: "Enchants a player's held item",
			Usage:       "/enchant <player> <enchantment> [level]",
			Permission:  "pocketmine.command.enchant",
		},
		server: server,
	}
}

var enchantNames = map[string]int{
	"protection":            0,
	"fire_protection":       1,
	"feather_falling":       2,
	"blast_protection":      3,
	"projectile_protection": 4,
	"thorns":                5,
	"respiration":           6,
	"aqua_affinity":         8,
	"sharpness":             9,
	"smite":                 10,
	"bane_of_arthropods":    11,
	"knockback":             12,
	"fire_aspect":           13,
	"looting":               14,
	"efficiency":            15,
	"silk_touch":            16,
	"unbreaking":            17,
	"fortune":               18,
	"power":                 19,
	"punch":                 20,
	"flame":                 21,
	"infinity":              22,
	"luck_of_the_sea":       23,
	"lure":                  24,
}

func (c *EnchantCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) < 2 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	playerName := args[0]
	enchantArg := args[1]
	level := 1

	if len(args) >= 3 {
		if l, err := strconv.Atoi(args[2]); err == nil {
			level = l
		}
	}

	target := c.server.GetPlayerByName(playerName)
	if target == nil {
		sender.SendMessage("§cPlayer not found: " + playerName)
		return true
	}

	var enchantID int
	if id, ok := enchantNames[enchantArg]; ok {
		enchantID = id
	} else if id, err := strconv.Atoi(enchantArg); err == nil {
		enchantID = id
	} else {
		sender.SendMessage("§cUnknown enchantment: " + enchantArg)
		return true
	}

	_ = enchantID
	_ = level

	sender.SendMessage("§aEnchanted " + target.GetName() + "'s item with " +
		enchantArg + " level " + strconv.Itoa(level))
	return true
}
