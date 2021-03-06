package eval

import (
	"chess/ai/gen"
	"math"
)

var (
	//Pawn Position value table
	//Mirrored for black pieces
	BlackPawnTable = [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		50, 50, 50, 50, 50, 50, 50, 50,
		10, 10, 20, 30, 30, 20, 10, 10,
		5, 5, 10, 27, 27, 10, 5, 5,
		0, 0, 0, 25, 25, 0, 0, 0,
		5, -5, -10, 0, 0, -10, -5, 5,
		5, 10, 10, -25, -25, 10, 10, 5,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	WhitePawnTable = [64]int{
		0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, -25, -25, 10, 10, 5,
		5, -5, -10, 0, 0, -10, -5, 5,
		0, 0, 0, 25, 25, 0, 0, 0,
		5, 5, 10, 27, 27, 10, 5, 5,
		10, 10, 20, 30, 30, 20, 10, 10,
		50, 50, 50, 50, 50, 50, 50, 50,
		0, 0, 0, 0, 0, 0, 0, 0,
	}

	//knight position value table
	//mirrored for black pieces
	KnightTable = []int{
		-50, -40, -30, -30, -30, -30, -40, -50,
		-40, -20, 0, 0, 0, 0, -20, -40,
		-30, 0, 10, 15, 15, 10, 0, -30,
		-30, 5, 15, 20, 20, 15, 5, -30,
		-30, 0, 15, 20, 20, 15, 0, -30,
		-30, 5, 10, 15, 15, 10, 5, -30,
		-40, -20, 0, 5, 5, 0, -20, -40,
		-50, -40, -20, -30, -30, -20, -40, -50,
	}

	//bishop position value table
	//mirrored for black pieces
	//emphasis on staying away from board edges
	BishopTable = []int{
		-20, -10, -10, -10, -10, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 5, 10, 10, 5, 0, -10,
		-10, 5, 5, 10, 10, 5, 5, -10,
		-10, 0, 10, 10, 10, 10, 0, -10,
		-10, 10, 10, 10, 10, 10, 10, -10,
		-10, 5, 0, 0, 0, 0, 5, -10,
		-20, -10, -40, -10, -10, -40, -10, -20,
	}

	//king position value table for early/mid game
	//Mirrored for black pieces
	//security found in corners on home rank
	KingTable = []int{
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-20, -30, -30, -40, -40, -30, -30, -20,
		-10, -20, -20, -20, -20, -20, -20, -10,
		20, 20, 0, 0, 0, 0, 20, 20,
		20, 30, 10, 0, 0, 10, 30, 20,
	}

	//King endgame position value
	//mirrored for black, though largely semetrical
	//best security is found in the centre
	KingEndGame = []int{
		-50, -40, -30, -20, -20, -30, -40, -50,
		-30, -20, -10, 0, 0, -10, -20, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -30, 0, 0, 0, 0, -30, -30,
		-50, -30, -30, -30, -30, -30, -30, -50,
	}

	//rook position table
	//mirrored for black
	RookTable = []int{
		0, 0, 0, 0, 0, 0, 0, 0,
		15, 15, 15, 15, 15, 15, 15, 15,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		-10, 0, 0, 10, 10, 0, 0, -10,
	}

	//queen position value table
	//mirrored black, though semetrical
	QueenTable = []int{
		-10, -10, -10, -10, -10, -10, -10, -10,
		-10, 0, 0, 0, 0, 0, 0, 0,
		-10, 0, 5, 5, 5, 5, 0, -10,
		-10, 0, 5, 10, 10, 5, 0, -10,
		-10, 0, 5, 10, 10, 5, 0, -10,
		-10, 0, 5, 5, 5, 5, 0, -10,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, -10, -10, -10, -10, -10, -10, -10,
	}

	//mirror provides reverse indexing for all
	//position value tables
	Mirror = []int{
		56, 57, 58, 59, 60, 61, 62, 63,
		48, 49, 50, 51, 52, 53, 54, 55,
		40, 41, 42, 43, 44, 45, 46, 47,
		32, 33, 34, 35, 36, 37, 38, 39,
		24, 25, 26, 27, 28, 29, 30, 31,
		16, 17, 18, 19, 20, 21, 22, 23,
		8, 9, 10, 11, 12, 13, 14, 15,
		0, 1, 2, 3, 4, 5, 6, 7,
	}

	//Distance from square to king square [sq][ksq]
	Distance [][]int

	//individual piece distances
	PawnOwnDistance      = []int{0, 8, 4, 2, 0, 0, 0, 0}
	PawnOpponentDistance = []int{0, 2, 1, 0, 0, 0, 0, 0}
	KnightDistance       = []int{0, 4, 4, 0, 0, 0, 0, 0}
	BishopDistance       = []int{0, 5, 4, 3, 2, 1, 0, 0}
	RookDistance         = []int{0, 7, 5, 4, 3, 0, 0, 0}
	QueenDistance        = []int{0, 10, 8, 5, 4, 0, 0, 0}

	//masks for calculating passed pawns
	//isolated pawns and backwards pawns on 
	//any square
	PassedWhite   []uint64
	PassedBlack   []uint64
	IsolatedWhite []uint64
	IsolatedBlack []uint64
	BackWhite     []uint64
	BackBlack     []uint64

	//king pawn weak/strong shielding 
	StrongShieldWhite []uint64
	StrongShieldBlack []uint64
	WeakShieldBlack   []uint64
	WeakShieldWhite   []uint64

	//Piece Values for evaluation
	PawnValue   = uint(10)
	KnightValue = uint(32)
	BishopValue = uint(35)
	RookValue   = uint(50)
	QueenValue  = uint(97)
	KingValue   = uint(99999)

	WhiteSquares = uint64(0xAAAAAAAAAAAAAAAA)
	BlackSquares = uint64(0x5555555555555555)

	//Bonuses
	PassedPawnBonus     = 15
	BishopPairBonus     = 45
	RookPassedPawnBonus = 30
	StrongShieldBonus   = 20
	WeakShieldBonus     = 10
	PinAndForkBonus     = 25
	DefenseBonus        = 15

	//Penalties
	DoublePawnPenalty   = 20
	IsolatedPawnPenalty = 20
	BackPawnPenalty     = 20

	//Constants
	CheckMate = uint(99999)
	Draw      = int(99990)
)

//initialiseGlobals initialises all data used during the evaluation process
func initialiseGlobals() {
	initDistanceArray()
	initialisePawnMasks()
	initialiseKingShields()
}

//initDistanceArray initialises the distance array with distances from square i
//to square sq
func initDistanceArray() {
	Distance = make([][]int, 64)
	for i := 0; i < 64; i++ {
		Distance[i] = make([]int, 64)
	}
	for i := 0; i < 64; i++ {
		for sq := 0; sq < 64; sq++ {
			if math.Abs(float64(gen.Ranks[i]-gen.Ranks[sq])) <
				math.Abs(float64(gen.Files[i]-gen.Files[sq])) {
				Distance[i][sq] = int(math.Abs(float64(gen.Ranks[i] - gen.Ranks[sq])))
			} else {
				Distance[i][sq] = int(math.Abs(float64(gen.Files[i] - gen.Files[sq])))
			}
		}
	}
}

//initialisePawnMasks initialises all pawn related masks
//first white id generated and then black is initialised with the mirror 
//of the white values
func initialisePawnMasks() {
	PassedWhite = make([]uint64, 64)
	PassedBlack = make([]uint64, 64)
	IsolatedBlack = make([]uint64, 64)
	IsolatedWhite = make([]uint64, 64)
	BackBlack = make([]uint64, 64)
	BackWhite = make([]uint64, 64)

	for i := 0; i < 64; i++ {
		//Passed White
		for rank := gen.Ranks[i] + 1; rank < 8; rank++ {
			if gen.Files[i]-1 > 0 {
				PassedWhite[i] ^= gen.BitSet[gen.Index[gen.Files[i]-1][rank]]
				PassedBlack[Mirror[i]] ^= gen.BitSet[Mirror[gen.Index[gen.Files[i]-1][rank]]]
			}
			PassedWhite[i] ^= gen.BitSet[gen.Index[gen.Files[i]][rank]]
			PassedBlack[Mirror[i]] ^= gen.BitSet[Mirror[gen.Index[gen.Files[i]][rank]]]
			if gen.Files[i]+1 <= 8 {
				PassedWhite[i] ^= gen.BitSet[gen.Index[gen.Files[i]+1][rank]]
				PassedBlack[i] ^= gen.BitSet[Mirror[gen.Index[gen.Files[i]+1][rank]]]
			}
		}
		//Isolated White
		for rank := 2; rank <= gen.Ranks[i]; rank++ {
			if gen.Files[i]-1 > 0 {
				IsolatedWhite[i] ^= gen.BitSet[gen.Index[gen.Files[i]-1][rank]]
				IsolatedBlack[Mirror[i]] ^= gen.BitSet[Mirror[gen.Index[gen.Files[i]-1][rank]]]
			}
			if gen.Files[i]+1 <= 8 {
				IsolatedWhite[i] ^= gen.BitSet[gen.Index[gen.Files[i]+1][rank]]
				IsolatedBlack[Mirror[i]] ^= gen.BitSet[Mirror[gen.Index[gen.Files[i]+1][rank]]]
			}
		}
		//backward white
		for rank := 2; rank <= gen.Ranks[i]; rank++ {
			if gen.Files[i]-1 > 0 {
				BackWhite[i] ^= gen.BitSet[gen.Index[gen.Files[i]-1][rank]]
				BackBlack[Mirror[i]] ^= gen.BitSet[Mirror[gen.Index[gen.Files[i]-1][rank]]]
			}
			if gen.Files[i]+1 <= 8 {
				BackWhite[i] ^= gen.BitSet[gen.Index[gen.Files[i]+1][rank]]
				BackBlack[Mirror[i]] ^= gen.BitSet[Mirror[gen.Index[gen.Files[i]+1][rank]]]
			}
		}
	}
}

//initialiseKingShields initialises all pawn shield masks
//white shields are generated first and then black shields are 
//initialised from the mirror of the white values
func initialiseKingShields() {
	StrongShieldBlack = make([]uint64, 64)
	StrongShieldWhite = make([]uint64, 64)
	WeakShieldWhite = make([]uint64, 64)
	WeakShieldBlack = make([]uint64, 64)

	for i := 0; i < 24; i++ {
		StrongShieldWhite[i] ^= gen.BitSet[i+8]
		WeakShieldWhite[i] ^= gen.BitSet[i+16]
		if gen.Files[i] > 1 {
			StrongShieldWhite[i] ^= gen.BitSet[i+7]
			WeakShieldWhite[i] ^= gen.BitSet[i+15]
		}
		if gen.Files[i] < 1 {
			StrongShieldWhite[i] ^= gen.BitSet[i+9]
			WeakShieldWhite[i] ^= gen.BitSet[i+17]
		}
		if gen.Files[i] == 1 {
			StrongShieldWhite[i] ^= gen.BitSet[i+10]
			WeakShieldWhite[i] ^= gen.BitSet[i+18]
		}
		if gen.Files[i] == 8 {
			StrongShieldWhite[i] ^= gen.BitSet[i+6]
			WeakShieldWhite[i] ^= gen.BitSet[i+14]
		}
	}

	for sq := 0; sq < 64; sq++ {
		if StrongShieldWhite[sq]&gen.BitSet[sq] != 0 {
			StrongShieldBlack[Mirror[sq]] |= gen.BitSet[Mirror[sq]]
		}
		if WeakShieldBlack[sq]&gen.BitSet[sq] != 0 {
			WeakShieldBlack[Mirror[sq]] |= gen.BitSet[Mirror[sq]]
		}
	}
}
