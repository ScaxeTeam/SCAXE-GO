package defaults

import (
	"fmt"
	"runtime"
	"time"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type StatusCommand struct {
	command.BaseCommand
	Server Server
}

func NewStatusCommand(s Server) *StatusCommand {
	return &StatusCommand{
		BaseCommand: command.BaseCommand{
			Name:        "status",
			Description: "Reads back the server's performance.",
			Usage:       "/status",
			Permission:  "pocketmine.command.status",
		},
		Server: s,
	}
}

func (c *StatusCommand) Execute(sender command.CommandSender, args []string) bool {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	sender.SendMessage("---- Server Status ----")

	uptime := time.Since(c.Server.GetStartTime())
	sender.SendMessage(fmt.Sprintf("Uptime: %s", uptime.Round(time.Second)))

	tps := c.Server.GetTPS()
	mspt := c.Server.GetMSPT()
	tpsPercent := (tps / 20.0) * 100
	sender.SendMessage(fmt.Sprintf("Current TPS: %.1f (%.0f%%)", tps, tpsPercent))
	sender.SendMessage(fmt.Sprintf("Average MSPT: %.2f ms", mspt))

	sender.SendMessage(fmt.Sprintf("Goroutines: %d", runtime.NumGoroutine()))
	sender.SendMessage(fmt.Sprintf("Allocated Memory: %.2f MB", float64(m.Alloc)/1024/1024))
	sender.SendMessage(fmt.Sprintf("Total Memory Allowed: %.2f MB", float64(m.Sys)/1024/1024))

	sender.SendMessage(fmt.Sprintf("RakNet Session Count: %d", c.Server.GetRakNetSessionCount()))

	return true
}
