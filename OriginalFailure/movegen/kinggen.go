package movegen

import (
	"chess/bb"
)

func WhiteKingMoves(state GameState) []GameState {
	king := state.White.King
	white := state.White.AllPieces()
	black := state.Black.AllPieces()
	states := make([]GameState, 0)

	lsb := bb.Lsb(king)
	for ; king != bb.Empty; {
		mask := bb.KingAttack(lsb)
		moves := mask &^ (black & white)
		move := bb.Lsb(moves)

		for ; move != bb.Empty; {
			newState := state
			newState.White.King = (newState.White.King &^ lsb) | move
			states = append(states, newState)
			moves = moves &^ move
			move = bb.Lsb(move)
		}
		
		king = king &^ lsb
		lsb = bb.Lsb(king)
	}
	return states
}

func BlackKingMoves(state GameState) []GameState {
	king := state.Black.King
	white := state.White.AllPieces()
	black := state.Black.AllPieces()
	states := make([]GameState, 0)

	lsb := bb.Lsb(king)
	for ; king != bb.Empty; {
		mask := bb.KingAttack(lsb)
		moves := mask &^ (black & white)
		move := bb.Lsb(moves)

		for ; move != bb.Empty; {
			newState := state
			newState.Black.King = (newState.Black.King &^ lsb) | move
			states = append(states, newState)
			moves = moves &^ move
			move = bb.Lsb(move)
		}
		
		king = king &^ lsb
		lsb = bb.Lsb(king)
	}
	return states
}

