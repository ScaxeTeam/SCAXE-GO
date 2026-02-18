package player

import (
	"testing"
)

func TestChunkHash(t *testing.T) {
	tests := []struct {
		x, z int32
	}{
		{0, 0},
		{1, 1},
		{-1, -1},
		{100, -200},
		{-32768, 32767},
	}

	for _, tt := range tests {
		hash := ChunkHash(tt.x, tt.z)
		gotX := int32(hash >> 32)
		gotZ := int32(hash & 0xFFFFFFFF)
		if gotX != tt.x || gotZ != tt.z {
			t.Errorf("ChunkHash(%d, %d) = %d, decoded (%d, %d), want (%d, %d)",
				tt.x, tt.z, hash, gotX, gotZ, tt.x, tt.z)
		}
	}
}

func TestQueueChunks(t *testing.T) {
	p := &Player{
		LoadedChunks:   make(map[int64]bool),
		usedChunks:     make(map[int64]bool),
		chunkLoadQueue: nil,
	}

	hashes := []int64{
		ChunkHash(0, 0),
		ChunkHash(1, 0),
		ChunkHash(0, 1),
	}
	p.QueueChunks(hashes)

	if len(p.chunkLoadQueue) != 3 {
		t.Errorf("Expected 3 queued chunks, got %d", len(p.chunkLoadQueue))
	}

	more := []int64{ChunkHash(2, 2)}
	p.QueueChunks(more)

	if len(p.chunkLoadQueue) != 4 {
		t.Errorf("Expected 4 queued chunks after append, got %d", len(p.chunkLoadQueue))
	}
}

func TestGetLoadedChunkList(t *testing.T) {
	p := &Player{
		LoadedChunks: make(map[int64]bool),
		usedChunks:   make(map[int64]bool),
	}

	p.MarkChunkLoaded(0, 0)
	p.MarkChunkLoaded(1, 1)
	p.MarkChunkLoaded(-1, -1)

	list := p.GetLoadedChunkList()
	if len(list) != 3 {
		t.Errorf("Expected 3 loaded chunks, got %d", len(list))
	}
}

func TestUnloadChunk(t *testing.T) {
	p := &Player{
		LoadedChunks: make(map[int64]bool),
		usedChunks:   make(map[int64]bool),
	}

	p.MarkChunkLoaded(5, 5)
	if !p.IsChunkLoaded(5, 5) {
		t.Error("Chunk (5,5) should be loaded after MarkChunkLoaded")
	}

	p.UnloadChunk(5, 5)
	if p.IsChunkLoaded(5, 5) {
		t.Error("Chunk (5,5) should not be loaded after UnloadChunk")
	}
}

func TestChunkRadius(t *testing.T) {
	p := &Player{
		LoadedChunks: make(map[int64]bool),
		usedChunks:   make(map[int64]bool),
	}

	p.SetChunkRadius(1)
	if p.GetChunkRadius() != 2 {
		t.Errorf("Expected min clamped radius 2, got %d", p.GetChunkRadius())
	}

	p.SetChunkRadius(20)
	if p.GetChunkRadius() != 16 {
		t.Errorf("Expected max clamped radius 16, got %d", p.GetChunkRadius())
	}

	p.SetChunkRadius(8)
	if p.GetChunkRadius() != 8 {
		t.Errorf("Expected radius 8, got %d", p.GetChunkRadius())
	}
}
