package protocol

const (
	EventSoundClick               int16 = 1000
	EventSoundClickFail           int16 = 1001
	EventSoundShoot               int16 = 1002
	EventSoundDoor                int16 = 1003
	EventSoundFizz                int16 = 1004
	EventSoundTNT                 int16 = 1005
	EventSoundGhast               int16 = 1007
	EventSoundGhastShoot          int16 = 1008
	EventSoundBlazeShoot          int16 = 1009
	EventSoundDoorBump            int16 = 1010
	EventSoundDoorCrash           int16 = 1012
	EventSoundBatFly              int16 = 1015
	EventSoundZombieInfect        int16 = 1016
	EventSoundZombieHeal          int16 = 1017
	EventSoundEndermanTeleport    int16 = 1018
	EventSoundAnvilBreak          int16 = 1020
	EventSoundAnvilUse            int16 = 1021
	EventSoundAnvilFall           int16 = 1022
	EventSoundPop                 int16 = 1030
	EventSoundThrowProjectile     int16 = 1031
	EventSoundItemFrameAddItem    int16 = 1040
	EventSoundItemFrameRemove     int16 = 1041
	EventSoundItemFramePlace      int16 = 1042
	EventSoundItemFrameRemoveItem int16 = 1043
	EventSoundItemFrameRotateItem int16 = 1044
	EventSoundCamera              int16 = 1050
	EventParticleShoot            int16 = 2000
	EventParticleDestroyBlock     int16 = 2001
	EventParticleSplash           int16 = 2002
	EventParticleEyeDespawn       int16 = 2003
	EventParticleSpawn            int16 = 2004
	EventParticleGuardianCurse    int16 = 2006
	EventParticleBlockForceField  int16 = 2008
	EventStartRain                int16 = 3001
	EventStartThunder             int16 = 3002
	EventStopRain                 int16 = 3003
	EventStopThunder              int16 = 3004
	EventSoundButtonClick         int16 = 3500
	EventCauldronExplode          int16 = 3501
	EventCauldronDyeArmor         int16 = 3502
	EventCauldronCleanArmor       int16 = 3503
	EventCauldronFillPotion       int16 = 3504
	EventCauldronTakePotion       int16 = 3505
	EventCauldronFillWater        int16 = 3506
	EventCauldronTakeWater        int16 = 3507
	EventCauldronAddDye           int16 = 3508
	EventSoundExplode             int16 = 3501
	EventSoundSpell               int16 = 3504
	EventSoundSplash              int16 = 3506
	EventSoundGraySplash          int16 = 3507
	EventSetData                  int16 = 4000
	EventPlayersSleeping          int16 = 9800
	EventAddParticleMask          int16 = 0x4000
)

type LevelEventPacket struct {
	BasePacket
	EventID uint16
	X       float32
	Y       float32
	Z       float32
	Data    int32
}

func NewLevelEventPacket() *LevelEventPacket {
	return &LevelEventPacket{
		BasePacket: BasePacket{PacketID: IDLevelEvent},
	}
}

func (p *LevelEventPacket) Name() string {
	return "LevelEventPacket"
}

func (p *LevelEventPacket) Encode(stream *BinaryStream) error {
	EncodeHeader(stream, p.ID())
	stream.WriteBEUShort(p.EventID)
	stream.WriteFloat(p.X)
	stream.WriteFloat(p.Y)
	stream.WriteFloat(p.Z)
	stream.WriteInt(p.Data)
	return nil
}

func (p *LevelEventPacket) Decode(stream *BinaryStream) error {
	DecodeHeader(stream, p.ID())
	var err error
	p.EventID, err = stream.ReadBEUShort()
	if err != nil {
		return err
	}
	p.X, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.Y, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.Z, err = stream.ReadFloat()
	if err != nil {
		return err
	}
	p.Data, err = stream.ReadInt()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	RegisterPacket(IDLevelEvent, func() DataPacket { return NewLevelEventPacket() })
}
