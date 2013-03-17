package search

import (
	"chess/ai/eval"
	"chess/ai/gen"
	"fmt"
	"math"
	"time"

//	"log"
)

//AlphaBeta structure used to conform to the Searcher interface and contains
//and eval.Evaluator for evaluating individual nodes at the bottom of the search
//it also counts the number of nodes it has searched
type AlphaBeta struct {
	Evaluator eval.Evaluator
	Nodes     int
	timeTaken time.Duration
}

//NewAlphaBeta returns an initialised AlphaBeta searcher
func NewAlphaBeta() *AlphaBeta {
	ab := new(AlphaBeta)
	ab.Evaluator = eval.BasicEvaluator
	return ab
}

//Search is the function to conform to the Searcher interface
//is the entry point to the alpha beta pruning algorithm
func (a *AlphaBeta) Search(b *gen.Board, depth int) (gen.Move, int) {
	a.Nodes = 0
	bestmove := gen.Move(0)
	bestscore := 0
	moves := gen.GenerateAllMoves(b)
	stop := make(chan struct{})

	t0 := time.Now()
	defer func() { a.timeTaken = time.Since(t0) }()

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
			return bestmove, score
		}
		if score > alpha {
			alpha = score
		}
	}
	return bestmove, 0
}

//alpha_beta is a standard alpha beta pruning algorithm. 
//it searches based on the current turn and based on narrowing the search window
//through the use of alpha beta values
func (s *AlphaBeta) alpha_beta(board *gen.Board, a, b, depth int, stop chan struct{}) int {
	s.Nodes++ //count all nodes
	if depth <= 0 {
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
			if !board.IsOtherPlayerChecked() {
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
			if !board.IsOtherPlayerChecked() {
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

func (a *AlphaBeta) PrintDebug() {
	fmt.Println("== Alpha Beta Debug ==")
	fmt.Println("Nodes searched:", a.Nodes)
	fmt.Printf("Time Taken: %v\n", a.timeTaken)
}
