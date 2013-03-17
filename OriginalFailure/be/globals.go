package be

import (
	//"fmt"
)

const (
	MaxMoves = 4096
	MaxPly = 64
	
	//Square Indexes
	A1 = 0; A2 = 8; A3 = 16; A4 = 24; A5 = 32; A6 = 40; A7 = 48; A8 = 56;
	B1 = 1; B2 = 9; B3 = 17; B4 = 25; B5 = 33; B6 = 41; B7 = 49; B8 = 57;
	C1 = 2; C2 = 10; C3 = 18; C4 = 26; C5 = 34; C6 = 42; C7 = 50; C8 = 58;
	D1 = 3; D2 = 11; D3 = 19; D4 = 27; D5 = 35; D6 = 43; D7 = 51; D8 = 59;
	E1 = 4; E2 = 12; E3 = 20; E4 = 28; E5 = 36; E6 = 44; E7 = 52; E8 = 60;
	F1 = 5; F2 = 13; F3 = 21; F4 = 29; F5 = 37; F6 = 45; F7 = 53; F8 = 61;
	G1 = 6; G2 = 14; G3 = 22; G4 = 30; G5 = 38; G6 = 46; G7 = 54; G8 = 62;
	H1 = 7; H2 = 15; H3 = 23; H4 = 31; H5 = 39; H6 = 47; H7 = 55; H8 = 63;

	//Move Identifiers
	WhiteMove = 0
	BlackMove = 1
	
/*
	Piece Identifiers
	4 bits per pieces
	Properties:
	 white = 0...
	 black = 1...
	 sliding = .1..
*/
	Empty = byte(0) //0000
	WhitePawn = byte(1) //0001
	WhiteKing = byte(2) //0010
	WhiteKnight = byte(3) //0011
	WhiteBishop = byte(5) //0101
	WhiteRook = byte(6) //0110
	WhiteQueen = byte(7) //0111

	BlackPawn = byte(9) //1001
	BlackKing = byte(10) //1010
	BlackKnight = byte(11) //1011
	BlackBishop = byte(13) //1101
	BlackRook = byte(14) //1110
	BlackQueen = byte(15) //1111	

	//Piece Values
	//used for evaluation in decision making
	PawnValue = 100
	KnightValue = 300
	BishopValue = 300
	RookValue = 500
	QueenValue = 900
	KingValue = 9999
	CheckMate = KingValue

	//Castling Constants
	CanCastleOO = byte(1)
	CanCastleOOO = byte(2)

)

var (
	//Rank and File Masks
	//Used for identifying rank / file from square indexes 0-63
	Files = []int {
		1,2,3,4,5,6,7,8,
		1,2,3,4,5,6,7,8,
		1,2,3,4,5,6,7,8,
		1,2,3,4,5,6,7,8,
		1,2,3,4,5,6,7,8,
		1,2,3,4,5,6,7,8,
		1,2,3,4,5,6,7,8,
		1,2,3,4,5,6,7,8,
	}

	Ranks = []int {
		1,1,1,1,1,1,1,1,
		2,2,2,2,2,2,2,2,
		3,3,3,3,3,3,3,3,
		4,4,4,4,4,4,4,4,
		5,5,5,5,5,5,5,5,
		6,6,6,6,6,6,6,6,
		7,7,7,7,7,7,7,7,
		8,8,8,8,8,8,8,8,
	}

	//Piece names for printing the board
	PieceNames = []string {
		"  ", "P ", "K ", "N ", "  ", "B ", "R ", "Q ",
		"  ", "P*", "K*", "N*", "  ", "B*", "R*", "Q*",
	}
	
	//array of bitmaps with individual bits set
	//used for fast lookup of positions
	BitSet = make([]uint64, 64)

	//Array of Most significant bits, used for Eugene Nalimov's
	// reverse bit scan
	MsbTable = make([]int, 256)

	//Indexes Ranks and Files to 0..63
	Index = [9][9]int {}

	//Colour Castling
	WhiteCastleOO uint
	WhiteCastleOOO uint
	BlackCastleOO uint
	BlackCastleOOO uint

	BlackSquares uint64
	WhiteSquares uint64

	PassedWhite []uint64
	PassedBlack []uint64
	IsolatedWhite []uint64
	IsolatedBlack []uint64
	BackWhite []uint64
	BackBlack []uint64
	
	KingShieldSWhite []uint64
	KingShieldSBlack []uint64
	KingShieldWWhite []uint64
	KingShieldWBlack []uint64
)

func InitAllData() {
	initBitSet()
	initBoardIndices()
	initMsbTable()
	InitialiseAllMasks() //from attackmasks.go
}

func initBitSet() {
	BitSet[0] = 0x1
	for i := 1; i < 64; i++ {
		BitSet[i] = BitSet[i - 1] << 1
	}
}

func initBoardIndices() {
	for rank := 0; rank < 9; rank++ {
		for file := 0; file < 9; file++ {
			Index[file][rank] = (rank - 1) * 8 + file - 1
		}
	}
}

func initMsbTable() {
	for i := 0; i < 256; i++ {
		if i > 127 {
			MsbTable[i] = 7
		} else if i > 63 {
			MsbTable[i] = 6
		} else if i > 31 {
			MsbTable[i] = 5
		} else if i > 15 {
			MsbTable[i] = 4
		} else if i > 7 {
			MsbTable[i] = 3
		} else if i > 3 {
			MsbTable[i] = 2
		} else if i > 1 {
			MsbTable[i] = 1
		} else {
			MsbTable[i] = 0
		}
	}
}

func initSquareMasks() {
	BlackSquares = 0
	for i := 0; i < 64; i++ {
		if (i + Ranks[i]) % 2 != 0 {
			BlackSquares ^= BitSet[i]
		}
	}
	WhiteSquares = Universe &^ BlackSquares
}