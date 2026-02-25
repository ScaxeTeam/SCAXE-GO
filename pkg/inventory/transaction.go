package inventory

import (
	"time"

	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/logger"
)

// ── 防刷物常量 ──────────────────────────────────────

const (
	// MaxTransactionsPerGroup 单次交易组最大交易数（防构造异常包）
	MaxTransactionsPerGroup = 50

	// TransactionTimeout 交易组超时时间（秒），超时后拒绝执行
	TransactionTimeout = 8.0
)

// ── Transaction ─────────────────────────────────────

// Transaction represents a single slot change in an inventory:
// sourceItem (what was in the slot) → targetItem (what should be placed).
type Transaction struct {
	inventory    Inventory
	slot         int
	sourceItem   item.Item
	targetItem   item.Item
	creationTime float64
}

// NewTransaction creates a new inventory transaction.
// 对应 PHP BaseTransaction::__construct()
func NewTransaction(inv Inventory, slot int, sourceItem, targetItem item.Item) *Transaction {
	return &Transaction{
		inventory:    inv,
		slot:         slot,
		sourceItem:   sourceItem,
		targetItem:   targetItem,
		creationTime: float64(time.Now().UnixNano()) / 1e9,
	}
}

func (t *Transaction) GetInventory() Inventory  { return t.inventory }
func (t *Transaction) GetSlot() int             { return t.slot }
func (t *Transaction) GetSourceItem() item.Item { return t.sourceItem }
func (t *Transaction) GetTargetItem() item.Item { return t.targetItem }
func (t *Transaction) GetCreationTime() float64 { return t.creationTime }

// ── TransactionGroup ────────────────────────────────

// TransactionGroup collects multiple transactions and executes them atomically
// if the item balance is valid (items in == items out).
// 对应 PHP SimpleTransactionGroup
type TransactionGroup struct {
	source       Viewer
	transactions []*Transaction
	inventories  map[Inventory]bool
	hasExecuted  bool
	creationTime float64

	// ── 事件钩子 ──
	// OnExecute 在 Execute 成功执行前调用。
	// 返回 false 则取消执行（对应 PHP InventoryTransactionEvent 的 cancel 机制）。
	OnExecute func(g *TransactionGroup) bool
}

// NewTransactionGroup creates a new transaction group for the given source player.
// 对应 PHP SimpleTransactionGroup::__construct()
func NewTransactionGroup(source Viewer) *TransactionGroup {
	return &TransactionGroup{
		source:       source,
		transactions: make([]*Transaction, 0),
		inventories:  make(map[Inventory]bool),
		creationTime: float64(time.Now().UnixNano()) / 1e9,
	}
}

func (g *TransactionGroup) GetSource() Viewer               { return g.source }
func (g *TransactionGroup) GetCreationTime() float64        { return g.creationTime }
func (g *TransactionGroup) GetTransactions() []*Transaction { return g.transactions }
func (g *TransactionGroup) HasExecuted() bool               { return g.hasExecuted }

// GetInventories returns all inventories involved in this transaction group.
func (g *TransactionGroup) GetInventories() []Inventory {
	invs := make([]Inventory, 0, len(g.inventories))
	for inv := range g.inventories {
		invs = append(invs, inv)
	}
	return invs
}

// ── 交易添加（含防刷物验证） ──────────────────────────

// AddTransaction adds a transaction to the group.
// If a transaction for the same inventory+slot already exists, the newer one wins.
// 对应 PHP SimpleTransactionGroup::addTransaction()
//
// 防刷物增强：
//   - 拒绝超过 MaxTransactionsPerGroup 的交易
//   - 验证槽位范围有效性
func (g *TransactionGroup) AddTransaction(tx *Transaction) bool {
	// 防刷物: 超过最大交易数
	if len(g.transactions) >= MaxTransactionsPerGroup {
		logger.Warn("TransactionGroup: too many transactions, rejecting",
			"count", len(g.transactions),
			"max", MaxTransactionsPerGroup)
		return false
	}

	// 防刷物: 验证槽位范围
	slot := tx.GetSlot()
	invSize := tx.GetInventory().GetSize()
	// 对 PlayerInventory 需要允许护甲栏 (36-39)
	maxSlot := invSize
	if _, ok := tx.GetInventory().(*PlayerInventory); ok {
		maxSlot = playerTotalSize // 40
	}
	if slot < 0 || slot >= maxSlot {
		logger.Warn("TransactionGroup: invalid slot",
			"slot", slot,
			"maxSlot", maxSlot,
			"inventory", tx.GetInventory().GetName())
		return false
	}

	// 去重: same inventory + same slot → keep the newer one
	// 对应 PHP SimpleTransactionGroup::addTransaction() L78-86
	for i, existing := range g.transactions {
		if existing.GetInventory() == tx.GetInventory() && existing.GetSlot() == tx.GetSlot() {
			if tx.GetCreationTime() >= existing.GetCreationTime() {
				// 替换旧交易
				g.transactions = append(g.transactions[:i], g.transactions[i+1:]...)
				break
			} else {
				// 现有的更新，跳过
				return false
			}
		}
	}

	g.transactions = append(g.transactions, tx)
	g.inventories[tx.GetInventory()] = true
	return true
}

