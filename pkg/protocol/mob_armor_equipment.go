package protocol

import "github.com/scaxe/scaxe-go/pkg/logger"

func init() {
	RegisterPacket(IDMobArmorEquipment, func() DataPacket { return NewMobArmorEquipmentPacket() })
}

type ArmorItem struct {
	ID    int16
	Count int8
	Meta  uint16
}

type MobArmorEquipmentPacket struct {
	BasePacket
	EntityID int64
	Slots    [4]ArmorItem
}

func NewMobArmorEquipmentPacket() *MobArmorEquipmentPacket {
	return &MobArmorEquipmentPacket{
		BasePacket: BasePacket{PacketID: IDMobArmorEquipment},
	}
}

func (p *MobArmorEquipmentPacket) ID() byte {
	return IDMobArmorEquipment
}

func (p *MobArmorEquipmentPacket) Name() string {
	return "MobArmorEquipmentPacket"
}

func writeSlotItem(stream *BinaryStream, item ArmorItem) {
	if item.ID <= 0 {
		stream.WriteBEShort(0)
		return
	}
	stream.WriteBEShort(item.ID)
	stream.WriteByte(uint8(item.Count))
	stream.WriteBEUShort(item.Meta)
	stream.WriteLEShort(0)
}

func readSlotItem(stream *BinaryStream) ArmorItem {
	item := ArmorItem{}
	id, err := stream.ReadBEShort()
	if err != nil {
		return item
	}
	item.ID = id
	if item.ID > 0 {
		count, _ := stream.ReadByte()
		item.Count = int8(count)
		item.Meta, _ = stream.ReadBEUShort()
		nbtLen, _ := stream.ReadLEShort()
		if nbtLen > 0 {
			stream.Skip(int(nbtLen))
		}
	}
	return item
}

func (p *MobArmorEquipmentPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteBELong(p.EntityID)
	for i := 0; i < 4; i++ {
		writeSlotItem(stream, p.Slots[i])
	}
	return nil
}

func (p *MobArmorEquipmentPacket) Decode(stream *BinaryStream) error {
	var err error
	p.EntityID, err = stream.ReadBELong()
	if err != nil {
		logger.Debug("MobArmorEquipment: failed to read EntityID", "error", err)
		return nil
	}
	for i := 0; i < 4; i++ {
		p.Slots[i] = readSlotItem(stream)
	}
	return nil
}
