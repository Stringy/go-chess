//eval package provides functionality to evaluate the value of a particular 
//board state.
package eval

import (
	"chess/ai/gen"
)

const ()

var (
	//Basic Evaluator contains the standard evaluation modules
	//one for each piece plus material
	//users can extend this, or create their own
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

//Evaluator.AddModule can be used to dynamically add modules at
//runtime
func (e Evaluator) AddModule(module Module) Evaluator {
	return append(e, module)
}

//Evaluator.Eval accumulates the scores from each of the modules 
//and returns a positive/negative score based on which side is to move
func (e Evaluator) Eval(board *gen.Board) int {
	score := 0
	for _, module := range e {
		score += module.eval(board)
	}
	if board.NextMove == gen.BlackMove {
		return -score
	}
	return score
}

//Evaluator.Debug prints debug information from each module
func (e Evaluator) Debug() {
	for _, module := range e {
		module.debug()
	}
}