// ── 物品平衡验证 ──────────────────────────────────

// matchItems verifies that all source items in slots match the actual current
// inventory contents, and collects "have" (source) and "need" (target) items.
// Returns true if source items are valid, false if there's a mismatch.
// 对应 PHP SimpleTransactionGroup::matchItems()
func (g *TransactionGroup) matchItems() (needItems, haveItems []item.Item, ok bool) {
	for _, tx := range g.transactions {
		targetItem := tx.GetTargetItem()
		if targetItem.ID != 0 {
			needItems = append(needItems, targetItem)
		}

		// 验证 source item 是否与背包中实际物品匹配
		checkItem := tx.GetInventory().GetItem(tx.GetSlot())
		sourceItem := tx.GetSourceItem()
		if !checkItem.Equals(sourceItem, true, true) || sourceItem.Count != checkItem.Count {
			logger.Debug("Transaction mismatch",
				"slot", tx.GetSlot(),
				"expected", sourceItem.ID,
				"expectedCount", sourceItem.Count,
				"actual", checkItem.ID,
				"actualCount", checkItem.Count)
			return nil, nil, false
		}

		if sourceItem.ID != 0 {
			haveItems = append(haveItems, sourceItem)
		}
	}

	// 平衡检查: need 和 have 必须完全匹配
	// 对应 PHP SimpleTransactionGroup::matchItems() L112-127
	for i := 0; i < len(needItems); i++ {
		for j := 0; j < len(haveItems); j++ {
			if needItems[i].Equals(haveItems[j], true, true) {
				amount := needItems[i].Count
				if haveItems[j].Count < amount {
					amount = haveItems[j].Count
				}
				needItems[i].Count -= amount
				haveItems[j].Count -= amount

				if haveItems[j].Count == 0 {
					haveItems = append(haveItems[:j], haveItems[j+1:]...)
					j--
				}
				if needItems[i].Count == 0 {
					needItems = append(needItems[:i], needItems[i+1:]...)
					i--
					break
				}
			}
		}
	}

	return needItems, haveItems, true
}

// ── 执行验证 ──────────────────────────────────────

// CanExecute checks if the transaction group is balanced, valid, and not expired.
// 对应 PHP SimpleTransactionGroup::canExecute()
//
// 防刷物增强：
//   - 检查交易组是否超时
func (g *TransactionGroup) CanExecute() bool {
	if len(g.transactions) == 0 {
		return false
	}

	// 防刷物: 超时检查
	now := float64(time.Now().UnixNano()) / 1e9
	if now-g.creationTime > TransactionTimeout {
		logger.Warn("TransactionGroup: expired",
			"age", now-g.creationTime,
			"timeout", TransactionTimeout)
		return false
	}

	needItems, haveItems, ok := g.matchItems()
	return ok && len(needItems) == 0 && len(haveItems) == 0
}

// Execute applies all transactions if the group is valid and balanced.
// Returns true if successful, false if rejected.
// 对应 PHP SimpleTransactionGroup::execute()
func (g *TransactionGroup) Execute() bool {
	if g.hasExecuted || !g.CanExecute() {
		g.SendInventories()
		return false
	}

	// 事件钩子: 对应 PHP InventoryTransactionEvent
	if g.OnExecute != nil {
		if !g.OnExecute(g) {
			logger.Debug("TransactionGroup: cancelled by event hook")
			g.SendInventories()
			return false
		}
	}

	// 应用所有交易
	for _, tx := range g.transactions {
		if err := tx.GetInventory().SetItem(tx.GetSlot(), tx.GetTargetItem()); err != nil {
			logger.Error("Transaction failed",
				"inventory", tx.GetInventory().GetName(),
				"slot", tx.GetSlot(),
				"item", tx.GetTargetItem().ID,
				"error", err)
		}
	}

	g.hasExecuted = true
	return true
}

// ── 重同步 ──────────────────────────────────────

// SendInventories re-sends all affected inventory contents to the source player,
// used when a transaction is rejected to resync the client.
// 对应 PHP SimpleTransactionGroup::sendInventories()
//
// 增强: 如果涉及 PlayerInventory，额外发送 ArmorContents
func (g *TransactionGroup) SendInventories() {
	if g.source == nil {
		return
	}
	for inv := range g.inventories {
		// PHP: if($inventory instanceof PlayerInventory) { $inventory->sendArmorContents($this->getSource()); }
		if playerInv, ok := inv.(*PlayerInventory); ok {
			playerInv.SendArmorContents(g.source)
		}
		inv.SendContents(g.source)
	}
}
