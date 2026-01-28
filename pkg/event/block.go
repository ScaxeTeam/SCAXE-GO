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
