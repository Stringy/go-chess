//eval package provides functionality to evaluate the value of a particular 
//board state.
package eval

import (
	"chess/ai/gen"

//	"fmt"
)

const ()

var (
	BasicEvaluator = Evaluator{
		new(MaterialModule),
		new(PawnModule),
		new(RookModule),
		new(KnightModule),
		new(BishopModule),
		new(QueenModule),
		new(KingModule),
	}
)

type Module interface {
	eval(*gen.Board) int
	debug()
}

type Evaluator []Module

//init initialises the global data for evaluation
func init() {
	initialiseGlobals()
}

func (e Evaluator) AddModule(module Module) Evaluator {
	return append(e, module)
}

func (e Evaluator) Eval(board *gen.Board) int {
	score := 0
	for _, module := range e {
		score += module.eval(board)
	}
	return score
}

func (e Evaluator) Debug() {
	for _, module := range e {
		module.debug()
	}
}
