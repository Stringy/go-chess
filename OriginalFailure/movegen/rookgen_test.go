package movegen

import (
	"testing"
	"chess/bb"
	"fmt"
)

func TestWhiteRookMoves(t *testing.T) {
	fmt.Println("Testing White Rook Generation")
	state := GameState { ColourState { 0, 0, 0, 0, 0, 0 }, ColourState { 0, 0, 0, 0, 0, 0 }}
	
	state.White.Rooks = bb.GetSquare(6, 6)
	fmt.Println("Square to be tested:\n")
	bb.PrintBitboard(state.White.Rooks)

	out := WhiteRookMoves(state)

	for _, val := range out {
		move := val.White.Rooks
		if move & bb.Rank7 != bb.Empty ||
			move & bb.FileG != bb.Empty {
			// do nothing expected
		} else {
			fmt.Println("Unexpected Move")
			t.Fail()
		}
	} 
}