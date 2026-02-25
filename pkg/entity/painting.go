package entity

// painting.go — 画实体
// 对应 PHP: entity/Painting.php, entity/Hanging.php
//
// Painting 是挂在墙上的装饰实体:
//   - NetworkID: 83
//   - Motive (画作主题) 从 NBT 加载
//   - 不能移动
//   - 受到任何伤害立即销毁，掉落画物品
//   - 使用 AddPaintingPacket 而非 AddEntityPacket 生成

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
)

const PaintingNetworkID = 83

// ============ 画作主题目录 ============

// PaintingMotive 画作主题信息
type PaintingMotive struct {
	Name   string // 画作名称（NBT "Motive" 值）
	Width  int    // 宽度（方块数）
	Height int    // 高度（方块数）
}

// AllPaintingMotives 所有可用画作（MCPE 0.14 版本）
var AllPaintingMotives = []PaintingMotive{
	// 1×1
	{Name: "Kebab", Width: 1, Height: 1},
	{Name: "Aztec", Width: 1, Height: 1},
	{Name: "Alban", Width: 1, Height: 1},
	{Name: "Aztec2", Width: 1, Height: 1},
	{Name: "Bomb", Width: 1, Height: 1},
	{Name: "Plant", Width: 1, Height: 1},
	{Name: "Wasteland", Width: 1, Height: 1},

	// 1×2
	{Name: "Wanderer", Width: 1, Height: 2},
	{Name: "Graham", Width: 1, Height: 2},

	// 2×1
	{Name: "Pool", Width: 2, Height: 1},
	{Name: "Courbet", Width: 2, Height: 1},
	{Name: "Sunset", Width: 2, Height: 1},
	{Name: "Sea", Width: 2, Height: 1},
	{Name: "Creebet", Width: 2, Height: 1},

	// 2×2
	{Name: "Match", Width: 2, Height: 2},
	{Name: "Bust", Width: 2, Height: 2},
	{Name: "Stage", Width: 2, Height: 2},
	{Name: "Void", Width: 2, Height: 2},
	{Name: "SkullAndRoses", Width: 2, Height: 2},
	{Name: "Wither", Width: 2, Height: 2},

	// 4×2
	{Name: "Fighters", Width: 4, Height: 2},

	// 4×3
	{Name: "Skeleton", Width: 4, Height: 3},
	{Name: "DonkeyKong", Width: 4, Height: 3},

	// 4×4
	{Name: "Pointer", Width: 4, Height: 4},
	{Name: "Pigscene", Width: 4, Height: 4},
	{Name: "BurningSkull", Width: 4, Height: 4},
}

// motivesByName 按名称查找画作
var motivesByName map[string]*PaintingMotive

func init() {
	motivesByName = make(map[string]*PaintingMotive, len(AllPaintingMotives))
	for i := range AllPaintingMotives {
		motivesByName[AllPaintingMotives[i].Name] = &AllPaintingMotives[i]
	}
}

// GetPaintingMotive 按名称获取画作信息
func GetPaintingMotive(name string) *PaintingMotive {
	return motivesByName[name]
}

// GetMotivesFittingSpace 获取能放入指定空间的画作列表
func GetMotivesFittingSpace(width, height int) []PaintingMotive {
	var result []PaintingMotive
	for _, m := range AllPaintingMotives {
		if m.Width <= width && m.Height <= height {
			result = append(result, m)
		}
	}
	return result
}

// ============ Painting 实体 ============

// Painting 画实体
type Painting struct {
	*Entity

	// Motive 画作主题名称
	Motive string

	// Direction 朝向 (0=南, 1=西, 2=北, 3=东)
	Direction int

	// BlockX/Y/Z 挂载方块坐标
	BlockX int
	BlockY int
	BlockZ int
}

// NewPainting 创建画实体
func NewPainting(motive string, direction int) *Painting {
	p := &Painting{
		Entity:    NewEntity(),
		Motive:    motive,
		Direction: direction,
	}

	p.Entity.NetworkID = PaintingNetworkID
	p.Entity.MaxHealth = 1
	p.Entity.Health = 1
	p.Entity.Width = 0
	p.Entity.Height = 0
	p.Entity.Gravity = 0 // 画不受重力
	p.Entity.Drag = 0

	return p
}

// GetMotive 获取画作主题
func (p *Painting) GetMotive() string {
	return p.Motive
}

// GetMotiveInfo 获取画作详细信息
func (p *Painting) GetMotiveInfo() *PaintingMotive {
	return GetPaintingMotive(p.Motive)
}

// GetDirection 获取朝向
func (p *Painting) GetDirection() int {
	return p.Direction
}

// GetName 获取实体名称
func (p *Painting) GetName() string {
	return "Painting"
}

// ============ NBT ============

// SavePaintingNBT 保存画 NBT
func (p *Painting) SavePaintingNBT() {
	p.Entity.SaveNBT()
	p.Entity.NamedTag.Set(nbt.NewStringTag("Motive", p.Motive))
	p.Entity.NamedTag.Set(nbt.NewByteTag("Direction", int8(p.Direction)))
}

// LoadPaintingFromNBT 从 NBT 加载画数据
func (p *Painting) LoadPaintingFromNBT() {
	if p.Entity.NamedTag == nil {
		return
	}
	motive := p.Entity.NamedTag.GetString("Motive")
	if motive != "" {
		p.Motive = motive
	}
	p.Direction = int(p.Entity.NamedTag.GetByte("Direction"))
}

// ============ 掉落 ============

// PaintingDropItemID 画的掉落物品ID
const PaintingDropItemID = 321

// GetDrops 画被破坏时的掉落物
func (p *Painting) GetDrops() (itemID, meta, count int) {
	return PaintingDropItemID, 0, 1
}
