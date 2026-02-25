package player

// inventory_windows.go — 玩家容器窗口管理系统
// 对应 PHP Player.php 中的 addWindow/removeWindow/openInventory/closeInventory
//
// MCPE 背包协议流程：
//   1. 玩家右键箱子/熔炉等 → 服务端调用 Player.OpenInventory(inv)
//   2. Player 分配一个 windowID → AddWindow
//   3. ContainerInventory.Open(player) → 发送 ContainerOpenPacket + ContainerSetContentPacket
//   4. 客户端显示背包 UI
//   5. 客户端拖动物品 → 发送 ContainerSetSlotPacket
//   6. 服务端验证并执行 Transaction
//   7. 玩家关闭背包 → 客户端发 ContainerClosePacket → HandleContainerClose

import (
	"sync"

	"github.com/scaxe/scaxe-go/pkg/inventory"
	"github.com/scaxe/scaxe-go/pkg/logger"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

// 窗口 ID 常量
const (
	WindowIDPlayer   byte = 0x00
	WindowIDArmor    byte = 0x78 // 120
	WindowIDCreative byte = 0x79 // 121

	// 容器窗口的起始/结束 ID（递增分配）
	WindowIDContainerMin byte = 0x01
	WindowIDContainerMax byte = 0x77 // 119
)

// InventoryWindows 管理玩家打开的容器窗口
type InventoryWindows struct {
	mu sync.RWMutex

	// windows: windowID → Inventory
	windows map[byte]inventory.Inventory

	// invToWindow: Inventory → windowID（反向映射）
	invToWindow map[inventory.Inventory]byte

	// windowCnt: 下一个可分配的窗口 ID
	windowCnt byte

	// currentWindow: 当前打开的容器窗口
	currentWindow inventory.Inventory
}

// NewInventoryWindows 创建窗口管理器
func NewInventoryWindows() *InventoryWindows {
	return &InventoryWindows{
		windows:     make(map[byte]inventory.Inventory),
		invToWindow: make(map[inventory.Inventory]byte),
		windowCnt:   WindowIDContainerMin,
	}
}

// ── inventory.Viewer 接口实现 ─────────────────────────

// 编译期验证 Player 实现了 inventory.Viewer
var _ inventory.Viewer = (*Player)(nil)

// GetWindowID 返回指定背包的窗口 ID。未找到返回 0xFF。
func (p *Player) GetWindowID(inv inventory.Inventory) byte {
	p.windows.mu.RLock()
	defer p.windows.mu.RUnlock()

	if id, ok := p.windows.invToWindow[inv]; ok {
		return id
	}
	return 0xFF
}

// SendDataPacket 向玩家发送数据包。
// inventory 包通过 interface{} 避免循环依赖。
func (p *Player) SendDataPacket(pk interface{}) {
	if dp, ok := pk.(protocol.DataPacket); ok {
		p.SendPacket(dp)
	}
}

// GetViewerID 返回玩家唯一标识符。
func (p *Player) GetViewerID() string {
	return p.Username
}

// 注意：IsSpawned() 已在 player.go 中定义。

// ── 窗口管理 ─────────────────────────────────

// AddWindow 向玩家注册一个背包窗口，返回分配的 windowID。
// 对应 PHP Player::addWindow()
func (p *Player) AddWindow(inv inventory.Inventory, forceID ...byte) byte {
	p.windows.mu.Lock()
	defer p.windows.mu.Unlock()

	// 已注册的背包直接返回现有 ID
	if existingID, ok := p.windows.invToWindow[inv]; ok {
		return existingID
	}

	var windowID byte
	if len(forceID) > 0 {
		windowID = forceID[0]
	} else {
		// 自动分配 ID
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

// RemoveWindow 移除一个背包窗口。
// 对应 PHP Player::removeWindow()
func (p *Player) RemoveWindow(inv inventory.Inventory) {
	p.windows.mu.Lock()
	defer p.windows.mu.Unlock()

	if windowID, ok := p.windows.invToWindow[inv]; ok {
		delete(p.windows.windows, windowID)
		delete(p.windows.invToWindow, inv)
	}
}

// GetWindowByID 根据窗口 ID 获取背包。
func (p *Player) GetWindowByID(windowID byte) inventory.Inventory {
	p.windows.mu.RLock()
	defer p.windows.mu.RUnlock()
	return p.windows.windows[windowID]
}

// ── 打开/关闭容器 ─────────────────────────────────

// OpenInventory 打开一个容器背包。
// 分配窗口ID → 调用 inv.Open(player) → 发送 OpenPacket + ContentsPacket
// 对应 PHP Player 打开容器的完整流程。
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

// CloseInventory 关闭当前打开的容器背包。
// 对应 PHP Player::removeWindow() + ContainerInventory::onClose()
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

// GetCurrentWindow 获取当前打开的容器窗口（nil = 无）。
func (p *Player) GetCurrentWindow() inventory.Inventory {
	p.windows.mu.RLock()
	defer p.windows.mu.RUnlock()
	return p.windows.currentWindow
}

// HandleContainerClose 处理客户端发来的 ContainerClosePacket。
// 对应 PHP Player 中对 ContainerClosePacket 的处理。
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
