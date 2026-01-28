package defaults

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/command"
	"github.com/scaxe/scaxe-go/pkg/entity"
	"github.com/scaxe/scaxe-go/pkg/level"
)

type SummonCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewSummonCommand(server ServerInterface) *SummonCommand {
	return &SummonCommand{
		BaseCommand: command.BaseCommand{
			Name:        "summon",
			Description: "Spawns an entity",
			Usage:       "/summon <entity> [x y z]",
			Permission:  "pocketmine.command.summon",
		},
		server: server,
	}
}

func (c *SummonCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) < 1 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	entityType := args[0]

	var x, y, z float64
	var ok bool

	if len(args) >= 4 {
		x, y, z, ok = parseLocation(sender, args, 1)
		if !ok {
			sender.SendMessage("§cInvalid coordinates")
			return true
		}
	} else if posSender, isPos := sender.(Positional); isPos {

		pPos := posSender.GetPosition()
		x, y, z = pPos.X, pPos.Y, pPos.Z
	} else {
		sender.SendMessage("§cPlease provide coordinates (sender is not positional)")
		return true
	}

	var targetLevel *level.Level
	if la, ok := sender.(LevelAware); ok {
		if lvl, ok := la.GetLevel().(*level.Level); ok {
			targetLevel = lvl
		}
	}
	if targetLevel == nil {
		if defaultLvl, ok := c.server.GetLevelManager().GetDefaultLevel().(*level.Level); ok {
			targetLevel = defaultLvl
		}
	}

	if targetLevel == nil {
		sender.SendMessage("§cInternal Error: No target level found.")
		return true
	}

	ent := entity.NewEntity()

	ent.Position.X = x
	ent.Position.Y = y
	ent.Position.Z = z
	ent.Level = targetLevel

	targetLevel.AddEntity(ent)

	sender.SendMessage(fmt.Sprintf("§aSummoned entity '%s' (Generic) at %.2f, %.2f, %.2f", entityType, x, y, z))
	return true
}
