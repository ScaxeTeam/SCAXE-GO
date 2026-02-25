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

type PlayerPreLoginEvent struct {
	*PlayerEvent
	Address     string
	Port        int
	KickMessage string
}

var playerPreLoginHandlers = NewHandlerList()

func NewPlayerPreLoginEvent(playerName string, playerID int64, address string, port int) *PlayerPreLoginEvent {
	return &PlayerPreLoginEvent{
		PlayerEvent: NewPlayerEvent("PlayerPreLoginEvent", playerName, playerID),
		Address:     address,
		Port:        port,
	}
}

func (e *PlayerPreLoginEvent) GetHandlers() *HandlerList { return playerPreLoginHandlers }
func (e *PlayerPreLoginEvent) GetKickMessage() string    { return e.KickMessage }
func (e *PlayerPreLoginEvent) SetKickMessage(msg string) { e.KickMessage = msg }

type PlayerLoginEvent struct {
	*PlayerEvent
	KickMessage string
}

var playerLoginHandlers = NewHandlerList()

func NewPlayerLoginEvent(playerName string, playerID int64) *PlayerLoginEvent {
	return &PlayerLoginEvent{
		PlayerEvent: NewPlayerEvent("PlayerLoginEvent", playerName, playerID),
	}
}

func (e *PlayerLoginEvent) GetHandlers() *HandlerList { return playerLoginHandlers }
func (e *PlayerLoginEvent) GetKickMessage() string    { return e.KickMessage }
func (e *PlayerLoginEvent) SetKickMessage(msg string) { e.KickMessage = msg }

type PlayerCreationEvent struct {
	*PlayerEvent
	Address string
	Port    int
}

var playerCreationHandlers = NewHandlerList()

func NewPlayerCreationEvent(address string, port int) *PlayerCreationEvent {
	return &PlayerCreationEvent{
		PlayerEvent: NewPlayerEvent("PlayerCreationEvent", "", 0),
		Address:     address,
		Port:        port,
	}
}

func (e *PlayerCreationEvent) GetHandlers() *HandlerList { return playerCreationHandlers }

type PlayerRespawnEvent struct {
	*PlayerEvent
	X, Y, Z float64
}

var playerRespawnHandlers = NewHandlerList()

func NewPlayerRespawnEvent(playerName string, playerID int64, x, y, z float64) *PlayerRespawnEvent {
	return &PlayerRespawnEvent{
		PlayerEvent: NewPlayerEvent("PlayerRespawnEvent", playerName, playerID),
		X:           x, Y: y, Z: z,
	}
}

func (e *PlayerRespawnEvent) GetHandlers() *HandlerList { return playerRespawnHandlers }

type PlayerKickEvent struct {
	*PlayerEvent
	Reason string
}

var playerKickHandlers = NewHandlerList()

func NewPlayerKickEvent(playerName string, playerID int64, reason string) *PlayerKickEvent {
	return &PlayerKickEvent{
		PlayerEvent: NewPlayerEvent("PlayerKickEvent", playerName, playerID),
		Reason:      reason,
	}
}

func (e *PlayerKickEvent) GetHandlers() *HandlerList { return playerKickHandlers }

type PlayerTransferEvent struct {
	*PlayerEvent
	Address string
	Port    int
}

var playerTransferHandlers = NewHandlerList()

func NewPlayerTransferEvent(playerName string, playerID int64, address string, port int) *PlayerTransferEvent {
	return &PlayerTransferEvent{
		PlayerEvent: NewPlayerEvent("PlayerTransferEvent", playerName, playerID),
		Address:     address,
		Port:        port,
	}
}

func (e *PlayerTransferEvent) GetHandlers() *HandlerList { return playerTransferHandlers }

type PlayerDropItemEvent struct {
	*PlayerEvent
	ItemID    int
	ItemMeta  int
	ItemCount int
}

var playerDropItemHandlers = NewHandlerList()

func NewPlayerDropItemEvent(playerName string, playerID int64, itemID, itemMeta, itemCount int) *PlayerDropItemEvent {
	return &PlayerDropItemEvent{
		PlayerEvent: NewPlayerEvent("PlayerDropItemEvent", playerName, playerID),
		ItemID:      itemID,
		ItemMeta:    itemMeta,
		ItemCount:   itemCount,
	}
}

func (e *PlayerDropItemEvent) GetHandlers() *HandlerList { return playerDropItemHandlers }

type PlayerItemHeldEvent struct {
	*PlayerEvent
	ItemID     int
	Slot       int
	HotbarSlot int
}

var playerItemHeldHandlers = NewHandlerList()

