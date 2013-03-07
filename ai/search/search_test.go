package search

import (
	"chess/ai/eval"
	"chess/ai/gen"
	"fmt"
	"testing"
	"time"
)

const max_depth = 10

func TestAlphaBeta(t *testing.T) {
	fmt.Println("NOT PARALLEL ALPHA BETA")

	for depth := 0; depth < max_depth; depth++ {
		t0 := time.Now()
		board := gen.NewBoard()
		ab := NewAlphaBeta()
		_, _ = ab.Search(board, depth)
		fmt.Printf("Searched %d nodes at depth %d in %v\n", ab.Nodes, depth, time.Since(t0))
	}
}

func TestPVSearch(t *testing.T) {
	fmt.Println("NOT PARALLEL PRINCIPAL VARIATION SEARCH")
	pvs := NewPVSearch(eval.BasicEvaluator)

	for depth := 0; depth < max_depth; depth++ {
		t0 := time.Now()
		board := gen.NewBoard()
		pvs = NewPVSearch(eval.BasicEvaluator)
		_, _ = pvs.Search(board, depth)
		fmt.Printf("Searched %d nodes at depth %d in %v\n", pvs.Nodes, depth, time.Since(t0))
	}
}
