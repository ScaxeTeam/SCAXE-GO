package protocol

type AdventureSettingsPacket struct {
	BasePacket
	Flags            int32
	UserPermission   int32
	GlobalPermission int32
}

const (
	AdventureFlagWorldImmutable   = 1 << 0
	AdventureFlagNoPlayerVsPlayer = 1 << 1
	AdventureFlagAutoJump         = 1 << 5
	AdventureFlagAllowFlight      = 1 << 6
	AdventureFlagNoClip           = 1 << 7
	AdventureFlagFlying           = 1 << 9
)

const (
	PermissionMember   = 0
	PermissionOperator = 1
	PermissionCustom   = 2
)

func NewAdventureSettingsPacket() *AdventureSettingsPacket {
	return &AdventureSettingsPacket{
		BasePacket: BasePacket{PacketID: IDAdventureSettings},
	}
}

func (p *AdventureSettingsPacket) ID() byte {
	return IDAdventureSettings
}

func (p *AdventureSettingsPacket) Name() string {
	return "AdventureSettingsPacket"
}

func (p *AdventureSettingsPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteInt(p.Flags)
	stream.WriteInt(p.UserPermission)
	stream.WriteInt(p.GlobalPermission)
	return nil
}

func (p *AdventureSettingsPacket) Decode(stream *BinaryStream) error {
	var err error
	p.Flags, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.UserPermission, err = stream.ReadInt()
	if err != nil {
		return err
	}
	p.GlobalPermission, err = stream.ReadInt()
	return err
}

func init() {
	RegisterPacket(IDAdventureSettings, func() DataPacket {
		return NewAdventureSettingsPacket()
	})
}
