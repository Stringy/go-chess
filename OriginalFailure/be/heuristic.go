package be

import (

)

var (
	//number of white pieces
	wpawns,
		wrooks,
		wknights,
		wbishops,
		wqueens uint

	wpassedpawns int

	//number of black pieces
	bpawns,
		brooks,
		bknights, 
		bbishops, 
		bqueens uint

	bpassedpawns int
	//king position
	wkingsq, 
		bkingsq uint
	
	//material values
	whitemat,
		blackmat uint
	
	//totals
	wtotal, 
		btotal uint

	endgame bool

	score int
)

const (
	EndGameLim = 150
)

func EvalPositionValue(b *Board) int {
	score = b.Material

	//init king positions
	wkingsq = FirstOne(b.WhiteKing)
	bkingsq = FirstOne(b.BlackKing)

	//init white material values
	wpawns = BitCount(b.WhitePawns)
	wrooks = BitCount(b.WhiteRooks)
	wknights = BitCount(b.WhiteKnights)
	wbishops = BitCount(b.WhiteBishops)
	wqueens = BitCount(b.WhiteQueens)
	whitemat = (RookValue * wrooks) +
		(KnightValue * wknights) + 
		(BishopValue * wbishops) +
		(QueenValue * wqueens)
	wtotal = wpawns + wrooks + wbishops + wknights + wqueens

	//init black material values
	bpawns = BitCount(b.WhitePawns)
	brooks = BitCount(b.WhiteRooks)
	bknights = BitCount(b.WhiteKnights)
	bbishops = BitCount(b.WhiteBishops)
	bqueens = BitCount(b.WhiteQueens)
	blackmat = (RookValue * brooks) +
		(KnightValue * bknights) + 
		(BishopValue * bbishops) +
		(QueenValue * bqueens)
	btotal = bpawns + brooks + bbishops + bknights + bqueens

	//endgame defined by material on the board
	endgame = whitemat < EndGameLim || blackmat < EndGameLim

	if evalDraw(b) {
		return 0
	}

	return score
}

func evalDraw(b *Board) bool {
	if wpawns == 0 || bpawns == 0 {		
		//king vs king 
		if wtotal == 0 && btotal == 0 {
			return true
		}
		
		//king and knight vs. king 
		wkandn := (wtotal == KnightValue && wknights == 1) && btotal == 0
		bkandn := (btotal == KnightValue && bknights == 1) && wtotal == 0
		if wkandn || bkandn {
			return true
		}

		//two kings with one or more bishops 
		if wbishops + bbishops > 0 {
			onlywbish := (wknights == 0 && wrooks == 0 && wqueens == 0)
			onlybbish := (bknights == 0 && brooks == 0 && bqueens == 0)
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

func evalMaterial(b *Board) {
	if wtotal + wpawns > btotal + bpawns {
		score += int(45 + 3 * wtotal - 6 * btotal)
	} else {
		score += int(-45 - int(3 * btotal + 6 * wtotal))
	}
}

func evalWhitePieces(b *Board) {

}

func evalWhitePawns(b *Board) {
	wpassedpawns = 0
	temp := b.WhitePawns
	for ; temp != 0; {
		sq := FirstOne(temp)
		score += PawnPosW[sq]
		score += PawnOppDist[Dist[sq][bkingsq]]
		
		if endgame != 0 {
			score += PawnOppDist[Dist[sq][wkingsq]]
		}
		if (PassedWhite[sq] & b.BlackPawns) == 0 {
			score += BonusPassedPawn
			wpassedpawns ^= BitSet[sq]
		}
	}
}
