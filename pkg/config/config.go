package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/scaxe/scaxe-go/pkg/logger"
)

type ServerConfig struct {
	ServerName string
	ServerPort int
	ServerIP   string
	MaxPlayers int
	MOTD       string

	Gamemode        int
	Difficulty      int
	LevelName       string
	LevelSeed       string
	LevelType       string
	SpawnProtection int

	OnlineMode  bool
	WhiteList   bool
	AllowFlight bool
	Hardcore    bool
	PvP         bool

	ViewDistance int
	TickRate     int

	DebugMode       bool
	DebugItemPickup bool

	Properties map[string]string
}

func DefaultConfig() *ServerConfig {
	logger.Debug("DefaultConfig", "action", "creating default configuration")
	return &ServerConfig{
		ServerName:      "Scaxe Go Server",
		ServerPort:      19132,
		ServerIP:        "0.0.0.0",
		MaxPlayers:      20,
		MOTD:            "A Scaxe Go Server",
		Gamemode:        0,
		Difficulty:      1,
		LevelName:       "world",
		LevelSeed:       "",
		LevelType:       "gorigional",
		SpawnProtection: 16,
		OnlineMode:      false,
		WhiteList:       false,
		AllowFlight:     false,
		Hardcore:        false,
		PvP:             true,
		ViewDistance:    8,
		TickRate:        20,
		DebugMode:       false,
		DebugItemPickup: false,
		Properties:      make(map[string]string),
	}
}

