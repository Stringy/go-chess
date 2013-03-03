package search

import (
	"chess/ai/eval"
	"chess/ai/gen"
)

type Update struct {
	alpha, beta, depth int
}

type Worker struct {
	Nodes   int
	update  chan Update
	friends []chan Update
	stop    chan struct{}
}

func (w *Worker) Process(b *Board, depth int) {

}

func (w *Worker) parallel_alpha_beta(b *Board, alpha, beta, depth int) int {
	if depth <= 0 {
		return eval.BasicEvaluator.Eval(b)
	}

	select {
	case ud := <-update:
		if ud.depth > this.depth {
			score := eval.BasicEvaluator.Eval(b)
			if score >= ud.beta {

				return beta
			}
			if score > ud.alpha {
				alpha = score
				for _, friend := range friend {
					friend <- Update{alpha, beta, depth}
				}
			}
		}
	default:
	}

	moves := gen.GenerateAllMoves(b)
	for _, move := range moves {
		b.MakeMove(&move)
		if !b.IsOtherKingAttacked() {

			if score >= beta {
				return beta
			}
			if score > alpha {
				alpha = score
			}
			b.UnmakeMove(&move)
		} else {
			b.UnmakeMove(&move)
		}
	}
}
