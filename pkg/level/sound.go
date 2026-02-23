package level

import "github.com/scaxe/scaxe-go/pkg/protocol"

func NewSoundPacket(x, y, z float32, soundID int16, pitch float32) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(soundID)
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = int32(pitch * 1000)
	return pk
}

func NewClickSound(x, y, z float32, pitch float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundClick, pitch)
}

func NewClickFailSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundClickFail, 0)
}

func NewShootSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundShoot, 0)
}

func NewDoorSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundDoor, 0)
}

func NewFizzSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundFizz, 0)
}

func NewTNTSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundTNT, 0)
}

func NewGhastSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundGhast, 0)
}

func NewGhastShootSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundGhastShoot, 0)
}

func NewBlazeShootSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundBlazeShoot, 0)
}

func NewDoorBumpSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundDoorBump, 0)
}

func NewDoorCrashSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundDoorCrash, 0)
}

func NewBatSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundBatFly, 0)
}

func NewZombieInfectSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundZombieInfect, 0)
}

func NewZombieHealSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundZombieHeal, 0)
}

func NewEndermanTeleportSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundEndermanTeleport, 0)
}

func NewAnvilBreakSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundAnvilBreak, 0)
}

func NewAnvilUseSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundAnvilUse, 0)
}

func NewAnvilFallSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundAnvilFall, 0)
}

func NewPopSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundPop, 0)
}

func NewLaunchSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundThrowProjectile, 0)
}

func NewItemFrameAddItemSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundItemFrameAddItem, 0)
}

func NewItemFrameRemoveSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundItemFrameRemove, 0)
}

func NewItemFramePlaceSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundItemFramePlace, 0)
}

func NewItemFrameRemoveItemSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundItemFrameRemoveItem, 0)
}

func NewItemFrameRotateItemSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundItemFrameRotateItem, 0)
}

func NewButtonClickSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundButtonClick, 0)
}

func NewExplodeSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundExplode, 0)
}

func NewSpellSound(x, y, z float32, pitch float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundSpell, pitch)
}

func NewSplashSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundSplash, 0)
}

func NewGraySplashSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundGraySplash, 0)
}

func NewTNTPrimeSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewTNTSound(x, y, z)
}

func NewNoteblockSound(x, y, z float32, instrument int, pitch int) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(protocol.EventSoundClick)
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = int32(instrument<<8 | pitch)
	return pk
}
