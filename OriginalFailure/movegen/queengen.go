package movegen 

import (
	"chess/bb"
)

func WhiteQueenMoves(state GameState) []GameState {
	queens := state.White.Queens
	black := state.Black.AllPieces()
	white := state.White.AllPieces()
	states := make([]GameState, 0)

	lsb := bb.Lsb(queens)

	for ; queens != bb.Empty; {
		mask := bb.QueenAttack(lsb)
		moves := mask &^ (black & white)
		move := bb.Lsb(moves)

		for ; moves != bb.Empty; {
			newState := state
			newState.White.Queens = (newState.White.Queens &^ lsb) | move
			states = append(states, newState)
			moves = moves &^ move
			move = bb.Lsb(moves)
		}

		queens = queens &^ lsb
		lsb = bb.Lsb(queens)
	}
	return states
}

func BlackQueenMoves(state GameState) []GameState {
	queens := state.Black.Queens
	black := state.Black.AllPieces()
	white := state.White.AllPieces()
	states := make([]GameState, 0)

	lsb := bb.Lsb(queens)

	for ; queens != bb.Empty; {
		mask := bb.QueenAttack(lsb)
		moves := mask &^ (black & white)
		move := bb.Lsb(moves)

		for ; moves != bb.Empty; {
			newState := state
			newState.Black.Queens = (newState.Black.Queens &^ lsb) | move
			states = append(states, newState)
			moves = moves &^ move
			move = bb.Lsb(moves)
		}

		queens = queens &^ lsb
		lsb = bb.Lsb(queens)
	}
	return states
}
