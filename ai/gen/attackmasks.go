package gen

import (
	"math"
)

const ()

var (
	//Pawn Attack / Move Masks
	WhitePawnAttacks     = make([]uint64, 64)
	WhitePawnMoves       = make([]uint64, 64)
	WhitePawnDoubleMoves = make([]uint64, 64)

	BlackPawnAttacks     = make([]uint64, 64)
	BlackPawnMoves       = make([]uint64, 64)
	BlackPawnDoubleMoves = make([]uint64, 64)

	//Other Piece Attacks
	KingAttacks     = make([]uint64, 64)
	KnightAttacks   = make([]uint64, 64)
	RankAttacks     = [64][64]uint64{} //for rooks
	FileAttacks     = [64][64]uint64{} //for rooks
	DiagA1H8Attacks = [64][64]uint64{} //for bishops
	DiagA8H1Attacks = [64][64]uint64{} //for bishops

	RankMask      = make([]uint64, 64)
	FileMask      = make([]uint64, 64)
	FileMagic     = make([]uint64, 64)
	DiagA8H1Mask  = make([]uint64, 64)
	DiagA8H1Magic = make([]uint64, 64)
	DiagA1H8Mask  = make([]uint64, 64)
	DiagA1H8Magic = make([]uint64, 64)

	GeneralSlidingAttacks = [8][64]byte{}

	MaskEG = make([]uint64, 2)
	MaskFG = make([]uint64, 2)
	MaskBD = make([]uint64, 2)
	MaskCE = make([]uint64, 2)

	CharBitSet = make([]byte, 8)

	RankShift = []int{
		1, 1, 1, 1, 1, 1, 1, 1,
		9, 9, 9, 9, 9, 9, 9, 9,
		17, 17, 17, 17, 17, 17, 17, 17,
		25, 25, 25, 25, 25, 25, 25, 25,
		33, 33, 33, 33, 33, 33, 33, 33,
		41, 41, 41, 41, 41, 41, 41, 41,
		49, 49, 49, 49, 49, 49, 49, 49,
		57, 57, 57, 57, 57, 57, 57, 57,
	}

	//File Magic Bitboard Mask
	FileMagicMask = []uint64{
		0x8040201008040200,
		0x4020100804020100,
		0x2010080402010080,
		0x1008040201008040,
		0x0804020100804020,
		0x0402010080402010,
		0x0201008040201008,
		0x0100804020100804,
	}

	//A8 H1 Diagonal Magic bitboard mask
	DiagA8H1MagicMask = []uint64{
		0x0,
		0x0,
		0x0101010101010100,
		0x0101010101010100,
		0x0101010101010100,
		0x0101010101010100,
		0x0101010101010100,
		0x0101010101010100,
		0x0080808080808080,
		0x0040404040404040,
		0x0020202020202020,
		0x0010101010101010,
		0x0008080808080808,
		0x0,
		0x0,
	}
	//A1 H8 Diagonal Magic bitboard mask	
	DiagA1H8MagicMask = []uint64{
		0x0,
		0x0,
		0x0101010101010100,
		0x0101010101010100,
		0x0101010101010100,
		0x0101010101010100,
		0x0101010101010100,
		0x0101010101010100,
		0x8080808080808000,
		0x4040404040400000,
		0x2020202020000000,
		0x1010101000000000,
		0x0808080000000000,
		0x0,
		0x0,
	}
)

