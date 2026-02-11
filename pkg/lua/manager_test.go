package lua

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scaxe/scaxe-go/pkg/command"
)

// Mock implementations for testing

type mockServerAPI struct {
	broadcastMsg string
	players      map[string]*mockPlayerAPI
	level        *mockLevelAPI
	commands     map[string]command.Command
	currentTick  int64
}

func newMockServerAPI() *mockServerAPI {
	return &mockServerAPI{
		players:  make(map[string]*mockPlayerAPI),
		level:    &mockLevelAPI{},
		commands: make(map[string]command.Command),
	}
}

func (m *mockServerAPI) BroadcastMessage(msg string) { m.broadcastMsg = msg }
func (m *mockServerAPI) GetOnlineCount() int         { return len(m.players) }
func (m *mockServerAPI) GetMaxPlayers() int          { return 20 }
func (m *mockServerAPI) GetTPS() float64             { return 20.0 }
func (m *mockServerAPI) GetServerName() string       { return "TestServer" }
func (m *mockServerAPI) GetPlayer(name string) PlayerAPI {
	p, ok := m.players[name]
	if !ok {
		return nil
	}
	return p
}
func (m *mockServerAPI) GetOnlinePlayers() []PlayerAPI {
	result := make([]PlayerAPI, 0)
	for _, p := range m.players {
		result = append(result, p)
	}
	return result
}
func (m *mockServerAPI) KickPlayer(name string, reason string) {}
func (m *mockServerAPI) GetLevel() LevelAPI                    { return m.level }
func (m *mockServerAPI) RegisterCommand(cmd command.Command) {
	m.commands[cmd.GetName()] = cmd
}
func (m *mockServerAPI) UnregisterCommand(name string) {
	delete(m.commands, name)
}
func (m *mockServerAPI) Stop()                 {}
func (m *mockServerAPI) GetCurrentTick() int64 { return m.currentTick }

type mockPlayerAPI struct {
	name     string
	x, y, z  float64
	health   int
	gamemode int
	op       bool
	messages []string
}

func (p *mockPlayerAPI) GetName() string                          { return p.name }
func (p *mockPlayerAPI) GetPosition() (float64, float64, float64) { return p.x, p.y, p.z }
func (p *mockPlayerAPI) SetPosition(x, y, z float64)              { p.x, p.y, p.z = x, y, z }
func (p *mockPlayerAPI) SendMessage(msg string)                   { p.messages = append(p.messages, msg) }
func (p *mockPlayerAPI) GetGamemode() int                         { return p.gamemode }
func (p *mockPlayerAPI) SetGamemode(mode int)                     { p.gamemode = mode }
func (p *mockPlayerAPI) IsOp() bool                               { return p.op }
func (p *mockPlayerAPI) Kick(reason string)                       {}
func (p *mockPlayerAPI) GetHealth() int                           { return p.health }
func (p *mockPlayerAPI) SetHealth(health int)                     { p.health = health }
func (p *mockPlayerAPI) GetEntityID() int64                       { return 1 }

type mockLevelAPI struct {
	blocks map[[3]int32][2]uint8
	time   int64
	seed   int64
}

func (l *mockLevelAPI) GetBlock(x, y, z int32) (uint8, uint8) {
	key := [3]int32{x, y, z}
	if b, ok := l.blocks[key]; ok {
		return b[0], b[1]
	}
	return 0, 0
}
func (l *mockLevelAPI) SetBlock(x, y, z int32, id, meta uint8) {
	if l.blocks == nil {
		l.blocks = make(map[[3]int32][2]uint8)
	}
	l.blocks[[3]int32{x, y, z}] = [2]uint8{id, meta}
}
func (l *mockLevelAPI) GetTime() int64                                { return l.time }
func (l *mockLevelAPI) SetTime(t int64)                               { l.time = t }
func (l *mockLevelAPI) GetSeed() int64                                { return l.seed }
func (l *mockLevelAPI) GetSpawnLocation() (float64, float64, float64) { return 0, 64, 0 }

// Tests

