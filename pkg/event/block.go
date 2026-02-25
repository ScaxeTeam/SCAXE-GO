package event

type BlockEvent struct {
	*BaseEvent
	BlockX, BlockY, BlockZ int
	BlockID                int
	BlockMeta              int
}

func NewBlockEvent(name string, x, y, z, id, meta int) *BlockEvent {
	return &BlockEvent{
		BaseEvent: NewBaseEvent(name),
		BlockX:    x, BlockY: y, BlockZ: z,
		BlockID: id, BlockMeta: meta,
	}
}

func (e *BlockEvent) GetBlockPosition() (x, y, z int) {
	return e.BlockX, e.BlockY, e.BlockZ
}

type BlockBreakEvent struct {
	*BlockEvent
	PlayerID  int64
	ItemID    int
	Drops     []interface{}
	ExpToDrop int
	FastBreak bool
}

var blockBreakHandlers = NewHandlerList()

func NewBlockBreakEvent(x, y, z, blockID, blockMeta int, playerID int64, itemID int) *BlockBreakEvent {
	return &BlockBreakEvent{
		BlockEvent: NewBlockEvent("BlockBreakEvent", x, y, z, blockID, blockMeta),
		PlayerID:   playerID,
		ItemID:     itemID,
		Drops:      []interface{}{},
		ExpToDrop:  0,
		FastBreak:  false,
	}
}

func (e *BlockBreakEvent) GetHandlers() *HandlerList {
	return blockBreakHandlers
}

func (e *BlockBreakEvent) GetPlayerID() int64 {
	return e.PlayerID
}

func (e *BlockBreakEvent) GetDrops() []interface{} {
	return e.Drops
}

func (e *BlockBreakEvent) SetDrops(drops []interface{}) {
	e.Drops = drops
}

type BlockPlaceEvent struct {
	*BlockEvent
	PlayerID     int64
	ItemID       int
	ReplacedID   int
	ReplacedMeta int
}

var blockPlaceHandlers = NewHandlerList()

func NewBlockPlaceEvent(x, y, z, blockID, blockMeta int, playerID int64, itemID, replacedID, replacedMeta int) *BlockPlaceEvent {
	return &BlockPlaceEvent{
		BlockEvent:   NewBlockEvent("BlockPlaceEvent", x, y, z, blockID, blockMeta),
		PlayerID:     playerID,
		ItemID:       itemID,
		ReplacedID:   replacedID,
		ReplacedMeta: replacedMeta,
	}
}

func (e *BlockPlaceEvent) GetHandlers() *HandlerList {
	return blockPlaceHandlers
}

type BlockUpdateEvent struct {
	*BlockEvent
}

var blockUpdateHandlers = NewHandlerList()

func NewBlockUpdateEvent(x, y, z, blockID, blockMeta int) *BlockUpdateEvent {
	return &BlockUpdateEvent{
		BlockEvent: NewBlockEvent("BlockUpdateEvent", x, y, z, blockID, blockMeta),
	}
}

func (e *BlockUpdateEvent) GetHandlers() *HandlerList {
	return blockUpdateHandlers
}

type SignChangeEvent struct {
	*BlockEvent
	PlayerID int64
	Lines    [4]string
}

var signChangeHandlers = NewHandlerList()

func NewSignChangeEvent(x, y, z, blockID, blockMeta int, playerID int64, lines [4]string) *SignChangeEvent {
	return &SignChangeEvent{
		BlockEvent: NewBlockEvent("SignChangeEvent", x, y, z, blockID, blockMeta),
		PlayerID:   playerID,
		Lines:      lines,
	}
}

func (e *SignChangeEvent) GetHandlers() *HandlerList {
	return signChangeHandlers
}

func (e *SignChangeEvent) GetLine(index int) string {
	if index >= 0 && index < 4 {
		return e.Lines[index]
	}
	return ""
}

func (e *SignChangeEvent) SetLine(index int, text string) {
	if index >= 0 && index < 4 {
		e.Lines[index] = text
	}
}

// ---------- BlockBurnEvent ----------
// Fired when a block is destroyed by fire.

type BlockBurnEvent struct {
	*BlockEvent
}

var blockBurnHandlers = NewHandlerList()

