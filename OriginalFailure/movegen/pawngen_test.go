package movegen

import (
	"testing"
	"chess/bb"
	"fmt"
	"unsafe"
)

func TestWhitePawnGen(t *testing.T) {
	bb.InitAttackMasks()
	fmt.Println("Testing White Move Generation")
	state := GameState { ColourState { 0, 0, 0, 0, 0, 0 }, ColourState { 0, 0, 0, 0, 0, 0 }} 

	state.White.Pawns = bb.Rank2

	fmt.Println("Pawn from which generation will take place:\n")
	bb.PrintBitboard(state.White.Pawns)
	
	fmt.Println("Size of GameState: ", unsafe.Sizeof(state))
	
	out := WhitePawnMoves(state)

	for _, val := range out {
		move := val.White.Pawns
		if move == bb.GetSquare(2, 4) ||
			move == bb.GetSquare(3, 4) ||
			move == bb.GetSquare(2, 3) ||
			move == bb.GetSquare(2, 5) {
			//do nothing, expected result
		} else {
			fmt.Println("Unexpected Move Generated")
			t.Fail()
		}
	}
}

func TestBlackPawnGen(t *testing.T) {
	bb.InitAttackMasks()
	fmt.Println("Testing Black Move Generation")
	state := GameState { ColourState { 0, 0, 0, 0, 0, 0 }, ColourState { 0, 0, 0, 0, 0, 0 }}

	state.Black.Pawns = bb.GetSquare(6, 4)

	fmt.Println("Pawn from which generation will take place:\n")
	bb.PrintBitboard(state.Black.Pawns)

	out := BlackPawnMoves(state)
	
	for _, val := range out {
		move := val.Black.Pawns
	//	bb.PrintBitboard(move)
		if move == bb.GetSquare(5, 4) ||
			move == bb.GetSquare(4, 4) ||
			move == bb.GetSquare(5, 3) ||
			move == bb.GetSquare(5, 5) {
			//Do nothing, expected result
		} else {
			fmt.Println("Unexpected Black Move Generated")
			bb.PrintBitboard(move)
			t.Fail()
		}
	}
}
