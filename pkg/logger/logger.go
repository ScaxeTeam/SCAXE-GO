package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"time"
)

const (
	CategoryServer = "SERVER"
	CategoryPlayer = "PLAYER"
	CategoryPacket = "PACKET"
	CategoryDebug  = "DEBUG"
	CategoryError  = "ERROR"
	CategoryWarn   = "WARN"
)

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorWhite   = "\033[37m"
	colorGray    = "\033[90m"
	colorBold    = "\033[1m"
)

var (
	mu            sync.Mutex
	output        io.Writer = os.Stdout
	logFile       *os.File
	fileOutput    io.Writer
	debugEnabled  bool = false
	packetEnabled bool = false
	colorEnabled  bool = true
	initialized   bool = false
	ansiRegex          = regexp.MustCompile(`\x1b\[[0-9;]*m`)
)

func Init(out io.Writer, debug bool) {
	mu.Lock()
	defer mu.Unlock()

	output = out
	debugEnabled = debug

	if runtime.GOOS == "windows" {
		enableWindowsVT()
	}

	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logs directory: %v\n", err)
	} else {

		timestamp := time.Now().Format("2006-01-02_15-04-05")
		logPath := filepath.Join(logsDir, fmt.Sprintf("server_%s.log", timestamp))
		var err error
		logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to create log file: %v\n", err)
		} else {
			fileOutput = logFile

			fmt.Fprintf(logFile, "=== SCAXE-GO Server Log ===\n")
			fmt.Fprintf(logFile, "Started: %s\n", time.Now().Format(time.RFC3339))
			fmt.Fprintf(logFile, "===========================\n\n")
		}
	}

	initialized = true
}

func SetDebug(enabled bool) {
	mu.Lock()
	defer mu.Unlock()
	debugEnabled = enabled
}

func SetPacketLogging(enabled bool) {
	mu.Lock()
	defer mu.Unlock()
	packetEnabled = enabled
}

func SetColor(enabled bool) {
	mu.Lock()
	defer mu.Unlock()
	colorEnabled = enabled
}

func Close() {
	mu.Lock()
	defer mu.Unlock()
	if logFile != nil {
		fmt.Fprintf(logFile, "\n===========================\n")
		fmt.Fprintf(logFile, "Stopped: %s\n", time.Now().Format(time.RFC3339))
		fmt.Fprintf(logFile, "=== End of Log ===\n")
		logFile.Close()
		logFile = nil
		fileOutput = nil
	}
}

func stripANSI(text string) string {
	return ansiRegex.ReplaceAllString(text, "")
}

func IsDebugEnabled() bool {
	return debugEnabled
}

func timestamp() string {
	return time.Now().Format("15:04:05")
}

func colorize(color, text string) string {
	if colorEnabled {
		return color + text + colorReset
	}
	return text
}

func formatLog(category, color, msg string, args ...any) string {
	ts := colorize(colorGray, timestamp())
	cat := colorize(color, fmt.Sprintf("[%s]", category))

	fullMsg := msg
	if len(args) > 0 {

		for i := 0; i+1 < len(args); i += 2 {
			key := fmt.Sprintf("%v", args[i])
			val := fmt.Sprintf("%v", args[i+1])
			fullMsg += fmt.Sprintf(" %s=%s", colorize(colorCyan, key), val)
		}

		if len(args)%2 == 1 {
			fullMsg += fmt.Sprintf(" %v", args[len(args)-1])
		}
	}

	return fmt.Sprintf("%s %s %s\n", ts, cat, fullMsg)
}

func write(text string) {
	mu.Lock()
	defer mu.Unlock()

	fmt.Fprint(output, text)

	if fileOutput != nil {
		fmt.Fprint(fileOutput, stripANSI(text))
	}
}

func Server(msg string, args ...any) {
	text := formatLog(CategoryServer, colorGreen+colorBold, msg, args...)
	write(text)
}

func Player(msg string, args ...any) {
	text := formatLog(CategoryPlayer, colorCyan+colorBold, msg, args...)
	write(text)
}

func PlayerJoin(username, address string, protocol int) {
	msg := fmt.Sprintf("%s joined", colorize(colorGreen+colorBold, username))
	text := formatLog(CategoryPlayer, colorCyan+colorBold, msg, "address", address, "protocol", protocol)
	write(text)
}

func PlayerLeave(username, reason string) {
	msg := fmt.Sprintf("%s left", colorize(colorYellow, username))
	text := formatLog(CategoryPlayer, colorCyan+colorBold, msg, "reason", reason)
	write(text)
}

func PacketIn(packetName string, from string, args ...any) {

	if !packetEnabled {
		return
	}
	arrow := colorize(colorMagenta, "←")
	msg := fmt.Sprintf("%s %s from %s", arrow, packetName, from)
	text := formatLog(CategoryPacket, colorMagenta, msg, args...)
	write(text)
}

func PacketOut(packetName string, to string, args ...any) {

	if !packetEnabled {
		return
	}
	arrow := colorize(colorBlue, "→")
	msg := fmt.Sprintf("%s %s to %s", arrow, packetName, to)
	text := formatLog(CategoryPacket, colorBlue, msg, args...)
	write(text)
}

func Info(msg string, args ...any) {
	text := formatLog("INFO", colorWhite, msg, args...)
	write(text)
}

func Warn(msg string, args ...any) {
	text := formatLog(CategoryWarn, colorYellow, msg, args...)
	write(text)
}

func Error(msg string, args ...any) {
	text := formatLog(CategoryError, colorRed+colorBold, msg, args...)
	write(text)
}

func Debug(msg string, args ...any) {
	if !debugEnabled {
		return
	}
	text := formatLog(CategoryDebug, colorGray, msg, args...)
	write(text)
}

func Banner(serverName, version, address string, maxPlayers int) {
	line := colorize(colorGreen+colorBold, "═══════════════════════════════════════════")

	write(fmt.Sprintf("\n%s\n", line))
	write(fmt.Sprintf("  %s\n", colorize(colorGreen+colorBold, "SCAXE-GO")))
	write(fmt.Sprintf("  %s\n", colorize(colorWhite, version)))
	write(fmt.Sprintf("%s\n", line))
	write(fmt.Sprintf("  Server:  %s\n", colorize(colorCyan, serverName)))
	write(fmt.Sprintf("  Address: %s\n", colorize(colorCyan, address)))
	write(fmt.Sprintf("  Max:     %s players\n", colorize(colorCyan, fmt.Sprintf("%d", maxPlayers))))
	write(fmt.Sprintf("%s\n\n", line))
}
