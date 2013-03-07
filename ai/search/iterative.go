package search

import (
	"chess/ai/eval"
	"chess/ai/gen"

	"fmt"
	"time"
)

const (
	search_depth = 15
)

type IterativeSearch struct {
	Nodes []int
	times []time.Duration
	Pvs   *PVSearch
	Stop  chan struct{}
}

func NewIterative(pvs *PVSearch, evaluator eval.Evaluator) *IterativeSearch {
	iter := new(IterativeSearch)
	iter.Stop = make(chan struct{})
	iter.Pvs = pvs
	iter.Nodes = make([]int, 0)
	iter.times = make([]time.Duration, 0)
	return iter
}

func (is *IterativeSearch) Search(b *gen.Board, depth int) (gen.Move, int) {
	move := gen.Move(0)
	score := 0
	is.Nodes = make([]int, 0, 0)
	is.times = make([]time.Duration, 0)
	sig := make(chan struct{})
	for depth = 0; depth < search_depth; depth++ {
		t0 := time.Now()
		go is.StartSearch(b.Clone(), depth, sig) //start new search
		for {                                    //wait for either user input or search result
			select {
			case _, ok := <-is.Stop:
				if !ok {
					return move, score
				}
			case <-sig:
				goto EXIT_LOOP
			default:
			}
		}
	EXIT_LOOP:
		is.Nodes = append(is.Nodes, is.Pvs.Nodes)
		is.times = append(is.times, time.Now().Sub(t0))
		//	fmt.Println("Nodes searched:", is.Pvs.Nodes)
		//	fmt.Println("\tCurrent best: \n\t\t", is.Pvs.PVSLine.Table[0][0].String())
		if score > (int(eval.CheckMate)-int(depth)) || score < -(int(eval.CheckMate)+int(depth)) {
			break
		}
	}
	return move, score
}

func (is *IterativeSearch) StartSearch(b *gen.Board, depth int, sig chan struct{}) {
	_, _ = is.Pvs.Search(b, depth)
	sig <- *new(struct{}) //send finished signal
}

func (iter *IterativeSearch) PrintDebug() {
	fmt.Println("\t=== Iterative Debug ===")
	for depth, nodes := range iter.Nodes {
		fmt.Print(depth, ":\t", nodes)
		fmt.Printf(", %v\n", iter.times[depth])
	}
}