//InitialiseAllMasks initialises all masks to be used in the move generation
//process 	
func InitialiseAllMasks() {
	initialiseCharBitSet()
	initialiseRankMask()
	initialiseFileMask()
	initialiseDiagA8H1()
	initialiseDiagA1H8Magic()
	initialiseFileMagic()
	initialiseGeneralSlidingAttacks()
	initialiseAttackMasks()

	//Castling Masks
	//[0] = white [1] = black
	MaskEG[0] = BitSet[SquareMap["E1"]] |
		BitSet[SquareMap["F1"]] |
		BitSet[SquareMap["G1"]]
	MaskEG[1] = BitSet[SquareMap["E8"]] |
		BitSet[SquareMap["F8"]] |
		BitSet[SquareMap["G8"]]

	MaskFG[0] = BitSet[SquareMap["G1"]] |
		BitSet[SquareMap["F1"]]
	MaskFG[1] = BitSet[SquareMap["F8"]] |
		BitSet[SquareMap["G8"]]

	MaskBD[0] = BitSet[SquareMap["B1"]] |
		BitSet[SquareMap["C1"]] |
		BitSet[SquareMap["D1"]]
	MaskBD[1] = BitSet[SquareMap["B8"]] |
		BitSet[SquareMap["C8"]] |
		BitSet[SquareMap["D8"]]

	MaskEG[0] = BitSet[SquareMap["C1"]] |
		BitSet[SquareMap["D1"]] |
		BitSet[SquareMap["E1"]]
	MaskEG[1] = BitSet[SquareMap["C8"]] |
		BitSet[SquareMap["D8"]] |
		BitSet[SquareMap["E8"]]

	//Pre generate Castling Moves
	var move Move
	move.Clear()
	move.SetCapt(Empty)
	move.SetPiece(WhiteKing)
	move.SetProm(WhiteKing)
	move.SetFrom(SquareMap["E1"])
	move.SetTo(SquareMap["G1"])
	WhiteCastleOO = move
	move.SetTo(SquareMap["C1"])
	WhiteCastleOOO = move

	move.SetPiece(BlackKing)
	move.SetProm(BlackKing)
	move.SetFrom(SquareMap["E8"])
	move.SetTo(SquareMap["G8"])
	BlackCastleOO = move
	move.SetTo(SquareMap["C8"])
	BlackCastleOOO = move
}

//initialiseRankMask inisitalises the RankMask array which contains
//Masks for all ranks, indexed by squares on the board
func initialiseRankMask() {
	for file := 1; file < 9; file++ {
		for rank := 1; rank < 9; rank++ {
			RankMask[Index[file][rank]] =
				BitSet[Index[2][rank]] |
					BitSet[Index[3][rank]] |
					BitSet[Index[4][rank]] |
					BitSet[Index[5][rank]] |
					BitSet[Index[6][rank]] |
					BitSet[Index[7][rank]]
		}
	}
}

func initialiseCharBitSet() {
	CharBitSet[0] = 1
	for i := 1; i < 8; i++ {
		CharBitSet[i] = CharBitSet[i-1] << 1
	}
}

//initialiseFileMask initialises the FileMask array with masks 
//indexed by squares on the board
func initialiseFileMask() {
	for file := 1; file < 9; file++ {
		for rank := 1; rank < 9; rank++ {
			FileMask[Index[file][rank]] =
				BitSet[Index[file][2]] |
					BitSet[Index[file][3]] |
					BitSet[Index[file][4]] |
					BitSet[Index[file][5]] |
					BitSet[Index[file][6]] |
					BitSet[Index[file][7]]
		}
	}
}

//initialiseDiagA8H1 initialises the DiagA8H1 array with masks indexed
//by square
func initialiseDiagA8H1() {
	for file := 1; file < 9; file++ {
		for rank := 1; rank < 9; rank++ {
			diag := file + rank
			DiagA8H1Magic[Index[file][rank]] = DiagA8H1MagicMask[diag-2]

			DiagA8H1Mask[Index[file][rank]] = 0x0
			if diag < 10 {
				for sq := 2; sq < diag-1; sq++ {
					DiagA8H1Mask[Index[file][rank]] |= BitSet[Index[sq][diag-sq]]
				}
			} else {
				for sq := 2; sq < 17-diag; sq++ {
					DiagA8H1Mask[Index[file][rank]] |= BitSet[Index[diag+sq-9][9-sq]]
				}
			}
		}
	}
}

