package object

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/level/generator/populator"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type Dungeon struct {
}

func NewDungeon() *Dungeon {
	return &Dungeon{}
}

func (d *Dungeon) Generate(w populator.ChunkManager, r *rand.Random, pos world.BlockPos) bool {

	j := r.NextBoundedInt(2) + 2
	k := -j - 1
	l := j + 1

	k1 := r.NextBoundedInt(2) + 2
	l1 := -k1 - 1
	i2 := k1 + 1
	j2 := 0

	for k2 := k; k2 <= l; k2++ {
		for l2 := -1; l2 <= 4; l2++ {
			for i3 := l1; i3 <= i2; i3++ {
				target := pos.Add(int32(k2), int32(l2), int32(i3))
				mat := w.GetBlockId(target.X(), target.Y(), target.Z())
				isSolid := mat != 0

				if l2 == -1 && !isSolid {
					return false
				}
				if l2 == 4 && !isSolid {
					return false
				}

				if (k2 == k || k2 == l || i3 == l1 || i3 == i2) && l2 == 0 && w.GetBlockId(target.X(), target.Y(), target.Z()) == 0 && w.GetBlockId(target.X(), target.Y()+1, target.Z()) == 0 {
					j2++
				}
			}
		}
	}

	if j2 >= 1 && j2 <= 5 {

		for k3 := k; k3 <= l; k3++ {
			for i4 := 3; i4 >= -1; i4-- {
				for k4 := l1; k4 <= i2; k4++ {
					target := pos.Add(int32(k3), int32(i4), int32(k4))

					if k3 != k && i4 != -1 && k4 != l1 && k3 != l && i4 != 4 && k4 != i2 {
						if w.GetBlockId(target.X(), target.Y(), target.Z()) != block.CHEST {
							w.SetBlock(target.X(), target.Y(), target.Z(), 0, 0, false)
						}
					} else if target.Y() >= 0 && w.GetBlockId(target.X(), target.Y()-1, target.Z()) == 0 {
						w.SetBlock(target.X(), target.Y(), target.Z(), 0, 0, false)
					} else if w.GetBlockId(target.X(), target.Y(), target.Z()) != 0 && w.GetBlockId(target.X(), target.Y(), target.Z()) != block.CHEST {
						if i4 == -1 && r.NextBoundedInt(4) != 0 {
							w.SetBlock(target.X(), target.Y(), target.Z(), block.MOSS_STONE, 0, false)
						} else {
							w.SetBlock(target.X(), target.Y(), target.Z(), block.COBBLESTONE, 0, false)
						}
					}
				}
			}
		}

		w.SetBlock(pos.X(), pos.Y(), pos.Z(), block.MONSTER_SPAWNER, 0, false)

		return true
	}
	return false
}
