//search package provides functionality for searching through a chess game tree
//to calculate the values of each node of the tree and come up with the best move
package search

import (
	"chess/ai/gen"
	"fmt"
	"runtime"
	"strings"
)

//Searcher interface is to provide polymorphism to allow for 
//different search algorithms to be easily plugged in higher up
type Searcher interface {
	Search(*gen.Board, int) (*gen.Move, int)
	PrintDebug()
}

type SearchResult struct {
	move  gen.Move
	score int
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
}

func IntToStr(n int) string {
	s := fmt.Sprintf("%d", n)
	l := len(s)
	x := l / 3
	if l%3 > 0 {
		x++
	}
	slice := make([]string, x)
	i := len(s) - 3
	for j := 1; i > 0; i -= 3 {
		slice[x-j] = s[i : i+3]
		j++
	}
	j := i + 3
	slice[0] = s[:j]
	return strings.Join(slice, ",")
}
