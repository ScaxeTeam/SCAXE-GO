package level

import "github.com/scaxe/scaxe-go/pkg/protocol"
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
func NewParticlePacket(x, y, z float32, particleType int, data int32) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(protocol.EventAddParticleMask) | uint16(particleType&0xFFF)
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = data
	return pk
}
func NewDestroyBlockParticle(x, y, z float32, blockID int, blockMeta int) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(protocol.EventParticleDestroyBlock)
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = int32(blockID) + int32(blockMeta<<12)
	return pk
}
func NewShootParticle(x, y, z float32, data int32) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(protocol.EventParticleShoot)
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = data
	return pk
}
func NewSplashParticle(x, y, z float32, data int32) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(protocol.EventParticleSplash)
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = data
	return pk
}
func NewSpawnParticle(x, y, z float32, data int32) *protocol.LevelEventPacket {
	pk := protocol.NewLevelEventPacket()
	pk.EventID = uint16(protocol.EventParticleSpawn)
	pk.X = x
	pk.Y = y
	pk.Z = z
	pk.Data = data
	return pk
}
func NewSmokeParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleSmoke, 0)
}
func NewFlameParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleFlame, 0)
}
func NewHeartParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleHeart, 0)
}
func NewCriticalParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleCritical, 0)
}
func NewLavaParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleLava, 0)
}
func NewExplodeParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleExplode, 0)
}
func NewHugeExplodeParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleHugeExplode, 0)
}
func NewBubbleParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleBubble, 0)
}
func NewPortalParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticlePortal, 0)
}
func NewInkParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleInk, 0)
}
func NewDustParticle(x, y, z float32, r, g, b, a int) *protocol.LevelEventPacket {
	data := int32((a&0xFF)<<24 | (r&0xFF)<<16 | (g&0xFF)<<8 | (b & 0xFF))
	return NewParticlePacket(x, y, z, ParticleDust, data)
}
func NewSpellParticle(x, y, z float32, r, g, b, a int) *protocol.LevelEventPacket {
	data := int32((a&0xFF)<<24 | (r&0xFF)<<16 | (g&0xFF)<<8 | (b & 0xFF))
	return NewParticlePacket(x, y, z, ParticleMobSpell, data)
}
func NewRedstoneParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleRedstone, 0)
}
func NewEnchantParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleEnchantmentTable, 0)
}
func NewItemBreakParticle(x, y, z float32, itemID, itemMeta int) *protocol.LevelEventPacket {
	data := int32(itemID) + int32(itemMeta<<16)
	return NewParticlePacket(x, y, z, ParticleItemBreak, data)
}
func NewTerrainParticle(x, y, z float32, blockID, blockMeta int) *protocol.LevelEventPacket {
	data := int32(blockID) | int32(blockMeta<<12)
	return NewParticlePacket(x, y, z, ParticleTerrain, data)
}
func NewWaterDripParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleDripWater, 0)
}
func NewLavaDripParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleDripLava, 0)
}
func NewAngryVillagerParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleVillagerAngry, 0)
}
func NewHappyVillagerParticle(x, y, z float32) *protocol.LevelEventPacket {
	return NewParticlePacket(x, y, z, ParticleVillagerHappy, 0)
}
func NewMobSpawnParticle(x, y, z float32, width, height float32) *protocol.LevelEventPacket {
	data := int32(int(width*256)&0xFFFF) | (int32(int(height*256)&0xFFFF) << 16)
	return NewParticlePacket(x, y, z, ParticleSnowballPoof, data)
}
