package search

import (
	"chess/ai/eval"
	"chess/ai/gen"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"
)

var (
	debugNodes     []int
	debugResults   []int
	evaluatedNodes int
	alphaCutoffs   int
	betaCutoffs    int
	nullMoves      int
	timeTaken      time.Duration
	NullReduction  = 4
	NullLim        = gen.KnightValue - 1
)

type PVSTable struct {
	Table [][]gen.Move
	mutex *sync.RWMutex
}

func (pvs *PVSTable) Update(ply int, move gen.Move) {
	pvs.mutex.Lock()
	pvs.Table[ply][ply] = move
	for j := ply + 1; j < len(pvs.Table[ply+1]); j++ {
		pvs.Table[ply][j] = pvs.Table[ply+1][j]
	}
	pvs.mutex.Unlock()
}

type PVSearch struct {
	Evaluator    eval.Evaluator
	Nodes        int
	Ply          int
	PVSLine      *PVSTable
	SearchPV     bool
	NullSearch   bool
	WhiteHistory [][]int
	BlackHistory [][]int
}

func NewPVSTable() *PVSTable {
	pvs := new(PVSTable)
	pvs.Table = make([][]gen.Move, 64)
	for i := 0; i < 64; i++ {
		pvs.Table[i] = make([]gen.Move, 64)
	}
	pvs.mutex = new(sync.RWMutex)
	return pvs
}

func NewPVSearch(evaluator eval.Evaluator) *PVSearch {
	searcher := new(PVSearch)
	searcher.SearchPV = true
	searcher.NullSearch = true
	searcher.Evaluator = evaluator
	searcher.PVSLine = NewPVSTable()
	searcher.Ply = 0
	searcher.WhiteHistory = make([][]int, 64)
	searcher.BlackHistory = make([][]int, 64)
	for i := 0; i < 64; i++ {
		searcher.WhiteHistory[i] = make([]int, 64)
		searcher.BlackHistory[i] = make([]int, 64)
	}
	return searcher
}

//pvs.Search is the entry point to the searching algorithm 
//satisfies search.Searcher interface
func (pvs *PVSearch) Search(b *gen.Board, depth int) (gen.Move, int) {
	stop := make(chan struct{})
	resetDebugSymbols()

	t0 := time.Now()
	pvs.Nodes = 0
	val := pvs.alpha_beta_pvs(b.Clone(), 0, depth, -int(eval.CheckMate), int(eval.CheckMate), stop)
	timeTaken = time.Since(t0)
	debugNodes = append(debugNodes, pvs.Nodes)
	debugResults = append(debugResults, val)
	return pvs.PVSLine.Table[0][0], val
}

func (pvs *PVSearch) ParallelSearch(b *gen.Board, depth int) (gen.Move, int) {
	res := make(chan SearchResult)
	stop := make(chan struct{})
	moves := gen.GenerateAllMoves(b)
	rescount := 0

	resetDebugSymbols()
	t0 := time.Now()
	pvs.sortMoves(b, moves, 0, depth, true)

	for _, move := range moves {
		b.MakeMove(&move)
		go func() {
			res <- SearchResult{move,
				pvs.alpha_beta_pvs(b.Clone(), 1, depth-1, -int(eval.CheckMate), int(eval.CheckMate), stop)}
		}()
		b.UnmakeMove(&move)
	}

	//	alpha := -math.MaxInt32
	beta := math.MaxInt32
	alpha := -math.MaxInt32
	bestscore := beta
	for {
		select {
		case r := <-res:
			rescount++
			debugNodes = append(debugNodes, pvs.Nodes)
			debugResults = append(debugResults, r.score)
			if r.score < bestscore {
				bestscore = r.score
			}
			if r.score > alpha {
				alpha = -r.score
				pvs.PVSLine.Update(0, r.move)
			}
			if alpha >= beta {
				close(stop)
				timeTaken = time.Since(t0)
				return pvs.PVSLine.Table[0][0], bestscore
			}
		default:

		}
		if rescount == len(moves) {
			timeTaken = time.Since(t0)
			return pvs.PVSLine.Table[0][0], bestscore
		}
	}
	return pvs.PVSLine.Table[0][0], 0 //never reached, compiler compliance
}

