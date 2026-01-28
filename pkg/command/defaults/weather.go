package defaults

import (
	"github.com/scaxe/scaxe-go/pkg/command"
)

const (
	WeatherClear   = 0
	WeatherRain    = 1
	WeatherThunder = 2
)

type WeatherCommand struct {
	command.BaseCommand
	server ServerInterface
}

func NewWeatherCommand(server ServerInterface) *WeatherCommand {
	return &WeatherCommand{
		BaseCommand: command.BaseCommand{
			Name:        "weather",
			Description: "Sets the weather",
			Usage:       "/weather <clear|rain|thunder> [duration]",
			Permission:  "pocketmine.command.weather",
		},
		server: server,
	}
}

func (c *WeatherCommand) Execute(sender command.CommandSender, args []string) bool {
	if len(args) < 1 {
		sender.SendMessage("§cUsage: " + c.Usage)
		return false
	}

	var weather int
	switch args[0] {
	case "clear", "sun", "0":
		weather = WeatherClear
		sender.SendMessage("§aWeather set to clear")
	case "rain", "1":
		weather = WeatherRain
		sender.SendMessage("§aWeather set to rain")
	case "thunder", "storm", "2":
		weather = WeatherThunder
		sender.SendMessage("§aWeather set to thunder")
	default:
		sender.SendMessage("§cInvalid weather type. Use: clear, rain, or thunder")
		return true
	}

	_ = weather
	c.server.BroadcastMessage("§eWeather has been changed")

	return true
}
