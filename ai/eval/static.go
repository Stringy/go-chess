package eval

import (
	"chess/ai/gen"
)

//StaticHeuristic is an evaluation method for static analysis on the game board
//the structure contains various information from the game board and is initialised
//upon each evaluation 
type StaticHeuristic struct {
	//number of each piece type
	wpawns,
	wrooks,
	wknights,
	wbishops,
	wqueens uint

	bpawns,
	brooks,
	bknights,
	bbishops,
	bqueens uint

	//king squares
	bkingsq, wkingsq uint

	//material for white and black
	wmat, bmat uint

	//total number of pieces
	wtotal, btotal uint

	//passed pawns
	wpassedpawns, bpassedpawns uint64
}

//Eval takes into account the following considerations during the 
//evaluation of the chess position
//    material value
//    king safety
//    bishop pair
//    +bishop val as num pawns decreases
//    -knight val as num pawns decreases
//    decrease score if double pawn
//    bonus for pawns close to promotion
// 
//    TODO:
//    lower queen value if moving early
//    bonus for mobility
//    bonus for threatening other pieces
//    bonus for protection of own peices
//    bonus for development speed of minor pieces in the opening
func (s *StaticHeuristic) eval(b *gen.Board) int {
	score := 0
	s.initFromBoard(b)

	if s.isDraw(b) {
		return score
	}

	score += s.materialValue()
	score += s.whiteEvaluation(b)
	score -= s.blackEvaluation(b)
	return score
}

//initFromBoard initialises the heurisitc structure with material values
//of the board
func (s *StaticHeuristic) initFromBoard(b *gen.Board) {
	s.wkingsq = uint(b.WKingSq)
	s.bkingsq = uint(b.BKingSq)

	s.wpawns = b.NumWPawns
	s.wrooks = b.NumWRooks
	s.wknights = b.NumWKnights
	s.wbishops = b.NumWBishops
	s.wqueens = b.NumWQueens
	s.wmat = uint(b.WMaterial)
	s.wtotal = s.wpawns + s.wrooks + s.wbishops + s.wknights + s.wqueens

	s.bpawns = b.NumBPawns
	s.brooks = b.NumBRooks
	s.bknights = b.NumBKnights
	s.bbishops = b.NumBBishops
	s.bqueens = b.NumBQueens
	s.bmat = uint(b.BMaterial)
	s.btotal = s.bpawns + s.brooks + s.bbishops + s.bknights + s.bqueens
}

//materialValue evaluates the material currently on the board
//it penalises loss of pieces in the late game and encourages
//taking of pieces for the current winner
func (s StaticHeuristic) materialValue() int {
	if s.wmat+s.wpawns > s.bmat+s.bpawns {
		return int(45 + 3*s.wtotal - 6*s.btotal)
	}
	return int(-45 - int(3*s.btotal+6*s.wtotal))
}

//isDraw checks the state for common draw positions
//    K vs. KN
//    K vs. KB
func (s StaticHeuristic) isDraw(b *gen.Board) bool {
	if s.wpawns == 0 || s.bpawns == 0 {
		//king vs. king
		if s.wtotal == 0 && s.btotal == 0 {
			return true
		}
		//king and knight vs. king
		wkandn := (s.wtotal == KnightValue && s.wknights == 1) && s.btotal == 0
		bkandn := (s.btotal == KnightValue && s.bknights == 1) && s.wtotal == 0
		if wkandn || bkandn {
			return true
		}

		//two kings and one or more bishops
		if s.wbishops+s.bbishops > 0 {
			onlywbish := (s.wknights == 0 && s.wrooks == 0 && s.wqueens == 0)
			onlybbish := (s.bknights == 0 && s.brooks == 0 && s.bqueens == 0)
			if onlywbish && onlybbish {
				n := ((b.WhiteBishops | b.BlackBishops) & WhiteSquares) == 0
				m := ((b.WhiteBishops | b.BlackBishops) & BlackSquares) == 0
				if n || m {
					return true
				}
			}
		}
	}
	return false
}