func NewPlayerItemHeldEvent(playerName string, playerID int64, itemID, slot, hotbarSlot int) *PlayerItemHeldEvent {
	return &PlayerItemHeldEvent{
		PlayerEvent: NewPlayerEvent("PlayerItemHeldEvent", playerName, playerID),
		ItemID:      itemID,
		Slot:        slot,
		HotbarSlot:  hotbarSlot,
	}
}

func (e *PlayerItemHeldEvent) GetHandlers() *HandlerList { return playerItemHeldHandlers }

type PlayerItemConsumeEvent struct {
	*PlayerEvent
	ItemID   int
	ItemMeta int
}

var playerItemConsumeHandlers = NewHandlerList()

func NewPlayerItemConsumeEvent(playerName string, playerID int64, itemID, itemMeta int) *PlayerItemConsumeEvent {
	return &PlayerItemConsumeEvent{
		PlayerEvent: NewPlayerEvent("PlayerItemConsumeEvent", playerName, playerID),
		ItemID:      itemID,
		ItemMeta:    itemMeta,
	}
}

func (e *PlayerItemConsumeEvent) GetHandlers() *HandlerList { return playerItemConsumeHandlers }

type PlayerCommandPreprocessEvent struct {
	*PlayerEvent
	Message string
}

var playerCommandPreprocessHandlers = NewHandlerList()

func NewPlayerCommandPreprocessEvent(playerName string, playerID int64, message string) *PlayerCommandPreprocessEvent {
	return &PlayerCommandPreprocessEvent{
		PlayerEvent: NewPlayerEvent("PlayerCommandPreprocessEvent", playerName, playerID),
		Message:     message,
	}
}

func (e *PlayerCommandPreprocessEvent) GetHandlers() *HandlerList {
	return playerCommandPreprocessHandlers
}
func (e *PlayerCommandPreprocessEvent) GetMessage() string    { return e.Message }
func (e *PlayerCommandPreprocessEvent) SetMessage(msg string) { e.Message = msg }

type PlayerGameModeChangeEvent struct {
	*PlayerEvent
	NewGameMode int
}

var playerGameModeChangeHandlers = NewHandlerList()

func NewPlayerGameModeChangeEvent(playerName string, playerID int64, newGameMode int) *PlayerGameModeChangeEvent {
	return &PlayerGameModeChangeEvent{
		PlayerEvent: NewPlayerEvent("PlayerGameModeChangeEvent", playerName, playerID),
		NewGameMode: newGameMode,
	}
}

func (e *PlayerGameModeChangeEvent) GetHandlers() *HandlerList { return playerGameModeChangeHandlers }

type PlayerToggleSneakEvent struct {
	*PlayerEvent
	IsSneaking bool
}

var playerToggleSneakHandlers = NewHandlerList()

func NewPlayerToggleSneakEvent(playerName string, playerID int64, isSneaking bool) *PlayerToggleSneakEvent {
	return &PlayerToggleSneakEvent{
		PlayerEvent: NewPlayerEvent("PlayerToggleSneakEvent", playerName, playerID),
		IsSneaking:  isSneaking,
	}
}

func (e *PlayerToggleSneakEvent) GetHandlers() *HandlerList { return playerToggleSneakHandlers }

type PlayerToggleSprintEvent struct {
	*PlayerEvent
	IsSprinting bool
}

var playerToggleSprintHandlers = NewHandlerList()

func NewPlayerToggleSprintEvent(playerName string, playerID int64, isSprinting bool) *PlayerToggleSprintEvent {
	return &PlayerToggleSprintEvent{
		PlayerEvent: NewPlayerEvent("PlayerToggleSprintEvent", playerName, playerID),
		IsSprinting: isSprinting,
	}
}

func (e *PlayerToggleSprintEvent) GetHandlers() *HandlerList { return playerToggleSprintHandlers }

type PlayerJumpEvent struct {
	*PlayerEvent
}

var playerJumpHandlers = NewHandlerList()

func NewPlayerJumpEvent(playerName string, playerID int64) *PlayerJumpEvent {
	return &PlayerJumpEvent{
		PlayerEvent: NewPlayerEvent("PlayerJumpEvent", playerName, playerID),
	}
}

func (e *PlayerJumpEvent) GetHandlers() *HandlerList { return playerJumpHandlers }

type PlayerAnimationEvent struct {
	*PlayerEvent
	AnimationType int32
}

var playerAnimationHandlers = NewHandlerList()

func NewPlayerAnimationEvent(playerName string, playerID int64, animationType int32) *PlayerAnimationEvent {
	return &PlayerAnimationEvent{
		PlayerEvent:   NewPlayerEvent("PlayerAnimationEvent", playerName, playerID),
		AnimationType: animationType,
	}
}

