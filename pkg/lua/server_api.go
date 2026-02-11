package lua

import "github.com/scaxe/scaxe-go/pkg/command"

type ServerAPI interface {
	BroadcastMessage(message string)
	GetOnlineCount() int
	GetMaxPlayers() int
	GetTPS() float64
	GetServerName() string
	GetPlayer(username string) PlayerAPI
	GetOnlinePlayers() []PlayerAPI
	KickPlayer(username string, reason string)
	GetLevel() LevelAPI
	RegisterCommand(cmd command.Command)
	UnregisterCommand(name string)
	Stop()
	GetCurrentTick() int64
}

type PlayerAPI interface {
	GetName() string
	GetPosition() (x, y, z float64)
	SetPosition(x, y, z float64)
	SendMessage(msg string)
	GetGamemode() int
	SetGamemode(mode int)
	IsOp() bool
	Kick(reason string)
	GetHealth() int
	SetHealth(health int)
	GetEntityID() int64
}

type LevelAPI interface {
	GetBlock(x, y, z int32) (id, meta uint8)
	SetBlock(x, y, z int32, id, meta uint8)
	GetTime() int64
	SetTime(time int64)
	GetSeed() int64
	GetSpawnLocation() (x, y, z float64)
}
