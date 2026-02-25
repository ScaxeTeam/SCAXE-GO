package tile

// brewing_stand.go — 酿造台 TileEntity
// 对应 PHP BrewingStand.php
// Spawnable + Container(4格) + Nameable + OnUpdate(酿造tick)
//
// 槽位布局:
//   0 = 原料 (ingredient)
//   1-3 = 药水瓶

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	// MaxBrewTime 酿造最大时间 (tick)
	MaxBrewTime = 400

	// BrewingSlots 酿造台槽位数 (1原料 + 3药水)
	BrewingSlots = 4
)

// 合法的酿造原料 (item ID -> true)
var brewingIngredients = map[int16]bool{
	372: true, // NETHER_WART
	348: true, // GLOWSTONE_DUST
	331: true, // REDSTONE
	376: true, // FERMENTED_SPIDER_EYE
	378: true, // MAGMA_CREAM
	353: true, // SUGAR
	382: true, // GLISTERING_MELON
	375: true, // SPIDER_EYE
	370: true, // GHAST_TEAR
	377: true, // BLAZE_POWDER
	396: true, // GOLDEN_CARROT
	462: true, // PUFFER_FISH
	414: true, // RABBIT_FOOT
	289: true, // GUNPOWDER
}

// BrewingStand 酿造台 TileEntity
type BrewingStand struct {
	SpawnableBase
	ContainerBase
	NameableBase

	CookTime   int16
	NeedUpdate bool
}

// NewBrewingStand 创建酿造台实例
func NewBrewingStand(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
	bs := &BrewingStand{}
	InitSpawnableBase(&bs.SpawnableBase, TypeBrewingStand, chunk, nbtData)
	InitContainerBase(&bs.ContainerBase, BrewingSlots)
	bs.NameableBase.LoadNameFromNBT(nbtData)

	// 从 NBT 加载物品
	bs.ContainerBase.LoadItemsFromNBT(nbtData)

	// 加载酿造时间
	bs.CookTime = nbtData.GetShort("CookTime")
	if bs.CookTime == 0 {
		bs.CookTime = MaxBrewTime
	}

	bs.NeedUpdate = true

	return bs
}

func (bs *BrewingStand) GetName() string {
	if bs.HasCustomName() {
		return bs.GetCustomName()
	}
	return "Brewing Stand"
}

// OnUpdate 酿造 tick 逻辑
// 对应 PHP BrewingStand::onUpdate()
// 完整的酿造配方匹配需要 CraftingManager，此处仅实现框架
func (bs *BrewingStand) OnUpdate() bool {
	if bs.IsClosed() {
		return false
	}
	if !bs.NeedUpdate {
		return true
	}

	// 检查是否有原料和药水瓶
	ingredient := bs.GetItem(0)
	canBrew := false

	// 检查槽位 1-3 是否有药水
	for i := 1; i <= 3; i++ {
		it := bs.GetItem(i)
		if it.ID == 373 || it.ID == 438 { // POTION || SPLASH_POTION
			canBrew = true
			break
		}
	}

	// 验证原料
	if canBrew && ingredient.ID != 0 && ingredient.Count > 0 {
		if _, ok := brewingIngredients[int16(ingredient.ID)]; !ok {
			canBrew = false
		}
	} else {
		canBrew = false
	}

	if canBrew {
		bs.CookTime--
		if bs.CookTime <= 0 {
			bs.CookTime = MaxBrewTime
			// TODO: 配方匹配 + 产出药水 (需要 CraftingManager)
		}
	} else {
		bs.CookTime = MaxBrewTime
	}

	bs.NeedUpdate = canBrew
	return true
}

// GetSpawnCompound 客户端渲染数据
func (bs *BrewingStand) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeBrewingStand))
	compound.Set(nbt.NewIntTag("x", bs.X))
	compound.Set(nbt.NewIntTag("y", bs.Y))
	compound.Set(nbt.NewIntTag("z", bs.Z))
	compound.Set(nbt.NewShortTag("CookTime", MaxBrewTime))

	if bs.HasCustomName() {
		compound.Set(nbt.NewStringTag("CustomName", bs.GetCustomName()))
	}
	return compound
}

func (bs *BrewingStand) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	return DefaultUpdateCompoundTag(nbtData)
}

func (bs *BrewingStand) SpawnTo(sender PacketSender) bool {
	return SpawnTo(bs, sender)
}

func (bs *BrewingStand) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(bs, broadcaster)
}

func (bs *BrewingStand) SaveNBT() {
	bs.SpawnableBase.SaveNBT()
	bs.ContainerBase.SaveItemsToNBT(bs.NBT)
	bs.NameableBase.SaveNameToNBT(bs.NBT)
	bs.NBT.Set(nbt.NewShortTag("CookTime", bs.CookTime))
}

// CheckIngredient 检查物品是否为合法酿造原料
func CheckIngredient(itemID int16) bool {
	_, ok := brewingIngredients[itemID]
	return ok
}

func init() {
	RegisterTile(TypeBrewingStand, NewBrewingStand)
}
