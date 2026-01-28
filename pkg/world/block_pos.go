package world

type BlockPos struct {
	x, y, z int32
}

func NewBlockPos(x, y, z int32) BlockPos {
	return BlockPos{x: x, y: y, z: z}
}

func (p BlockPos) Add(x, y, z int32) BlockPos {
	return BlockPos{x: p.x + x, y: p.y + y, z: p.z + z}
}

func (p BlockPos) X() int32 {
	return p.x
}

func (p BlockPos) Y() int32 {
	return p.y
}

func (p BlockPos) Z() int32 {
	return p.z
}

func (p BlockPos) Up(n int32) BlockPos {
	return p.Add(0, n, 0)
}

func (p BlockPos) Down() BlockPos {
	return p.Add(0, -1, 0)
}

func (p BlockPos) North() BlockPos {
	return p.Add(0, 0, -1)
}

func (p BlockPos) South() BlockPos {
	return p.Add(0, 0, 1)
}

func (p BlockPos) East() BlockPos {
	return p.Add(1, 0, 0)
}

func (p BlockPos) West() BlockPos {
	return p.Add(-1, 0, 0)
}
