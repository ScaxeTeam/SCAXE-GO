package event

type InventoryEvent struct {
	*BaseEvent
	InventoryType int
}

func NewInventoryEvent(name string, invType int) *InventoryEvent {
	return &InventoryEvent{
		BaseEvent:     NewBaseEvent(name),
		InventoryType: invType,
	}
}

type InventoryOpenEvent struct {
	*InventoryEvent
	PlayerID int64
}

var inventoryOpenHandlers = NewHandlerList()

func NewInventoryOpenEvent(playerID int64, invType int) *InventoryOpenEvent {
	return &InventoryOpenEvent{
		InventoryEvent: NewInventoryEvent("InventoryOpenEvent", invType),
		PlayerID:       playerID,
	}
}

func (e *InventoryOpenEvent) GetHandlers() *HandlerList {
	return inventoryOpenHandlers
}

type InventoryCloseEvent struct {
	*InventoryEvent
	PlayerID int64
}

var inventoryCloseHandlers = NewHandlerList()

func NewInventoryCloseEvent(playerID int64, invType int) *InventoryCloseEvent {
	return &InventoryCloseEvent{
		InventoryEvent: NewInventoryEvent("InventoryCloseEvent", invType),
		PlayerID:       playerID,
	}
}

func (e *InventoryCloseEvent) GetHandlers() *HandlerList {
	return inventoryCloseHandlers
}

type InventoryTransactionEvent struct {
	*InventoryEvent
	PlayerID  int64
	Slot      int
	OldItemID int
	NewItemID int
}

var inventoryTransactionHandlers = NewHandlerList()

func NewInventoryTransactionEvent(playerID int64, slot, oldItem, newItem int) *InventoryTransactionEvent {
	return &InventoryTransactionEvent{
		InventoryEvent: NewInventoryEvent("InventoryTransactionEvent", 0),
		PlayerID:       playerID,
		Slot:           slot,
		OldItemID:      oldItem,
		NewItemID:      newItem,
	}
}

func (e *InventoryTransactionEvent) GetHandlers() *HandlerList {
	return inventoryTransactionHandlers
}

type CraftItemEvent struct {
	*InventoryEvent
	PlayerID int64
	ResultID int
}

var craftItemHandlers = NewHandlerList()

func NewCraftItemEvent(playerID int64, resultID int) *CraftItemEvent {
	return &CraftItemEvent{
		InventoryEvent: NewInventoryEvent("CraftItemEvent", 0),
		PlayerID:       playerID,
		ResultID:       resultID,
	}
}

func (e *CraftItemEvent) GetHandlers() *HandlerList {
	return craftItemHandlers
}

type FurnaceSmeltEvent struct {
	*InventoryEvent
	BlockX, BlockY, BlockZ int
	SourceID               int
	ResultID               int
}

var furnaceSmeltHandlers = NewHandlerList()

func NewFurnaceSmeltEvent(bx, by, bz, sourceID, resultID int) *FurnaceSmeltEvent {
	return &FurnaceSmeltEvent{
		InventoryEvent: NewInventoryEvent("FurnaceSmeltEvent", 0),
		BlockX:         bx, BlockY: by, BlockZ: bz,
		SourceID: sourceID,
		ResultID: resultID,
	}
}

func (e *FurnaceSmeltEvent) GetHandlers() *HandlerList {
	return furnaceSmeltHandlers
}

type FurnaceBurnEvent struct {
	*InventoryEvent
	BlockX, BlockY, BlockZ int
	FuelID                 int
	BurnTime               int
}

var furnaceBurnHandlers = NewHandlerList()

func NewFurnaceBurnEvent(bx, by, bz, fuelID, burnTime int) *FurnaceBurnEvent {
	return &FurnaceBurnEvent{
		InventoryEvent: NewInventoryEvent("FurnaceBurnEvent", 0),
		BlockX:         bx, BlockY: by, BlockZ: bz,
		FuelID:   fuelID,
		BurnTime: burnTime,
	}
}

func (e *FurnaceBurnEvent) GetHandlers() *HandlerList {
	return furnaceBurnHandlers
}
