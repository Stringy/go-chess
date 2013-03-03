package search

import (
	"chess/ai/eval"
	"chess/ai/gen"
)

func quiescenceSearch(b *gen.Board, evaluator eval.Evaluator, alpha, beta int) int {
	bestval := evaluator.Eval(b)
	if bestval >= beta {
		return bestval
	}
	if bestval > alpha {
		alpha = bestval
	}

	moves := gen.GenerateAllMoves(b)
	for _, move := range moves {
		if !move.IsCapture() {
			continue
		} else {
			b.MakeMove(&move)
			score := -quiescenceSearch(b, evaluator, -beta, -alpha)
			b.UnmakeMove(&move)
			if score >= beta {
				return beta
			}
			if score > alpha {
				alpha = score
			}
		}
	}
	return alpha
}

func evalMove(b *gen.Board, move *gen.Move, evaluator eval.Evaluator) int {
	if b.IsCheck() {

	}
	if move.IsCapture() && move.IsPromotion() {
		return evaluator.Eval(b) - int(eval.PawnValue/2+eval.Value(move.GetCapt())+eval.Value(move.GetProm()))
	}
	if move.IsCapture() {
		return evaluator.Eval(b) + int(eval.PawnValue/2+eval.Value(move.GetCapt()))
	}
	if move.IsPromotion() {
		return evaluator.Eval(b) - int(eval.PawnValue/2+eval.Value(move.GetProm()))
	}
	return 0
}
