package movegen

import (
//	"fmt"
)

const (
	WHITE = 1
	BLACK = 0
)

type ColourState struct {
	Pawns uint64
	Rooks uint64
	Knights uint64
	Bishops uint64
	Queens uint64
	King uint64
}

type GameState struct {
	White ColourState
	Black ColourState
}

func (c *ColourState) AllPieces() uint64 {
	var union = uint64(0)
	union |= c.Pawns
	union |= c.Rooks
	union |= c.Knights
	union |= c.Bishops
	union |= c.Queens
	union |= c.King
	return union
}

func generateTree(head *Tree, depth uint, col int) {
	if depth == 0 {
		return
	}
	
	var nextCol int
	gen := []func(GameState) []GameState { GenerateBlackMoves, GenerateWhiteMoves }

	if col == BLACK { 
		nextCol = WHITE
	} else {
		nextCol = col
	}
	
	if len(head.Children) == 0 {
		moves := gen[col](head.State)
		AddChildren(head, GameStatesToTrees(moves)...)
	}

	for _, child := range head.Children {
		generateTree(child, depth - 1, nextCol)
	}
}

func GenerateFullGameTree(state GameState, depth uint) Tree {
	tree := Tree { state, make([]*Tree, 0) }
	generateTree(&tree, depth, BLACK)
	return tree
}

func GenerateMovesFor(state GameState, col int) Tree {
	tree := Tree { state, make([]*Tree, 0) }
	var moves []GameState
	
	if col == WHITE {
		moves = GenerateWhiteMoves(state)
	} else if col == BLACK {
		moves = GenerateBlackMoves(state)
	}

	AddChildren(&tree, GameStatesToTrees(moves)...)

	return tree
}

func GenerateWhiteMoves(state GameState) []GameState {
	moves := make([]GameState, 0)

	moves = append(moves, WhitePawnMoves(state)...)
	moves = append(moves, WhiteRookMoves(state)...)
	moves = append(moves, WhiteKnightMoves(state)...)
	moves = append(moves, WhiteBishopMoves(state)...)
	moves = append(moves, WhiteQueenMoves(state)...)
	moves = append(moves, WhiteKingMoves(state)...)

	return moves
}

func GenerateBlackMoves(state GameState) []GameState {
	moves := make([]GameState, 0)

	moves = append(moves, BlackPawnMoves(state)...)
	moves = append(moves, BlackRookMoves(state)...)
	moves = append(moves, BlackKnightMoves(state)...)
	moves = append(moves, BlackBishopMoves(state)...)
	moves = append(moves, BlackQueenMoves(state)...)
	moves = append(moves, BlackKingMoves(state)...)

	return moves
}