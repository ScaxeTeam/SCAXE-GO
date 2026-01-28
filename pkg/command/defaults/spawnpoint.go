package defaults

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/command"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

type SpawnpointCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewSpawnpointCommand(server ServerInterface) *SpawnpointCommand {
	return &SpawnpointCommand{
		BaseCommand: command.BaseCommand{
			Name:        "spawnpoint",
			Description: "Sets the spawn point for a player",
			Usage:       "/spawnpoint [player] [x y z]",
			Permission:  "pocketmine.command.spawnpoint",
		},
		server: server,
	}
}

func (c *SpawnpointCommand) Execute(sender command.CommandSender, args []string) bool {
	var target command.PlayerSender
	var x, y, z float64
	var ok bool

	if len(args) == 0 {

		if p, ok := sender.(command.PlayerSender); ok {
			target = p
			if posSender, isPos := sender.(Positional); isPos {
				pos := posSender.GetPosition()
				x, y, z = pos.X, pos.Y, pos.Z
			} else {
				sender.SendMessage("§cYou are not in the world.")
				return true
			}
		} else {
			sender.SendMessage("§cPlease specify a player.")
			return true
		}
	} else {

		p := c.server.GetPlayerByName(args[0])
		if p != nil {
			target = p
			if len(args) >= 4 {
				x, y, z, ok = parseLocation(sender, args, 1)
				if !ok {
					sender.SendMessage("§cInvalid coordinates")
					return true
				}
			} else {

				if posSender, isPos := p.(Positional); isPos {
					pos := posSender.GetPosition()
					x, y, z = pos.X, pos.Y, pos.Z
				} else {
					sender.SendMessage("§cTarget is not in the world.")
					return true
				}
			}
		} else {

			if self, isPlayer := sender.(command.PlayerSender); isPlayer {
				target = self
				if len(args) >= 3 {
					x, y, z, ok = parseLocation(sender, args, 0)
					if !ok {
						sender.SendMessage("§cPlayer not found or invalid coordinates")
						return true
					}
				} else {
					sender.SendMessage("§cPlayer not found: " + args[0])
					return true
				}
			} else {
				sender.SendMessage("§cPlayer not found: " + args[0])
				return true
			}
		}
	}

	pk := protocol.NewSetSpawnPositionPacket()
	pk.X = int32(x)
	pk.Y = int32(y)
	pk.Z = int32(z)

	sender.SendMessage(fmt.Sprintf("§aSet %s's spawn point to %.1f, %.1f, %.1f", target.GetName(), x, y, z))

	return true
}