//initialiseDiagA8H1Magic initialises the DiagA8HMagic array with masks indexed
//by square
func initialiseDiagA1H8Magic() {
	for file := 1; file < 9; file++ {
		for rank := 1; rank < 9; rank++ {
			diag := file - rank
			DiagA1H8Magic[Index[file][rank]] = DiagA1H8MagicMask[diag+7]

			DiagA1H8Mask[Index[file][rank]] = 0x0
			if diag > -1 {
				for sq := 2; sq < 8-diag; sq++ {
					DiagA1H8Mask[Index[file][rank]] |= BitSet[Index[diag+sq][sq]]
				}
			} else {
				for sq := 2; sq < 8+diag; sq++ {
					DiagA1H8Mask[Index[file][rank]] |= BitSet[Index[sq][sq-diag]]
				}
			}
		}
	}
}

//initialiseFileMagic initialises the FileMagic array with masks
//indexed by square
func initialiseFileMagic() {
	for file := 1; file < 9; file++ {
		for rank := 1; rank < 9; rank++ {
			FileMagic[Index[file][rank]] = FileMagicMask[file-1]
		}
	}
}

//initialiseGeneralSlidingAttacks initialises the GeneralSlidingAttacks array
//indexed by square
func initialiseGeneralSlidingAttacks() {
	for sq := 0; sq < 8; sq++ {
		for state6 := 0; state6 < 64; state6++ {
			state8 := state6 << 1
			attBit := byte(0)
			if sq < 7 {
				attBit |= CharBitSet[sq+1]
			}
			slide := sq + 2
			for slide < 8 {
				if (int(CharBitSet[slide-1]) &^ state8) != 0 {
					attBit |= CharBitSet[slide]
				} else {
					break
				}
				slide++
			}

			if sq > 0 {
				attBit |= CharBitSet[sq-1]
			}
			slide = sq - 2
			for slide >= 0 {
				if (int(CharBitSet[slide+1]) &^ state8) != 0 {
					attBit |= CharBitSet[slide]
				} else {
					break
				}
				slide--
			}
			GeneralSlidingAttacks[sq][state6] = attBit
		}
	}
}

//inisitaliseAttackMasks initialises all attack masks	
func initialiseAttackMasks() {
	for sq := 0; sq < 64; sq++ {
		KnightAttacks[sq] = 0x0
		KingAttacks[sq] = 0x0
		WhitePawnAttacks[sq] = 0x0
		WhitePawnMoves[sq] = 0x0
		WhitePawnDoubleMoves[sq] = 0x0
		BlackPawnMoves[sq] = 0x0
		BlackPawnAttacks[sq] = 0x0
		BlackPawnDoubleMoves[sq] = 0x0

		initialiseWhitePawnMovesAndAttacks(sq)
		initialiseBlackPawnMovesAndAttacks(sq)
		initialiseKnightAttacks(sq)
		initialiseKingAttacks(sq)

		for state := 0; state < 64; state++ {
			RankAttacks[sq][state] = 0x0
			FileAttacks[sq][state] = 0x0
			DiagA8H1Attacks[sq][state] = 0x0
			DiagA1H8Attacks[sq][state] = 0x0

			initialiseRankAttacks(sq, state)
			initialiseFileAttacks(sq, state)
			initialiseA8H1Masks(sq, state)
			initialiseA1H8Masks(sq, state)
		}
	}
}

func initialiseWhitePawnMovesAndAttacks(sq int) {
	file := Files[sq]
	rank := Ranks[sq]

	afile := file - 1
	arank := rank + 1

	//attacks 
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		WhitePawnAttacks[sq] |= BitSet[Index[afile][arank]]
	}
	afile = file + 1
	arank = rank + 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		WhitePawnAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	//moves
	afile = file
	arank = rank + 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		WhitePawnMoves[sq] |= BitSet[Index[afile][arank]]
	}
	if rank == 2 {
		afile = file
		arank = rank + 2
		if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
			WhitePawnDoubleMoves[sq] |= BitSet[Index[afile][arank]]
		}
	}
}

