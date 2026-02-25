package command

import "github.com/scaxe/scaxe-go/pkg/event"

type CommandSender interface {
	SendMessage(message string)
	GetName() string
	IsOp() bool
	HasPermission(name string) bool
}

type PlayerSender interface {
	CommandSender
	GetEntityID() int64
}

type Command interface {
	Execute(sender CommandSender, args []string) bool
	GetName() string
	GetDescription() string
	GetUsage() string
	GetPermission() string
}

type BaseCommand struct {
	Name        string
	Description string
	Usage       string
	Permission  string
}

func (c *BaseCommand) GetName() string {
	return c.Name
}

func (c *BaseCommand) GetDescription() string {
	return c.Description
}

func (c *BaseCommand) GetUsage() string {
	return c.Usage
}

func (c *BaseCommand) GetPermission() string {
	return c.Permission
}

type CommandMap struct {
	commands map[string]Command
}

func NewCommandMap() *CommandMap {
	return &CommandMap{
		commands: make(map[string]Command),
	}
}

func (m *CommandMap) Register(cmd Command) {
	m.commands[cmd.GetName()] = cmd
}

func (m *CommandMap) Unregister(name string) {
	delete(m.commands, name)
}

func (m *CommandMap) RegisterAlias(alias string, targetName string) {
	if cmd, ok := m.commands[targetName]; ok {
		m.commands[alias] = cmd
	}
}

func (m *CommandMap) Dispatch(sender CommandSender, cmdLine string) bool {
	if cmdLine == "" {
		return false
	}

	// Fire PlayerCommandPreprocessEvent for player senders
	if ps, ok := sender.(PlayerSender); ok {
		ppEvt := event.NewPlayerCommandPreprocessEvent(sender.GetName(), ps.GetEntityID(), cmdLine)
		event.Call(ppEvt)
		if ppEvt.IsCancelled() {
			return true
		}
		cmdLine = ppEvt.GetMessage() // allow plugins to modify the command
	}

	// Fire CommandEvent
	cmdEvt := event.NewCommandEvent(cmdLine, sender.GetName())
	event.Call(cmdEvt)
	if cmdEvt.IsCancelled() {
		return true
	}
	cmdLine = cmdEvt.GetCommand()

	var args []string
	currentArg := ""
	for _, char := range cmdLine {
		if char == ' ' {
			if currentArg != "" {
				args = append(args, currentArg)
				currentArg = ""
			}
		} else {
			currentArg += string(char)
		}
	}
	if currentArg != "" {
		args = append(args, currentArg)
	}

	if len(args) == 0 {
		return false
	}

	label := args[0]

	args = args[1:]

	if cmd, ok := m.commands[label]; ok {
		if cmd.GetPermission() != "" && !sender.HasPermission(cmd.GetPermission()) {
			sender.SendMessage("§cYou do not have permission to use this command.")
			return true
		}
		return cmd.Execute(sender, args)
	}

	return false
}