func Load(path string) (*ServerConfig, error) {
	logger.Debug("Config.Load", "path", path)

	cfg := DefaultConfig()

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warn("Config.Load", "warning", "config file not found, creating default", "path", path)
			if err := cfg.Save(path); err != nil {
				logger.Error("Config.Load", "error", "failed to create default config", "err", err)
			}
			return cfg, nil
		}
		logger.Error("Config.Load", "error", err, "path", path)
		return nil, fmt.Errorf("failed to open config: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			logger.Warn("Config.Load", "warning", "invalid line format", "line", lineNum, "content", line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		cfg.Properties[key] = value

		switch key {
		case "server-name":
			cfg.ServerName = value
			logger.Debug("Config.Load", "key", key, "value", value)
		case "server-port":
			if v, err := strconv.Atoi(value); err == nil {
				cfg.ServerPort = v
				logger.Debug("Config.Load", "key", key, "value", v)
			}
		case "server-ip":
			cfg.ServerIP = value
			logger.Debug("Config.Load", "key", key, "value", value)
		case "max-players":
			if v, err := strconv.Atoi(value); err == nil {
				cfg.MaxPlayers = v
				logger.Debug("Config.Load", "key", key, "value", v)
			}
		case "motd":
			cfg.MOTD = value
			logger.Debug("Config.Load", "key", key, "value", value)
		case "gamemode":
			if v, err := strconv.Atoi(value); err == nil {
				cfg.Gamemode = v
				logger.Debug("Config.Load", "key", key, "value", v)
			}
		case "difficulty":
			if v, err := strconv.Atoi(value); err == nil {
				cfg.Difficulty = v
				logger.Debug("Config.Load", "key", key, "value", v)
			}
		case "level-name":
			cfg.LevelName = value
			logger.Debug("Config.Load", "key", key, "value", value)
		case "level-seed":
			cfg.LevelSeed = value
			logger.Debug("Config.Load", "key", key, "value", value)
		case "level-type":
			cfg.LevelType = value
			logger.Debug("Config.Load", "key", key, "value", value)
		case "spawn-protection":
			if v, err := strconv.Atoi(value); err == nil {
				cfg.SpawnProtection = v
				logger.Debug("Config.Load", "key", key, "value", v)
			}
		case "online-mode":
			cfg.OnlineMode = parseBool(value)
			logger.Debug("Config.Load", "key", key, "value", cfg.OnlineMode)
		case "white-list":
			cfg.WhiteList = parseBool(value)
			logger.Debug("Config.Load", "key", key, "value", cfg.WhiteList)
		case "allow-flight":
			cfg.AllowFlight = parseBool(value)
			logger.Debug("Config.Load", "key", key, "value", cfg.AllowFlight)
		case "hardcore":
			cfg.Hardcore = parseBool(value)
			logger.Debug("Config.Load", "key", key, "value", cfg.Hardcore)
		case "pvp":
			cfg.PvP = parseBool(value)
			logger.Debug("Config.Load", "key", key, "value", cfg.PvP)
		case "view-distance":
			if v, err := strconv.Atoi(value); err == nil {
				cfg.ViewDistance = v
				logger.Debug("Config.Load", "key", key, "value", v)
			}
		case "debug":
			cfg.DebugMode = parseBool(value)
			logger.Debug("Config.Load", "key", key, "value", cfg.DebugMode)
		case "debug-item-pickup":
			cfg.DebugItemPickup = parseBool(value)
			logger.Debug("Config.Load", "key", key, "value", cfg.DebugItemPickup)
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Config.Load", "error", "scanner error", "err", err)
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	logger.Info("Config.Load", "status", "configuration loaded", "properties", len(cfg.Properties))

	if err := cfg.Save(path); err != nil {
		logger.Warn("Config.Load", "warning", "failed to update config file", "err", err)
	}

	return cfg, nil
}

func (c *ServerConfig) Save(path string) error {
	logger.Debug("Config.Save", "path", path)

	file, err := os.Create(path)
	if err != nil {
		logger.Error("Config.Save", "error", err, "path", path)
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	lines := []string{
		"# Scaxe Go Server Properties",
		"# Generated automatically",
		"",
		fmt.Sprintf("server-name=%s", c.ServerName),
		fmt.Sprintf("server-port=%d", c.ServerPort),
		fmt.Sprintf("server-ip=%s", c.ServerIP),
		fmt.Sprintf("max-players=%d", c.MaxPlayers),
		fmt.Sprintf("motd=%s", c.MOTD),
		fmt.Sprintf("gamemode=%d", c.Gamemode),
		fmt.Sprintf("difficulty=%d", c.Difficulty),
		fmt.Sprintf("level-name=%s", c.LevelName),
		fmt.Sprintf("level-seed=%s", c.LevelSeed),
		fmt.Sprintf("level-type=%s", c.LevelType),
		fmt.Sprintf("spawn-protection=%d", c.SpawnProtection),
		fmt.Sprintf("online-mode=%t", c.OnlineMode),
		fmt.Sprintf("white-list=%t", c.WhiteList),
		fmt.Sprintf("allow-flight=%t", c.AllowFlight),
		fmt.Sprintf("hardcore=%t", c.Hardcore),
		fmt.Sprintf("pvp=%t", c.PvP),
		fmt.Sprintf("view-distance=%d", c.ViewDistance),
		fmt.Sprintf("debug=%t", c.DebugMode),
		fmt.Sprintf("debug-item-pickup=%t", c.DebugItemPickup),
	}

	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			logger.Error("Config.Save", "error", err)
			return fmt.Errorf("failed to write config: %w", err)
		}
	}

	logger.Info("Config.Save", "status", "configuration saved", "path", path)
	return nil
}

func (c *ServerConfig) Get(key, defaultValue string) string {
	if v, ok := c.Properties[key]; ok {
		return v
	}
	return defaultValue
}

func (c *ServerConfig) GetInt(key string, defaultValue int) int {
	if v, ok := c.Properties[key]; ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}

func (c *ServerConfig) GetBool(key string, defaultValue bool) bool {
	if v, ok := c.Properties[key]; ok {
		return parseBool(v)
	}
	return defaultValue
}

func parseBool(value string) bool {
	v := strings.ToLower(strings.TrimSpace(value))
	return v == "true" || v == "on" || v == "yes" || v == "1"
}