//blackEvaluation is a wrapper function for all black piece evaluation functions
func (s StaticHeuristic) blackEvaluation(b *gen.Board) int {
	score := s.evalBlackPawns(b) +
		s.evalBlackRooks(b) +
		s.evalBlackKnights(b) +
		s.evalBlackBishops(b) +
		s.evalBlackQueens(b)
	return score
}

//whiteEvaluation is a wrapper function for all white piece evaluation functions
func (s StaticHeuristic) whiteEvaluation(b *gen.Board) int {
	score := s.evalWhitePawns(b) +
		s.evalWhiteRooks(b) +
		s.evalWhiteKnights(b) +
		s.evalWhiteBishops(b) +
		s.evalWhiteQueens(b) +
		s.evalWhiteKing(b)
	return score
}

//evalWhitePawns evaluates the position and strength of the white pawns
//on the board
//including passed pawns, backwards pawns, and doubled pawns
func (s *StaticHeuristic) evalWhitePawns(b *gen.Board) int {
	twpawns := b.WhitePawns
	tscore := 0
	for twpawns != 0 {
		sq := gen.FirstOne(twpawns)
		//passed pawns
		tscore += WhitePawnTable[sq]
		tscore += PawnOpponentDistance[Distance[sq][s.bkingsq]]
		if (PassedWhite[sq] & b.BlackPawns) != 0 {
			tscore += PassedPawnBonus
			s.wpassedpawns ^= gen.BitSet[sq]
		}
		//doubled pawns
		if (b.WhitePawns^gen.BitSet[sq])&gen.FileMask[sq] == 0 {
			tscore -= DoublePawnPenalty
		}
		//isolated pawns TODO
		if IsolatedWhite[sq]&b.WhitePawns == 0 {
			tscore -= IsolatedPawnPenalty
		} else {
			if gen.WhitePawnAttacks[sq+8]&b.BlackPawns == 0 {
				if BackWhite[sq]&b.WhitePawns != 0 {
					tscore -= BackPawnPenalty
				}
			}
		}
		twpawns ^= gen.BitSet[sq]
	}
	return tscore
}

//evalWhiteKnights evaluates the position and strength of all white knights on the board
func (s StaticHeuristic) evalWhiteKnights(b *gen.Board) int {
	twknights := b.WhiteKnights
	tscore := s.evalGenericWhitePiece(twknights, KnightTable, KnightDistance)
	return tscore
}

//evalWhiteBishops evaluates the position and strength of all white bishops
//including a bonus for having the bishop pair
func (s StaticHeuristic) evalWhiteBishops(b *gen.Board) int {
	twbishops := b.WhiteBishops
	tscore := s.evalGenericWhitePiece(twbishops, BishopTable, BishopDistance)
	if twbishops != 0 {
		if gen.BitCount(twbishops) == 2 {
			tscore += BishopPairBonus
		}
	}
	return tscore
}

//evalWhiteRooks evaluates the position and value of all white rooks
//including a bonus for passed pawns
func (s StaticHeuristic) evalWhiteRooks(b *gen.Board) int {
	twrooks := b.WhiteRooks
	tscore := s.evalGenericWhitePiece(twrooks, RookTable, RookDistance)
	for twrooks != 0 {
		sq := gen.FirstOne(twrooks)
		if gen.FileMask[sq]&s.wpassedpawns != 0 {
			if sq < gen.LastOne(gen.FileMask[sq]&s.wpassedpawns) {
				tscore += RookPassedPawnBonus
			}
		}
		twrooks ^= gen.BitSet[sq]
	}
	return tscore
}

//evalWhiteQueens evaluates the position and strength of all white queens
//currently on the board
func (s StaticHeuristic) evalWhiteQueens(b *gen.Board) int {
	twqueens := b.WhiteQueens
	tscore := s.evalGenericWhitePiece(twqueens, QueenTable, QueenDistance)
	return tscore
}

//evalWhiteKing evaluates the position of the white king
//including PawnShield and endgame positioning
func (s StaticHeuristic) evalWhiteKing(b *gen.Board) int {
	return 0
}

