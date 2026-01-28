package gorigional

import (
	"runtime"
	"sync"

	"github.com/scaxe/scaxe-go/pkg/level/generator/gorigional/noise"
	"github.com/scaxe/scaxe-go/pkg/math/rand"
	"github.com/scaxe/scaxe-go/pkg/world"
)

type ParallelChunkGenerator struct {
	baseSeed   int64
	workers    int
	generators []*ChunkGeneratorOverworld
	mu         sync.Mutex
}

func NewParallelChunkGenerator(seed int64, workers int) *ParallelChunkGenerator {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	pg := &ParallelChunkGenerator{
		baseSeed:   seed,
		workers:    workers,
		generators: make([]*ChunkGeneratorOverworld, workers),
	}

	for i := 0; i < workers; i++ {
		pg.generators[i] = NewChunkGeneratorOverworld(seed)
	}

	return pg
}

type ChunkCoord struct {
	X, Z int32
}

func (pg *ParallelChunkGenerator) GenerateChunksParallel(coords []ChunkCoord) []*world.Chunk {
	results := make([]*world.Chunk, len(coords))
	var wg sync.WaitGroup

	sem := make(chan int, pg.workers)

	for i, coord := range coords {
		wg.Add(1)
		sem <- i % pg.workers

		go func(idx int, cx, cz int32, workerID int) {
			defer wg.Done()
			defer func() { <-sem }()

			gen := pg.generators[workerID]

			chunk := pg.generateChunkDeterministic(gen, cx, cz)
			results[idx] = chunk
		}(i, coord.X, coord.Z, i%pg.workers)
	}

	wg.Wait()
	return results
}

func (pg *ParallelChunkGenerator) generateChunkDeterministic(gen *ChunkGeneratorOverworld, cx, cz int32) *world.Chunk {
	x, z := int(cx), int(cz)

	chunkSeed := int64(x)*341873128712 + int64(z)*132897987541
	localRand := rand.NewRandom(chunkSeed)

	c := world.NewChunk(cx, cz)

	gen.SetBlocksInChunk(x, z, c)
	gen.replaceBiomeBlocks(x, z, c, localRand)
	gen.caves.GenerateChunk(cx, cz, c)
	gen.ravines.GenerateChunk(cx, cz, c)

	_ = localRand

	return c
}

func (pg *ParallelChunkGenerator) GenerateAndPopulateParallel(coords []ChunkCoord, chunkManager interface {
	SetChunk(x, z int32, c *world.Chunk)
	GetChunk(x, z int32, generate bool) *world.Chunk
}) {

	chunks := pg.GenerateChunksParallel(coords)

	for i, coord := range coords {
		if chunks[i] != nil {
			chunkManager.SetChunk(coord.X, coord.Z, chunks[i])
		}
	}

	var wg sync.WaitGroup
	sem := make(chan int, pg.workers)

	for _, coord := range coords {
		wg.Add(1)
		sem <- 1

		go func(cx, cz int32) {
			defer wg.Done()
			defer func() { <-sem }()

			pg.mu.Lock()
			pg.generators[0].PopulateChunk(cx, cz)
			pg.mu.Unlock()
		}(coord.X, coord.Z)
	}

	wg.Wait()
}

type WorkerPoolNoise struct {
	octaves []*noise.OctavesNoise
	current int
	mu      sync.Mutex
}

func (wp *WorkerPoolNoise) GetNoise() *noise.OctavesNoise {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	n := wp.octaves[wp.current]
	wp.current = (wp.current + 1) % len(wp.octaves)
	return n
}
