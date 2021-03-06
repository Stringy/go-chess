package gen

import (
	"fmt"
)

//Board structure used to contain all information regarding the current board state
//includes bitboards for every piece type along with current turn and castle status
//of each player
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

	//all pieces masks
	WhitePieces     uint64
	BlackPieces     uint64
	OccupiedSquares uint64

	NextMove    byte //white or black to move next
	WhiteCastle byte //white castle status
	BlackCastle byte //black castle status
	Ep          int  //enpassant target
	NumMoves    int  //number of moves (used to count fifty move rule)

	Material int //value of all pieces on board

	Squares []byte //current occupancy of all squares

	//number of white pieces on the board
	NumWPawns,
	NumWRooks,
	NumWKnights,
	NumWBishops,
	NumWQueens uint

	//number of black pieces on the board
	NumBPawns,
	NumBRooks,
	NumBKnights,
	NumBBishops,
	NumBQueens uint

	//record of king square positions
	//used for evaluation
	WKingSq byte
	BKingSq byte

	//Colour specific material values
	WMaterial int
	BMaterial int
}

//New Board returns a pointer to an initialised board
//with the standard chess start position
func NewBoard() *Board {
	var board Board
	board.Init()
	return &board
}

//Board.Init initialises a new board state with initial position of
//all the pieces along with material values and castling statuses
func (b *Board) Init() {
	b.Squares = make([]byte, 64)

	for i := 0; i < 64; i++ {
		b.Squares[i] = Empty
	}

	b.Squares[SquareMap["E1"]] = WhiteKing
	b.WKingSq = SquareMap["E1"]
	b.Squares[SquareMap["D1"]] = WhiteQueen
	b.Squares[SquareMap["A1"]] = WhiteRook
	b.Squares[SquareMap["H1"]] = WhiteRook
	b.Squares[SquareMap["B1"]] = WhiteKnight
	b.Squares[SquareMap["G1"]] = WhiteKnight
	b.Squares[SquareMap["C1"]] = WhiteBishop
	b.Squares[SquareMap["F1"]] = WhiteBishop
	b.Squares[SquareMap["A2"]] = WhitePawn
	b.Squares[SquareMap["B2"]] = WhitePawn
	b.Squares[SquareMap["C2"]] = WhitePawn
	b.Squares[SquareMap["D2"]] = WhitePawn
	b.Squares[SquareMap["E2"]] = WhitePawn
	b.Squares[SquareMap["F2"]] = WhitePawn
	b.Squares[SquareMap["G2"]] = WhitePawn
	b.Squares[SquareMap["H2"]] = WhitePawn

	b.Squares[SquareMap["E8"]] = BlackKing
	b.BKingSq = SquareMap["E8"]
	b.Squares[SquareMap["D8"]] = BlackQueen
	b.Squares[SquareMap["A8"]] = BlackRook
	b.Squares[SquareMap["H8"]] = BlackRook
	b.Squares[SquareMap["B8"]] = BlackKnight
	b.Squares[SquareMap["G8"]] = BlackKnight
	b.Squares[SquareMap["C8"]] = BlackBishop
	b.Squares[SquareMap["F8"]] = BlackBishop
	b.Squares[SquareMap["A7"]] = BlackPawn
	b.Squares[SquareMap["B7"]] = BlackPawn
	b.Squares[SquareMap["C7"]] = BlackPawn
	b.Squares[SquareMap["D7"]] = BlackPawn
	b.Squares[SquareMap["E7"]] = BlackPawn
	b.Squares[SquareMap["F7"]] = BlackPawn
	b.Squares[SquareMap["G7"]] = BlackPawn
	b.Squares[SquareMap["H7"]] = BlackPawn

	b.InitFromGrid(b.Squares, WhiteMove, 0, CanCastleOO+CanCastleOOO, CanCastleOO+CanCastleOOO, 0)
}