//evalBlackPawns evaluates the position and strength of the black pawns
//on the board
//including passed pawns, backwards pawns, and doubled pawns
func (s *StaticHeuristic) evalBlackPawns(b *gen.Board) int {
	tbpawns := b.BlackPawns
	tscore := 0
	for tbpawns != 0 {
		sq := gen.FirstOne(tbpawns)
		//passed pawns
		tscore += BlackPawnTable[sq]
		tscore += PawnOpponentDistance[Distance[sq][s.wkingsq]]
		if (PassedBlack[sq] & b.WhitePawns) != 0 {
			tscore -= PassedPawnBonus
			s.bpassedpawns ^= gen.BitSet[sq]
		}
		//doubled pawns
		if (b.BlackPawns^gen.BitSet[sq])&gen.FileMask[sq] == 0 {
			tscore += DoublePawnPenalty
		}
		//isolated pawns TODO
		if IsolatedBlack[sq]&b.BlackPawns != 0 {
			tscore += IsolatedPawnPenalty
		} else {
			if gen.BlackPawnAttacks[sq-8]&b.WhitePawns == 0 {
				if BackBlack[sq]&b.BlackPawns != 0 {
					tscore += BackPawnPenalty
				}
			}
		}
		tbpawns ^= gen.BitSet[sq]
	}
	return tscore
}

//evalBlackKnights evaluates the position and strength of the black knights
func (s StaticHeuristic) evalBlackKnights(b *gen.Board) int {
	tbknights := b.BlackKnights
	tscore := s.evalGenericBlackPiece(tbknights, KnightTable, KnightDistance)
	return tscore
}

//evalBlackBishops evaluates the position and strength of the black bishops
//including bonus for having the bishop pair
func (s StaticHeuristic) evalBlackBishops(b *gen.Board) int {
	tbbishops := b.BlackBishops
	tscore := s.evalGenericBlackPiece(tbbishops, BishopTable, BishopDistance)
	if tbbishops != 0 {
		if gen.BitCount(tbbishops) == 2 {
			tscore += BishopPairBonus
		}
	}
	return tscore
}

//evalBlackRooks evaluates the position and strength of the black rooks
//including a bonus for the passed pawns
func (s StaticHeuristic) evalBlackRooks(b *gen.Board) int {
	tbrooks := b.BlackRooks
	tscore := s.evalGenericBlackPiece(tbrooks, RookTable, RookDistance)
	for tbrooks != 0 {
		sq := gen.FirstOne(tbrooks)
		if gen.FileMask[sq]&s.bpassedpawns != 0 {
			if sq < gen.LastOne(gen.FileMask[sq]&s.bpassedpawns) {
				tscore += RookPassedPawnBonus
			}
		}
		tbrooks ^= gen.BitSet[sq]
	}
	return tscore
}

//evalBlackQueens evaluates the position and strength of the black queens
func (s StaticHeuristic) evalBlackQueens(b *gen.Board) int {
	tbqueens := b.BlackQueens
	tscore := s.evalGenericBlackPiece(tbqueens, QueenTable, QueenDistance)
	return tscore
}

//evalBlackKing evaluates the position and value of the black king
//including endgame positioning and Pawn Shielding
func (s StaticHeuristic) evalBlackKing(b *gen.Board) int {
	return 0
}

//evalGenericWhitePieces provides evaluation for repeated functionality used for
//almost all pieces. 
//this includes basic position value, from pregenerated tables and
//distance from the black king
func (s StaticHeuristic) evalGenericWhitePiece(pieces uint64, values, distances []int) int {
	temp := pieces
	score := 0
	for temp != 0 {
		sq := gen.FirstOne(temp)
		score += values[sq]
		score += distances[Distance[sq][s.bkingsq]]
		temp ^= gen.BitSet[sq]
	}
	return score
}

//evalGenericBlackPieces provides evaluation for repeated functionality used for
//almost all pieces. 
//this includes basic position value, from pregenerated tables and
//distance from the white king
func (s StaticHeuristic) evalGenericBlackPiece(pieces uint64, values, distances []int) int {
	temp := pieces
	score := 0
	for temp != 0 {
		sq := gen.FirstOne(temp)
		score += values[Mirror[sq]]
		score += distances[Distance[sq][s.wkingsq]]
		temp ^= gen.BitSet[sq]
	}
	return score
}
 
 
 
 