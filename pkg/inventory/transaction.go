package inventory

import (
	"time"

	"github.com/scaxe/scaxe-go/pkg/item"
	"github.com/scaxe/scaxe-go/pkg/logger"
)

const (
	MaxTransactionsPerGroup = 50
	TransactionTimeout = 8.0
)
type Transaction struct {
	inventory    Inventory
	slot         int
	sourceItem   item.Item
	targetItem   item.Item
	creationTime float64
}
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
type TransactionGroup struct {
	source       Viewer
	transactions []*Transaction
	inventories  map[Inventory]bool
	hasExecuted  bool
	creationTime float64
	OnExecute func(g *TransactionGroup) bool
}
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
func (g *TransactionGroup) GetInventories() []Inventory {
	invs := make([]Inventory, 0, len(g.inventories))
	for inv := range g.inventories {
		invs = append(invs, inv)
	}
	return invs
}
func (g *TransactionGroup) AddTransaction(tx *Transaction) bool {
	if len(g.transactions) >= MaxTransactionsPerGroup {
		logger.Warn("TransactionGroup: too many transactions, rejecting",
			"count", len(g.transactions),
			"max", MaxTransactionsPerGroup)
		return false
	}
	slot := tx.GetSlot()
	invSize := tx.GetInventory().GetSize()
	maxSlot := invSize
	if _, ok := tx.GetInventory().(*PlayerInventory); ok {
		maxSlot = playerTotalSize
	}
	if slot < 0 || slot >= maxSlot {
		logger.Warn("TransactionGroup: invalid slot",
			"slot", slot,
			"maxSlot", maxSlot,
			"inventory", tx.GetInventory().GetName())
		return false
	}
	for i, existing := range g.transactions {
		if existing.GetInventory() == tx.GetInventory() && existing.GetSlot() == tx.GetSlot() {
			if tx.GetCreationTime() >= existing.GetCreationTime() {
				g.transactions = append(g.transactions[:i], g.transactions[i+1:]...)
				break
			} else {
				return false
			}
		}
	}

	g.transactions = append(g.transactions, tx)
	g.inventories[tx.GetInventory()] = true
	return true
}
func (g *TransactionGroup) matchItems() (needItems, haveItems []item.Item, ok bool) {
	for _, tx := range g.transactions {
		targetItem := tx.GetTargetItem()
		if targetItem.ID != 0 {
			needItems = append(needItems, targetItem)
		}
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
func (g *TransactionGroup) CanExecute() bool {
	if len(g.transactions) == 0 {
		return false
	}
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
func (g *TransactionGroup) Execute() bool {
	if g.hasExecuted || !g.CanExecute() {
		g.SendInventories()
		return false
	}
	if g.OnExecute != nil {
		if !g.OnExecute(g) {
			logger.Debug("TransactionGroup: cancelled by event hook")
			g.SendInventories()
			return false
		}
	}
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
func (g *TransactionGroup) SendInventories() {
	if g.source == nil {
		return
	}
	for inv := range g.inventories {
		if playerInv, ok := inv.(*PlayerInventory); ok {
			playerInv.SendArmorContents(g.source)
		}
		inv.SendContents(g.source)
	}
}
