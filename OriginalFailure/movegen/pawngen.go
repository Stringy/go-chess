package movegen

import (
	"chess/bb"
//	"fmt"
)

func WhitePawnMoves(state GameState) []GameState {
	pawns := state.White.Pawns
	black := state.Black.AllPieces()
	white := state.White.AllPieces()
	var newState GameState

	var states = make([]GameState, 0)
	
	for lsb := bb.Lsb(pawns); pawns != bb.Empty; pawns = pawns &^ lsb {
		//check if can double move
		if (lsb & bb.Rank2) != bb.Empty {
			dbm := bb.SouthOne(bb.SouthOne(lsb))
			//check blocking black pieces
			if (dbm & white) == bb.Empty || (dbm & black) == bb.Empty {
				newState = state
				newState.White.Pawns = (newState.White.Pawns &^ lsb) | dbm
				//out <- newState
				states = append(states, newState)
			}
		}
		if m := bb.SouthEastOne(lsb); m & black == bb.Empty { //attack
			if m & white == bb.Empty {
				newState = state
				newState.White.Pawns = (newState.White.Pawns &^ lsb) | m
				states = append(states, newState)
			}
		}
		if m := bb.SouthWestOne(lsb); m & black == bb.Empty { //attack
			if m & white == bb.Empty {
				newState = state
				newState.White.Pawns = (newState.White.Pawns &^ lsb) | m
				states = append(states, newState)
			}
		}
		if m := bb.SouthOne(lsb); m & black == bb.Empty { //normal move
			if m & white == bb.Empty {
				newState = state
				newState.White.Pawns = (newState.White.Pawns &^ lsb) | m
				states = append(states, newState)
			}
		}
	}
	return states
}

func BlackPawnMoves(state GameState) []GameState {
	pawns := state.Black.Pawns
	black := state.Black.AllPieces()
	white := state.White.AllPieces()
	var newState GameState
	
	states := make([]GameState, 0)

	for lsb := bb.Lsb(pawns); pawns != bb.Empty; pawns = pawns &^ lsb {
		//check if can double move
		if (lsb & bb.Rank2) != bb.Empty {
			dbm := bb.NorthOne(bb.NorthOne(lsb))
			//check blocking black pieces
			if (dbm & white) == bb.Empty || (dbm & black) == bb.Empty {
				newState = state
				newState.Black.Pawns = (newState.Black.Pawns &^ lsb) | dbm
				//out <- newState
				states = append(states, newState)
			}
		}
		if m := bb.NorthEastOne(lsb); m & black == bb.Empty { //attack
			if m & white == bb.Empty {
				newState = state
				newState.Black.Pawns = (newState.Black.Pawns &^ lsb) | m
				states = append(states, newState)
			}
		}
		if m := bb.NorthWestOne(lsb); m & black == bb.Empty { //attack
			if m & white == bb.Empty {
				newState = state
				newState.Black.Pawns = (newState.Black.Pawns &^ lsb) | m
				states = append(states, newState)
			}
		}
		if m := bb.NorthOne(lsb); m & black == bb.Empty { //normal move
			if m & white == bb.Empty {
				newState = state
				newState.Black.Pawns = (newState.Black.Pawns &^ lsb) | m
				states = append(states, newState)
			}
		}
	}

	return states
}