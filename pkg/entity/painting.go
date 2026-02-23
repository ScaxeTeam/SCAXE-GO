package entity

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
)

const PaintingNetworkID = 83

type PaintingMotive struct {
	Name   string
	Width  int
	Height int
}

var AllPaintingMotives = []PaintingMotive{

	{Name: "Kebab", Width: 1, Height: 1},
	{Name: "Aztec", Width: 1, Height: 1},
	{Name: "Alban", Width: 1, Height: 1},
	{Name: "Aztec2", Width: 1, Height: 1},
	{Name: "Bomb", Width: 1, Height: 1},
	{Name: "Plant", Width: 1, Height: 1},
	{Name: "Wasteland", Width: 1, Height: 1},

	{Name: "Wanderer", Width: 1, Height: 2},
	{Name: "Graham", Width: 1, Height: 2},

	{Name: "Pool", Width: 2, Height: 1},
	{Name: "Courbet", Width: 2, Height: 1},
	{Name: "Sunset", Width: 2, Height: 1},
	{Name: "Sea", Width: 2, Height: 1},
	{Name: "Creebet", Width: 2, Height: 1},

	{Name: "Match", Width: 2, Height: 2},
	{Name: "Bust", Width: 2, Height: 2},
	{Name: "Stage", Width: 2, Height: 2},
	{Name: "Void", Width: 2, Height: 2},
	{Name: "SkullAndRoses", Width: 2, Height: 2},
	{Name: "Wither", Width: 2, Height: 2},

	{Name: "Fighters", Width: 4, Height: 2},

	{Name: "Skeleton", Width: 4, Height: 3},
	{Name: "DonkeyKong", Width: 4, Height: 3},

	{Name: "Pointer", Width: 4, Height: 4},
	{Name: "Pigscene", Width: 4, Height: 4},
	{Name: "BurningSkull", Width: 4, Height: 4},
}

var motivesByName map[string]*PaintingMotive

func init() {
	motivesByName = make(map[string]*PaintingMotive, len(AllPaintingMotives))
	for i := range AllPaintingMotives {
		motivesByName[AllPaintingMotives[i].Name] = &AllPaintingMotives[i]
	}
}

func GetPaintingMotive(name string) *PaintingMotive {
	return motivesByName[name]
}

func GetMotivesFittingSpace(width, height int) []PaintingMotive {
	var result []PaintingMotive
	for _, m := range AllPaintingMotives {
		if m.Width <= width && m.Height <= height {
			result = append(result, m)
		}
	}
	return result
}

type Painting struct {
	*Entity

	Motive string

	Direction int

	BlockX int
	BlockY int
	BlockZ int
}

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
	p.Entity.Gravity = 0
	p.Entity.Drag = 0

	return p
}

func (p *Painting) GetMotive() string {
	return p.Motive
}

func (p *Painting) GetMotiveInfo() *PaintingMotive {
	return GetPaintingMotive(p.Motive)
}

func (p *Painting) GetDirection() int {
	return p.Direction
}

func (p *Painting) GetName() string {
	return "Painting"
}

func (p *Painting) SavePaintingNBT() {
	p.Entity.SaveNBT()
	p.Entity.NamedTag.Set(nbt.NewStringTag("Motive", p.Motive))
	p.Entity.NamedTag.Set(nbt.NewByteTag("Direction", int8(p.Direction)))
}

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

const PaintingDropItemID = 321

func (p *Painting) GetDrops() (itemID, meta, count int) {
	return PaintingDropItemID, 0, 1
}
