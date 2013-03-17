package be

import (
	"math"
//	"fmt"
)

const (
	MaxPlayer = BlackMove
	Depth = 5
)

func AlphaBeta(b *Board) int {
	return alpha_beta(b, math.MinInt32, math.MaxInt32, Depth)
}

func alpha_beta(board *Board, a, b, depth int) int {
//	fmt.Print(depth)
	if depth < 1 {
//		fmt.Println("depth == 0")
		return EvalPositionValue(board)
	}
//	fmt.Println("Not At Depth 0")
	if board.NextMove == MaxPlayer { //black move
//		fmt.Println("Generating All Black Moves", depth)
		moves := GenerateAllMoves(board)
//		fmt.Println("Generated Black Moves", depth)
		for _, move := range moves {
			board.MakeMove(&move)
			a = int(math.Max(float64(a), float64(alpha_beta(board, a, b, depth - 1))))
			board.UnmakeMove(&move)
			if b <= a {
				break
			}
		}
		return a
	} else { //white move
//		fmt.Println("Generating All White Moves", depth)
		moves := GenerateAllMoves(board)
//		fmt.Println("Generated All White Moves", depth)
		for _, move := range moves {
			board.MakeMove(&move)
			b = int(math.Min(float64(a), float64(alpha_beta(board, a, b, depth - 1))))
			board.UnmakeMove(&move)
			if b <= a {
				break
			}
		}
		return b
	}
	return 0
}