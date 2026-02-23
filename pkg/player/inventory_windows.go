package player

import (
	"sync"

	"github.com/scaxe/scaxe-go/pkg/inventory"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

const (
	WindowIDPlayer   byte = 0x00
	WindowIDArmor    byte = 0x78
	WindowIDCreative byte = 0x79

	WindowIDContainerMin byte = 0x01
	WindowIDContainerMax byte = 0x77
)

type InventoryWindows struct {
	mu sync.RWMutex

	windows map[byte]inventory.Inventory

	invToWindow map[inventory.Inventory]byte

	windowCnt byte

	currentWindow inventory.Inventory
}

func NewInventoryWindows() *InventoryWindows {
	return &InventoryWindows{
		windows:     make(map[byte]inventory.Inventory),
		invToWindow: make(map[inventory.Inventory]byte),
		windowCnt:   WindowIDContainerMin,
	}
}

var _ inventory.Viewer = (*Player)(nil)

func (p *Player) GetWindowID(inv inventory.Inventory) byte {
	p.windows.mu.RLock()
	defer p.windows.mu.RUnlock()

	if id, ok := p.windows.invToWindow[inv]; ok {
		return id
	}
	return 0xFF
}

func (p *Player) SendDataPacket(pk interface{}) {
	if dp, ok := pk.(protocol.DataPacket); ok {
		p.SendPacket(dp)
	}
}

func (p *Player) GetViewerID() string {
	return p.Username
}

func (p *Player) AddWindow(inv inventory.Inventory, forceID ...byte) byte {
	p.windows.mu.Lock()
	defer p.windows.mu.Unlock()

	if existingID, ok := p.windows.invToWindow[inv]; ok {
		return existingID
	}

	var windowID byte
	if len(forceID) > 0 {
		windowID = forceID[0]
	} else {

		windowID = p.windows.windowCnt
		p.windows.windowCnt++
		if p.windows.windowCnt > WindowIDContainerMax {
			p.windows.windowCnt = WindowIDContainerMin
		}
	}

	p.windows.windows[windowID] = inv
	p.windows.invToWindow[inv] = windowID

	return windowID
}

func (p *Player) RemoveWindow(inv inventory.Inventory) {
	p.windows.mu.Lock()
	defer p.windows.mu.Unlock()

	if windowID, ok := p.windows.invToWindow[inv]; ok {
		delete(p.windows.windows, windowID)
		delete(p.windows.invToWindow, inv)
	}
}

func (p *Player) GetWindowByID(windowID byte) inventory.Inventory {
	p.windows.mu.RLock()
	defer p.windows.mu.RUnlock()
	return p.windows.windows[windowID]
}

func (p *Player) OpenInventory(inv inventory.Inventory) bool {
	if p.windows.currentWindow != nil {
		p.CloseInventory()
	}

	windowID := p.AddWindow(inv)
	logger.Debug("OpenInventory",
		"player", p.Username,
		"windowID", windowID,
		"type", inv.GetType().GetDefaultTitle())

	p.windows.mu.Lock()
	p.windows.currentWindow = inv
	p.windows.mu.Unlock()

	return inv.Open(p)
}

func (p *Player) CloseInventory() {
	p.windows.mu.Lock()
	current := p.windows.currentWindow
	p.windows.currentWindow = nil
	p.windows.mu.Unlock()

	if current == nil {
		return
	}

	logger.Debug("CloseInventory",
		"player", p.Username,
		"type", current.GetType().GetDefaultTitle())

	current.Close(p)
	p.RemoveWindow(current)
}

func (p *Player) GetCurrentWindow() inventory.Inventory {
	p.windows.mu.RLock()
	defer p.windows.mu.RUnlock()
	return p.windows.currentWindow
}

func (p *Player) HandleContainerClose(windowID byte) {
	if windowID == WindowIDPlayer {
		return
	}

	inv := p.GetWindowByID(windowID)
	if inv == nil {
		return
	}

	p.windows.mu.Lock()
	if p.windows.currentWindow == inv {
		p.windows.currentWindow = nil
	}
	p.windows.mu.Unlock()

	inv.Close(p)
	p.RemoveWindow(inv)
}
