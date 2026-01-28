package defaults

import (
	"time"

	"github.com/scaxe/scaxe-go/pkg/command"
	"github.com/scaxe/scaxe-go/pkg/player"
	"github.com/scaxe/scaxe-go/pkg/protocol"
)

type ServerInterface interface {
	GetOnlinePlayers() []*player.Player
	GetMaxPlayers() int
	GetStartTime() time.Time
	GetRakNetSessionCount() int

	GetTPS() float64
	GetAverageTPS() float64
	GetMSPT() float64

	AddOp(name string, cid int64)
	RemoveOp(name string)
	IsOp(name string) bool

	Stop()

	GetPlayerByName(name string) command.PlayerSender

	BroadcastMessage(message string)
	BroadcastPacket(pk protocol.DataPacket)

	GetTime() int
	SetTime(time int)

	GetDifficulty() int
	SetDifficulty(difficulty int)

	GetSeed() int64

	SetPlayerGamemode(playerName string, gamemode int)

	GetLevelManager() LevelManager
}

type LevelManager interface {
	GetLevelNames() []string
	GetLevel(name string) interface{}
	GetDefaultLevel() interface{}
	LoadLevel(name string) (interface{}, error)
	GenerateLevel(name string, generatorName string, seed int64) (interface{}, error)
	UnloadLevel(name string) bool
}

type Server = ServerInterface