//Board.InitFromGrid initialises the board from the array of bytes containing the
//position of the pieces. Initialises material values and castle statuses
func (b *Board) InitFromGrid(grid []byte, moves, ep int, cw, cb, next byte) {
	if len(b.Squares) != 64 {
		b.Squares = make([]byte, 64)
	}
	for i := 0; i < 64; i++ {
		b.Squares[i] = grid[i]
		switch b.Squares[i] {
		case WhiteKing:
			b.WhiteKing |= BitSet[i]
		case WhiteQueen:
			b.WhiteQueens |= BitSet[i]
		case WhiteBishop:
			b.WhiteBishops |= BitSet[i]
		case WhiteKnight:
			b.WhiteKnights |= BitSet[i]
		case WhiteRook:
			b.WhiteRooks |= BitSet[i]
		case WhitePawn:
			b.WhitePawns |= BitSet[i]

		case BlackKing:
			b.BlackKing |= BitSet[i]
		case BlackQueen:
			b.BlackQueens |= BitSet[i]
		case BlackBishop:
			b.BlackBishops |= BitSet[i]
		case BlackKnight:
			b.BlackKnights |= BitSet[i]
		case BlackRook:
			b.BlackRooks |= BitSet[i]
		case BlackPawn:
			b.BlackPawns |= BitSet[i]
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

	b.NumBPawns = 8
	b.NumBRooks = 2
	b.NumBKnights = 2
	b.NumBBishops = 2
	b.NumBQueens = 1

	b.NumWPawns = 8
	b.NumWRooks = 2
	b.NumWKnights = 2
	b.NumWBishops = 2
	b.NumWQueens = 1

	b.WMaterial = int(
		BitCount(b.WhitePawns)*PawnValue +
			BitCount(b.WhiteRooks)*RookValue +
			BitCount(b.WhiteKnights)*KnightValue +
			BitCount(b.WhiteBishops)*BishopValue +
			BitCount(b.WhiteQueens)*QueenValue)

	b.BMaterial = -int(
		BitCount(b.BlackPawns)*PawnValue +
			BitCount(b.BlackRooks)*RookValue +
			BitCount(b.BlackKnights)*KnightValue +
			BitCount(b.BlackBishops)*BishopValue +
			BitCount(b.BlackQueens)*QueenValue)

	b.Material = b.WMaterial - b.BMaterial
}

//Board.PrintBoard prints a representation of the board
func (b *Board) PrintBoard() {
	fmt.Println("    a   b   c   d   e   f   g   h")
	for rank := 8; rank > 0; rank-- {
		fmt.Println("  +---+---+---+---+---+---+---+---+")
		fmt.Print("  |")
		for file := 1; file <= 8; file++ {
			fmt.Print(" ", PieceNames[b.Squares[Index[file][rank]]]+"|")
		}
		fmt.Println(rank)
	}
	fmt.Println("  +---+---+---+---+---+---+---+---+")
}

//Board.Clone creates a copy of the board
//This is primarily used to pass distinct copies to go routines
//for safer searching.
func (b *Board) Clone() *Board {
	clone := new(Board)

	clone.WhitePawns = b.WhitePawns
	clone.WhiteRooks = b.WhiteRooks
	clone.WhiteKnights = b.WhiteKnights
	clone.WhiteBishops = b.WhiteBishops
	clone.WhiteQueens = b.WhiteQueens
	clone.WhiteKing = b.WhiteKing

	clone.BlackPawns = b.BlackPawns
	clone.BlackRooks = b.BlackRooks
	clone.BlackKnights = b.BlackKnights
	clone.BlackBishops = b.BlackBishops
	clone.BlackQueens = b.BlackQueens
	clone.BlackKing = b.BlackKing

	clone.Material = b.Material
	clone.WhitePieces = b.WhitePieces
	clone.BlackPieces = b.BlackPieces
	clone.OccupiedSquares = b.OccupiedSquares

	clone.NextMove = b.NextMove
	clone.WhiteCastle = b.WhiteCastle
	clone.BlackCastle = b.BlackCastle
	clone.Ep = b.Ep
	clone.NumMoves = b.NumMoves

	clone.Squares = make([]byte, len(b.Squares))
	copy(clone.Squares, b.Squares)

	return clone
}

//IsOtherKingAttacked returns true if the opponent king
//is being attacked by any piece. 
func (b *Board) IsOtherPlayerChecked() bool {
	if b.NextMove == BlackMove {
		return b.isAttacked(b.WhiteKing, BlackMove)
	}
	return b.isAttacked(b.BlackKing, WhiteMove)
}

//IsCheck determined whether the player to move is in check or not
func (b *Board) IsCheck() bool {
	if b.NextMove == BlackMove {
		return b.isAttacked(b.BlackKing, WhiteMove)
	}
	return b.isAttacked(b.WhiteKing, BlackMove)
}

//IsCheckMate returns true if the current player to move is checkmated, 
//else returns false
func (b *Board) IsCheckMate() bool {
	if b.NextMove == BlackMove {
		if !b.isAttacked(b.BlackKing, WhiteMove) {
			return false
		} else {
			kingMoves := generateBlackKingMoves(b)
			whiteMoves := generateAllWhiteMoves(b)
			newboard := b.Clone()
			newboard.MakeAllMoves(append(kingMoves, whiteMoves...))
			for _, piece := range newboard.Squares {
				if piece == BlackKing {
					return false
				}
			}
			return true
		}
	} else {
		if !b.isAttacked(b.WhiteKing, BlackMove) {
			return false
		} else {
			kingMoves := generateWhiteKingMoves(b)
			blackMoves := generateAllBlackMoves(b)
			newboard := b.Clone()
			newboard.MakeAllMoves(append(kingMoves, blackMoves...))
			for _, piece := range newboard.Squares {
				if piece == BlackKing {
					return false
				}
			}
			return true
		}
	}
	return false
}

//isAttacked works out if a particular square is attacked by a specific piece type. 
//
//returns true if attacked, false otherwise
func (b *Board) isAttacked(bits uint64, attacker byte) bool {
	targ := bits
	if attacker == BlackMove {
		moves := generateAllBlackMoves(b)
		for targ != 0 {
			to := LSB(targ)
			for _, move := range moves {
				if uint(move.GetTo()) == to {
					return true
				}
			}
			targ ^= BitSet[to]
		}
	} else {
		moves := generateAllWhiteMoves(b)
		for targ != 0 {
			to := LSB(targ)
			for _, move := range moves {
				if uint(move.GetTo()) == to {
					return true
				}
			}
			targ ^= BitSet[to]
		}
	}
	return false
}

//NullMove makes a null move on the board, that is changes the turn
func (b *Board) NullMove() {
	if b.NextMove == WhiteMove {
		b.NextMove = BlackMove
		return
	}
	b.NextMove = WhiteMove
}
