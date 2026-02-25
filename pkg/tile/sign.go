package tile

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)
type Sign struct {
	SpawnableBase

	text [4]string
}
func NewSign(chunk *world.Chunk, nbtData *nbt.CompoundTag) *Sign {
	s := &Sign{}
	for i, key := range signTextKeys {
		if val := nbtData.GetString(key); val != "" {
			s.text[i] = val
		} else {
			nbtData.Set(nbt.NewStringTag(key, ""))
		}
	}

	InitSpawnableBase(&s.SpawnableBase, TypeSign, chunk, nbtData)
	return s
}
var signTextKeys = [4]string{"Text1", "Text2", "Text3", "Text4"}
func (s *Sign) GetText() [4]string {
	return s.text
}
func (s *Sign) GetLine(index int) string {
	if index < 0 || index > 3 {
		return ""
	}
	return s.text[index]
}
func (s *Sign) SetText(line1, line2, line3, line4 string) {
	s.text[0] = line1
	s.text[1] = line2
	s.text[2] = line3
	s.text[3] = line4
	for i, key := range signTextKeys {
		s.NBT.Set(nbt.NewStringTag(key, s.text[i]))
	}
}
func (s *Sign) SetLine(index int, line string) {
	if index < 0 || index > 3 {
		return
	}
	s.text[index] = line
	s.NBT.Set(nbt.NewStringTag(signTextKeys[index], line))
}
func (s *Sign) GetSpawnCompound() *nbt.CompoundTag {
	compound := nbt.NewCompoundTag("")
	compound.Set(nbt.NewStringTag("id", TypeSign))
	compound.Set(nbt.NewIntTag("x", s.X))
	compound.Set(nbt.NewIntTag("y", s.Y))
	compound.Set(nbt.NewIntTag("z", s.Z))
	for i, key := range signTextKeys {
		compound.Set(nbt.NewStringTag(key, s.text[i]))
	}
	return compound
}
func (s *Sign) UpdateCompoundTag(nbtData *nbt.CompoundTag) bool {
	if nbtData.GetString("id") != TypeSign {
		return false
	}

	var lines [4]string
	for i, key := range signTextKeys {
		lines[i] = nbtData.GetString(key)
	}

	s.SetText(lines[0], lines[1], lines[2], lines[3])
	return true
}
func (s *Sign) SpawnTo(sender PacketSender) bool {
	return SpawnTo(s, sender)
}
func (s *Sign) SpawnToAll(broadcaster ChunkBroadcaster) {
	SpawnToAll(s, broadcaster)
}
func (s *Sign) SaveNBT() {
	s.BaseTile.SaveNBT()
	for i, key := range signTextKeys {
		s.NBT.Set(nbt.NewStringTag(key, s.text[i]))
	}
	s.NBT.Remove("Creator")
}

func init() {
	RegisterTile(TypeSign, func(chunk *world.Chunk, nbtData *nbt.CompoundTag) Tile {
		return NewSign(chunk, nbtData)
	})
}
