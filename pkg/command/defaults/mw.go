package defaults

import (
	"strconv"
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type MWCommand struct {
	server LevelManagerProvider
}

type PlayerSwitcher interface {
	SwitchLevel(lvl interface{}) bool
}

type LevelManagerProvider interface {
	GetLevelManager() LevelManager
}

func NewMWCommand(server LevelManagerProvider) *MWCommand {
	return &MWCommand{server: server}
}

func (c *MWCommand) GetName() string {
	return "mw"
}

func (c *MWCommand) GetDescription() string {
	return "Multi-world management (list/create/load/unload/tp/info)"
}

func (c *MWCommand) GetUsage() string {
	return "/mw <list|create|load|unload|tp|info> [args...]"
}

func (c *MWCommand) GetPermission() string {
	return "scaxe.command.mw"
}

func (c *MWCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) == 0 {
		sender.SendMessage("§eMulti-World Commands:")
		sender.SendMessage("§7/mw list §f- List all loaded worlds")
		sender.SendMessage("§7/mw create <name> [generator] [seed] §f- Create new world")
		sender.SendMessage("§7/mw load <name> §f- Load existing world")
		sender.SendMessage("§7/mw unload <name> §f- Unload world")
		sender.SendMessage("§7/mw tp <world> §f- Teleport to world spawn")
		sender.SendMessage("§7/mw info [world] §f- Show world info")
		return true
	}

	subCmd := strings.ToLower(args[0])
	lm := c.server.GetLevelManager()

	switch subCmd {
	case "list":
		return c.cmdList(sender, lm)
	case "create":
		return c.cmdCreate(sender, lm, args[1:])
	case "load":
		return c.cmdLoad(sender, lm, args[1:])
	case "unload":
		return c.cmdUnload(sender, lm, args[1:])
	case "tp":
		return c.cmdTeleport(sender, lm, args[1:])
	case "info":
		return c.cmdInfo(sender, lm, args[1:])
	default:
		sender.SendMessage("§cUnknown subcommand: " + subCmd)
		return false
	}
}

func (c *MWCommand) cmdList(sender command.CommandSender, lm LevelManager) bool {
	names := lm.GetLevelNames()
	if len(names) == 0 {
		sender.SendMessage("§eNo worlds loaded.")
		return true
	}

	sender.SendMessage("§aLoaded Worlds §7(" + strconv.Itoa(len(names)) + ")§a:")
	defaultLvl := lm.GetDefaultLevel()
	for _, name := range names {
		lvl := lm.GetLevel(name)
		marker := ""
		if lvl == defaultLvl {
			marker = " §6[Default]"
		}
		sender.SendMessage("§7- §f" + name + marker)
	}
	return true
}

func (c *MWCommand) cmdCreate(sender command.CommandSender, lm LevelManager, args []string) bool {
	if len(args) < 1 {
		sender.SendMessage("§cUsage: /mw create <name> [generator] [seed]")
		return false
	}

	name := args[0]
	generator := "flat"
	var seed int64 = 0

	if len(args) >= 2 {
		generator = strings.ToLower(args[1])
	}
	if len(args) >= 3 {
		if s, err := strconv.ParseInt(args[2], 10, 64); err == nil {
			seed = s
		}
	}

	sender.SendMessage("§eCreating world '" + name + "' with generator '" + generator + "'...")

	_, err := lm.GenerateLevel(name, generator, seed)
	if err != nil {
		sender.SendMessage("§cFailed to create world: " + err.Error())
		return false
	}

	sender.SendMessage("§aWorld '" + name + "' created successfully!")
	return true
}

func (c *MWCommand) cmdLoad(sender command.CommandSender, lm LevelManager, args []string) bool {
	if len(args) < 1 {
		sender.SendMessage("§cUsage: /mw load <name>")
		return false
	}

	name := args[0]

	if lm.GetLevel(name) != nil {
		sender.SendMessage("§eWorld '" + name + "' is already loaded.")
		return true
	}

	sender.SendMessage("§eLoading world '" + name + "'...")

	_, err := lm.LoadLevel(name)
	if err != nil {
		sender.SendMessage("§cFailed to load world: " + err.Error())
		return false
	}

	sender.SendMessage("§aWorld '" + name + "' loaded successfully!")
	return true
}

func (c *MWCommand) cmdUnload(sender command.CommandSender, lm LevelManager, args []string) bool {
	if len(args) < 1 {
		sender.SendMessage("§cUsage: /mw unload <name>")
		return false
	}

	name := args[0]

	if lm.GetLevel(name) == nil {
		sender.SendMessage("§cWorld '" + name + "' is not loaded.")
		return false
	}

	if !lm.UnloadLevel(name) {
		sender.SendMessage("§cFailed to unload world '" + name + "'. It may be the default world.")
		return false
	}

	sender.SendMessage("§aWorld '" + name + "' unloaded successfully!")
	return true
}

func (c *MWCommand) cmdTeleport(sender command.CommandSender, lm LevelManager, args []string) bool {

	ps, ok := sender.(PlayerSwitcher)
	if !ok {
		sender.SendMessage("§cThis command can only be used by players.")
		return false
	}

	if len(args) < 1 {
		sender.SendMessage("§cUsage: /mw tp <world>")
		return false
	}

	name := args[0]
	lvl := lm.GetLevel(name)
	if lvl == nil {
		sender.SendMessage("§cWorld '" + name + "' is not loaded.")
		return false
	}

	sender.SendMessage("§aTeleporting to world '" + name + "'...")

	if ps.SwitchLevel(lvl) {
		sender.SendMessage("§aTeleported to world spawn.")
	} else {
		sender.SendMessage("§cFailed to switch level.")
	}
	return true
}

func (c *MWCommand) cmdInfo(sender command.CommandSender, lm LevelManager, args []string) bool {
	var name string
	if len(args) >= 1 {
		name = args[0]
	} else {
		defaultLvl := lm.GetDefaultLevel()
		if defaultLvl == nil {
			sender.SendMessage("§cNo default world set.")
			return false
		}
		sender.SendMessage("§eShowing info for default world...")
		return true
	}

	lvl := lm.GetLevel(name)
	if lvl == nil {
		sender.SendMessage("§cWorld '" + name + "' is not loaded.")
		return false
	}

	sender.SendMessage("§a=== World Info: " + name + " ===")
	sender.SendMessage("§7Status: §aLoaded")
	return true
}
