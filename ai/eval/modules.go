package eval

import (
	"chess/ai/gen"
	"fmt"
)

type MaterialModule struct {
	score int
}

type KingModule struct {
	score int
}

type QueenModule struct {
	score int
}

type BishopModule struct {
	score int
}

type KnightModule struct {
	score int
}

type RookModule struct {
	score            int
	num_passed_white int
	num_passed_black int
}

type PawnModule struct {
	score           int
	isolated_white  int
	backwards_white int
	isolated_black  int
	backwards_black int
	double_white    int
	double_black    int
}

//MaterialModule.eval evaluates the current material value on the board
//quite simply
func (m *MaterialModule) eval(b *gen.Board) int {
	wtotal := b.NumWPawns +
		b.NumWRooks +
		b.NumWBishops +
		b.NumWKnights +
		b.NumWQueens
	btotal := b.NumBPawns +
		b.NumBRooks +
		b.NumBBishops +
		b.NumBKnights +
		b.NumBQueens

	if b.WMaterial+int(b.NumWPawns) > b.BMaterial+int(b.NumBPawns) {
		m.score = int(45 + 3*wtotal - 6*btotal)
	} else {
		m.score = int(-45 - int(3*btotal+6*wtotal))
	}
	return m.score
}

func (m *MaterialModule) debug() {
	fmt.Println("\t=== MATERIAL EVALUATION RESULTS ===")
	fmt.Println("Score:", m.score)
}

//PawnModule.eval evaluates the pawn situation for both sides and returns the
//cumulative score of both. i.e. A positive value will be the result of a 
//good white pawn state, negative for a better black position
func (p *PawnModule) eval(b *gen.Board) int {
	//white eval
	twpawns := b.WhitePawns

	p.isolated_black = 0
	p.isolated_white = 0
	p.backwards_white = 0
	p.backwards_black = 0
	p.double_black = 0
	p.double_white = 0

	score := 0

	for twpawns != 0 {
		sq := gen.LSB(twpawns)
		score += WhitePawnTable[sq]
		score += PawnOpponentDistance[Distance[sq][b.BKingSq]]
		score += numAttackedPieces(b.WhitePieces, gen.WhitePawnAttacks[sq]) * DefenseBonus //defending pieces
		if (PassedWhite[sq] & b.BlackPawns) != 0 {
			score += PassedPawnBonus
		}
		if (b.WhitePawns^gen.BitSet[sq])&gen.FileMask[sq] != 0 {
			p.double_white++
			score -= DoublePawnPenalty
		}
		if IsolatedWhite[sq]&b.WhitePawns == 0 {
			p.isolated_white++
			score -= IsolatedPawnPenalty
		} else {
			if (gen.WhitePawnAttacks[sq+8] & b.BlackPawns) != 0 {
				if BackWhite[sq]&b.WhitePawns != 0 {
					p.backwards_white++
					score -= BackPawnPenalty
				}
			}
		}
		twpawns ^= gen.BitSet[sq]
	}

	//black eval
	tbpawns := b.BlackPawns
	for tbpawns != 0 {
		sq := gen.LSB(tbpawns)
		score -= BlackPawnTable[sq]
		score -= PawnOpponentDistance[Distance[sq][b.WKingSq]]
		score -= numAttackedPieces(b.BlackPieces, gen.BlackPawnAttacks[sq]) * DefenseBonus
		if PassedBlack[sq]&b.WhitePawns != 0 {
			score -= PassedPawnBonus
		}
		if (b.BlackPawns^gen.BitSet[sq])&gen.FileMask[sq] != 0 {
			p.double_black++
			score += DoublePawnPenalty
		}
		if IsolatedBlack[sq]&b.BlackPawns == 0 {
			p.isolated_black++
			score += IsolatedPawnPenalty
		} else {
			if gen.BlackPawnAttacks[sq-8]&b.WhitePawns != 0 {
				if BackBlack[sq]&b.BlackPawns != 0 {
					p.backwards_black++
					score += BackPawnPenalty
				}
			}
		}
		tbpawns ^= gen.BitSet[sq]
	}
	p.score = score
	return score
}

func (p *PawnModule) debug() {
	fmt.Println("\t=== PAWN EVALUATION RESULTS ===")
	fmt.Println("Score =", p.score)
	fmt.Println("Number of isolated pawns:\n\tBlack:", p.isolated_black, "\n\tWhite:", p.isolated_white)
	fmt.Println("Number of backwards pawns:\n\tBlack:", p.backwards_black, "\n\tWhite:", p.backwards_white)
	fmt.Println("Number of double pawns:\n\tBlack:", p.double_black, "\n\tWhite:", p.double_white)
}

