package permission

import (
	"encoding/json"
	"os"
	"strings"
	"sync"

	"github.com/scaxe/scaxe-go/pkg/logger"
)

type OpEntry struct {
	Name     string `json:"name"`
	ClientID int64  `json:"cid"`
}

type OpManager struct {
	mu       sync.RWMutex
	ops      map[string]int64
	filePath string
}

func NewOpManager(path string) *OpManager {
	return &OpManager{
		ops:      make(map[string]int64),
		filePath: path,
	}
}

func (m *OpManager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, err := os.ReadFile(m.filePath)
	if err != nil {
		if os.IsNotExist(err) {

			m.ops = make(map[string]int64)
			return m.saveInternal()
		}
		return err
	}

	var entries []OpEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		logger.Error("Failed to parse ops.json", "error", err)
		return err
	}

	m.ops = make(map[string]int64)
	for _, entry := range entries {
		m.ops[strings.ToLower(entry.Name)] = entry.ClientID
	}

	logger.Info("Loaded OPs", "count", len(m.ops))
	return nil
}

func (m *OpManager) Save() error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.saveInternal()
}

func (m *OpManager) saveInternal() error {
	entries := []OpEntry{}
	for name, cid := range m.ops {
		entries = append(entries, OpEntry{Name: name, ClientID: cid})
	}

	data, err := json.MarshalIndent(entries, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(m.filePath, data, 0644)
}

func (m *OpManager) IsOp(name string, cid int64) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	expectedCid, ok := m.ops[strings.ToLower(name)]
	if !ok {
		return false
	}
	return expectedCid == cid
}

func (m *OpManager) AddOp(name string, cid int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ops[strings.ToLower(name)] = cid
	m.saveInternal()
}

func (m *OpManager) RemoveOp(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.ops, strings.ToLower(name))
	m.saveInternal()
}

func (m *OpManager) IsOpByName(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.ops[strings.ToLower(name)]
	return ok
}
