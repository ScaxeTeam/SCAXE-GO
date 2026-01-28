package anvil

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/scaxe/scaxe-go/pkg/nbt"
	"github.com/scaxe/scaxe-go/pkg/world"
)

const (
	SectorSize = 4096
	Headers    = 2

	CompressionGzip = 1
	CompressionZlib = 2
)

type RegionLoader struct {
	x, z int32
	path string
	file *os.File
	mu   sync.Mutex

	locationTable [1024][2]uint32
	timestamps    [1024]uint32

	lastSector uint32
}

func NewRegionLoader(basePath string, x, z int32) (*RegionLoader, error) {
	path := filepath.Join(basePath, "region", fmt.Sprintf("r.%d.%d.mca", x, z))

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	rl := &RegionLoader{
		x:    x,
		z:    z,
		path: path,
		file: f,
	}

	stat, err := f.Stat()
	if err == nil && stat.Size() == 0 {
		rl.createBlank()
	} else {
		if err := rl.loadLocationTable(); err != nil {
			f.Close()
			return nil, err
		}
	}

	return rl, nil
}

func (rl *RegionLoader) Close() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	if rl.file != nil {

		return rl.file.Close()
	}
	return nil
}

func (rl *RegionLoader) createBlank() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	header := make([]byte, SectorSize*2)

	if _, err := rl.file.Write(header); err != nil {
		return err
	}
	rl.lastSector = 1
	return nil
}

func (rl *RegionLoader) loadLocationTable() error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.file.Seek(0, 0)
	header := make([]byte, SectorSize*2)
	if _, err := io.ReadFull(rl.file, header); err != nil {
		return err
	}

	locations := header[:SectorSize]
	for i := 0; i < 1024; i++ {
		off := i * 4
		val := binary.BigEndian.Uint32(locations[off : off+4])
		offset := val >> 8
		count := val & 0xFF
		rl.locationTable[i] = [2]uint32{offset, count}

		if offset+count-1 > rl.lastSector {
			rl.lastSector = offset + count - 1
		}
	}

	timestamps := header[SectorSize:]
	for i := 0; i < 1024; i++ {
		off := i * 4
		rl.timestamps[i] = binary.BigEndian.Uint32(timestamps[off : off+4])
	}

	if rl.lastSector < 1 {
		rl.lastSector = 1
	}

	return nil
}

func (rl *RegionLoader) ReadChunk(cx, cz int32) (*world.Chunk, error) {
	localX := cx & 31
	localZ := cz & 31
	index := localX + (localZ * 32)

	rl.mu.Lock()
	loc := rl.locationTable[index]
	rl.mu.Unlock()

	offset := loc[0]
	count := loc[1]

	if offset == 0 || count == 0 {
		return nil, nil
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	if _, err := rl.file.Seek(int64(offset)*SectorSize, 0); err != nil {
		return nil, err
	}

	var length uint32
	if err := binary.Read(rl.file, binary.BigEndian, &length); err != nil {
		return nil, err
	}

	if length == 0 {
		return nil, nil
	}

	var compression byte
	if err := binary.Read(rl.file, binary.BigEndian, &compression); err != nil {
		return nil, err
	}

	data := make([]byte, length-1)
	if _, err := io.ReadFull(rl.file, data); err != nil {
		return nil, err
	}

	var r io.Reader
	switch compression {
	case CompressionGzip:
		var err error
		r, err = gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
	case CompressionZlib:
		var err error
		r, err = zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown compression type: %d", compression)
	}

	decompressed, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if closer, ok := r.(io.Closer); ok {
		closer.Close()
	}

	n := nbt.New(nbt.BigEndian)
	if err := n.Read(decompressed); err != nil {
		return nil, err
	}

	root := n.GetData()

	var levelTag *nbt.CompoundTag
	if root.Has("Level") {
		levelTag = root.GetCompound("Level")
	} else if root.Name() == "Level" {
		levelTag = root
	} else {

		if root.Has("Sections") {
			levelTag = root
		}
	}

	if levelTag == nil {
		return nil, fmt.Errorf("invalid chunk NBT structure: missing 'Level' tag")
	}

	chunk := world.ChunkFromNBT(levelTag)
	if chunk != nil {
		chunk.X = cx
		chunk.Z = cz
	}

	return chunk, nil
}

func (rl *RegionLoader) WriteChunk(chunk *world.Chunk) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	nbtTag := chunk.ToNBT()
	n := nbt.New(nbt.BigEndian)
	n.SetData(nbtTag)

	compressedData, err := n.WriteZlib()
	if err != nil {
		return err
	}

	payloadLen := len(compressedData) + 1
	totalLen := 4 + payloadLen

	sectorsNeeded := (totalLen + SectorSize - 1) / SectorSize

	localX := chunk.X & 31
	localZ := chunk.Z & 31
	index := int(localX + (localZ * 32))

	offset := rl.locationTable[index][0]
	currentSectors := rl.locationTable[index][1]

	if offset == 0 || currentSectors < uint32(sectorsNeeded) {

		offset = rl.lastSector + 1
		rl.lastSector += uint32(sectorsNeeded)

	} else {

		currentSectors = uint32(sectorsNeeded)
	}

	if _, err := rl.file.Seek(int64(offset)*SectorSize, 0); err != nil {
		return err
	}

	if err := binary.Write(rl.file, binary.BigEndian, uint32(payloadLen)); err != nil {
		return err
	}

	if _, err := rl.file.Write([]byte{CompressionZlib}); err != nil {
		return err
	}

	if _, err := rl.file.Write(compressedData); err != nil {
		return err
	}

	written := 4 + 1 + len(compressedData)
	padding := int(uint32(sectorsNeeded)*SectorSize) - written
	if padding > 0 {
		pad := make([]byte, padding)
		if _, err := rl.file.Write(pad); err != nil {
			return err
		}
	}

	rl.locationTable[index][0] = offset
	rl.locationTable[index][1] = uint32(sectorsNeeded)
	rl.timestamps[index] = uint32(time.Now().Unix())

	if _, err := rl.file.Seek(int64(index)*4, 0); err != nil {
		return err
	}
	locVal := (offset << 8) | (uint32(sectorsNeeded) & 0xFF)
	if err := binary.Write(rl.file, binary.BigEndian, locVal); err != nil {
		return err
	}

	if _, err := rl.file.Seek(int64(SectorSize+(index*4)), 0); err != nil {
		return err
	}
	if err := binary.Write(rl.file, binary.BigEndian, rl.timestamps[index]); err != nil {
		return err
	}

	return nil
}