//pvs.alpha_beta_search is the actual search algorithm. It is a recursive algorithm which keeps a record of the 
//principal variation of the previous iterations
//moves are sorted before recursive searches to speed up the search so better nodes are found early. 
func (pvs *PVSearch) alpha_beta_pvs(b *gen.Board, ply, depth, alpha, beta int, stop chan struct{}) int {
	if depth <= 0 {
		pvs.SearchPV = false
		evaluatedNodes++
		return pvs.Evaluator.Eval(b)
	}

	//handle timeout or stop signal
	select {
	case _, ok := <-stop:
		if !ok {
			pvs.SearchPV = false
			return pvs.Evaluator.Eval(b)
		}
	default:
	}

	//null move pruning
	if (!pvs.SearchPV) && pvs.NullSearch {
		if (b.NextMove == gen.BlackMove && (b.BMaterial > NullLim)) ||
			(b.NextMove == gen.WhiteMove && (b.WMaterial > NullLim)) {
			if !b.IsCheck() {
				pvs.NullSearch = false //prevent two null moves in a row
				b.NullMove()           //make null move
				nullMoves++
				val := -pvs.alpha_beta_pvs(b, ply, depth-NullReduction, -beta, -beta+1, stop)
				b.NullMove() //unmake null move
				pvs.NullSearch = true
				if val >= beta {
					return val
				}
			}
		}
	}
	pvs.NullSearch = true

	score := 0
	moves := gen.GenerateAllMoves(b)
	movesfound := 0
	pvs.sortMoves(b, moves, ply, depth, pvs.SearchPV)
	//	fmt.Println(moves)
	for _, move := range moves {
		b.MakeMove(&move)
		if !b.IsOtherKingAttacked() {
			pvs.Nodes++
			if movesfound == 0 {
				score = -pvs.alpha_beta_pvs(b, ply+1, depth-1, -alpha-1, -alpha, stop)
				if score > alpha && score < beta {
					score = -pvs.alpha_beta_pvs(b, ply+1, depth-1, -beta, -alpha, stop)
				}
			} else {
				score = -pvs.alpha_beta_pvs(b, ply+1, depth-1, -beta, -alpha, stop)
			}
			b.UnmakeMove(&move)
			if score >= beta {
				//record history of beta cutoffs
				pvs.updateHistory(b, move.GetTo(), move.GetFrom(), depth)
				betaCutoffs++
				return beta
			}
			if score > alpha {
				movesfound++
				alpha = score
				alphaCutoffs++
				// record principal variation
				pvs.updateHistory(b, move.GetTo(), move.GetFrom(), depth)
				pvs.PVSLine.Update(ply, move)
			}
		} else {
			b.UnmakeMove(&move)
		}
	}
	return alpha
}

func (pvs *PVSearch) updateHistory(b *gen.Board, to, from byte, depth int) {
	if b.NextMove == gen.BlackMove {
		pvs.BlackHistory[from][to] += depth * depth
	} else {
		pvs.WhiteHistory[from][to] += depth * depth
	}
}

//pvs.sortMoves sorts the moves before they are analysed
//first, it'll check for a principal variation node, and moves that to the front and returns
//if pv is not being searched, it'll sort the moves based on previous discovery of beta cutoffs
func (pvs *PVSearch) sortMoves(b *gen.Board, moves []gen.Move, ply, depth int, pv bool) {
	if pv && depth > 1 {
		for i, move := range moves {
			if move == pvs.PVSLine.Table[ply][0] {
				moves = toFront(move, moves, i)
				return
			}
		}
	}

	if b.NextMove == gen.BlackMove {
		best := pvs.BlackHistory[moves[0].GetFrom()][moves[0].GetTo()]
		for i, move := range moves {
			if pvs.BlackHistory[move.GetFrom()][move.GetTo()] > best {
				best = pvs.BlackHistory[move.GetFrom()][move.GetTo()]
				moves = toFront(move, moves, i)
			}
		}
	} else {
		for i, move := range moves {
			best := pvs.WhiteHistory[moves[0].GetFrom()][moves[0].GetTo()]
			if pvs.WhiteHistory[move.GetFrom()][move.GetTo()] > best {
				best = pvs.WhiteHistory[move.GetFrom()][move.GetTo()]
				moves = toFront(move, moves, i)
			}
		}
	}
}

func resetDebugSymbols() {
	debugNodes = make([]int, 0)
	debugResults = make([]int, 0)
	evaluatedNodes = 0
	alphaCutoffs = 0
	betaCutoffs = 0
}

func (pvs *PVSearch) PrintDebug() {
	total := 0
	for _, num := range debugNodes {
		total += num
	}

	fmt.Println("Total nodes searched:\t", IntToStr(int(total)))
	if len(debugNodes) != 0 {
		fmt.Println("Each routine:")
		for i, nodes := range debugNodes {
			fmt.Println("                   ", i, "\t", IntToStr(nodes), ":", debugResults[i])
		}
	} else {
		fmt.Println("No routine specific data")
	}
	fmt.Printf("Time taken: %v\n", timeTaken)
	fmt.Printf("Nodes/Second: %f\n", float64(total)/timeTaken.Seconds())
	fmt.Printf("\tExact Nodes: %d\n\tAlpha Nodes: %d\n\tBeta Nodes: %d\n", evaluatedNodes, alphaCutoffs, betaCutoffs)
	fmt.Printf("\tNull Moves: %d\n", nullMoves)
	fmt.Println("Current principal variation:")
	for _, move := range pvs.PVSLine.Table[0] {
		if move.GetPiece() != 0 {
			fmt.Print(move.Print())
		}
	}
}

func IntToStr(n int) string {
	s := fmt.Sprintf("%d", n)
	l := len(s)
	x := l / 3
	if l%3 > 0 {
		x++
	}
	slice := make([]string, x)
	i := len(s) - 3
	for j := 1; i > 0; i -= 3 {
		slice[x-j] = s[i : i+3]
		j++
	}
	j := i + 3
	slice[0] = s[:j]
	return strings.Join(slice, ",")
}

func toFront(move gen.Move, moves []gen.Move, i int) []gen.Move {
	newslice := make([]gen.Move, len(moves))
	newslice[0] = move
	copy(newslice[1:], moves[:i])
	copy(newslice[i+1:], moves[i+1:])
	return newslice
}