func TestPluginMetaParsing(t *testing.T) {
	dir := t.TempDir()

	metaContent := `name: TestPlugin
version: "1.0.0"
author: Tester
description: "A test plugin"
main: main.lua
`
	if err := os.WriteFile(filepath.Join(dir, "plugin.yml"), []byte(metaContent), 0644); err != nil {
		t.Fatal(err)
	}

	meta, err := loadPluginMeta(dir)
	if err != nil {
		t.Fatal(err)
	}

	if meta.Name != "TestPlugin" {
		t.Errorf("Expected name 'TestPlugin', got '%s'", meta.Name)
	}
	if meta.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", meta.Version)
	}
	if meta.Author != "Tester" {
		t.Errorf("Expected author 'Tester', got '%s'", meta.Author)
	}
	if meta.Main != "main.lua" {
		t.Errorf("Expected main 'main.lua', got '%s'", meta.Main)
	}
}

func TestPluginMetaDefaultMain(t *testing.T) {
	dir := t.TempDir()

	metaContent := `name: NoMain
version: "1.0"
`
	if err := os.WriteFile(filepath.Join(dir, "plugin.yml"), []byte(metaContent), 0644); err != nil {
		t.Fatal(err)
	}

	meta, err := loadPluginMeta(dir)
	if err != nil {
		t.Fatal(err)
	}

	if meta.Main != "main.lua" {
		t.Errorf("Expected default main 'main.lua', got '%s'", meta.Main)
	}
}

