package world

func (c *Chunk) RecalculateColumn(x, z int) {
	if x < 0 || x >= 16 || z < 0 || z >= 16 {
		return
	}
	for y := 127; y >= 0; y-- {
		id := c.GetBlockId(x, y, z)
		if id != 0 {
			c.HeightMap[(z<<4)|x] = byte(y + 1)
			return
		}
	}
	c.HeightMap[(z<<4)|x] = 0
}

func (c *Chunk) RecalculateHeightMap() {
	for x := 0; x < 16; x++ {
		for z := 0; z < 16; z++ {
			c.RecalculateColumn(x, z)
		}
	}
}
