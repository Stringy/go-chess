package movegen

import (
	"chess/bb"
)

func WhiteRookMoves(state GameState) []GameState {
	black := state.Black.AllPieces()
	white := state.White.AllPieces()
	states := make([]GameState, 0)
	rooks := state.White.Rooks

	lsb := bb.Lsb(rooks)
	for ; rooks != bb.Empty; {

		mask := bb.RookAttack(lsb)
		moves := mask &^ (black & white)
		move := bb.Lsb(moves)

		for ; moves != bb.Empty; { 
			newState := state
			newState.White.Rooks = (newState.White.Rooks &^ lsb) | move
			states = append(states, newState)
			moves = moves &^ move
			move = bb.Lsb(moves)
		}
		
		rooks = rooks &^ lsb
		lsb = bb.Lsb(rooks)
	}
	return states
}

func BlackRookMoves(state GameState) []GameState {
	black := state.Black.AllPieces()
	white := state.White.AllPieces()
	states := make([]GameState, 0)
	rooks := state.Black.Rooks

	lsb := bb.Lsb(rooks)
	
	for ; rooks != bb.Empty; {

		mask := bb.RookAttack(lsb)
		moves := mask &^ (black & white)
		move := bb.Lsb(moves)

		for ; moves != bb.Empty; {
			newState := state
			newState.White.Rooks = (newState.White.Rooks &^ lsb) | move
			states = append(states, newState)
			moves = moves &^ move
			move = bb.Lsb(moves)
		}
		
		rooks = rooks &^ lsb
		lsb = bb.Lsb(rooks)
	}
	return states
}

