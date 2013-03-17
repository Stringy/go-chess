package search

import (
	"chess/ai/eval"
	"chess/ai/gen"
	"fmt"
	"sync"
	"time"
)

var (
	debugNodes     []int //debug node count for latest search
	debugResults   []int
	evaluatedNodes int                   //number of nodes evaluated
	alphaCutoffs   int                   //number of alpha cut offs
	betaCutoffs    int                   //number of beta nodes
	nullMoves      int                   //number of null moves 
	timeTaken      time.Duration         //time of a single search
	NullReduction  = 4                   //depth reduction from null moves
	NullLim        = gen.KnightValue - 1 //material limit under which no null moves will be made
)

//PVSTable is a thread safe structure for remembering
//the principal variation at various depths
type PVSTable struct {
	Table [][]gen.Move
	mutex *sync.RWMutex
}

//PVSTable.Update updates the best moves found depending on 
//the ply at which a new best move has been moves
func (pvs *PVSTable) Update(ply int, move gen.Move) {
	pvs.mutex.Lock()
	pvs.Table[ply][ply] = move
	for j := ply + 1; j < len(pvs.Table[ply+1]); j++ {
		pvs.Table[ply][j] = pvs.Table[ply+1][j]
	}
	pvs.mutex.Unlock()
}

//PVSearch is a structure which conforms to the Searcher interface
type PVSearch struct {
	Evaluator    eval.Evaluator //for evaluating nodes
	Nodes        int            //number of nodes searched
	Ply          int            //current ply of the search in progress
	PVSLine      *PVSTable      //table containing best moves indexed by ply
	SearchPV     bool           //flag for searching the principal variation
	NullSearch   bool           //flag for null moves
	WhiteHistory [][]int        //White history heuristic
	BlackHistory [][]int        //Black history heuristic
}

//NewPVSTable initialises a new table for the principal variation
func NewPVSTable() *PVSTable {
	pvs := new(PVSTable)
	pvs.Table = make([][]gen.Move, 64)
	for i := 0; i < 64; i++ {
		pvs.Table[i] = make([]gen.Move, 64)
	}
	pvs.mutex = new(sync.RWMutex)
	return pvs
}

//NewPVSearch initialises and new principal variation search
// structure
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
func (pvs *PVSearch) Search(b *gen.Board, depth int) (*gen.Move, int) {
	stop := make(chan struct{})
	resetDebugSymbols()

	t0 := time.Now()
	pvs.Nodes = 0
	val := pvs.alpha_beta_pvs(b.Clone(), 0, depth, -int(eval.CheckMate), int(eval.CheckMate), stop)
	timeTaken = time.Since(t0)
	debugNodes = append(debugNodes, pvs.Nodes)
	debugResults = append(debugResults, val)
	return &pvs.PVSLine.Table[0][0], val
}

//pvs.alpha_beta_search is the actual search algorithm. It is a recursive algorithm which keeps a record of the 
//principal variation of the previous iterations
//moves are sorted before recursive searches to speed up the search so better nodes are found early. 
//null moves are made to attempt to discover beta cut offs early.
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
				pvs.Nodes++
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
	for _, move := range moves {
		b.MakeMove(&move)
		if !b.IsOtherPlayerChecked() {
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

//resetDebugSymbols resets all debug information to zero
//used before any new search
func resetDebugSymbols() {
	debugNodes = make([]int, 0)
	debugResults = make([]int, 0)
	evaluatedNodes = 0
	alphaCutoffs = 0
	betaCutoffs = 0
}

//PVSearch.PrintDebug prints debug information about the latest search
//including total nodes, alpha/beta/absolute nodes and null moves
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

//toFront takes a slice of moves, a move and its index in the slice and brings the move
//to the front, as a way of ordering the slice
func toFront(move gen.Move, moves []gen.Move, i int) []gen.Move {
	newslice := make([]gen.Move, len(moves))
	newslice[0] = move
	copy(newslice[1:], moves[:i])
	copy(newslice[i+1:], moves[i+1:])
	return newslice
}