func (e *PlayerAnimationEvent) GetHandlers() *HandlerList { return playerAnimationHandlers }

type PlayerBucketEmptyEvent struct {
	*PlayerEvent
	BucketID               int
	BlockX, BlockY, BlockZ int
	Face                   int
}

var playerBucketEmptyHandlers = NewHandlerList()

func NewPlayerBucketEmptyEvent(playerName string, playerID int64, bucketID, bx, by, bz, face int) *PlayerBucketEmptyEvent {
	return &PlayerBucketEmptyEvent{
		PlayerEvent: NewPlayerEvent("PlayerBucketEmptyEvent", playerName, playerID),
		BucketID:    bucketID,
		BlockX:      bx, BlockY: by, BlockZ: bz,
		Face: face,
	}
}

func (e *PlayerBucketEmptyEvent) GetHandlers() *HandlerList { return playerBucketEmptyHandlers }

type PlayerBucketFillEvent struct {
	*PlayerEvent
	BucketID               int
	BlockX, BlockY, BlockZ int
	Face                   int
}

var playerBucketFillHandlers = NewHandlerList()

func NewPlayerBucketFillEvent(playerName string, playerID int64, bucketID, bx, by, bz, face int) *PlayerBucketFillEvent {
	return &PlayerBucketFillEvent{
		PlayerEvent: NewPlayerEvent("PlayerBucketFillEvent", playerName, playerID),
		BucketID:    bucketID,
		BlockX:      bx, BlockY: by, BlockZ: bz,
		Face: face,
	}
}

func (e *PlayerBucketFillEvent) GetHandlers() *HandlerList { return playerBucketFillHandlers }

type PlayerGlassBottleEvent struct {
	*PlayerEvent
	BlockX, BlockY, BlockZ int
}

var playerGlassBottleHandlers = NewHandlerList()

func NewPlayerGlassBottleEvent(playerName string, playerID int64, bx, by, bz int) *PlayerGlassBottleEvent {
	return &PlayerGlassBottleEvent{
		PlayerEvent: NewPlayerEvent("PlayerGlassBottleEvent", playerName, playerID),
		BlockX:      bx, BlockY: by, BlockZ: bz,
	}
}

func (e *PlayerGlassBottleEvent) GetHandlers() *HandlerList { return playerGlassBottleHandlers }

type PlayerBedEnterEvent struct {
	*PlayerEvent
	BlockX, BlockY, BlockZ int
}

var playerBedEnterHandlers = NewHandlerList()

func NewPlayerBedEnterEvent(playerName string, playerID int64, bx, by, bz int) *PlayerBedEnterEvent {
	return &PlayerBedEnterEvent{
		PlayerEvent: NewPlayerEvent("PlayerBedEnterEvent", playerName, playerID),
		BlockX:      bx, BlockY: by, BlockZ: bz,
	}
}

func (e *PlayerBedEnterEvent) GetHandlers() *HandlerList { return playerBedEnterHandlers }

type PlayerBedLeaveEvent struct {
	*PlayerEvent
	BlockX, BlockY, BlockZ int
}

var playerBedLeaveHandlers = NewHandlerList()

func NewPlayerBedLeaveEvent(playerName string, playerID int64, bx, by, bz int) *PlayerBedLeaveEvent {
	return &PlayerBedLeaveEvent{
		PlayerEvent: NewPlayerEvent("PlayerBedLeaveEvent", playerName, playerID),
		BlockX:      bx, BlockY: by, BlockZ: bz,
	}
}

func (e *PlayerBedLeaveEvent) GetHandlers() *HandlerList { return playerBedLeaveHandlers }

type PlayerExhaustEvent struct {
	*PlayerEvent
	Amount float64
	Cause  int
}

const (
	ExhaustCauseAttack     = 0
	ExhaustCauseDamage     = 1
	ExhaustCauseMining     = 2
	ExhaustCauseSprint     = 3
	ExhaustCauseJump       = 4
	ExhaustCauseSwim       = 5
	ExhaustCauseSprintJump = 6
	ExhaustCauseRegenerate = 7
)

var playerExhaustHandlers = NewHandlerList()

func NewPlayerExhaustEvent(playerName string, playerID int64, amount float64, cause int) *PlayerExhaustEvent {
	return &PlayerExhaustEvent{
		PlayerEvent: NewPlayerEvent("PlayerExhaustEvent", playerName, playerID),
		Amount:      amount,
		Cause:       cause,
	}
}

func (e *PlayerExhaustEvent) GetHandlers() *HandlerList { return playerExhaustHandlers }
func (e *PlayerExhaustEvent) SetAmount(a float64)       { e.Amount = a }

