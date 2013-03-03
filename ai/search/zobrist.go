package search

import (
	"chess/ai/gen"
	"math/rand"
	"sync"
)

const (
// z_wPawn   = 0
// z_wRook   = 1
// z_wKnight = 2
// z_wBishop = 3
// z_wQueen  = 4
// z_wKing   = 5

// z_bPawn   = 6
// z_bRook   = 7
// z_bKnight = 8
// z_bBishop = 9
// z_bQueen  = 10
// z_bKing   = 11
)

var (
	bitTable [][]int64

	AlphaNode = byte(1)
	BetaNode  = byte(2)
	PvsNode   = byte(3)
	ExactNode = byte(4)
)

type TransInfo struct {
	hash     int64
	depth    int
	score    int
	nodeType byte
	old      bool
}

type TransTable struct {
	table map[int64]TransInfo
	mutex *sync.RWMutex
}

func NewTransTable() *TransTable {
	table := new(TransTable)
	table.table = make(map[int64]TransInfo)
	table.mutex = new(sync.RWMutex)
	return table
}

func (t *TransTable) Cleanup() {
	for key, entry := range t.table {
		if entry.old {
			delete(t.table, key)
		}
	}
}

func InitialiseZobrist() {
	bitTable = make([][]int64, 64)
	for i := 0; i < 64; i++ {
		bitTable[i] = make([]int64, 16)
		for j := 0; j < 16; j++ {
			bitTable[i][j] = rand.Int63()
		}
	}
}

func Hash(board *gen.Board) int64 {
	hash := int64(0)
	for i, piece := range board.Squares {
		if piece != gen.Empty {
			hash ^= bitTable[i][piece]
		}
	}
	return hash
}
