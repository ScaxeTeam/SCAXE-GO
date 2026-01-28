package gorigional

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestVisualizeBiomeDiffs(t *testing.T) {

	javaPath := "../../../../../java_world.bin"
	goPath := "../../../../../go_world.bin"

	dumpChunk(javaPath, "Java", 0, 0)
	dumpChunk(goPath, "Go", 0, 0)
}

func dumpChunk(path string, label string, targetCx, targetCz int32) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", path, err)
		return
	}
	defer f.Close()

	chunkSize := 50440
	buf := make([]byte, chunkSize)

	fmt.Printf("Scanning %s for chunk (%d, %d)...\n", label, targetCx, targetCz)

	header := make([]byte, 12)
	if _, err := io.ReadFull(f, header); err != nil {
		fmt.Printf("Error reading header: %v\n", err)
		return
	}

	count := 0
	for {
		_, err := io.ReadFull(f, buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading file: %v\n", err)
			return
		}

		cx := int32(binary.BigEndian.Uint32(buf[0:4]))
		cz := int32(binary.BigEndian.Uint32(buf[4:8]))

		if cx == targetCx && cz == targetCz {
			fmt.Printf("Found Chunk (%d, %d)\n", cx, cz)

			biomeOffset := 49160
			biomeData := buf[biomeOffset : biomeOffset+256]

			fmt.Printf("--- %s Biome Grid (16x16) ---\n", label)
			for z := 0; z < 16; z++ {
				for x := 0; x < 16; x++ {
					id := biomeData[z*16+x]
					fmt.Printf("%3d ", id)
				}
				fmt.Println()
			}
			return
		}
		count++
	}
	fmt.Printf("Chunk (%d, %d) not found in %s after %d chunks\n", targetCx, targetCz, label, count)
}
