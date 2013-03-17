package ai

import (
	"chess/ai/eval"
	"chess/ai/gen"
	"chess/ai/search"
)

var ()

//Player interface
//serves as a standardisation for human/ai players
//for simplifying the game loop
type Player interface {
	GetBestMove(*gen.Board, int) *gen.Move
	PonderUntilInput(*gen.Board, chan struct{}, *gen.Move)
	Debug()
}

//AI player struct conforms to Player interface
type AIPlayer struct {
	pvs  *search.PVSearch
	iter *search.IterativeSearch
}

//NewAIPlayer creates a new AI with initialised pvs and iterative 
// deepening searchers.
func NewAIPlayer() *AIPlayer {
	player := new(AIPlayer)
	player.pvs = search.NewPVSearch(eval.BasicEvaluator)
	player.iter = search.NewIterative(player.pvs, eval.BasicEvaluator)
	return player
}

//AIPlayer.GetBestMove returns the best move found from a principal variation
//search
func (ai *AIPlayer) GetBestMove(b *gen.Board, depth int) *gen.Move {
	bestmove, _ := ai.pvs.Search(b, depth)
	return bestmove
}

//AIPlayer.PonderUntilInput performs an iterative deepening search on the board
//until it discovers a solution to the game (checkmate) or the stop channel is closed
func (ai *AIPlayer) PonderUntilInput(b *gen.Board, stop chan struct{}, move *gen.Move) {
	ai.iter.Stop = stop
	move, _ = ai.iter.Search(b.Clone(), 0)
}

//AIPlayer.Debug prints the debug information from the 
//latest iterative and principal variation searches.
func (ai *AIPlayer) Debug() {
	ai.pvs.PrintDebug()
	ai.iter.PrintDebug()
}
