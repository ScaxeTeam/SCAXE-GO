package protocol

import (
	"github.com/scaxe/scaxe-go/pkg/item"
)

type UseItemPacket struct {
	BasePacket
	X, Y, Z          int32
	Face             byte
	FX, FY, FZ       float32
	PosX, PosY, PosZ float32
	Slot             int32
	Item             item.Item
}

func NewUseItemPacket() *UseItemPacket {
	return &UseItemPacket{
		BasePacket: BasePacket{PacketID: IDUseItem},
	}
}

func (p *UseItemPacket) Name() string {
	return "UseItemPacket"
}

func (p *UseItemPacket) Encode(stream *BinaryStream) error {
	return nil
}

func (p *UseItemPacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())
	var err error
	p.X, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.Y, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.Z, err = stream.ReadInt()
	if err != nil {
		return err
	}

	p.Face, err = stream.ReadByte()
	if err != nil {
		return err
	}

	p.FX, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.FY, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.FZ, err = stream.ReadFloat()
	if err != nil {
		return err
	}

	p.PosX, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.PosY, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.PosZ, err = stream.ReadFloat()
	if err != nil {
		return err
	}

	if ProtocolCurrent >= 70 {
		p.Slot, err = stream.ReadInt()
		if err != nil {
			return err
		}
	}

	p.Item, err = readSlot(stream)
	if err != nil {
		return err
	}

	return nil
}

func readSlot(stream *BinaryStream) (item.Item, error) {
	id, err := stream.ReadShort()
	if err != nil {
		return item.Item{}, err
	}
	if id <= 0 {
		return item.Item{ID: 0}, nil
	}

	cnt, err := stream.ReadByte()
	if err != nil {
		return item.Item{}, err
	}

	damage, err := stream.ReadShort()
	if err != nil {
		return item.Item{}, err
	}

	nbtLen, err := stream.ReadLShort()
	if err != nil {
		return item.Item{}, err
	}

	if nbtLen > 0 {
		_, err = stream.ReadBytes(int(nbtLen))
		if err != nil {
			return item.Item{}, err
		}

	}

	return item.NewItem(int(id), int(damage), int(cnt)), nil
}
