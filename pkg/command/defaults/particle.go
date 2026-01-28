package defaults

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/command"
	"github.com/scaxe/scaxe-go/pkg/level"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

type ParticleCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewParticleCommand(server ServerInterface) *ParticleCommand {
	return &ParticleCommand{
		BaseCommand: command.BaseCommand{
			Name:        "particle",
			Description: "Spawns particles",
			Usage:       "/particle <name> <x> <y> <z> [data]",
			Permission:  "pocketmine.command.particle",
		},
		server: server,
	}
}

func (c *ParticleCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) < 4 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	name := args[0]
	x, y, z, ok := parseLocation(sender, args, 1)
	if !ok {
		sender.SendMessage("§cInvalid coordinates")
		return true
	}

	data := 0
	if len(args) >= 5 {

		fmt.Sscanf(args[4], "%d", &data)
	}

	particleID := 0

	switch name {
	case "explode":
		particleID = 1
	case "smoke":
		particleID = 2

	default:

		fmt.Sscanf(name, "%d", &particleID)
	}

	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(protocol.EventParticleSpawn)
	pk.X = float32(x)
	pk.Y = float32(y)
	pk.Z = float32(z)
	pk.Data = int32(particleID)

	var targetLevel *level.Level
	if la, ok := sender.(LevelAware); ok {
		if lvl, ok := la.GetLevel().(*level.Level); ok {
			targetLevel = lvl
		}
	} else {

		if defaultLvl, ok := c.server.GetLevelManager().GetDefaultLevel().(*level.Level); ok {
			targetLevel = defaultLvl
		}
	}

	if targetLevel == nil {
		sender.SendMessage("§cInternal Error: Level not found.")
		return true
	}

	sender.SendMessage(fmt.Sprintf("§aSpawned particle %s (ID %d) at %.1f %.1f %.1f", name, particleID, x, y, z))
	return true
}
