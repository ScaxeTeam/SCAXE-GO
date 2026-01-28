package defaults

import (
	"fmt"
	"runtime"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type TpsCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewTpsCommand(server ServerInterface) *TpsCommand {
	return &TpsCommand{
		BaseCommand: command.BaseCommand{
			Name:        "tps",
			Description: "Shows server TPS",
			Usage:       "/tps",
			Permission:  "pocketmine.command.tps",
		},
		server: server,
	}
}

func (c *TpsCommand) Execute(sender command.CommandSender, args []string) bool {
	tps := c.server.GetTPS()
	avgTps := c.server.GetAverageTPS()

	var tpsColor string
	if tps >= 20 {
		tpsColor = "§a"
	} else if tps >= 15 {
		tpsColor = "§e"
	} else if tps >= 10 {
		tpsColor = "§6"
	} else {
		tpsColor = "§c"
	}

	sender.SendMessage(fmt.Sprintf("§aCurrent TPS: %s%.2f§a (Average: %.2f)", tpsColor, tps, avgTps))

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	sender.SendMessage(fmt.Sprintf("§aMemory: %.2f MB / %.2f MB allocated",
		float64(m.Alloc)/1024/1024,
		float64(m.Sys)/1024/1024))

	return true
}
