package search

import (
	"chess/ai/eval"
	"chess/ai/gen"
	"fmt"
	"testing"
	"time"
)

func TestAlphaBeta(t *testing.T) {
	fmt.Println("NOT PARALLEL ALPHA BETA")

	for depth := 0; depth < 7; depth++ {
		t0 := time.Now()
		board := gen.NewBoard()
		ab := NewAlphaBeta()
		_ = ab.Search(board, depth)
		fmt.Printf("Searched %d nodes at depth %d in %v\n", ab.Nodes, depth, time.Since(t0))
	}
}

func TestCCAlphaBeta(t *testing.T) {
	fmt.Println("PARALLEL ALPHA BETA")
	stop := make(chan struct{})

	for depth := 0; depth < 7; depth++ {
		t0 := time.Now()
		ab := NewAlphaBeta()
		board := gen.NewBoard()
		_ = ab.ParallelSearch(board, depth, stop)
		fmt.Printf("Searched %d nodes at depth %d in %v\n", ab.Nodes, depth, time.Since(t0))
	}
}

func TestPVSearch(t *testing.T) {
	fmt.Println("NOT PARALLEL PRINCIPAL VARIATION SEARCH")
	pvs := NewPVSearch(&eval.StaticHeuristic{})

	for depth := 0; depth < 7; depth++ {
		t0 := time.Now()
		board := gen.NewBoard()
		_ = pvs.Search(board, depth)
		fmt.Printf("Searched %d nodes at depth %d in %v\n", pvs.Nodes, depth, time.Since(t0))
	}
}

func TestCCPVSearch(t *testing.T) {
	fmt.Println("PARALLEL PRINCIPAL VARIATION SEARCH")
	pvs := NewPVSearch(&eval.StaticHeuristic{})

	stop := make(chan struct{})
	for depth := 0; depth < 7; depth++ {
		t0 := time.Now()
		board := gen.NewBoard()
		_ = pvs.ParallelSearch(board, depth, stop)
		fmt.Printf("Searched %d nodes at depth %d in %v\n", pvs.Nodes, depth, time.Since(t0))
	}
}
