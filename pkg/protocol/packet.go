package protocol

import (
	"fmt"

	"github.com/scaxe/scaxe-go/pkg/logger"
)

const (
	ProtocolCurrent = 70
	ProtocolMax     = 84

	Protocol_0_12_First = 34
	Protocol_0_13_First = 38
	Protocol_0_14_First = 41
	Protocol_0_15_First = 81
)

var AcceptedProtocols = []int{

	34,

	37, 38, 39,

	41, 42, 43, 44, 45, 46, 60, 70,

	81, 82, 83, 84,
}

func IsProtocolSupported(protocol int) bool {
	for _, p := range AcceptedProtocols {
		if p == protocol {
			return true
		}
	}
	return false
}

const (
	IDLogin                  byte = 0x8f
	IDPlayStatus             byte = 0x90
	IDDisconnect             byte = 0x91
	IDBatch                  byte = 0x92
	IDText                   byte = 0x93
	IDSetTime                byte = 0x94
	IDStartGame              byte = 0x95
	IDAddPlayer              byte = 0x96
	IDRemovePlayer           byte = 0x97
	IDAddEntity              byte = 0x98
	IDRemoveEntity           byte = 0x99
	IDAddItemEntity          byte = 0x9a
	IDTakeItemEntity         byte = 0x9b
	IDMoveEntity             byte = 0x9c
	IDMovePlayer             byte = 0x9d
	IDRemoveBlock            byte = 0x9e
	IDUpdateBlock            byte = 0x9f
	IDAddPainting            byte = 0xa0
	IDExplode                byte = 0xa1
	IDLevelEvent             byte = 0xa2
	IDBlockEvent             byte = 0xa3
	IDEntityEvent            byte = 0xa4
	IDMobEffect              byte = 0xa5
	IDUpdateAttributes       byte = 0xa6
	IDMobEquipment           byte = 0xa7
	IDMobArmorEquipment      byte = 0xa8
	IDInteract               byte = 0xa9
	IDUseItem                byte = 0xaa
	IDPlayerAction           byte = 0xab
	IDHurtArmor              byte = 0xac
	IDSetEntityData          byte = 0xad
	IDSetEntityMotion        byte = 0xae
	IDSetEntityLink          byte = 0xaf
	IDSetHealth              byte = 0xb0
	IDSetSpawnPosition       byte = 0xb1
	IDAnimate                byte = 0xb2
	IDRespawn                byte = 0xb3
	IDDropItem               byte = 0xb4
	IDContainerOpen          byte = 0xb5
	IDContainerClose         byte = 0xb6
	IDContainerSetSlot       byte = 0xb7
	IDContainerSetData       byte = 0xb8
	IDContainerSetContent    byte = 0xb9
	IDCraftingData           byte = 0xba
	IDCraftingEvent          byte = 0xbb
	IDAdventureSettings      byte = 0xbc
	IDBlockEntityData        byte = 0xbd
	IDPlayerInput            byte = 0xbe
	IDFullChunkData          byte = 0xbf
	IDSetDifficulty          byte = 0xc0
	IDChangeDimension        byte = 0xc1
	IDSetPlayerGameType      byte = 0xc2
	IDPlayerList             byte = 0xc3
	IDTelemetryEvent         byte = 0xc4
	IDClientboundMapItemData byte = 0xc6
	IDMapInfoRequest         byte = 0xc7
	IDRequestChunkRadius     byte = 0xc8
	IDChunkRadiusUpdated     byte = 0xc9
	IDItemFrameDropItem      byte = 0xca
	IDReplaceSelectedItem    byte = 0xcb
)

const RakNetWrapperID byte = 0x8e

type DataPacket interface {
	ID() byte
	Name() string
	Encode(stream *BinaryStream) error
	Decode(stream *BinaryStream) error
}

type BasePacket struct {
	PacketID byte
}

func (p *BasePacket) ID() byte {
	return p.PacketID
}

func EncodeHeader(stream *BinaryStream, id byte) {
	stream.WriteByte(id)
	logger.Debug("EncodeHeader", "packetID", id)
}

func DecodeHeader(stream *BinaryStream, expectedID byte) error {
	id, err := stream.ReadByte()
	if err != nil {
		return err
	}
	if id != expectedID {
		logger.Warn("DecodeHeader", "expected", expectedID, "got", id)
	}
	logger.Debug("DecodeHeader", "packetID", id)
	return nil
}

var PacketPool = make(map[byte]func() DataPacket)

func RegisterPacket(id byte, constructor func() DataPacket) {
	PacketPool[id] = constructor
	logger.Debug("RegisterPacket", "id", id, "hex", fmt.Sprintf("0x%02x", id))
}

func init() {
	RegisterPacket(IDRespawn, func() DataPacket { return &RespawnPacket{} })
	RegisterPacket(IDDropItem, func() DataPacket { return &DropItemPacket{} })
}

func GetPacket(id byte) DataPacket {
	constructor, ok := PacketPool[id]
	if !ok {
		logger.Warn("GetPacket", "id", id, "hex", fmt.Sprintf("0x%02x", id), "error", "unknown packet ID")
		return nil
	}
	return constructor()
}
