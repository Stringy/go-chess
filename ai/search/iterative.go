package search

import (
	"chess/ai/eval"
	"chess/ai/gen"

	"fmt"
	"time"
)

const (
	search_depth = 15 //max search depth for iterative deepening
)

//IterativeSearch struct, conforms to Searcher interface
//contains a PVSearch object, for searching along with a stop
//channel to interrupt the search upon user input/timeout
type IterativeSearch struct {
	Nodes []int
	times []time.Duration
	Pvs   *PVSearch
	Stop  chan struct{}
}

//NewIterative returns a new IterativeSearch with initialised data
func NewIterative(pvs *PVSearch, evaluator eval.Evaluator) *IterativeSearch {
	iter := new(IterativeSearch)
	iter.Stop = make(chan struct{})
	iter.Pvs = pvs
	iter.Nodes = make([]int, 0)
	iter.times = make([]time.Duration, 0)
	return iter
}

//IterativeSearch.Seach performs iterative deepening on game until it can find a solution
//if a solution isn't found, it returns nil
func (is *IterativeSearch) Search(b *gen.Board, depth int) (*gen.Move, int) {
	score := 0
	is.Nodes = make([]int, 0, 0)
	is.times = make([]time.Duration, 0)
	sig := make(chan int)
	for depth = 0; depth < search_depth; depth++ {
		t0 := time.Now()
		go is.StartSearch(b.Clone(), depth, sig) //start new search
		for {                                    //wait for either user input or search result
			select {
			case _, ok := <-is.Stop:
				if !ok {
					return nil, 0
				}
			case score = <-sig:
				goto EXIT_LOOP
			default:
			}
		}
	EXIT_LOOP:
		is.Nodes = append(is.Nodes, is.Pvs.Nodes) //record node debug info
		is.times = append(is.times, time.Now().Sub(t0)) //record time taken

		if score > (int(eval.CheckMate)-int(depth)) || score < -(int(eval.CheckMate)+int(depth)) {
			break
		}
	}
	return &is.Pvs.PVSLine.Table[0][0], score
}

func (is *IterativeSearch) StartSearch(b *gen.Board, depth int, sig chan int) {
	score := 0
	_, score = is.Pvs.Search(b, depth)
	sig <- score //send finished signal
}

func (iter *IterativeSearch) PrintDebug() {
	fmt.Println("\t=== Iterative Debug ===")
	for depth, nodes := range iter.Nodes {
		fmt.Print(depth, ":\t", nodes)
		fmt.Printf(", %v\n", iter.times[depth])
	}
}
