package event

type PlayerEvent struct {
	*BaseEvent
	PlayerName string
	PlayerID   int64
}

func NewPlayerEvent(name string, playerName string, playerID int64) *PlayerEvent {
	return &PlayerEvent{
		BaseEvent:  NewBaseEvent(name),
		PlayerName: playerName,
		PlayerID:   playerID,
	}
}

func (e *PlayerEvent) GetPlayerName() string {
	return e.PlayerName
}

func (e *PlayerEvent) GetPlayerID() int64 {
	return e.PlayerID
}

type PlayerJoinEvent struct {
	*PlayerEvent
	JoinMessage string
}

var playerJoinHandlers = NewHandlerList()

func NewPlayerJoinEvent(playerName string, playerID int64, joinMessage string) *PlayerJoinEvent {
	return &PlayerJoinEvent{
		PlayerEvent: NewPlayerEvent("PlayerJoinEvent", playerName, playerID),
		JoinMessage: joinMessage,
	}
}

func (e *PlayerJoinEvent) GetHandlers() *HandlerList {
	return playerJoinHandlers
}

func (e *PlayerJoinEvent) GetJoinMessage() string {
	return e.JoinMessage
}

func (e *PlayerJoinEvent) SetJoinMessage(msg string) {
	e.JoinMessage = msg
}

type PlayerQuitEvent struct {
	*PlayerEvent
	QuitMessage string
	Reason      string
}

var playerQuitHandlers = NewHandlerList()

func NewPlayerQuitEvent(playerName string, playerID int64, quitMessage, reason string) *PlayerQuitEvent {
	return &PlayerQuitEvent{
		PlayerEvent: NewPlayerEvent("PlayerQuitEvent", playerName, playerID),
		QuitMessage: quitMessage,
		Reason:      reason,
	}
}

func (e *PlayerQuitEvent) GetHandlers() *HandlerList {
	return playerQuitHandlers
}

func (e *PlayerQuitEvent) GetQuitMessage() string {
	return e.QuitMessage
}

func (e *PlayerQuitEvent) SetQuitMessage(msg string) {
	e.QuitMessage = msg
}

type PlayerChatEvent struct {
	*PlayerEvent
	Message    string
	Format     string
	Recipients []int64
}

var playerChatHandlers = NewHandlerList()

func NewPlayerChatEvent(playerName string, playerID int64, message string, recipients []int64) *PlayerChatEvent {
	return &PlayerChatEvent{
		PlayerEvent: NewPlayerEvent("PlayerChatEvent", playerName, playerID),
		Message:     message,
		Format:      "<%s> %s",
		Recipients:  recipients,
	}
}

func (e *PlayerChatEvent) GetHandlers() *HandlerList {
	return playerChatHandlers
}

func (e *PlayerChatEvent) GetMessage() string {
	return e.Message
}

func (e *PlayerChatEvent) SetMessage(msg string) {
	e.Message = msg
}

func (e *PlayerChatEvent) GetFormat() string {
	return e.Format
}

func (e *PlayerChatEvent) SetFormat(format string) {
	e.Format = format
}

type PlayerMoveEvent struct {
	*PlayerEvent
	FromX, FromY, FromZ float64
	ToX, ToY, ToZ       float64
}

var playerMoveHandlers = NewHandlerList()

func NewPlayerMoveEvent(playerName string, playerID int64, fromX, fromY, fromZ, toX, toY, toZ float64) *PlayerMoveEvent {
	return &PlayerMoveEvent{
		PlayerEvent: NewPlayerEvent("PlayerMoveEvent", playerName, playerID),
		FromX:       fromX, FromY: fromY, FromZ: fromZ,
		ToX: toX, ToY: toY, ToZ: toZ,
	}
}

func (e *PlayerMoveEvent) GetHandlers() *HandlerList {
	return playerMoveHandlers
}

type PlayerDeathEvent struct {
	*PlayerEvent
	DeathMessage   string
	KeepInventory  bool
	KeepExperience bool
}

var playerDeathHandlers = NewHandlerList()

func NewPlayerDeathEvent(playerName string, playerID int64, deathMessage string) *PlayerDeathEvent {
	return &PlayerDeathEvent{
		PlayerEvent:    NewPlayerEvent("PlayerDeathEvent", playerName, playerID),
		DeathMessage:   deathMessage,
		KeepInventory:  false,
		KeepExperience: false,
	}
}

func (e *PlayerDeathEvent) GetHandlers() *HandlerList {
	return playerDeathHandlers
}

func (e *PlayerDeathEvent) GetDeathMessage() string {
	return e.DeathMessage
}

func (e *PlayerDeathEvent) SetDeathMessage(msg string) {
	e.DeathMessage = msg
}

const (
	ActionLeftClickBlock  = 0
	ActionRightClickBlock = 1
	ActionLeftClickAir    = 2
	ActionRightClickAir   = 3
)

type PlayerInteractEvent struct {
	*PlayerEvent
	Action                 int
	BlockX, BlockY, BlockZ int
	Face                   int
	ItemID                 int
}

var playerInteractHandlers = NewHandlerList()

func NewPlayerInteractEvent(playerName string, playerID int64, action int, bx, by, bz, face, itemID int) *PlayerInteractEvent {
	return &PlayerInteractEvent{
		PlayerEvent: NewPlayerEvent("PlayerInteractEvent", playerName, playerID),
		Action:      action,
		BlockX:      bx, BlockY: by, BlockZ: bz,
		Face:   face,
		ItemID: itemID,
	}
}

func (e *PlayerInteractEvent) GetHandlers() *HandlerList {
	return playerInteractHandlers
}