func NewBlockBurnEvent(x, y, z, blockID, blockMeta int) *BlockBurnEvent {
	return &BlockBurnEvent{
		BlockEvent: NewBlockEvent("BlockBurnEvent", x, y, z, blockID, blockMeta),
	}
}

func (e *BlockBurnEvent) GetHandlers() *HandlerList { return blockBurnHandlers }

// ---------- BlockFormEvent ----------
// Fired when a block forms naturally (e.g. snow layer from snowfall, obsidian from lava+water).

type BlockFormEvent struct {
	*BlockEvent
	NewBlockID   int
	NewBlockMeta int
}

var blockFormHandlers = NewHandlerList()

func NewBlockFormEvent(x, y, z, blockID, blockMeta, newBlockID, newBlockMeta int) *BlockFormEvent {
	return &BlockFormEvent{
		BlockEvent:   NewBlockEvent("BlockFormEvent", x, y, z, blockID, blockMeta),
		NewBlockID:   newBlockID,
		NewBlockMeta: newBlockMeta,
	}
}

func (e *BlockFormEvent) GetHandlers() *HandlerList { return blockFormHandlers }

// ---------- BlockGrowEvent ----------
// Fired when a crop or plant block grows.

type BlockGrowEvent struct {
	*BlockEvent
	NewBlockID   int
	NewBlockMeta int
}

var blockGrowHandlers = NewHandlerList()

func NewBlockGrowEvent(x, y, z, blockID, blockMeta, newBlockID, newBlockMeta int) *BlockGrowEvent {
	return &BlockGrowEvent{
		BlockEvent:   NewBlockEvent("BlockGrowEvent", x, y, z, blockID, blockMeta),
		NewBlockID:   newBlockID,
		NewBlockMeta: newBlockMeta,
	}
}

func (e *BlockGrowEvent) GetHandlers() *HandlerList { return blockGrowHandlers }

// ---------- BlockSpreadEvent ----------
// Fired when a block spreads to another position (e.g. fire, grass, mycelium).

type BlockSpreadEvent struct {
	*BlockEvent
	SourceX, SourceY, SourceZ int
	NewBlockID                int
	NewBlockMeta              int
}

var blockSpreadHandlers = NewHandlerList()

func NewBlockSpreadEvent(x, y, z, blockID, blockMeta, srcX, srcY, srcZ, newBlockID, newBlockMeta int) *BlockSpreadEvent {
	return &BlockSpreadEvent{
		BlockEvent:   NewBlockEvent("BlockSpreadEvent", x, y, z, blockID, blockMeta),
		SourceX:      srcX,
		SourceY:      srcY,
		SourceZ:      srcZ,
		NewBlockID:   newBlockID,
		NewBlockMeta: newBlockMeta,
	}
}

func (e *BlockSpreadEvent) GetHandlers() *HandlerList { return blockSpreadHandlers }

// ---------- LeavesDecayEvent ----------
// Fired when a leaf block decays naturally.

type LeavesDecayEvent struct {
	*BlockEvent
}

var leavesDecayHandlers = NewHandlerList()

func NewLeavesDecayEvent(x, y, z, blockID, blockMeta int) *LeavesDecayEvent {
	return &LeavesDecayEvent{
		BlockEvent: NewBlockEvent("LeavesDecayEvent", x, y, z, blockID, blockMeta),
	}
}

func (e *LeavesDecayEvent) GetHandlers() *HandlerList { return leavesDecayHandlers }

// ---------- ItemFrameDropItemEvent ----------
// Fired when an item is removed from an item frame.

type ItemFrameDropItemEvent struct {
	*BlockEvent
	PlayerID int64
	ItemID   int
	ItemMeta int
}

var itemFrameDropItemHandlers = NewHandlerList()

func NewItemFrameDropItemEvent(x, y, z, blockID, blockMeta int, playerID int64, itemID, itemMeta int) *ItemFrameDropItemEvent {
	return &ItemFrameDropItemEvent{
		BlockEvent: NewBlockEvent("ItemFrameDropItemEvent", x, y, z, blockID, blockMeta),
		PlayerID:   playerID,
		ItemID:     itemID,
		ItemMeta:   itemMeta,
	}
}

func (e *ItemFrameDropItemEvent) GetHandlers() *HandlerList { return itemFrameDropItemHandlers }
