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
	fmt.Println("ALPHA BETA")
	tot := 0
	t1 := time.Now()
	for depth := 0; depth < max_depth; depth++ {
		t0 := time.Now()
		board := gen.NewBoard()
		ab := NewAlphaBeta()
		_, _ = ab.Search(board, 0)
		tot += ab.Nodes
		fmt.Printf("Searched %d nodes at depth %d in %v\n", ab.Nodes, depth, time.Since(t0))
	}
	fmt.Printf("Average N/s = %f\n", float64(tot)/time.Since(t1).Seconds())
}

func TestPVSearch(t *testing.T) {
	fmt.Println("PRINCIPAL VARIATION SEARCH")
	pvs := NewPVSearch(eval.BasicEvaluator)
	t1 := time.Now()
	tot := 0
	for depth := 0; depth < max_depth; depth++ {
		t0 := time.Now()
		board := gen.NewBoard()
		pvs = NewPVSearch(eval.BasicEvaluator)
		_, _ = pvs.Search(board, 8)
		tot += pvs.Nodes
		fmt.Printf("Searched %d nodes at depth %d in %v\n", pvs.Nodes, depth, time.Since(t0))
	}
	fmt.Printf("Average N/s = %f\n", float64(tot)/time.Since(t1).Seconds())
}
