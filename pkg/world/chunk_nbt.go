package world

import (
	"github.com/scaxe/scaxe-go/pkg/nbt"
)

func (c *Chunk) ToNBT() *nbt.CompoundTag {
	nbtData := nbt.NewCompoundTag("Level")

	nbtData.Set(nbt.NewIntTag("xPos", c.X))
	nbtData.Set(nbt.NewIntTag("zPos", c.Z))
	nbtData.Set(nbt.NewByteTag("V", 1))
	nbtData.Set(nbt.NewByteTag("LightPopulated", int8(boolToByte(c.LightPopulated))))
	nbtData.Set(nbt.NewByteTag("TerrainPopulated", int8(boolToByte(c.Populated))))

	sectionsList := nbt.NewListTag("Sections", nbt.TagCompound)
	for _, section := range c.Sections {
		if section != nil && !section.IsEmpty() {
			secTag := nbt.NewCompoundTag("")
			secTag.Set(nbt.NewByteTag("Y", int8(section.Y)))
			secTag.Set(nbt.NewByteArrayTag("Blocks", section.Blocks))
			secTag.Set(nbt.NewByteArrayTag("Data", section.Data))
			secTag.Set(nbt.NewByteArrayTag("BlockLight", section.BlockLight))
			secTag.Set(nbt.NewByteArrayTag("SkyLight", section.SkyLight))
			sectionsList.Add(secTag)
		}
	}
	nbtData.Set(sectionsList)

	biomes := make([]int32, 256)
	for i, v := range c.BiomeColors {
		biomes[i] = int32(v)
	}
	nbtData.Set(nbt.NewIntArrayTag("BiomeColors", biomes))

	heightMap := make([]int32, 256)
	for i, v := range c.HeightMap {
		heightMap[i] = int32(v)
	}
	nbtData.Set(nbt.NewIntArrayTag("HeightMap", heightMap))

	entitiesList := nbt.NewListTag("Entities", nbt.TagCompound)
	for _, e := range c.Entities {
		entitiesList.Add(e)
	}
	nbtData.Set(entitiesList)

	tilesList := nbt.NewListTag("TileEntities", nbt.TagCompound)
	for _, t := range c.Tiles {
		tilesList.Add(t)
	}
	nbtData.Set(tilesList)

	return nbtData
}

func ChunkFromNBT(levelTag *nbt.CompoundTag) *Chunk {
	if levelTag == nil {
		return nil
	}

	x := levelTag.GetInt("xPos")
	z := levelTag.GetInt("zPos")

	c := NewChunk(x, z)

	c.LightPopulated = levelTag.GetByte("LightPopulated") > 0
	c.Populated = levelTag.GetByte("TerrainPopulated") > 0
	c.Generated = true

	if sectionsList, ok := levelTag.Get("Sections").(*nbt.ListTag); ok {
		for i := 0; i < sectionsList.Len(); i++ {
			if secTag, ok := sectionsList.Get(i).(*nbt.CompoundTag); ok {
				y := secTag.GetByte("Y")
				if y >= 0 && y < SectionCount {
					sec := c.getSection(int(y))

					blocks := secTag.GetByteArray("Blocks")
					if len(blocks) == 4096 {
						copy(sec.Blocks, blocks)
					}

					data := secTag.GetByteArray("Data")
					if len(data) == 2048 {
						copy(sec.Data, data)
					}

					blockLight := secTag.GetByteArray("BlockLight")
					if len(blockLight) == 2048 {
						copy(sec.BlockLight, blockLight)
					}

					skyLight := secTag.GetByteArray("SkyLight")
					if len(skyLight) == 2048 {
						copy(sec.SkyLight, skyLight)
					}
				}
			}
		}
	}

	if biomes := levelTag.GetIntArray("BiomeColors"); len(biomes) == 256 {
		for i, v := range biomes {
			c.BiomeColors[i] = uint32(v)
		}
	}

	if heightMap := levelTag.GetIntArray("HeightMap"); len(heightMap) == 256 {
		for i, v := range heightMap {
			c.HeightMap[i] = byte(v)
		}
	}

	if entitiesList, ok := levelTag.Get("Entities").(*nbt.ListTag); ok {
		for i := 0; i < entitiesList.Len(); i++ {
			if tag, ok := entitiesList.Get(i).(*nbt.CompoundTag); ok {
				c.Entities = append(c.Entities, tag)
			}
		}
	}

	if tilesList, ok := levelTag.Get("TileEntities").(*nbt.ListTag); ok {
		for i := 0; i < tilesList.Len(); i++ {
			if tag, ok := tilesList.Get(i).(*nbt.CompoundTag); ok {
				c.Tiles = append(c.Tiles, tag)
			}
		}
	}

	return c
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}
