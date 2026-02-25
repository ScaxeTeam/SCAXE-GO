package level

import "github.com/scaxe/scaxe-go/pkg/protocol"

// 粒子类型常量（对应 PHP Particle::TYPE_* ）
// 使用方法：通过 NewParticlePacket 生成 LevelEventPacket 发送给客户端
const (
	ParticleBubble                int = 1
	ParticleCritical              int = 2
	ParticleSmoke                 int = 3
	ParticleExplode               int = 4
	ParticleWhiteSmoke            int = 5
	ParticleFlame                 int = 6
	ParticleLava                  int = 7
	ParticleLargeSmoke            int = 8
	ParticleRedstone              int = 9
	ParticleItemBreak             int = 10
	ParticleSnowballPoof          int = 11
	ParticleLargeExplode          int = 12
	ParticleHugeExplode           int = 13
	ParticleMobFlame              int = 14
	ParticleHeart                 int = 15
	ParticleTerrain               int = 16
	ParticleTownAura              int = 17
	ParticlePortal                int = 18
	ParticleWaterSplash           int = 19
	ParticleWaterWake             int = 20
	ParticleDripWater             int = 21
	ParticleDripLava              int = 22
	ParticleDust                  int = 23
	ParticleMobSpell              int = 24
	ParticleMobSpellAmbient       int = 25
	ParticleMobSpellInstantaneous int = 26
	ParticleInk                   int = 27
	ParticleSlime                 int = 28
	ParticleRainSplash            int = 29
	ParticleVillagerAngry         int = 30
	ParticleVillagerHappy         int = 31
	ParticleEnchantmentTable      int = 32
)

// NewParticlePacket 创建通用粒子效果包
// 对应 PHP GenericParticle::encode()
// particleType 为上方 Particle* 常量，data 为附加数据（如颜色）
func NewParticlePacket(x, y, z float32, particleType int, data int32) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(protocol.EventAddParticleMask) | uint16(particleType&0xFFF)
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = data
	return pk
}

// NewDestroyBlockParticle 创建方块破坏粒子效果包
// 对应 PHP DestroyBlockParticle::encode()
// blockID 和 blockMeta 用于确定碎片外观
func NewDestroyBlockParticle(x, y, z float32, blockID int, blockMeta int) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(protocol.EventParticleDestroyBlock)
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = int32(blockID) + int32(blockMeta<<12)
	return pk
}

// NewShootParticle 创建射击粒子效果包
func NewShootParticle(x, y, z float32, data int32) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(protocol.EventParticleShoot)
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = data
	return pk
}

// NewSplashParticle 创建飞溅粒子效果包
func NewSplashParticle(x, y, z float32, data int32) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(protocol.EventParticleSplash)
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = data
	return pk
}

// NewSpawnParticle 创建生成粒子效果包（用于刷怪笼、传送门等）
func NewSpawnParticle(x, y, z float32, data int32) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(protocol.EventParticleSpawn)
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = data
	return pk
}

// ---- 便利函数：常用特定粒子 ----

// NewSmokeParticle 创建烟雾粒子
func NewSmokeParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleSmoke, 0)
}

// NewFlameParticle 创建火焰粒子
func NewFlameParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleFlame, 0)
}

// NewHeartParticle 创建爱心粒子
func NewHeartParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleHeart, 0)
}

// NewCriticalParticle 创建暴击粒子
func NewCriticalParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleCritical, 0)
}

// NewLavaParticle 创建岩浆粒子
func NewLavaParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleLava, 0)
}

// NewExplodeParticle 创建爆炸粒子
func NewExplodeParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleExplode, 0)
}

// NewHugeExplodeParticle 创建大型爆炸粒子
func NewHugeExplodeParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleHugeExplode, 0)
}

// NewBubbleParticle 创建气泡粒子
func NewBubbleParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleBubble, 0)
}

// NewPortalParticle 创建传送门粒子
func NewPortalParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticlePortal, 0)
}

// NewInkParticle 创建墨汁粒子（鱿鱼）
func NewInkParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleInk, 0)
}

// NewDustParticle 创建灰尘粒子（带颜色）
// r, g, b, a 范围 0-255
func NewDustParticle(x, y, z float32, r, g, b, a int) *protocol.LevelEventPacket {
	data := int32((a&0xFF)<<24 | (r&0xFF)<<16 | (g&0xFF)<<8 | (b & 0xFF))
	return NewParticlePacket(x, y, z, ParticleDust, data)
}

// NewSpellParticle 创建药水效果粒子（带颜色）
func NewSpellParticle(x, y, z float32, r, g, b, a int) *protocol.LevelEventPacket {
	data := int32((a&0xFF)<<24 | (r&0xFF)<<16 | (g&0xFF)<<8 | (b & 0xFF))
	return NewParticlePacket(x, y, z, ParticleMobSpell, data)
}

// NewRedstoneParticle 创建红石粉粒子
func NewRedstoneParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleRedstone, 0)
}

// NewEnchantParticle 创建附魔粒子
func NewEnchantParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleEnchantmentTable, 0)
}

// NewItemBreakParticle 创建物品破碎粒子
// itemID 和 itemMeta 决定碎片外观
func NewItemBreakParticle(x, y, z float32, itemID, itemMeta int) *protocol.LevelEventPacket {
	data := int32(itemID) + int32(itemMeta<<16)
	return NewParticlePacket(x, y, z, ParticleItemBreak, data)
}

// NewTerrainParticle 创建地形粒子
// blockID 和 blockMeta 决定碎片外观
func NewTerrainParticle(x, y, z float32, blockID, blockMeta int) *protocol.LevelEventPacket {
	data := int32(blockID) | int32(blockMeta<<12)
	return NewParticlePacket(x, y, z, ParticleTerrain, data)
}

// NewWaterDripParticle 创建水滴粒子
func NewWaterDripParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleDripWater, 0)
}

// NewLavaDripParticle 创建岩浆滴粒子
func NewLavaDripParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleDripLava, 0)
}

// NewAngryVillagerParticle 创建生气村民粒子
func NewAngryVillagerParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleVillagerAngry, 0)
}

// NewHappyVillagerParticle 创建开心村民粒子
func NewHappyVillagerParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleVillagerHappy, 0)
}

// NewMobSpawnParticle 创建生物生成粒子
// width 和 height 为生物碰撞箱尺寸
func NewMobSpawnParticle(x, y, z float32, width, height float32) *protocol.LevelEventPacket {
	data := int32(int(width*256)&0xFFFF) | (int32(int(height*256)&0xFFFF) << 16)
	return NewParticlePacket(x, y, z, ParticleSnowballPoof, data)
}
