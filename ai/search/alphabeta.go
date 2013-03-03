package search

import (
	"chess/ai/eval"
	"chess/ai/gen"
	//	"fmt"
	"math"

//	"log"
)

//AlphaBeta structure used to conform to the Searcher interface and contains
//and eval.Evaluator for evaluating individual nodes at the bottom of the search
//it also counts the number of nodes it has searched
type AlphaBetaAlgo struct {
	Evaluator eval.Evaluator
	Nodes     int
}

func NewAlphaBeta() *AlphaBetaAlgo {
	return &AlphaBetaAlgo{eval.BasicEvaluator, 0}
}

//Search is the function to conform to the Searcher interface
//is the entry point to the alpha beta pruning algorithm
func (a AlphaBetaAlgo) Search(b *gen.Board, depth int) gen.Move {
	a.Nodes = 0
	bestmove := gen.Move(0)
	bestscore := 0
	moves := gen.GenerateAllMoves(b)
	stop := make(chan struct{})

	alpha := -int(eval.CheckMate)
	beta := int(eval.CheckMate)
	for _, move := range moves {
		b.MakeMove(&move)
		score := a.alpha_beta(b, alpha, beta, depth, stop)
		b.UnmakeMove(&move)
		if score > bestscore {
			bestscore = score
			bestmove = move
		}
		if score >= beta {
			return bestmove
		}
		if score > alpha {
			alpha = score
		}
	}
	return bestmove
}

func (a AlphaBetaAlgo) ParallelSearch(b *gen.Board, depth int) gen.Move {
	res := make(chan SearchResult)
	stop := make(chan struct{})
	moves := gen.GenerateAllMoves(b)

	for _, move := range moves {
		b.MakeMove(&move)
		go a.ccSearch(b.Clone(), move, depth-1, res, stop)
		b.UnmakeMove(&move)
	}

	bestscore := int(eval.CheckMate)
	bestmove := gen.Move(0)
	alpha := -int(eval.CheckMate)
	beta := int(eval.CheckMate)
	rescount := 0
	//	fmt.Println(len(moves))

	for {
		select {
		case r := <-res:
			rescount++
			if -r.score < bestscore {
				bestscore = -r.score
				bestmove = r.move
			}
			alpha = int(math.Max(float64(alpha), float64(r.score)))
			if alpha >= beta {
				close(stop)
				return bestmove
			}
		default:
			if rescount == len(moves) {
				return bestmove
			}
		}
	}
	return bestmove
}

func (a AlphaBetaAlgo) ccSearch(b *gen.Board, move gen.Move, depth int, res chan SearchResult, stop chan struct{}) {
	res <- SearchResult{move, a.alpha_beta(b, -int(eval.CheckMate), int(eval.CheckMate), depth, stop)}
}

//alpha_beta is a standard alpha beta pruning algorithm. 
//it searches based on the current turn and based on narrowing the search window
//through the use of alpha beta values
func (s *AlphaBetaAlgo) alpha_beta(board *gen.Board, a, b, depth int, stop chan struct{}) int {
	s.Nodes++ //count all nodes
	if depth < 1 {
		return s.Evaluator.Eval(board)
	}

	select {
	case _, ok := <-stop:
		if !ok {
			return s.Evaluator.Eval(board)
		}
	default:
	}

	if board.NextMove == gen.BlackMove { //black move
		moves := gen.GenerateAllMoves(board)
		for _, move := range moves {
			board.MakeMove(&move)
			if !board.IsOtherKingAttacked() {
				a = int(math.Max(float64(a), float64(s.alpha_beta(board, a, b, depth-1, stop))))
				board.UnmakeMove(&move)
				if b <= a {
					break
				}
			} else {
				board.UnmakeMove(&move)
			}
		}
		return a
	} else { //white move
		moves := gen.GenerateAllMoves(board)
		for _, move := range moves {
			board.MakeMove(&move)
			if !board.IsOtherKingAttacked() {
				b = int(math.Min(float64(b), float64(s.alpha_beta(board, a, b, depth-1, stop))))
				board.UnmakeMove(&move)
				if b <= a {
					break
				}
			} else {
				board.UnmakeMove(&move)
			}
		}
		return b
	}
	return 0 //this is never reached, conforms to compiler return checking
}