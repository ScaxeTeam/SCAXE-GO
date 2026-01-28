package defaults

import (
	"strconv"
	"strings"

	"github.com/scaxe/scaxe-go/pkg/command"
)

type GamemodeCommand struct {
	server ServerInterface
}

func NewGamemodeCommand(server ServerInterface) *GamemodeCommand {
	return &GamemodeCommand{server: server}
}

func (c *GamemodeCommand) GetName() string {
	return "gamemode"
}

func (c *GamemodeCommand) GetDescription() string {
	return "Changes a player's game mode"
}

func (c *GamemodeCommand) GetUsage() string {
	return "/gamemode <0|1|2|3|s|c|a|sp> [player]"
}

func (c *GamemodeCommand) GetPermission() string {
	return "scaxe.command.gamemode"
}

const (
	GamemodeSurvival  = 0
	GamemodeCreative  = 1
	GamemodeAdventure = 2
	GamemodeSpectator = 3
)

func (c *GamemodeCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) < 1 {
		sender.SendMessage("§cUsage: " + c.GetUsage())
		return true
	}

	var targetPlayerName string
	var isSenderPlayer bool

	if player, ok := sender.(command.PlayerSender); ok {
		isSenderPlayer = true
		if !c.server.IsOp(player.GetName()) {
			sender.SendMessage("§cYou don't have permission to use this command")
			return true
		}
		targetPlayerName = player.GetName()
	}

	gamemode := parseGamemode(args[0])
	if gamemode == -1 {
		sender.SendMessage("§cInvalid game mode: " + args[0])
		sender.SendMessage("§7Use: 0/s (survival), 1/c (creative), 2/a (adventure), 3/sp (spectator)")
		return true
	}

	if len(args) >= 2 {
		targetPlayerName = args[1]
		found := c.server.GetPlayerByName(targetPlayerName)
		if found == nil {
			sender.SendMessage("§cPlayer not found: " + targetPlayerName)
			return true
		}
	} else if !isSenderPlayer {

		sender.SendMessage("§cUsage: /gamemode <mode> <player>")
		return true
	}

	if targetPlayerName == "" {
		sender.SendMessage("§cNo target player specified")
		return true
	}

	c.server.SetPlayerGamemode(targetPlayerName, gamemode)

	modeName := getGamemodeName(gamemode)
	if targetPlayerName == sender.GetName() {
		sender.SendMessage("§aYour game mode has been set to " + modeName)
	} else {
		sender.SendMessage("§aSet " + targetPlayerName + "'s game mode to " + modeName)

		target := c.server.GetPlayerByName(targetPlayerName)
		if target != nil {
			target.SendMessage("§aYour game mode has been set to " + modeName)
		}
	}
	return true
}

func parseGamemode(input string) int {
	input = strings.ToLower(strings.TrimSpace(input))

	if num, err := strconv.Atoi(input); err == nil {
		if num >= 0 && num <= 3 {
			return num
		}
		return -1
	}

	switch input {
	case "s", "survival":
		return GamemodeSurvival
	case "c", "creative":
		return GamemodeCreative
	case "a", "adventure":
		return GamemodeAdventure
	case "sp", "spectator", "v", "view":
		return GamemodeSpectator
	default:
		return -1
	}
}

func getGamemodeName(mode int) string {
	switch mode {
	case GamemodeSurvival:
		return "Survival"
	case GamemodeCreative:
		return "Creative"
	case GamemodeAdventure:
		return "Adventure"
	case GamemodeSpectator:
		return "Spectator"
	default:
		return "Unknown"
	}
}