func initialiseBlackPawnMovesAndAttacks(sq int) {
	file := Files[sq]
	rank := Ranks[sq]

	afile := file - 1
	arank := rank - 1

	//attacks 
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		BlackPawnAttacks[sq] |= BitSet[Index[afile][arank]]
	}
	afile = file + 1
	arank = rank - 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		BlackPawnAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	//moves
	afile = file
	arank = rank - 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		BlackPawnMoves[sq] |= BitSet[Index[afile][arank]]
	}
	if rank == 7 {
		afile = file
		arank = rank - 2
		if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
			BlackPawnDoubleMoves[sq] |= BitSet[Index[afile][arank]]
		}
	}
}

func initialiseKnightAttacks(sq int) {
	file := Files[sq]
	rank := Ranks[sq]

	afile := file - 2
	arank := rank + 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KnightAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file - 1
	arank = rank + 2
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KnightAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file + 1
	arank = rank + 2
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KnightAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file + 2
	arank = rank + 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KnightAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file + 2
	arank = rank - 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KnightAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file + 1
	arank = rank - 2
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KnightAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file - 1
	arank = rank - 2
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KnightAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file - 2
	arank = rank - 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KnightAttacks[sq] |= BitSet[Index[afile][arank]]
	}
}

func initialiseKingAttacks(sq int) {
	file := Files[sq]
	rank := Ranks[sq]

	afile := file - 1
	arank := rank
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KingAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file - 1
	arank = rank + 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KingAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file
	arank = rank + 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KingAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file + 1
	arank = rank + 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KingAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file + 1
	arank = rank
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KingAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file + 1
	arank = rank - 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KingAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file
	arank = rank - 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KingAttacks[sq] |= BitSet[Index[afile][arank]]
	}

	afile = file - 1
	arank = rank - 1
	if (afile >= 1) && (afile <= 8) && (arank >= 1) && (arank <= 8) {
		KingAttacks[sq] |= BitSet[Index[afile][arank]]
	}
}

func initialiseRankAttacks(sq, state int) {
	RankAttacks[sq][state] = 0
	RankAttacks[sq][state] |= uint64(GeneralSlidingAttacks[Files[sq]-1][state]) << uint(RankShift[sq]-1)
}

func initialiseFileAttacks(sq, state int) {
	FileAttacks[sq][state] = 0x0
	for attbit := 0; attbit < 8; attbit++ {
		if (GeneralSlidingAttacks[8-Ranks[sq]][state] & CharBitSet[attbit]) != 0 {
			FileAttacks[sq][state] |= BitSet[Index[Files[sq]][8-attbit]]
		}
	}
}

func initialiseA1H8Masks(sq, state int) {
	DiagA1H8Attacks[sq][state] = 0x0
	for attbit := 0; attbit < 8; attbit++ {
		pos := int(math.Min(float64(Files[sq]-1), float64(Ranks[sq]-1)))
		if (GeneralSlidingAttacks[pos][state] & CharBitSet[attbit]) != 0 {
			diag := Files[sq] - Ranks[sq]
			var file, rank int
			if diag < 0 {
				file = attbit + 1
				rank = file - diag
			} else {
				rank = attbit + 1
				file = diag + rank
			}
			if (file > 0) && (file < 9) && (rank > 0) && (rank < 9) {
				DiagA1H8Attacks[sq][state] |= BitSet[Index[file][rank]]
			}
		}
	}
}

func initialiseA8H1Masks(sq, state int) {
	DiagA8H1Attacks[sq][state] = 0x0
	for attbit := 0; attbit < 8; attbit++ {
		pos := int(math.Min(float64(8-Ranks[sq]), float64(Files[sq]-1)))
		if (GeneralSlidingAttacks[pos][state] & CharBitSet[attbit]) != 0 {
			diag := Files[sq] + Ranks[sq]
			var file, rank int
			if diag < 10 {
				file = attbit + 1
				rank = diag - file
			} else {
				rank = 8 - attbit
				file = diag - rank
			}
			if (file > 0) && (file < 9) && (rank > 0) && (rank < 9) {
				DiagA8H1Attacks[sq][state] |= BitSet[Index[file][rank]]
			}
		}
	}
}
