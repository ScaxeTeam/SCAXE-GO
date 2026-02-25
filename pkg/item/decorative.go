package item

// decorative.go — 装饰性物品
// 染料子类型, 旗帜, 唱片, 盔甲架等

// ============ 染料子类型 ============
// DYE (351) 通过 meta 区分子类型

const (
	DyeInkSac          = 0  // 墨囊
	DyeRoseRed         = 1  // 玫瑰红
	DyeCactusGreen     = 2  // 仙人掌绿
	DyeCocoaBeans      = 3  // 可可豆
	DyeLapisLazuli     = 4  // 青金石
	DyePurple          = 5  // 紫色染料
	DyeCyan            = 6  // 青色染料
	DyeLightGray       = 7  // 淡灰色染料
	DyeGray            = 8  // 灰色染料
	DyePink            = 9  // 粉红色染料
	DyeLime            = 10 // 黄绿色染料
	DyeDandelionYellow = 11 // 蒲公英黄
	DyeLightBlue       = 12 // 淡蓝色染料
	DyeMagenta         = 13 // 品红色染料
	DyeOrange          = 14 // 橙色染料
	DyeBoneMeal        = 15 // 骨粉
)

// DyeNames 染料 meta → 名称映射
var DyeNames = [16]string{
	"Ink Sac", "Rose Red", "Cactus Green", "Cocoa Beans",
	"Lapis Lazuli", "Purple Dye", "Cyan Dye", "Light Gray Dye",
	"Gray Dye", "Pink Dye", "Lime Dye", "Dandelion Yellow",
	"Light Blue Dye", "Magenta Dye", "Orange Dye", "Bone Meal",
}

// GetDyeName 获取染料名称
func GetDyeName(meta int) string {
	if meta < 0 || meta > 15 {
		return "Unknown Dye"
	}
	return DyeNames[meta]
}

// ============ 头颅子类型 ============
// SKULL/MOB_HEAD (397) 通过 meta 区分

const (
	SkullSkeleton       = 0
	SkullWitherSkeleton = 1
	SkullZombie         = 2
	SkullSteve          = 3
	SkullCreeper        = 4
)

var SkullNames = [5]string{
	"Skeleton Skull", "Wither Skeleton Skull",
	"Zombie Head", "Steve Head", "Creeper Head",
}

func GetSkullName(meta int) string {
	if meta < 0 || meta > 4 {
		return "Mob Head"
	}
	return SkullNames[meta]
}

// ============ 唱片 ============
// MCPE 唱片物品 ID 范围: 500-511

const (
	RECORD_13      = 500
	RECORD_CAT     = 501
	RECORD_BLOCKS  = 502
	RECORD_CHIRP   = 503
	RECORD_FAR     = 504
	RECORD_MALL    = 505
	RECORD_MELLOHI = 506
	RECORD_STAL    = 507
	RECORD_STRAD   = 508
	RECORD_WARD    = 509
	RECORD_11      = 510
	RECORD_WAIT    = 511
)

// RecordInfo 唱片信息
type RecordInfo struct {
	ID   int
	Name string
}

var records = []RecordInfo{
	{RECORD_13, "Music Disc - 13"},
	{RECORD_CAT, "Music Disc - cat"},
	{RECORD_BLOCKS, "Music Disc - blocks"},
	{RECORD_CHIRP, "Music Disc - chirp"},
	{RECORD_FAR, "Music Disc - far"},
	{RECORD_MALL, "Music Disc - mall"},
	{RECORD_MELLOHI, "Music Disc - mellohi"},
	{RECORD_STAL, "Music Disc - stal"},
	{RECORD_STRAD, "Music Disc - strad"},
	{RECORD_WARD, "Music Disc - ward"},
	{RECORD_11, "Music Disc - 11"},
	{RECORD_WAIT, "Music Disc - wait"},
}

// IsRecord 判断是否为唱片
func IsRecord(id int) bool {
	return id >= RECORD_13 && id <= RECORD_WAIT
}

// GetRecordName 获取唱片名称
func GetRecordName(id int) string {
	if id < RECORD_13 || id > RECORD_WAIT {
		return ""
	}
	return records[id-RECORD_13].Name
}

// ============ 旗帜/盔甲架 ============

const (
	BANNER      = 446
	ARMOR_STAND = 425
)

// BannerBaseColors 旗帜底色 (meta 0-15, 与染料对应)
func GetBannerColorName(meta int) string {
	if meta < 0 || meta > 15 {
		return "White Banner"
	}
	return DyeNames[15-meta] + " Banner" // 旗帜颜色与染料反序
}
