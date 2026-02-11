package lua

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/scaxe/scaxe-go/pkg/logger"
	lua "github.com/yuin/gopher-lua"
)

type PluginManager struct {
	mu        sync.RWMutex
	server    ServerAPI
	plugins   map[string]*Plugin
	pluginDir string
}

func NewPluginManager(server ServerAPI, dir string) *PluginManager {
	return &PluginManager{
		server:    server,
		plugins:   make(map[string]*Plugin),
		pluginDir: dir,
	}
}

func (pm *PluginManager) LoadAll() error {
	if err := os.MkdirAll(pm.pluginDir, 0755); err != nil {
		return fmt.Errorf("failed to create plugins directory: %w", err)
	}

	entries, err := os.ReadDir(pm.pluginDir)
	if err != nil {
		return fmt.Errorf("failed to read plugins directory: %w", err)
	}

	loaded := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pluginDir := filepath.Join(pm.pluginDir, entry.Name())
		metaPath := filepath.Join(pluginDir, "plugin.yml")
		if _, err := os.Stat(metaPath); os.IsNotExist(err) {
			continue
		}

		if err := pm.LoadPlugin(entry.Name()); err != nil {
			logger.Error("Failed to load plugin", "name", entry.Name(), "error", err)
		} else {
			loaded++
		}
	}

	if loaded > 0 {
		logger.Server("Plugins loaded", "count", loaded)
	} else {
		logger.Server("No plugins found in", "dir", pm.pluginDir)
	}

	return nil
}

func (pm *PluginManager) LoadPlugin(dirName string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pluginDir := filepath.Join(pm.pluginDir, dirName)

	meta, err := loadPluginMeta(pluginDir)
	if err != nil {
		return err
	}

	if _, exists := pm.plugins[meta.Name]; exists {
		return fmt.Errorf("plugin %s is already loaded", meta.Name)
	}

	plugin := newPlugin(meta, pluginDir)
	plugin.createState(pm.server)

	if err := plugin.enable(); err != nil {
		plugin.disable()
		return err
	}

	pm.plugins[meta.Name] = plugin
	logger.Server("Plugin enabled", "name", meta.Name, "version", meta.Version, "author", meta.Author)
	return nil
}

func (pm *PluginManager) UnloadPlugin(name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	plugin, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	plugin.disable()
	delete(pm.plugins, name)
	logger.Server("Plugin disabled", "name", name)
	return nil
}

func (pm *PluginManager) ReloadPlugin(name string) error {
	pm.mu.Lock()
	plugin, exists := pm.plugins[name]
	if !exists {
		pm.mu.Unlock()
		return fmt.Errorf("plugin %s not found", name)
	}

	dir := plugin.Dir
	plugin.disable()
	delete(pm.plugins, name)
	pm.mu.Unlock()

	dirName := filepath.Base(dir)
	return pm.LoadPlugin(dirName)
}

func (pm *PluginManager) GetPlugin(name string) *Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.plugins[name]
}

func (pm *PluginManager) GetPlugins() []*Plugin {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	result := make([]*Plugin, 0, len(pm.plugins))
	for _, p := range pm.plugins {
		result = append(result, p)
	}
	return result
}

func (pm *PluginManager) DisableAll() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	for name, plugin := range pm.plugins {
		plugin.disable()
		logger.Server("Plugin disabled", "name", name)
	}
	pm.plugins = make(map[string]*Plugin)
}

func (pm *PluginManager) GetPluginNames() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	names := make([]string, 0, len(pm.plugins))
	for name := range pm.plugins {
		names = append(names, name)
	}
	return names
}

func (pm *PluginManager) EnablePlugin(name string) error {
	return pm.LoadPlugin(name)
}

func (pm *PluginManager) DisablePlugin(name string) error {
	return pm.UnloadPlugin(name)
}

func (pm *PluginManager) Tick(currentTick int64) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	for _, plugin := range pm.plugins {
		plugin.tick(currentTick)
	}
}

func (pm *PluginManager) CallEvent(eventName string, data map[string]interface{}) bool {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	cancelled := false
	for _, plugin := range pm.plugins {
		if !plugin.Enabled || plugin.State == nil {
			continue
		}

		eventTable := mapToLuaTable(plugin.State, data)
		if plugin.callEvent(eventName, eventTable) {
			cancelled = true
		}
	}
	return cancelled
}

func mapToLuaTable(L *lua.LState, data map[string]interface{}) *lua.LTable {
	tbl := L.NewTable()
	for k, v := range data {
		switch val := v.(type) {
		case string:
			tbl.RawSetString(k, lua.LString(val))
		case int:
			tbl.RawSetString(k, lua.LNumber(val))
		case int32:
			tbl.RawSetString(k, lua.LNumber(val))
		case int64:
			tbl.RawSetString(k, lua.LNumber(val))
		case float32:
			tbl.RawSetString(k, lua.LNumber(val))
		case float64:
			tbl.RawSetString(k, lua.LNumber(val))
		case bool:
			if val {
				tbl.RawSetString(k, lua.LTrue)
			} else {
				tbl.RawSetString(k, lua.LFalse)
			}
		}
	}
	return tbl
}
