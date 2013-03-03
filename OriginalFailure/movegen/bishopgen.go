package movegen

import (
	"chess/bb"
)

func WhiteBishopMoves(state GameState) []GameState {
	bishops := state.White.Bishops
	black := state.Black.AllPieces()
	white := state.White.AllPieces()

	states := make([]GameState, 0)

	lsb := bb.Lsb(bishops)
	for ; bishops != bb.Empty; {
		
		mask := bb.BishopAttack(lsb)
		moves := mask &^ (black & white)
		move := bb.Lsb(moves)

		for ; moves != bb.Empty; {
			newState := state
			newState.White.Bishops = (newState.White.Bishops &^ lsb) | move
			states = append(states, newState)
			moves = moves &^ move
			move = bb.Lsb(moves)
		}
		
		bishops = bishops &^ lsb
		lsb = bb.Lsb(bishops)
	}
	return states
}

func BlackBishopMoves(state GameState) []GameState {
	bishops := state.Black.Bishops
	black := state.Black.AllPieces()
	white := state.White.AllPieces()

	states := make([]GameState, 0)

	lsb := bb.Lsb(bishops)
	for ; bishops != bb.Empty; {
		
		mask := bb.BishopAttack(lsb)
		moves := mask &^ (black & white)
		move := bb.Lsb(moves)

		for ; moves != bb.Empty; {
			newState := state
			newState.Black.Bishops = (newState.Black.Bishops &^ lsb) | move
			states = append(states, newState)
			moves = moves &^ move
			move = bb.Lsb(moves)
		}
		
		bishops = bishops &^ lsb
		lsb = bb.Lsb(bishops)
	}
	return states
}