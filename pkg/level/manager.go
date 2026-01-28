package level

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/scaxe/scaxe-go/pkg/level/generator"
	"github.com/scaxe/scaxe-go/pkg/logger"
)

type LevelManager struct {
	mu sync.RWMutex

	levels       map[string]*Level
	defaultLevel *Level
	basePath     string
}

func NewLevelManager(basePath string) *LevelManager {

	worldsPath := filepath.Join(basePath, "worlds")
	if _, err := os.Stat(worldsPath); os.IsNotExist(err) {
		os.MkdirAll(worldsPath, 0755)
	}

	return &LevelManager{
		levels:   make(map[string]*Level),
		basePath: basePath,
	}
}

func (m *LevelManager) LoadLevel(name string) (*Level, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if lvl, exists := m.levels[name]; exists {
		return lvl, nil
	}

	worldPath := filepath.Join(m.basePath, "worlds", name)
	if _, err := os.Stat(worldPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("world '%s' does not exist", name)
	}

	provider, err := NewAnvilProvider(worldPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider for '%s': %v", name, err)
	}

	lvl := NewLevel(name, worldPath, provider, "normal")
	m.levels[name] = lvl

	logger.Info("Loaded world", "name", name)
	return lvl, nil
}

func (m *LevelManager) GenerateLevel(name string, generatorName string, seed int64) (*Level, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.levels[name]; exists {
		return nil, fmt.Errorf("world '%s' already loaded", name)
	}

	worldPath := filepath.Join(m.basePath, "worlds", name)
	if _, err := os.Stat(worldPath); !os.IsNotExist(err) {
		return nil, fmt.Errorf("world directory '%s' already exists", name)
	}

	if err := os.MkdirAll(worldPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create world directory: %v", err)
	}

	provider, err := NewAnvilProvider(worldPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %v", err)
	}

	lvl := NewLevel(name, worldPath, provider, generatorName)
	lvl.Seed = seed

	gen := generator.GetGenerator(generatorName, nil)
	if gen != nil {
		lvl.Generator = gen
		gen.Init(lvl, seed)
	}

	m.levels[name] = lvl
	logger.Info("Generated new world", "name", name, "generator", generatorName, "seed", seed)
	return lvl, nil
}

func (m *LevelManager) UnloadLevel(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	lvl, exists := m.levels[name]
	if !exists {
		return false
	}

	if lvl == m.defaultLevel {
		logger.Warn("Cannot unload default level", "name", name)
		return false
	}

	lvl.Close()
	delete(m.levels, name)

	logger.Info("Unloaded world", "name", name)
	return true
}

func (m *LevelManager) GetLevel(name string) *Level {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.levels[name]
}

func (m *LevelManager) GetLevels() []*Level {
	m.mu.RLock()
	defer m.mu.RUnlock()

	levels := make([]*Level, 0, len(m.levels))
	for _, lvl := range m.levels {
		levels = append(levels, lvl)
	}
	return levels
}

func (m *LevelManager) GetLevelNames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.levels))
	for name := range m.levels {
		names = append(names, name)
	}
	return names
}

func (m *LevelManager) GetDefaultLevel() *Level {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.defaultLevel
}

func (m *LevelManager) SetDefaultLevel(lvl *Level) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.defaultLevel = lvl
}

func (m *LevelManager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, lvl := range m.levels {
		lvl.Close()
		logger.Info("Closed world", "name", name)
	}
	m.levels = make(map[string]*Level)
	m.defaultLevel = nil
}

func (m *LevelManager) LevelCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.levels)
}

func NewAnvilProvider(path string) (Provider, error) {

	return nil, fmt.Errorf("use anvil.NewAnvilProvider directly - this is a placeholder")
}