func TestPluginMetaMissingName(t *testing.T) {
	dir := t.TempDir()

	metaContent := `version: "1.0"
author: Tester
`
	if err := os.WriteFile(filepath.Join(dir, "plugin.yml"), []byte(metaContent), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := loadPluginMeta(dir)
	if err == nil {
		t.Error("Expected error for missing plugin name")
	}
}

func createTestPlugin(t *testing.T, dir, name, luaCode string) {
	t.Helper()
	pluginDir := filepath.Join(dir, name)
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatal(err)
	}
	meta := `name: ` + name + `
version: "1.0.0"
main: main.lua
`
	if err := os.WriteFile(filepath.Join(pluginDir, "plugin.yml"), []byte(meta), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pluginDir, "main.lua"), []byte(luaCode), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestPluginManagerLoadAll(t *testing.T) {
	dir := t.TempDir()
	server := newMockServerAPI()

	createTestPlugin(t, dir, "TestPlugin1", `
function onEnable()
end
`)
	createTestPlugin(t, dir, "TestPlugin2", `
function onEnable()
end
`)

	pm := NewPluginManager(server, dir)
	if err := pm.LoadAll(); err != nil {
		t.Fatal(err)
	}

	plugins := pm.GetPlugins()
	if len(plugins) != 2 {
		t.Errorf("Expected 2 plugins, got %d", len(plugins))
	}
}

func TestPluginManagerReload(t *testing.T) {
	dir := t.TempDir()
	server := newMockServerAPI()

	createTestPlugin(t, dir, "ReloadTest", `
function onEnable()
end
function onDisable()
end
`)

	pm := NewPluginManager(server, dir)
	if err := pm.LoadPlugin("ReloadTest"); err != nil {
		t.Fatal(err)
	}

	if err := pm.ReloadPlugin("ReloadTest"); err != nil {
		t.Fatal(err)
	}

	p := pm.GetPlugin("ReloadTest")
	if p == nil {
		t.Error("Plugin should exist after reload")
	}
	if !p.Enabled {
		t.Error("Plugin should be enabled after reload")
	}
}

func TestPluginManagerUnload(t *testing.T) {
	dir := t.TempDir()
	server := newMockServerAPI()

	createTestPlugin(t, dir, "UnloadTest", `
function onEnable()
end
`)

	pm := NewPluginManager(server, dir)
	if err := pm.LoadPlugin("UnloadTest"); err != nil {
		t.Fatal(err)
	}

	if err := pm.UnloadPlugin("UnloadTest"); err != nil {
		t.Fatal(err)
	}

	if pm.GetPlugin("UnloadTest") != nil {
		t.Error("Plugin should not exist after unload")
	}
}

func TestEventHandling(t *testing.T) {
	dir := t.TempDir()
	server := newMockServerAPI()

	createTestPlugin(t, dir, "EventTest", `
events.listen("PlayerJoinEvent", function(e)
    server.broadcast("Welcome " .. e.playerName)
end)
`)

	pm := NewPluginManager(server, dir)
	if err := pm.LoadPlugin("EventTest"); err != nil {
		t.Fatal(err)
	}

	pm.CallEvent("PlayerJoinEvent", map[string]interface{}{
		"playerName": "Steve",
	})

	if server.broadcastMsg != "Welcome Steve" {
		t.Errorf("Expected broadcast 'Welcome Steve', got '%s'", server.broadcastMsg)
	}
}

func TestEventCancellation(t *testing.T) {
	dir := t.TempDir()
	server := newMockServerAPI()

	createTestPlugin(t, dir, "CancelTest", `
events.listen("PlayerChatEvent", function(e)
    e.cancelled = true
end)
`)

	pm := NewPluginManager(server, dir)
	if err := pm.LoadPlugin("CancelTest"); err != nil {
		t.Fatal(err)
	}

	cancelled := pm.CallEvent("PlayerChatEvent", map[string]interface{}{
		"playerName": "Steve",
		"message":    "hello",
	})

	if !cancelled {
		t.Error("Event should be cancelled")
	}
}

func TestCommandRegistration(t *testing.T) {
	dir := t.TempDir()
	server := newMockServerAPI()

	createTestPlugin(t, dir, "CmdTest", `
commands.register({
    name = "testcmd",
    description = "Test command",
    callback = function(sender, args)
        -- do nothing
    end
})
`)

	pm := NewPluginManager(server, dir)
	if err := pm.LoadPlugin("CmdTest"); err != nil {
		t.Fatal(err)
	}

	if _, ok := server.commands["testcmd"]; !ok {
		t.Error("Command 'testcmd' should be registered")
	}
}

func TestSchedulerDelayed(t *testing.T) {
	dir := t.TempDir()
	server := newMockServerAPI()
	server.currentTick = 0

	createTestPlugin(t, dir, "SchedulerTest", `
scheduler.delayed(5, function()
    server.broadcast("delayed!")
end)
`)

	pm := NewPluginManager(server, dir)
	if err := pm.LoadPlugin("SchedulerTest"); err != nil {
		t.Fatal(err)
	}

	// Not yet triggered
	pm.Tick(3)
	if server.broadcastMsg != "" {
		t.Errorf("Scheduler should not fire yet, got '%s'", server.broadcastMsg)
	}

	// Should trigger at tick 5
	pm.Tick(5)
	if server.broadcastMsg != "delayed!" {
		t.Errorf("Expected 'delayed!', got '%s'", server.broadcastMsg)
	}
}

func TestServerAPIFromLua(t *testing.T) {
	dir := t.TempDir()
	server := newMockServerAPI()

	createTestPlugin(t, dir, "APITest", `
function onEnable()
    local count = server.getOnlineCount()
    local max = server.getMaxPlayers()
    local name = server.getServerName()
    local tps = server.getTPS()
    -- If we get here without error, APIs work
    logger.info("APIs OK: " .. name .. " " .. count .. "/" .. max .. " TPS=" .. tps)
end
`)

	pm := NewPluginManager(server, dir)
	if err := pm.LoadPlugin("APITest"); err != nil {
		t.Fatal(err)
	}

	if pm.GetPlugin("APITest") == nil {
		t.Error("Plugin should be loaded")
	}
}

func TestDisableAll(t *testing.T) {
	dir := t.TempDir()
	server := newMockServerAPI()

	createTestPlugin(t, dir, "DA1", `function onEnable() end`)
	createTestPlugin(t, dir, "DA2", `function onEnable() end`)

	pm := NewPluginManager(server, dir)
	pm.LoadAll()

	if len(pm.GetPlugins()) != 2 {
		t.Errorf("Expected 2 plugins, got %d", len(pm.GetPlugins()))
	}

	pm.DisableAll()

	if len(pm.GetPlugins()) != 0 {
		t.Errorf("Expected 0 plugins after DisableAll, got %d", len(pm.GetPlugins()))
	}
}

func TestEmptyPluginDir(t *testing.T) {
	dir := t.TempDir()
	server := newMockServerAPI()

	pm := NewPluginManager(server, dir)
	if err := pm.LoadAll(); err != nil {
		t.Fatal(err)
	}

	if len(pm.GetPlugins()) != 0 {
		t.Error("Should have no plugins in empty dir")
	}
}
