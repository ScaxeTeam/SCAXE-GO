package protocol

import "github.com/scaxe/scaxe-go/pkg/logger"

func init() {
	RegisterPacket(IDMobEquipment, func() DataPacket { return NewMobEquipmentPacket() })
}

type MobEquipmentPacket struct {
	BasePacket
	EntityID     int64
	ItemID       int16
	ItemCount    int8
	ItemMeta     uint16
	Slot         uint8
	SelectedSlot uint8
}

func NewMobEquipmentPacket() *MobEquipmentPacket {
	return &MobEquipmentPacket{
		BasePacket: BasePacket{PacketID: IDMobEquipment},
	}
}

func (p *MobEquipmentPacket) ID() byte {
	return IDMobEquipment
}

func (p *MobEquipmentPacket) Name() string {
	return "MobEquipmentPacket"
}

func (p *MobEquipmentPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())

	stream.WriteBELong(p.EntityID)

	if p.ItemID <= 0 {
		stream.WriteBEShort(0)
	} else {
		stream.WriteBEShort(p.ItemID)
		stream.WriteByte(uint8(p.ItemCount))
		stream.WriteBEUShort(p.ItemMeta)
		stream.WriteLEShort(0)
	}
	stream.WriteByte(p.Slot)
	stream.WriteByte(p.SelectedSlot)
	return nil
}

func (p *MobEquipmentPacket) Decode(stream *BinaryStream) error {
	var err error

	p.EntityID, err = stream.ReadBELong()
	if err != nil {
		logger.Debug("MobEquipment: failed to read EntityID", "error", err)
		return nil
	}

	p.ItemID, err = stream.ReadBEShort()
	if err != nil {
		logger.Debug("MobEquipment: failed to read ItemID", "error", err)
		return nil
	}

	if p.ItemID > 0 {
		count, _ := stream.ReadByte()
		p.ItemCount = int8(count)
		p.ItemMeta, _ = stream.ReadBEUShort()

		nbtLen, _ := stream.ReadLEShort()
		if nbtLen > 0 {
			stream.Skip(int(nbtLen))
		}
	}

	p.Slot, _ = stream.ReadByte()
	p.SelectedSlot, _ = stream.ReadByte()

	return nil
}
