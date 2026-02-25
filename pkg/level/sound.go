package level

import "github.com/scaxe/scaxe-go/pkg/protocol"

// NewSoundPacket 创建通用声音效果包
// 对应 PHP GenericSound::encode()
// soundID 使用 protocol.EventSound* 常量, pitch 为音调 (0.0 ~ 1.0+)
func NewSoundPacket(x, y, z float32, soundID int16, pitch float32) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(soundID)
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = int32(pitch * 1000) // PHP: $this->pitch = (float) $pitch * 1000
	return pk
}

// ---- 便利函数：常用声音 ----

// NewClickSound 创建点击声音
func NewClickSound(x, y, z float32, pitch float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundClick, pitch)
}

// NewClickFailSound 创建点击失败声音
func NewClickFailSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundClickFail, 0)
}

// NewShootSound 创建射击声音（弓箭、投掷等）
func NewShootSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundShoot, 0)
}

// NewDoorSound 创建门开关声音
func NewDoorSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundDoor, 0)
}

// NewFizzSound 创建嘶嘶声（水浇岩浆等）
func NewFizzSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundFizz, 0)
}

// NewTNTSound 创建 TNT 点燃声音
func NewTNTSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundTNT, 0)
}

// NewGhastSound 创建恶魂声音
func NewGhastSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundGhast, 0)
}

// NewGhastShootSound 创建恶魂射击声音
func NewGhastShootSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundGhastShoot, 0)
}

// NewBlazeShootSound 创建烈焰人射击声音
func NewBlazeShootSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundBlazeShoot, 0)
}

// NewDoorBumpSound 创建门碰撞声音（僵尸破门）
func NewDoorBumpSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundDoorBump, 0)
}

// NewDoorCrashSound 创建门破碎声音
func NewDoorCrashSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundDoorCrash, 0)
}

// NewBatSound 创建蝙蝠声音
func NewBatSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundBatFly, 0)
}

// NewZombieInfectSound 创建僵尸感染声音（村民被感染）
func NewZombieInfectSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundZombieInfect, 0)
}

// NewZombieHealSound 创建僵尸治愈声音（僵尸村民治愈）
func NewZombieHealSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundZombieHeal, 0)
}

// NewEndermanTeleportSound 创建末影人传送声音
func NewEndermanTeleportSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundEndermanTeleport, 0)
}

// NewAnvilBreakSound 创建铁砧碎裂声音
func NewAnvilBreakSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundAnvilBreak, 0)
}

// NewAnvilUseSound 创建铁砧使用声音
func NewAnvilUseSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundAnvilUse, 0)
}

// NewAnvilFallSound 创建铁砧掉落声音
func NewAnvilFallSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundAnvilFall, 0)
}

// NewPopSound 创建弹出声音（物品拾取等）
func NewPopSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundPop, 0)
}

// NewLaunchSound 创建投射声音
func NewLaunchSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundThrowProjectile, 0)
}

// NewItemFrameAddItemSound 创建物品展示框添加物品声音
func NewItemFrameAddItemSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundItemFrameAddItem, 0)
}

// NewItemFrameRemoveSound 创建物品展示框移除声音
func NewItemFrameRemoveSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundItemFrameRemove, 0)
}

// NewItemFramePlaceSound 创建物品展示框放置声音
func NewItemFramePlaceSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundItemFramePlace, 0)
}

// NewItemFrameRemoveItemSound 创建物品展示框取出物品声音
func NewItemFrameRemoveItemSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundItemFrameRemoveItem, 0)
}

// NewItemFrameRotateItemSound 创建物品展示框旋转物品声音
func NewItemFrameRotateItemSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundItemFrameRotateItem, 0)
}

// NewButtonClickSound 创建按钮点击声音
func NewButtonClickSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundButtonClick, 0)
}

// NewExplodeSound 创建爆炸声音
func NewExplodeSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundExplode, 0)
}

// NewSpellSound 创建施法声音
func NewSpellSound(x, y, z float32, pitch float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundSpell, pitch)
}

// NewSplashSound 创建飞溅声音（水花）
func NewSplashSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundSplash, 0)
}

// NewGraySplashSound 创建灰色飞溅声音
func NewGraySplashSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewSoundPacket(x, y, z, protocol.EventSoundGraySplash, 0)
}

// NewTNTPrimeSound 创建 TNT 引燃声音
func NewTNTPrimeSound(x, y, z float32) *protocol.LevelEventPacket {
	return NewTNTSound(x, y, z)
}

// NewNoteblockSound 创建音符盒声音
// instrument: 0=钢琴 1=低音鼓 2=鼓 3=棍棒 4=低音吉他
// pitch: 音高 0-24
func NewNoteblockSound(x, y, z float32, instrument int, pitch int) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(protocol.EventSoundClick) // 音符盒通过 BlockEventPacket 处理，这里用 LevelEvent 作为备用
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = int32(instrument<<8 | pitch)
	return pk
}
