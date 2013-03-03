package be

import (
	"fmt"
)

type Board struct {
	//white piece masks
	WhitePawns,
	WhiteRooks,
	WhiteKnights,
	WhiteBishops,
	WhiteQueens,
	WhiteKing uint64

	//black piece masks
	BlackPawns,
	BlackRooks,
	BlackKnights,
	BlackBishops,
	BlackQueens,
	BlackKing uint64

	WhitePieces uint64
	BlackPieces uint64
	OccupiedSquares uint64

	NextMove byte //white or black to move next
	WhiteCastle byte //white castle status
	BlackCastle byte //black castle status
	Ep int //enpassant target
	NumMoves int //number of moves (used to count fifty move rule)

	Material int //value of all pieces on board

	Squares []byte //current occupancy of all squares
}

func (b *Board) Init() {
	b.Squares = make([]byte, 64)

	for i := 0; i < 64; i++ {
		b.Squares[i] = Empty
	}

   b.Squares[E1] = WhiteKing;
   b.Squares[D1] = WhiteQueen;
   b.Squares[A1] = WhiteRook;
   b.Squares[H1] = WhiteRook;
   b.Squares[B1] = WhiteKnight;
   b.Squares[G1] = WhiteKnight;
   b.Squares[C1] = WhiteBishop;
   b.Squares[F1] = WhiteBishop;
   b.Squares[A2] = WhitePawn;
   b.Squares[B2] = WhitePawn;
   b.Squares[C2] = WhitePawn;
   b.Squares[D2] = WhitePawn;
   b.Squares[E2] = WhitePawn;
   b.Squares[F2] = WhitePawn;
   b.Squares[G2] = WhitePawn;
   b.Squares[H2] = WhitePawn;
	
   b.Squares[E8] = BlackKing;
   b.Squares[D8] = BlackQueen;
   b.Squares[A8] = BlackRook;
   b.Squares[H8] = BlackRook;
   b.Squares[B8] = BlackKnight;
   b.Squares[G8] = BlackKnight;
   b.Squares[C8] = BlackBishop;
   b.Squares[F8] = BlackBishop;
   b.Squares[A7] = BlackPawn;
   b.Squares[B7] = BlackPawn;
   b.Squares[C7] = BlackPawn;
   b.Squares[D7] = BlackPawn;
   b.Squares[E7] = BlackPawn;
   b.Squares[F7] = BlackPawn;
   b.Squares[G7] = BlackPawn;
   b.Squares[H7] = BlackPawn;
	
	b.InitFromGrid(b.Squares, WhiteMove, 0, CanCastleOO + CanCastleOOO, CanCastleOO + CanCastleOOO, 0)
}

func (b *Board) InitFromGrid(grid []byte, moves, ep int, cw, cb, next byte) {
	for i := 0; i < 64; i++ {
		b.Squares[i] = grid[i]
		switch b.Squares[i] {
		case WhiteKing: b.WhiteKing |= BitSet[i]
		case WhiteQueen: b.WhiteQueens |= BitSet[i]
		case WhiteBishop: b.WhiteBishops |= BitSet[i]
		case WhiteKnight: b.WhiteKnights |= BitSet[i]
		case WhiteRook: b.WhiteRooks |= BitSet[i]
		case WhitePawn: b.WhitePawns |= BitSet[i]

		case BlackKing: b.BlackKing |= BitSet[i]
		case BlackQueen: b.BlackQueens |= BitSet[i]
		case BlackBishop: b.BlackBishops |= BitSet[i]
		case BlackKnight: b.BlackKnights |= BitSet[i]
		case BlackRook: b.BlackRooks |= BitSet[i]
		case BlackPawn: b.BlackPawns |= BitSet[i]
		}
	}
	b.WhitePieces = 
		b.WhiteKing |
		b.WhiteQueens |
		b.WhiteBishops |
		b.WhiteKnights |
		b.WhiteRooks |
		b.WhitePawns

	b.BlackPieces = 
		b.BlackKing |
		b.BlackQueens |
		b.BlackBishops |
		b.BlackKnights |
		b.BlackRooks |
		b.BlackPawns

	b.NumMoves = moves
	b.OccupiedSquares = b.BlackPieces | b.WhitePieces
	b.NextMove = next
	b.WhiteCastle = cw
	b.BlackCastle = cb
	b.Ep = ep

	b.Material = int(
		BitCount(b.WhitePawns) * PawnValue +
		BitCount(b.WhiteRooks) * RookValue +
		BitCount(b.WhiteKnights) * KnightValue +
		BitCount(b.WhiteBishops) * BishopValue +
		BitCount(b.WhiteQueens) * QueenValue)

	b.Material -= int(
		BitCount(b.BlackPawns) * PawnValue +
		BitCount(b.BlackRooks) * RookValue +
		BitCount(b.BlackKnights) * KnightValue +
		BitCount(b.BlackBishops) * BishopValue +
		BitCount(b.BlackQueens) * QueenValue)

}

func (b *Board) PrintBoard() {
	fmt.Println("    h   g   f   e   d   c   b   a")
	for rank := 1; rank < 9; rank++ {
		fmt.Println("  +---+---+---+---+---+---+---+---+")
		fmt.Print("  |")
		for file := 8; file > 0; file-- {
			fmt.Print(" ", PieceNames[b.Squares[Index[file][rank]]] + "|")
		}
		fmt.Println(rank)
	}
	fmt.Println("  +---+---+---+---+---+---+---+---+")
}