type PlayerExperienceChangeEvent struct {
	*PlayerEvent
	OldLevel int
	OldExp   float64
	NewLevel int
	NewExp   float64
}

var playerExperienceChangeHandlers = NewHandlerList()

func NewPlayerExperienceChangeEvent(playerName string, playerID int64, oldLevel int, oldExp float64, newLevel int, newExp float64) *PlayerExperienceChangeEvent {
	return &PlayerExperienceChangeEvent{
		PlayerEvent: NewPlayerEvent("PlayerExperienceChangeEvent", playerName, playerID),
		OldLevel:    oldLevel, OldExp: oldExp,
		NewLevel: newLevel, NewExp: newExp,
	}
}

func (e *PlayerExperienceChangeEvent) GetHandlers() *HandlerList {
	return playerExperienceChangeHandlers
}

type PlayerHungerChangeEvent struct {
	*PlayerEvent
	OldHunger int
	NewHunger int
}

var playerHungerChangeHandlers = NewHandlerList()

func NewPlayerHungerChangeEvent(playerName string, playerID int64, oldHunger, newHunger int) *PlayerHungerChangeEvent {
	return &PlayerHungerChangeEvent{
		PlayerEvent: NewPlayerEvent("PlayerHungerChangeEvent", playerName, playerID),
		OldHunger:   oldHunger,
		NewHunger:   newHunger,
	}
}

func (e *PlayerHungerChangeEvent) GetHandlers() *HandlerList { return playerHungerChangeHandlers }

type PlayerAchievementAwardedEvent struct {
	*PlayerEvent
	Achievement string
}

var playerAchievementAwardedHandlers = NewHandlerList()

func NewPlayerAchievementAwardedEvent(playerName string, playerID int64, achievement string) *PlayerAchievementAwardedEvent {
	return &PlayerAchievementAwardedEvent{
		PlayerEvent: NewPlayerEvent("PlayerAchievementAwardedEvent", playerName, playerID),
		Achievement: achievement,
	}
}

func (e *PlayerAchievementAwardedEvent) GetHandlers() *HandlerList {
	return playerAchievementAwardedHandlers
}

type PlayerPickupExpOrbEvent struct {
	*PlayerEvent
	ExpAmount int
}

var playerPickupExpOrbHandlers = NewHandlerList()

func NewPlayerPickupExpOrbEvent(playerName string, playerID int64, expAmount int) *PlayerPickupExpOrbEvent {
	return &PlayerPickupExpOrbEvent{
		PlayerEvent: NewPlayerEvent("PlayerPickupExpOrbEvent", playerName, playerID),
		ExpAmount:   expAmount,
	}
}

func (e *PlayerPickupExpOrbEvent) GetHandlers() *HandlerList { return playerPickupExpOrbHandlers }

type PlayerTextPreSendEvent struct {
	*PlayerEvent
	Message string
	Type    int
}

var playerTextPreSendHandlers = NewHandlerList()

func NewPlayerTextPreSendEvent(playerName string, playerID int64, message string, textType int) *PlayerTextPreSendEvent {
	return &PlayerTextPreSendEvent{
		PlayerEvent: NewPlayerEvent("PlayerTextPreSendEvent", playerName, playerID),
		Message:     message,
		Type:        textType,
	}
}

func (e *PlayerTextPreSendEvent) GetHandlers() *HandlerList { return playerTextPreSendHandlers }

type PlayerFishEvent struct {
	*PlayerEvent
	ItemID    int
	ItemMeta  int
	ItemCount int
}

var playerFishHandlers = NewHandlerList()

func NewPlayerFishEvent(playerName string, playerID int64, itemID, itemMeta, itemCount int) *PlayerFishEvent {
	return &PlayerFishEvent{
		PlayerEvent: NewPlayerEvent("PlayerFishEvent", playerName, playerID),
		ItemID:      itemID,
		ItemMeta:    itemMeta,
		ItemCount:   itemCount,
	}
}

func (e *PlayerFishEvent) GetHandlers() *HandlerList { return playerFishHandlers }

const (
	FishingRodActionCast = 0
	FishingRodActionReel = 1
)

type PlayerUseFishingRodEvent struct {
	*PlayerEvent
	Action int
}

var playerUseFishingRodHandlers = NewHandlerList()

func NewPlayerUseFishingRodEvent(playerName string, playerID int64, action int) *PlayerUseFishingRodEvent {
	return &PlayerUseFishingRodEvent{
		PlayerEvent: NewPlayerEvent("PlayerUseFishingRodEvent", playerName, playerID),
		Action:      action,
	}
}

func (e *PlayerUseFishingRodEvent) GetHandlers() *HandlerList { return playerUseFishingRodHandlers }
