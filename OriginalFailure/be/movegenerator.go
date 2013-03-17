package be

import (
	"math"
	"runtime"
	"fmt"
)

var (
	target uint64 //used as bitboard containing all squares 
	              //available for each piece eg ~BlackPieces
	count int
	totMoves []Move
)

const (
	averageNumMoves = 45
)

func InitMoveGen() {

}

func myappend(total []Move, n *int, moves...Move) []Move {
	if len(moves) > len(total) - *n {
		total = append(total, make([]Move, 0, (len(moves) * averageNumMoves))...)
		runtime.GC()
	}

	fmt.Println(*n)

	for _, move := range moves {
		total[*n] = move
		*n++
	}
	return total
}

func Generate(depth int, b *Board) []Move {
	numofmoves := func() int {
		return int(math.Pow(float64(averageNumMoves), float64(depth)))
	}

	totMoves = make([]Move, numofmoves())
	count = 0
	GenerateToDepth(depth, b)
	return totMoves
}

func GenerateToDepth(depth int, b *Board) {
	if depth == 0 {
		return
	}

	moves := GenerateAllMoves(b)
	totMoves = myappend(totMoves, &count, moves...)
	for _, move := range moves {
		b.MakeMove(&move)
		GenerateToDepth(depth - 1, b)
		b.UnmakeMove(&move)
	}
}

//Generate All Moves for a particular board
func GenerateAllMoves(b *Board) []Move {
	moves := make([]Move, 0)
//	n := 0
	if b.NextMove == BlackMove {
		target = Universe &^ b.BlackPieces //can't move to own pieces
		moves = append(moves, generateBlackPawnMoves(b)...)
		moves = append(moves, generateKnightMoves(b, BlackKnight, b.BlackKnights)...)
		moves = append(moves, generateSlidingMoves(b, BlackRook, b.BlackRooks, rookMoves)...)
		moves = append(moves, generateSlidingMoves(b, BlackBishop, b.BlackBishops, bishopMoves)...)
		moves = append(moves, generateSlidingMoves(b, BlackQueen, b.BlackQueens, queenMoves)...)
		moves = append(moves, generateBlackKingMoves(b)...)
	} else { //white move
		target = Universe &^ b.WhitePieces //can't move to own pieces
		moves = append(moves, generateWhitePawnMoves(b)...)
		moves = append(moves, generateKnightMoves(b, WhiteKnight, b.WhiteKnights)...)
		moves = append(moves, generateSlidingMoves(b, WhiteRook, b.WhiteRooks, rookMoves)...)
		moves = append(moves, generateSlidingMoves(b, WhiteBishop, b.WhiteBishops, bishopMoves)...)
		moves = append(moves, generateSlidingMoves(b, WhiteQueen, b.WhiteQueens, queenMoves)...)
		moves = append(moves, generateWhiteKingMoves(b)...)
	}
	return moves
}

func generateBlackPawnMoves(b *Board) []Move {
	movearr := make([]Move, 0)
	move := Move { 0 }
	pawns := b.BlackPawns
	free := Universe &^ b.OccupiedSquares	

	move.SetPiece(BlackPawn)
	for pawns != 0 { //for all pawns currently in play
		from := FirstOne(pawns) // get first pawn position (0-63)
		move.SetFrom(byte(from)) 
		moves := BlackPawnMoves[from] & free //normal moves
		if Ranks[from] == 7 && moves != 0 {
			moves |= BlackPawnDoubleMoves[from] & free //double moves
		}
		moves |= BlackPawnAttacks[from] & b.WhitePieces //captures
		for moves != 0 { //trim off all moves
			to := FirstOne(moves) // get first move
			move.SetTo(byte(to))
			move.SetCapt(b.Squares[to])
			if Ranks[to] == 1 { // promotion
				move.SetProm(BlackQueen)
				movearr = append(movearr, move)

				move.SetProm(BlackRook)
				movearr = append(movearr, move)
				
				move.SetProm(BlackBishop)
				movearr = append(movearr, move)
				
				move.SetProm(BlackKnight)
				movearr = append(movearr, move)
				
				move.SetProm(Empty)
			} else {
				movearr = append(movearr, move)
			}
			moves ^= BitSet[to] //clear bit for next loop
		}
		
		if b.Ep != 0 { //en-passant captures
			if BlackPawnAttacks[from] & BitSet[b.Ep] != 0 {
				move.SetProm(BlackPawn)
				move.SetCapt(WhitePawn)
				move.SetTo(byte(b.Ep))
				movearr = append(movearr, move)
			}
		}
		pawns ^= BitSet[from]
		move.SetProm(Empty) //reset promotion bits
	}
	return movearr
}

