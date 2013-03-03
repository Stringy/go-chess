package search

import (
	"chess/ai/eval"
	"chess/ai/gen"
	"math"

//	"fmt"
)

type Negascout struct {
	Evaluator eval.Evaluator
}

func (n Negascout) Search(b *gen.Board, depth int, stop chan struct{}) int {
	score := n.negascout(b, math.MinInt32, math.MaxInt32, depth)
	//	fmt.Println(count)
	return score
}

func (n Negascout) negascout(board *gen.Board, alpha, beta, depth int) int {
	if depth < 1 {
		return n.Evaluator.Eval(board)
	}
	b := beta
	a := alpha
	moves := gen.GenerateAllMoves(board)
	for i, move := range moves {
		board.MakeMove(&move)
		score := -n.negascout(board, -b, -a, depth-1)
		if a < score && score < beta && i != 0 {
			a = -n.negascout(board, -beta, -score, depth-1)
		}
		alpha = int(math.Max(float64(alpha), float64(score)))
		if a >= beta {
			return a
		}
		b = a + 1
		board.UnmakeMove(&move)
	}
	return alpha
}
