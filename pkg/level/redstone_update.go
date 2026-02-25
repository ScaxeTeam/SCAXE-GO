package level

import (
	"github.com/scaxe/scaxe-go/pkg/block"
)

func (l *Level) GetRedstonePower(x, y, z int32, face int) int {
	bid := l.GetBlockId(x, y, z)
	if bid == 0 {
		return 0
	}
	behavior := block.Registry.GetBehavior(bid)
	if behavior == nil {
		return 0
	}
	if block.Registry.IsSolid(bid) {
		return l.GetStrongPowerTo(x, y, z)
	}
	meta := l.GetBlockData(x, y, z)
	return behavior.GetWeakPower(face, meta)
}

func (l *Level) GetStrongPowerTo(x, y, z int32) int {
	maxPower := 0
	offsets := [6][3]int32{
		{0, -1, 0}, {0, 1, 0},
		{0, 0, -1}, {0, 0, 1},
		{-1, 0, 0}, {1, 0, 0},
	}
	faces := [6]int{0, 1, 2, 3, 4, 5}

	for i, off := range offsets {
		nx, ny, nz := x+off[0], y+off[1], z+off[2]
		if ny < YMin || ny >= YMax {
			continue
		}
		nbid := l.GetBlockId(nx, ny, nz)
		if nbid == 0 {
			continue
		}
		behavior := block.Registry.GetBehavior(nbid)
		if behavior == nil {
			continue
		}
		nmeta := l.GetBlockData(nx, ny, nz)
		p := behavior.GetStrongPower(faces[i], nmeta)
		if p > maxPower {
			maxPower = p
		}
	}
	return maxPower
}

func (l *Level) IsBlockPowered(x, y, z int32) bool {
	offsets := [6][3]int32{
		{0, -1, 0}, {0, 1, 0},
		{0, 0, -1}, {0, 0, 1},
		{-1, 0, 0}, {1, 0, 0},
	}
	faces := [6]int{0, 1, 2, 3, 4, 5}

	for i, off := range offsets {
		nx, ny, nz := x+off[0], y+off[1], z+off[2]
		if ny < YMin || ny >= YMax {
			continue
		}
		if l.GetRedstonePower(nx, ny, nz, faces[i]) > 0 {
			return true
		}
	}
	return false
}

func (l *Level) IsBlockIndirectlyPowered(x, y, z int32) bool {
	offsets := [6][3]int32{
		{0, -1, 0}, {0, 1, 0},
		{0, 0, -1}, {0, 0, 1},
		{-1, 0, 0}, {1, 0, 0},
	}

	for _, off := range offsets {
		nx, ny, nz := x+off[0], y+off[1], z+off[2]
		if ny < YMin || ny >= YMax {
			continue
		}
		nbid := l.GetBlockId(nx, ny, nz)
		if nbid == 0 {
			continue
		}
		behavior := block.Registry.GetBehavior(nbid)
		if behavior == nil {
			continue
		}
		if behavior.IsPowerSource() {
			nmeta := l.GetBlockData(nx, ny, nz)
			for f := 0; f < 6; f++ {
				if behavior.GetWeakPower(f, nmeta) > 0 {
					return true
				}
			}
		}
		if block.Registry.IsSolid(nbid) && l.GetStrongPowerTo(nx, ny, nz) > 0 {
			return true
		}
	}
	return false
}

func (l *Level) UpdateRedstoneWire(x, y, z int32) {
	oldMeta := l.GetBlockData(x, y, z)
	newPower := l.calculateWirePower(x, y, z)

	if int(oldMeta) == newPower {
		return
	}

	l.SetBlock(x, y, z, block.REDSTONE_WIRE, byte(newPower), false)
	l.PendingBlockUpdates = append(l.PendingBlockUpdates, PendingBlockUpdate{
		X: x, Y: y, Z: z,
		ID: block.REDSTONE_WIRE, Meta: byte(newPower),
	})

	hOffsets := [4][2]int32{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, off := range hOffsets {
		nx, nz := x+off[0], z+off[1]
		l.notifyWireNeighbor(nx, y, nz)
		l.notifyWireNeighbor(nx, y-1, nz)
		l.notifyWireNeighbor(nx, y+1, nz)
	}
	l.notifyWireNeighbor(x, y-1, z)
	l.notifyWireNeighbor(x, y+1, z)
}

func (l *Level) notifyWireNeighbor(x, y, z int32) {
	if y < YMin || y >= YMax {
		return
	}
	bid := l.GetBlockId(x, y, z)
	if bid == block.REDSTONE_WIRE {
		l.UpdateRedstoneWire(x, y, z)
		return
	}
	behavior := block.Registry.GetBehavior(bid)
	if behavior != nil {
		bs := l.GetBlock(x, y, z)
		ctx := &block.BlockContext{
			X: int(x), Y: int(y), Z: int(z),
			Meta:    bs.Meta,
			Powered: l.getBlockPowered(bs.ID, bs.Meta, x, y, z),
		}
		behavior.OnUpdate(ctx, BlockUpdateNormal)
		l.applyBlockContextResult(ctx, x, y, z)
	}
}

func (l *Level) calculateWirePower(x, y, z int32) int {
	maxPower := 0

	offsets := [6][3]int32{
		{0, -1, 0}, {0, 1, 0},
		{0, 0, -1}, {0, 0, 1},
		{-1, 0, 0}, {1, 0, 0},
	}
	for _, off := range offsets {
		nx, ny, nz := x+off[0], y+off[1], z+off[2]
		if ny < YMin || ny >= YMax {
			continue
		}
		nbid := l.GetBlockId(nx, ny, nz)
		if nbid == 0 {
			continue
		}
		behavior := block.Registry.GetBehavior(nbid)
		if behavior == nil {
			continue
		}
		if nbid != block.REDSTONE_WIRE && behavior.IsPowerSource() {
			nmeta := l.GetBlockData(nx, ny, nz)
			for f := 0; f < 6; f++ {
				p := behavior.GetWeakPower(f, nmeta)
				if p > maxPower {
					maxPower = p
				}
			}
		}
		if block.Registry.IsSolid(nbid) {
			sp := l.GetStrongPowerTo(nx, ny, nz)
			if sp > maxPower {
				maxPower = sp
			}
		}
	}

	hOffsets := [4][2]int32{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, off := range hOffsets {
		nx, nz := x+off[0], z+off[1]
		wp := l.getWirePowerAt(nx, y, nz)
		if wp-1 > maxPower {
			maxPower = wp - 1
		}
		aboveSolid := block.Registry.IsSolid(l.GetBlockId(x, y+1, z))
		if !aboveSolid {
			wp = l.getWirePowerAt(nx, y+1, nz)
			if wp-1 > maxPower {
				maxPower = wp - 1
			}
		}
		belowID := l.GetBlockId(nx, y, nz)
		if block.Registry.IsSolid(belowID) {
			wp = l.getWirePowerAt(nx, y-1, nz)
			if wp-1 > maxPower {
				maxPower = wp - 1
			}
		}
	}

	if maxPower < 0 {
		maxPower = 0
	}
	if maxPower > 15 {
		maxPower = 15
	}
	return maxPower
}

func (l *Level) getWirePowerAt(x, y, z int32) int {
	if y < YMin || y >= YMax {
		return 0
	}
	if l.GetBlockId(x, y, z) == block.REDSTONE_WIRE {
		return int(l.GetBlockData(x, y, z))
	}
	return 0
}
