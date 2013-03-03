package movegen

import (
	"chess/bb"
)

func WhiteKnightMoves(state GameState) []GameState {
	knights := state.White.Knights 
	black := state.Black.AllPieces()
	white := state.White.AllPieces()
	states := make([]GameState, 0)

	lsb := bb.Lsb(knights)
	for ; knights != bb.Empty; {
		mask := bb.KnightAttack(lsb)
		moves := mask &^ (black & white)
		move := bb.Lsb(moves)
		
		for ; moves != bb.Empty; {
			newState := state
			newState.White.Knights = (newState.White.Knights &^ lsb) | move
			states = append(states, newState)
			moves = moves &^ move
			move = bb.Lsb(moves)
		}

		knights = knights &^ lsb
		lsb = bb.Lsb(knights)
	}
	return states
}

func BlackKnightMoves(state GameState) []GameState {
	knights := state.Black.Knights 
	black := state.Black.AllPieces()
	white := state.White.AllPieces()
	states := make([]GameState, 0)

	lsb := bb.Lsb(knights)
	for ; knights != bb.Empty; {
		mask := bb.KnightAttack(lsb)
		moves := mask &^ (black & white)
		move := bb.Lsb(moves)
		
		for ; moves != bb.Empty; {
			newState := state
			newState.Black.Knights = (newState.Black.Knights &^ lsb) | move
			states = append(states, newState)
			moves = moves &^ move
			move = bb.Lsb(moves)
		}

		knights = knights &^ lsb
		lsb = bb.Lsb(knights)
	}
	return states
}