//RookModule.eval evaluates the state of the all rooks currently in play
//including passed pawn bonuses
func (r *RookModule) eval(b *gen.Board) int {
	twrooks := b.WhiteRooks

	r.num_passed_black = 0
	r.num_passed_white = 0

	score := 0

	bpassedpawns := uint64(0)
	wpassedpawns := uint64(0)
	twpawns := b.WhitePawns
	tbpawns := b.BlackPawns

	for twpawns != 0 {
		sq := gen.LSB(twpawns)
		if PassedWhite[sq]&b.BlackPawns == 0 {
			r.num_passed_white++
			wpassedpawns ^= gen.BitSet[sq]
		}
		twpawns ^= gen.BitSet[sq]
	}

	for tbpawns != 0 {
		sq := gen.LSB(tbpawns)
		if PassedBlack[sq]&b.WhitePawns == 0 {
			r.num_passed_black++
			bpassedpawns ^= gen.BitSet[sq]
		}
		tbpawns ^= gen.BitSet[sq]
	}

	for twrooks != 0 {
		sq := gen.LSB(twrooks)
		score += RookTable[sq]
		score += RookDistance[Distance[sq][b.BKingSq]]
		score += numAttackedPieces(b.BlackPieces, gen.RookMoves(b, sq)) * PinAndForkBonus
		score += numAttackedPieces(b.WhitePieces, gen.RookMoves(b, sq)) * DefenseBonus
		pawnsOnFile := gen.FileMask[sq] & wpassedpawns
		if pawnsOnFile != 0 {
			if sq < gen.MSB(pawnsOnFile) {
				score += RookPassedPawnBonus
			}
		}
		twrooks ^= gen.BitSet[sq]
	}

	tbrooks := b.BlackRooks
	for tbrooks != 0 {
		sq := gen.LSB(tbrooks)
		score -= RookTable[Mirror[sq]]
		score -= RookDistance[Distance[sq][b.WKingSq]]
		score -= numAttackedPieces(b.WhitePieces, gen.RookMoves(b, sq)) * PinAndForkBonus
		score -= numAttackedPieces(b.BlackPieces, gen.RookMoves(b, sq)) * DefenseBonus
		pawnsOnFile := gen.FileMask[sq] & bpassedpawns
		if pawnsOnFile != 0 {
			if sq < gen.MSB(pawnsOnFile) {
				score -= RookPassedPawnBonus
			}
		}
		tbrooks ^= gen.BitSet[sq]
	}
	r.score = score
	return score
}

func (r *RookModule) debug() {
	fmt.Println("\t=== ROOK EVALUATION RESULTS ===")
	fmt.Println("Score:", r.score)
	fmt.Println("Number of passed pawns:\n\tBlack", r.num_passed_black, "\n\tWhite", r.num_passed_white)
}

//BishopModule.eval evaluates the bishop position for either side and returns
//the cumulative score of both
func (bish *BishopModule) eval(b *gen.Board) int {
	twbishops := b.WhiteBishops
	score := 0

	if gen.BitCount(twbishops) == 2 {
		score += BishopPairBonus
	}
	for twbishops != 0 {
		sq := gen.LSB(twbishops)
		score += BishopTable[sq]
		score += BishopDistance[Distance[sq][b.BKingSq]]
		score += numAttackedPieces(b.BlackPieces, gen.BishopMoves(b, sq)) * PinAndForkBonus
		score += numAttackedPieces(b.WhitePieces, gen.BishopMoves(b, sq)) * DefenseBonus
		twbishops ^= gen.BitSet[sq]
	}

	tbbishops := b.BlackBishops
	if gen.BitCount(tbbishops) == 2 {
		score -= BishopPairBonus
	}
	for tbbishops != 0 {
		sq := gen.LSB(tbbishops)
		score -= BishopTable[Mirror[sq]]
		score -= BishopDistance[Distance[sq][b.WKingSq]]
		score -= numAttackedPieces(b.WhitePieces, gen.BishopMoves(b, sq)) * PinAndForkBonus
		score -= numAttackedPieces(b.BlackPieces, gen.BishopMoves(b, sq)) * DefenseBonus
		tbbishops ^= gen.BitSet[sq]
	}
	bish.score = score
	return score
}

func (b *BishopModule) debug() {
	fmt.Println("\t=== BISHOP EVALUATION RESULTS ===")
	fmt.Println("Score:", b.score)
}

