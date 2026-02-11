package defaults

import (
	"fmt"
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type PluginManagerInterface interface {
	GetPluginNames() []string
	ReloadPlugin(name string) error
	EnablePlugin(name string) error
	DisablePlugin(name string) error
}

var pluginManager PluginManagerInterface

func SetPluginManager(pm PluginManagerInterface) {
	pluginManager = pm
}

type LuaPluginCommand struct {
	command.BaseCommand
}

func NewLuaPluginCommand() *LuaPluginCommand {
	return &LuaPluginCommand{
		BaseCommand: command.BaseCommand{
			Name:        "luaplugin",
			Description: "Manage Lua plugins",
			Usage:       "/luaplugin <list|reload|disable|enable> [name]",
			Permission:  "scaxe.command.luaplugin",
		},
	}
}

func (c *LuaPluginCommand) Execute(sender command.CommandSender, args []string) bool {
	if pluginManager == nil {
		sender.SendMessage("§cLua plugin system is not initialized")
		return true
	}

	if len(args) == 0 {
		sender.SendMessage("§eUsage: /luaplugin <list|reload|disable|enable> [name]")
		return true
	}

	subCmd := strings.ToLower(args[0])

	switch subCmd {
	case "list":
		names := pluginManager.GetPluginNames()
		if len(names) == 0 {
			sender.SendMessage("§eNo Lua plugins loaded")
		} else {
			sender.SendMessage(fmt.Sprintf("§aLua Plugins (%d): §f%s", len(names), strings.Join(names, ", ")))
		}

	case "reload":
		if len(args) < 2 {
			sender.SendMessage("§cUsage: /luaplugin reload <name>")
			return true
		}
		name := args[1]
		if err := pluginManager.ReloadPlugin(name); err != nil {
			sender.SendMessage(fmt.Sprintf("§cFailed to reload plugin '%s': %v", name, err))
		} else {
			sender.SendMessage(fmt.Sprintf("§aPlugin '%s' reloaded successfully", name))
		}

	case "disable":
		if len(args) < 2 {
			sender.SendMessage("§cUsage: /luaplugin disable <name>")
			return true
		}
		name := args[1]
		if err := pluginManager.DisablePlugin(name); err != nil {
			sender.SendMessage(fmt.Sprintf("§cFailed to disable plugin '%s': %v", name, err))
		} else {
			sender.SendMessage(fmt.Sprintf("§aPlugin '%s' disabled", name))
		}

	case "enable":
		if len(args) < 2 {
			sender.SendMessage("§cUsage: /luaplugin enable <name>")
			return true
		}
		name := args[1]
		if err := pluginManager.EnablePlugin(name); err != nil {
			sender.SendMessage(fmt.Sprintf("§cFailed to enable plugin '%s': %v", name, err))
		} else {
			sender.SendMessage(fmt.Sprintf("§aPlugin '%s' enabled", name))
		}

	default:
		sender.SendMessage("§cUnknown subcommand. Use: list, reload, disable, enable")
	}

	return true
}
