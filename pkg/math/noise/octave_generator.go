package noise

type OctaveGenerator interface {
	SetXScale(scale float64)
	SetYScale(scale float64)
	SetZScale(scale float64)
	SetScale(scale float64)
	GetSizeX() int
	GetSizeY() int
	GetSizeZ() int
}

type BaseOctaveGenerator struct {
	Octaves int
	SizeX   int
	SizeY   int
	SizeZ   int
	XScale  float64
	YScale  float64
	ZScale  float64
}

func (g *BaseOctaveGenerator) SetXScale(scale float64) {
	g.XScale = scale
}

func (g *BaseOctaveGenerator) SetYScale(scale float64) {
	g.YScale = scale
}

func (g *BaseOctaveGenerator) SetZScale(scale float64) {
	g.ZScale = scale
}

func (g *BaseOctaveGenerator) SetScale(scale float64) {
	g.XScale = scale
	g.YScale = scale
	g.ZScale = scale
}

func (g *BaseOctaveGenerator) GetSizeX() int {
	return g.SizeX
}

func (g *BaseOctaveGenerator) GetSizeY() int {
	return g.SizeY
}

func (g *BaseOctaveGenerator) GetSizeZ() int {
	return g.SizeZ
}