func generateWhitePawnMoves(b *Board) []Move {
	movearr := make([]Move, 0)
	move := Move { 0 }
	pawns := b.WhitePawns
	free := Universe &^ b.OccupiedSquares
	
	move.SetPiece(WhitePawn)
	for pawns != 0 { //for all pawns currently in play
		from := FirstOne(pawns) // get first pawn position (0-63)
		move.SetFrom(byte(from)) 
		moves := WhitePawnMoves[from] & free //normal moves
		if Ranks[from] == 2 && moves != 0 {
			moves |= WhitePawnDoubleMoves[from] & free //double moves
		}
		moves |= WhitePawnAttacks[from] & b.BlackPieces //captures
		for moves != 0 { //trim off all moves
			to := FirstOne(moves) // get first move
			move.SetTo(byte(to))
			move.SetCapt(b.Squares[to])
			if Ranks[to] == 1 { // promotion
				move.SetProm(WhiteQueen)
				movearr = append(movearr, move)

				move.SetProm(WhiteRook)
				movearr = append(movearr, move)
				
				move.SetProm(WhiteBishop)
				movearr = append(movearr, move)
				
				move.SetProm(WhiteKnight)
				movearr = append(movearr, move)
				
				move.SetProm(Empty)
			} else {
				movearr = append(movearr, move)
			}
			moves ^= BitSet[to] //clear bit for next loop
		}
		
		if b.Ep != 0 { //en-passant captures
			if WhitePawnAttacks[from] & BitSet[b.Ep] != 0 {
				move.SetProm(WhitePawn)
				move.SetCapt(BlackPawn)
				move.SetTo(byte(b.Ep))
				movearr = append(movearr, move)
			}
		}
		pawns ^= BitSet[from]
		move.SetProm(Empty) //reset promotion bits
	}
	return movearr
}


//Separate Generators for black and white kings
//due to castling special rule
func generateBlackKingMoves(b *Board) []Move {
	movearr := make([]Move, 0)
	move := Move { 0 }
	move.SetPiece(BlackKing)
	king := b.BlackKing
	
	for king != 0 {
		from := FirstOne(king)
		move.SetFrom(byte(from))
		moves := KingAttacks[from] & target
		for moves != 0 {
			to := FirstOne(moves)
			move.SetTo(byte(to))
			move.SetCapt(b.Squares[to])
			movearr = append(movearr, move)
			moves ^= BitSet[to]
		}

		//white 00 castling
		if (b.BlackCastle & CanCastleOO) != 0 {
			if (MaskFG[0] & b.OccupiedSquares) == 0 {
				if !isAttacked(b, MaskEG[BlackMove], WhiteMove) {
					move.Val = BlackCastleOO
					movearr = append(movearr, move) 
				}
			}
		}
		//white 000 castling
		if (b.BlackCastle & CanCastleOOO) != 0 {
			if (MaskBD[0] & b.OccupiedSquares) == 0 {
				if !isAttacked(b, MaskCE[BlackMove], WhiteMove) {
					move.Val = BlackCastleOOO
					movearr = append(movearr, move)
				}
			}
		}
		king ^= BitSet[from]
		move.SetProm(Empty)
	}	
	return movearr
}

func generateWhiteKingMoves(b *Board) []Move {
	movearr := make([]Move, 0)
	move := Move { 0 }
	move.SetPiece(WhiteKing)
	king := b.WhiteKing
	
	for king != 0 {
		from := FirstOne(king)
		move.SetFrom(byte(from))
		moves := KingAttacks[from] & target
		for moves != 0 {
			to := FirstOne(moves)
			move.SetTo(byte(to))
			move.SetCapt(b.Squares[to])
			movearr = append(movearr, move)
			moves ^= BitSet[to]
		}

		//white 00 castling
		if (b.WhiteCastle & CanCastleOO) != 0 {
			if (MaskFG[0] & b.OccupiedSquares) == 0 {
				if !isAttacked(b, MaskEG[WhiteMove], BlackMove) {
					move.Val = WhiteCastleOO
					movearr = append(movearr, move) 
				}
			}
		}
		//white 000 castling
		if (b.WhiteCastle & CanCastleOOO) != 0 {
			if (MaskBD[0] & b.OccupiedSquares) == 0 {
				if !isAttacked(b, MaskCE[WhiteMove], BlackMove) {
					move.Val = WhiteCastleOOO
					movearr = append(movearr, move)
				}
			}
		}
		king ^= BitSet[from]
		move.SetProm(Empty)
	}
	return movearr
}

func generateKnightMoves(b *Board, id byte, pieces uint64) []Move {
	movearr := make([]Move, 0)
	move := Move { 0 }
	move.SetPiece(id)
	knights := pieces
	
	for knights != 0 {
		from := FirstOne(knights)
		move.SetFrom(byte(from))
		moves := KnightAttacks[from] & target
		for moves != 0 {
			to := FirstOne(moves) 
			move.SetTo(byte(to))
			move.SetCapt(b.Squares[to])
			movearr = append(movearr, move)
			moves ^= BitSet[to]
		}
		knights ^= BitSet[from]
	}
	return movearr
}

