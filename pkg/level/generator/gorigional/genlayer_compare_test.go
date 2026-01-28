package gorigional

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/scaxe/scaxe-go/pkg/level/generator/gorigional/layer"
)

func TestDumpGenLayerParity(t *testing.T) {
	seed := int64(114514)

	layers := layer.InitializeAll(seed)
	riverMix := layers[0]

	targetX := -3*4 - 2
	targetZ := -3*4 - 2
	width := 10
	height := 10

	t.Logf("Dumping GenLayer (riverMix) at %d, %d (%dx%d)", targetX, targetZ, width, height)

	data := riverMix.GetInts(targetX, targetZ, width, height)

	fileName := "../../../../../dump_go_genlayer_parity.csv"
	f, err := os.Create(fileName)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	defer f.Close()

	var sb strings.Builder
	for i, val := range data {
		sb.WriteString(strconv.Itoa(val))
		sb.WriteString(",")
		if (i+1)%width == 0 {
			sb.WriteString("\n")
		}
	}
	f.WriteString(sb.String())
	t.Logf("Dumped parity data to %s", fileName)

	fmt.Println("=== Go GenLayer Dump ===")
	for i, val := range data {
		fmt.Printf("%d,", val)
		if (i+1)%width == 0 {
			fmt.Println()
		}
	}
}
