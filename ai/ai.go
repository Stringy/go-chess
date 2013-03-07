package ai

import (
	"chess/ai/eval"
	"chess/ai/gen"
	"chess/ai/search"
)

var (
	evaluator  = eval.BasicEvaluator
	ccPVSTable = search.NewPVSTable()
	pvs        *search.PVSearch
	iter       *search.IterativeSearch
)

type result struct {
	move  *gen.Move
	nodes int
	score int
}

func init() {
	pvs = search.NewPVSearch(evaluator)
	pvs.PVSLine = ccPVSTable
	iter = search.NewIterative(pvs, evaluator)
}

func GetBestMove(b gen.Board, depth int) *gen.Move {
	bestmove, _ := pvs.Search(&b, depth)
	return &bestmove
}

func PonderUntilInput(b *gen.Board, stop chan struct{}) {
	iter.Stop = stop
	_, _ = iter.Search(b.Clone(), 0)
}

func PrintDebug() {
	pvs.PrintDebug()
	iter.PrintDebug()
}