func generateSlidingMoves(b *Board, id byte, pieces uint64, moveGen func(*Board, uint) uint64) []Move {
	movearr := make([]Move, 0)
	move := Move { 0 }
	move.SetPiece(id)
	tempPieces := pieces

	for tempPieces != 0 {
		from := FirstOne(tempPieces)
		move.SetFrom(byte(from))
		moves := moveGen(b, from)
		for moves != 0 {
			to := FirstOne(moves)
			move.SetTo(byte(to))
			move.SetCapt(b.Squares[to])
			movearr = append(movearr, move)
			moves ^= BitSet[to]
		}
		tempPieces ^= BitSet[from]
	}
	return movearr
}

func rankMoves(b *Board, from uint) uint64 {
	n := (b.OccupiedSquares & RankMask[from]) >> uint(RankShift[from])
	return RankAttacks[from][n] & target
}

func fileMoves(b *Board, from uint) uint64 {
	n := ((b.OccupiedSquares & FileMask[from]) * FileMagic[from]) >> 57
	return FileAttacks[from][n] & target
}

func diagA8H1Moves(b *Board, from uint) uint64 {
	n := ((b.OccupiedSquares & DiagA8H1Mask[from]) * DiagA8H1Magic[from]) >> 57
	return DiagA8H1Attacks[from][n]
}

func diagA1H8Moves(b *Board, from uint) uint64 {
	n := ((b.OccupiedSquares & DiagA1H8Mask[from]) * DiagA1H8Magic[from]) >> 57
	return DiagA1H8Attacks[from][n]
}

func bishopMoves(b *Board, from uint) uint64 {
	return diagA1H8Moves(b, from) | diagA8H1Moves(b, from)
}

func rookMoves(b *Board, from uint) uint64 {
	return rankMoves(b, from) | fileMoves(b, from)	
}

func queenMoves(b *Board, from uint) uint64 {
	return bishopMoves(b, from) | rookMoves(b, from)
}

func isAttacked(b *Board, bits uint64, attacker byte) bool {
	targ := bits
	if attacker == BlackMove {
		for targ != 0 {
			to := FirstOne(targ)
			if b.BlackPawns & WhitePawnAttacks[to] != 0 {
				return true
			}
			if b.BlackKnights & KnightAttacks[to] != 0 {
				return true
			}
			if  b.BlackKing & KingAttacks[to] != 0 {
				return true
			}

			slidingAtt := b.BlackQueens | b.BlackRooks
			if slidingAtt != 0 {
				i := (b.OccupiedSquares & RankMask[to]) >> uint(RankShift[to])
				if RankAttacks[to][i] & slidingAtt != 0 {
					return true
				}
				i = ((b.OccupiedSquares & FileMask[to]) * FileMagic[to]) >> 57
				if FileAttacks[to][i] & slidingAtt != 0 {
					return true
				}
			}

			slidingAtt = b.BlackQueens | b.BlackBishops
			if slidingAtt != 0 {
				i := ((b.OccupiedSquares & DiagA8H1Mask[to]) * DiagA8H1Magic[to]) >> 57
				if DiagA8H1Attacks[to][i] & slidingAtt != 0 {
					return true
				}
				i = ((b.OccupiedSquares & DiagA1H8Mask[to]) * DiagA1H8Magic[to]) >> 57
				if DiagA1H8Attacks[to][i] & slidingAtt != 0 {
					return true
				}
			}
			targ ^= BitSet[to]
		}
	} else { //White Attacker
		for targ != 0 {
			to := FirstOne(targ)
			if b.WhitePawns & WhitePawnAttacks[to] != 0 {
				return true
			}
			if b.WhiteKnights & KnightAttacks[to] != 0 {
				return true
			}
			if  b.WhiteKing & KingAttacks[to] != 0 {
				return true
			}

			slidingAtt := b.WhiteQueens | b.WhiteRooks
			if slidingAtt != 0 {
				i := (b.OccupiedSquares & RankMask[to]) >> uint(RankShift[to])
				if RankAttacks[to][i] & slidingAtt != 0 {
					return true
				}
				i = ((b.OccupiedSquares & FileMask[to]) * FileMagic[to]) >> 57
				if FileAttacks[to][i] & slidingAtt != 0 {
					return true
				}
			}

			slidingAtt = b.WhiteQueens | b.WhiteBishops
			if slidingAtt != 0 {
				i := ((b.OccupiedSquares & DiagA8H1Mask[to]) * DiagA8H1Magic[to]) >> 57
				if DiagA8H1Attacks[to][i] & slidingAtt != 0 {
					return true
				}
				i = ((b.OccupiedSquares & DiagA1H8Mask[to]) * DiagA1H8Magic[to]) >> 57
				if DiagA1H8Attacks[to][i] & slidingAtt != 0 {
					return true
				}
			}
			targ ^= BitSet[to]
		}
	}
	return false
}