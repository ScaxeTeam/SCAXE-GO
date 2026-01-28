package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/scaxe/scaxe-go/internal/version"
	"github.com/scaxe/scaxe-go/pkg/config"

	_ "github.com/scaxe/scaxe-go/pkg/level/generator"
	_ "github.com/scaxe/scaxe-go/pkg/level/generator/gorigional"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/server"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("\n========== SERVER CRASHED ==========")
			fmt.Printf("Panic: %v\n", r)
			fmt.Println("\nStack trace:")
			debug.PrintStack()
			fmt.Println("=====================================")
			fmt.Println("Press Enter to exit...")
			bufio.NewReader(os.Stdin).ReadBytes('\n')
		}
	}()

	showVersion := flag.Bool("version", false, "Show version information and exit")
	showHelp := flag.Bool("help", false, "Show help message and exit")
	configPath := flag.String("config", "server.properties", "Path to server configuration file")
	debug := flag.Bool("debug", false, "Enable debug logging")
	noColor := flag.Bool("no-color", false, "Disable colored output")

	flag.Parse()

	if *showVersion {
		fmt.Println(version.Full())
		os.Exit(0)
	}

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	logger.Init(os.Stdout, *debug)
	if *noColor {
		logger.SetColor(false)
	}

	logger.Server("Loading configuration", "path", *configPath)

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	if cfg.DebugMode {
		logger.SetDebug(true)
		logger.SetPacketLogging(true)
	}

	logger.Server("Configuration loaded",
		"serverName", cfg.ServerName,
		"maxPlayers", cfg.MaxPlayers,
		"gamemode", cfg.Gamemode)

	srv := server.NewServer(cfg)

	if err := srv.Start(); err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if !srv.IsRunning() {
				return
			}
			srv.HandleConsoleCommand(scanner.Text())
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Server("Press Ctrl+C to stop the server")

	select {
	case sig := <-sigChan:
		logger.Server("Received signal", "signal", sig.String())
	case <-srv.StopChan():
		logger.Server("Server shutdown requested via command")
	}

	srv.Stop()
	logger.Close()
}

func printHelp() {
	fmt.Println(version.Full())
	fmt.Println()
	fmt.Println("Usage: server [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --version        Show version information and exit")
	fmt.Println("  --help           Show this help message and exit")
	fmt.Println("  --config PATH    Path to server.properties (default: server.properties)")
	fmt.Println("  --debug          Enable debug logging (shows packet details)")
	fmt.Println("  --no-color       Disable colored output")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  server                         Start with default config")
	fmt.Println("  server --config my.properties  Start with custom config")
	fmt.Println("  server --debug                 Start with debug output")
}
