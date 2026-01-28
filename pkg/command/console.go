package command

import "github.com/scaxe/scaxe-go/pkg/logger"

type ConsoleCommandSender struct{}

func (c *ConsoleCommandSender) SendMessage(message string) {
	logger.Info(message)
}

func (c *ConsoleCommandSender) GetName() string {
	return "CONSOLE"
}

func (c *ConsoleCommandSender) IsOp() bool {
	return true
}
