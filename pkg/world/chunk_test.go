package world

import (
	"bytes"
	"testing"
)

func TestChunkSection_SetGetBlockId(t *testing.T) {
	s := NewChunkSection(0)

	s.SetBlockId(0, 0, 0, 1)
	if id := s.GetBlockId(0, 0, 0); id != 1 {
		t.Errorf("Expected block ID 1, got %d", id)
	}

	s.SetBlockId(15, 15, 15, 255)
	if id := s.GetBlockId(15, 15, 15); id != 255 {
		t.Errorf("Expected block ID 255, got %d", id)
	}
}

func TestChunkSection_SetGetBlockData(t *testing.T) {
	s := NewChunkSection(0)

	s.SetBlockData(0, 0, 0, 0x0A)
	s.SetBlockData(1, 0, 0, 0x0B)

	if data := s.GetBlockData(0, 0, 0); data != 0x0A {
		t.Errorf("Expected block data 0x0A at 0,0,0, got 0x%X", data)
	}
	if data := s.GetBlockData(1, 0, 0); data != 0x0B {
		t.Errorf("Expected block data 0x0B at 1,0,0, got 0x%X", data)
	}

	rawByte := s.Data[0]
	if rawByte != 0xBA {
		t.Errorf("Expected raw byte 0xBA, got 0x%X", rawByte)
	}
}

func TestChunk_SetGetBlock(t *testing.T) {
	c := NewChunk(0, 0)

	c.SetBlock(0, 0, 0, 1, 0)
	id, meta := c.GetBlock(0, 0, 0)
	if id != 1 || meta != 0 {
		t.Errorf("Expected 1:0, got %d:%d", id, meta)
	}

	c.SetBlock(15, 127, 15, 7, 7)
	id, meta = c.GetBlock(15, 127, 15)
	if id != 7 || meta != 7 {
		t.Errorf("Expected 7:7, got %d:%d", id, meta)
	}
}

func TestChunk_ToPacketBytes(t *testing.T) {
	c := NewChunk(1, 2)

	c.SetBlock(0, 0, 0, 1, 0)

	data := c.ToPacketBytes()

	expectedLen := 4 + 4 + 1 + (1 + 4096 + 2048*3) + 256 + 1024 + 1
	if len(data) != expectedLen {
		t.Errorf("Expected length %d, got %d", expectedLen, len(data))
	}

	buf := bytes.NewReader(data)

	header := make([]byte, 9)
	buf.Read(header)

	if header[3] != 1 {
		t.Errorf("Expected X=1 (BE), got byte %d", header[3])
	}
	if header[7] != 2 {
		t.Errorf("Expected Z=2 (BE), got byte %d", header[7])
	}
	if header[8] != 1 {
		t.Errorf("Expected SectionCount=1, got %d", header[8])
	}
}
