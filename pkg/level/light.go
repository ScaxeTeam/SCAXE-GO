package level

import (
	"github.com/scaxe/scaxe-go/pkg/block"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type LightQueueNode struct {
	XPos, YPos, ZPos int32
	LightLevel       uint8
}

func (l *Level) GetBlockSkyLightAt(x, y, z int32) uint8 {
	if y < YMin || y >= YMax {
		if y >= YMax {
			return 15
		}
		return 0
	}
	chunk := l.GetChunk(x>>4, z>>4, false)
	if chunk == nil {
		return 15
	}
	return chunk.GetSkyLight(int(x&0x0f), int(y), int(z&0x0f))
}

func (l *Level) SetBlockSkyLightAt(x, y, z int32, level uint8) {
	if y < YMin || y >= YMax {
		return
	}
	chunk := l.GetChunk(x>>4, z>>4, true)
	if chunk != nil {
		chunk.SetSkyLight(int(x&0x0f), int(y), int(z&0x0f), level)
	}
}

func (l *Level) GetBlockLightAt(x, y, z int32) uint8 {
	if y < YMin || y >= YMax {
		return 0
	}
	chunk := l.GetChunk(x>>4, z>>4, false)
	if chunk == nil {
		return 0
	}
	return chunk.GetBlockLight(int(x&0x0f), int(y), int(z&0x0f))
}

func (l *Level) SetBlockLightAt(x, y, z int32, level uint8) {
	if y < YMin || y >= YMax {
		return
	}
	chunk := l.GetChunk(x>>4, z>>4, true)
	if chunk != nil {
		chunk.SetBlockLight(int(x&0x0f), int(y), int(z&0x0f), level)
	}
}

func (l *Level) GetHighestAdjacentBlockLight(x, y, z int32) uint8 {
	return max8(
		l.GetBlockLightAt(x+1, y, z),
		l.GetBlockLightAt(x-1, y, z),
		l.GetBlockLightAt(x, y+1, z),
		l.GetBlockLightAt(x, y-1, z),
		l.GetBlockLightAt(x, y, z+1),
		l.GetBlockLightAt(x, y, z-1),
	)
}

func (l *Level) GetHighestAdjacentBlockSkyLight(x, y, z int32) uint8 {
	return max8(
		l.GetBlockSkyLightAt(x+1, y, z),
		l.GetBlockSkyLightAt(x-1, y, z),
		l.GetBlockSkyLightAt(x, y+1, z),
		l.GetBlockSkyLightAt(x, y-1, z),
		l.GetBlockSkyLightAt(x, y, z+1),
		l.GetBlockSkyLightAt(x, y, z-1),
	)
}

func (l *Level) UpdateBlockLight(x, y, z int32, newLevel int) {
	lightPropagationQueue := make([]world.BlockPos, 0)
	lightRemovalQueue := make([]LightQueueNode, 0)
	visited := make(map[world.BlockPos]bool)
	removalVisited := make(map[world.BlockPos]bool)

	oldLevel := int(l.GetBlockLightAt(x, y, z))

	if newLevel == -1 {
		id := l.GetBlockId(x, y, z)
		emission := block.GetProperty(id).LightLevel
		filter := block.GetProperty(id).LightFilter
		decay := int(filter)
		if decay == 0 {
			decay = 1
		}

		adjacent := l.GetHighestAdjacentBlockLight(x, y, z)

		val := int(adjacent) - decay
		if val < 0 {
			val = 0
		}

		lightFromSource := int(emission)
		if lightFromSource > val {
			newLevel = lightFromSource
		} else {
			newLevel = val
		}
	}

	if oldLevel != newLevel {
		l.SetBlockLightAt(x, y, z, uint8(newLevel))

		pos := world.NewBlockPos(x, y, z)

		if newLevel < oldLevel {
			removalVisited[pos] = true
			lightRemovalQueue = append(lightRemovalQueue, LightQueueNode{x, y, z, uint8(oldLevel)})
		} else {
			visited[pos] = true
			lightPropagationQueue = append(lightPropagationQueue, pos)
		}
	}

	idxR := 0
	for idxR < len(lightRemovalQueue) {
		node := lightRemovalQueue[idxR]
		idxR++

		l.computeRemoveBlockLight(node.XPos-1, node.YPos, node.ZPos, node.LightLevel, &lightRemovalQueue, &lightPropagationQueue, visited, removalVisited)
		l.computeRemoveBlockLight(node.XPos+1, node.YPos, node.ZPos, node.LightLevel, &lightRemovalQueue, &lightPropagationQueue, visited, removalVisited)
		l.computeRemoveBlockLight(node.XPos, node.YPos-1, node.ZPos, node.LightLevel, &lightRemovalQueue, &lightPropagationQueue, visited, removalVisited)
		l.computeRemoveBlockLight(node.XPos, node.YPos+1, node.ZPos, node.LightLevel, &lightRemovalQueue, &lightPropagationQueue, visited, removalVisited)
		l.computeRemoveBlockLight(node.XPos, node.YPos, node.ZPos-1, node.LightLevel, &lightRemovalQueue, &lightPropagationQueue, visited, removalVisited)
		l.computeRemoveBlockLight(node.XPos, node.YPos, node.ZPos+1, node.LightLevel, &lightRemovalQueue, &lightPropagationQueue, visited, removalVisited)
	}

	idxP := 0
	for idxP < len(lightPropagationQueue) {
		pos := lightPropagationQueue[idxP]
		idxP++

		lightLevel := l.GetBlockLightAt(pos.X(), pos.Y(), pos.Z())

		if lightLevel >= 1 {
			l.computeSpreadBlockLight(pos.X()-1, pos.Y(), pos.Z(), lightLevel, &lightPropagationQueue, visited)
			l.computeSpreadBlockLight(pos.X()+1, pos.Y(), pos.Z(), lightLevel, &lightPropagationQueue, visited)
			l.computeSpreadBlockLight(pos.X(), pos.Y()-1, pos.Z(), lightLevel, &lightPropagationQueue, visited)
			l.computeSpreadBlockLight(pos.X(), pos.Y()+1, pos.Z(), lightLevel, &lightPropagationQueue, visited)
			l.computeSpreadBlockLight(pos.X(), pos.Y(), pos.Z()-1, lightLevel, &lightPropagationQueue, visited)
			l.computeSpreadBlockLight(pos.X(), pos.Y(), pos.Z()+1, lightLevel, &lightPropagationQueue, visited)
		}
	}
}

func (l *Level) computeRemoveBlockLight(x, y, z int32, currentLight uint8, removeQueue *[]LightQueueNode, spreadQueue *[]world.BlockPos, visited, removalVisited map[world.BlockPos]bool) {
	if y < YMin || y >= YMax {
		return
	}
	current := l.GetBlockLightAt(x, y, z)

	if current != 0 && current < currentLight {
		l.SetBlockLightAt(x, y, z, 0)

		pos := world.NewBlockPos(x, y, z)
		if !removalVisited[pos] {
			removalVisited[pos] = true
			if current > 1 {
				*removeQueue = append(*removeQueue, LightQueueNode{x, y, z, current})
			}
		}
	} else if current >= currentLight {
		pos := world.NewBlockPos(x, y, z)
		if !visited[pos] {
			visited[pos] = true
			*spreadQueue = append(*spreadQueue, pos)
		}
	}
}

func (l *Level) computeSpreadBlockLight(x, y, z int32, currentLight uint8, spreadQueue *[]world.BlockPos, visited map[world.BlockPos]bool) {
	if y < YMin || y >= YMax {
		return
	}

	current := l.GetBlockLightAt(x, y, z)
	id := l.GetBlockId(x, y, z)
	filter := block.GetProperty(id).LightFilter

	decay := int(filter)
	if decay == 0 {
		decay = 1
	}

	if int(currentLight) <= decay {
		return
	}
	potentialLight := currentLight - uint8(decay)

	if current < potentialLight {
		l.SetBlockLightAt(x, y, z, potentialLight)

		pos := world.NewBlockPos(x, y, z)
		if !visited[pos] {
			visited[pos] = true
			if potentialLight > 1 {
				*spreadQueue = append(*spreadQueue, pos)
			}
		}
	}
}

func (l *Level) UpdateBlockSkyLight(x, y, z int32) {
	chunk := l.GetChunk(x>>4, z>>4, false)
	if chunk == nil {
		return
	}

	lx, lz := x&0x0F, z&0x0F
	oldHeight := int32(chunk.HeightMap[(lz<<4)|lx])

	sourceId := l.GetBlockId(x, y, z)
	filter := block.GetProperty(sourceId).LightFilter

	_ = filter

	chunk.RecalculateColumn(int(lx), int(lz))
	newHeight := int32(chunk.HeightMap[(lz<<4)|lx])

	if newHeight > oldHeight {

		for i := y; i >= newHeight; i-- {

		}
	}

}

func max8(nums ...uint8) uint8 {
	var m uint8 = 0
	for _, n := range nums {
		if n > m {
			m = n
		}
	}
	return m
}
