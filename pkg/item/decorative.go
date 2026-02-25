package item

const (
	DyeInkSac          = 0
	DyeRoseRed         = 1
	DyeCactusGreen     = 2
	DyeCocoaBeans      = 3
	DyeLapisLazuli     = 4
	DyePurple          = 5
	DyeCyan            = 6
	DyeLightGray       = 7
	DyeGray            = 8
	DyePink            = 9
	DyeLime            = 10
	DyeDandelionYellow = 11
	DyeLightBlue       = 12
	DyeMagenta         = 13
	DyeOrange          = 14
	DyeBoneMeal        = 15
)
var DyeNames = [16]string{
	"Ink Sac", "Rose Red", "Cactus Green", "Cocoa Beans",
	"Lapis Lazuli", "Purple Dye", "Cyan Dye", "Light Gray Dye",
	"Gray Dye", "Pink Dye", "Lime Dye", "Dandelion Yellow",
	"Light Blue Dye", "Magenta Dye", "Orange Dye", "Bone Meal",
}
func GetDyeName(meta int) string {
	if meta < 0 || meta > 15 {
		return "Unknown Dye"
	}
	return DyeNames[meta]
}

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
func IsRecord(id int) bool {
	return id >= RECORD_13 && id <= RECORD_WAIT
}
func GetRecordName(id int) string {
	if id < RECORD_13 || id > RECORD_WAIT {
		return ""
	}
	return records[id-RECORD_13].Name
}

const (
	BANNER      = 446
	ARMOR_STAND = 425
)
func GetBannerColorName(meta int) string {
	if meta < 0 || meta > 15 {
		return "White Banner"
	}
	return DyeNames[15-meta] + " Banner"
}
