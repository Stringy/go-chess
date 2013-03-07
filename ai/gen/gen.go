//Package gen is responsible for generating a list of moves for a particular game state.
//It uses bitboards to contain the information regarding each piece type and a single 32 bit integer
//to contain movement information. It is only designed to generate from a single state, to a depth of one.
//Any further than this must be extended outside of this package. 
package gen

import ()

var ()

// init Initialises all global data including attack masks
func init() {
	InitAllData()        //init globals
	InitialiseAllMasks() //init attack masks
}

func Perft(b *Board, depth int) int {
	nodes := 0
	if depth == 0 {
		return 1
	}
	moves := GenerateAllMoves(b)
	for _, move := range moves {
		b.MakeMove(&move)
		nodes += Perft(b, depth-1)
		b.UnmakeMove(&move)
	}
	return nodes
}