//KnightModule.eval evaluates the knight position for both sides 
//and returns the cumulative score of both
func (k *KnightModule) eval(b *gen.Board) int {
	twknights := b.WhiteKnights
	score := 0
	for twknights != 0 {
		sq := gen.LSB(twknights)
		score += KnightTable[sq] + int(b.NumWPawns)
		score += KnightDistance[Distance[sq][b.BKingSq]]
		score += numAttackedPieces(b.BlackPieces, gen.KnightAttacks[sq]) * PinAndForkBonus
		score += numAttackedPieces(b.WhitePieces, gen.KnightAttacks[sq]) * DefenseBonus
		twknights ^= gen.BitSet[sq]
	}

	tbknights := b.BlackKnights
	for tbknights != 0 {
		sq := gen.LSB(tbknights)
		score -= KnightTable[Mirror[sq]] + int(b.NumBPawns)
		score -= KnightDistance[Distance[sq][b.BKingSq]]
		score -= numAttackedPieces(b.WhitePieces, gen.KnightAttacks[sq]) * PinAndForkBonus
		score -= numAttackedPieces(b.BlackPieces, gen.KnightAttacks[sq]) * DefenseBonus
		tbknights ^= gen.BitSet[sq]
	}
	k.score = score
	return score
}

func (k *KnightModule) debug() {
	fmt.Println("\t=== KNIGHT EVALUATION RESULTS ===")
	fmt.Println("Score:", k.score)
}

//QueenModule.eval evaluates the queen position for either side and 
//returns the cumulative score for both
func (q *QueenModule) eval(b *gen.Board) int {
	twqueens := b.WhiteQueens
	score := 0
	for twqueens != 0 {
		sq := gen.LSB(twqueens)
		score += QueenTable[sq]
		score += QueenDistance[Distance[sq][b.BKingSq]]
		mask := gen.QueenMoves(b, sq)
		score += numAttackedPieces(b.BlackPieces, mask) * PinAndForkBonus
		score += numAttackedPieces(b.WhitePieces, mask) * DefenseBonus
		twqueens ^= gen.BitSet[sq]
	}

	tbqueens := b.BlackQueens
	for tbqueens != 0 {
		sq := gen.LSB(tbqueens)
		score -= QueenTable[Mirror[sq]]
		score -= QueenDistance[Distance[sq][b.BKingSq]]
		mask := gen.QueenMoves(b, sq)
		score -= numAttackedPieces(b.WhitePieces, mask) * PinAndForkBonus
		score -= numAttackedPieces(b.BlackPieces, mask) * DefenseBonus
		tbqueens ^= gen.BitSet[sq]
	}
	q.score = score
	return score
}

func (q *QueenModule) debug() {
	fmt.Println("\t=== QUEEN EVALUATION RESULTS ===")
	fmt.Println("Score:", q.score)
}

func (k *KingModule) eval(b *gen.Board) int {
	score := 0
	if b.IsCheckMate() && b.NextMove == gen.WhiteMove {
		return -int(CheckMate)
	} else if b.IsCheckMate() && b.NextMove == gen.BlackMove {
		return int(CheckMate)
	}

	endgame := b.WMaterial < 150 || b.BMaterial < 150

	if endgame {
		score += KingEndGame[b.WKingSq]
		score -= KingEndGame[Mirror[b.BKingSq]]
	} else {
		score += KingTable[b.WKingSq]
		score -= KingTable[Mirror[b.BKingSq]]
	}

	shieldingWPawns := gen.BitCount(StrongShieldWhite[b.WKingSq] & b.WhitePawns)
	shieldingBPawns := gen.BitCount(StrongShieldBlack[b.BKingSq] & b.BlackPawns)
	score += int(shieldingWPawns * uint(StrongShieldBonus))
	score -= int(shieldingBPawns * uint(StrongShieldBonus))

	shieldingWPawns = gen.BitCount(WeakShieldWhite[b.WKingSq] & b.WhitePawns)
	shieldingBPawns = gen.BitCount(WeakShieldBlack[b.BKingSq] & b.BlackPawns)
	score += int(shieldingWPawns * uint(WeakShieldBonus))
	score -= int(shieldingBPawns * uint(WeakShieldBonus))

	k.score = score
	return score
}

func (k *KingModule) debug() {
	fmt.Println("\t=== KING EVALUATION RESULTS ===")
	fmt.Println("Score:", k.score)
}

//numAttackedPieces returns the number of pieces attacked on the target
//bitmap used for determining pins/forks and number of defended/attacked pieces
func numAttackedPieces(target uint64, mask uint64) int {
	return int(gen.BitCount(mask & target))
}
