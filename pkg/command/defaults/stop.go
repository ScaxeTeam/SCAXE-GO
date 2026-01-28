package defaults

import (
	"github.com/scaxe/scaxe-go/pkg/command"
	"github.com/scaxe/scaxe-go/pkg/logger"
)

type StopCommand struct {
	server ServerInterface
}

func NewStopCommand(server ServerInterface) *StopCommand {
	return &StopCommand{server: server}
}

func (c *StopCommand) GetName() string {
	return "stop"
}

func (c *StopCommand) GetDescription() string {
	return "Stops the server"
}

func (c *StopCommand) GetUsage() string {
	return "/stop"
}

func (c *StopCommand) GetPermission() string {
	return "scaxe.command.stop"
}

func (c *StopCommand) Execute(sender command.CommandSender, args []string) bool {

	if _, isConsole := sender.(*command.ConsoleCommandSender); !isConsole {
		if player, ok := sender.(command.PlayerSender); ok {
			if !c.server.IsOp(player.GetName()) {
				sender.SendMessage("§cYou don't have permission to use this command")
				return true
			}
		}
	}

	logger.Server("Server shutdown initiated by " + sender.GetName())
	sender.SendMessage("§eStopping the server...")

	c.server.Stop()
	return true
}
