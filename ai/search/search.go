//search package provides functionality for searching through a chess game tree
//to calculate the values of each node of the tree and come up with the best move
package search

import (
	"chess/ai/gen"
	"runtime"
)

//Searcher interface is to provide polymorphism to allow for 
//different search algorithms to be easily plugged in higher up
type Searcher interface {
	Search(*gen.Board, int) (gen.Move, int)
	ParallelSearch(*gen.Board, int) (gen.Move, int)
}

type SearchResult struct {
	move  gen.Move
	score int
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	InitialiseZobrist()
}
