package be

import (
	"fmt"
	"math"
)

type eval struct {
	move *Move
	val int
}

func GetBestMove(b *Board) Move {
	reschan := make(chan eval)
//	stop := make(chan struct{})
	fmt.Println("Generating first moves")
	moves := GenerateAllMoves(b)
	fmt.Println("Finished first Generating Moves")
	results := make([]eval, len(moves))
	rescount := 0

	fmt.Println("Starting Go Routines")
	for _, move := range moves {
		go concurrentAlphaBeta(*b, &move, reschan)
	}
	fmt.Println("Moves", len(moves))
	fmt.Println("Finished Starting Routines\nWaiting for Results...")
	for {
		select {
		case res := <-reschan:
			results[rescount] = res
			rescount++
		default:

		}
		if rescount == len(moves) {
			break
		}
	}
	
	bestMove := &Move { 0 }
	bestScore := math.MinInt32
	for _, evaled := range results {
		if evaled.val > bestScore {
			bestMove = evaled.move
			bestScore = evaled.val
		}
	}
	return *bestMove
}

func concurrentAlphaBeta(b Board, move *Move, res chan eval) {
	b.MakeMove(move)
	res <- eval { move, AlphaBeta(&b) }
	b.UnmakeMove(move)
}