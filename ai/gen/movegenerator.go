package gen

import (
//"fmt"
)

var (
	//used as bitboard containing all squares
	//available for each piece eg ~BlackPieces
	target uint64
)

//Generate All Moves takes a board state and generates all the valid 
//moves possible on the board, for the side whose turn it currently is.
//Averages between 20-25 moves per state
func GenerateAllMoves(b *Board) []Move {
	var moves []Move
	//	n := 0
	if b.NextMove == BlackMove {
		moves = generateAllBlackMoves(b)
	} else { //white move
		moves = generateAllWhiteMoves(b)
	}
	return moves
}

func generateAllWhiteMoves(b *Board) []Move {
	target = Universe &^ b.WhitePieces //can't move to own pieces
	moves := generateWhitePawnMoves(b)
	moves = append(moves, generateKnightMoves(b, WhiteKnight, b.WhiteKnights)...)
	moves = append(moves, generateSlidingMoves(b, WhiteRook, b.WhiteRooks, RookMoves)...)
	moves = append(moves, generateSlidingMoves(b, WhiteBishop, b.WhiteBishops, BishopMoves)...)
	moves = append(moves, generateSlidingMoves(b, WhiteQueen, b.WhiteQueens, QueenMoves)...)
	moves = append(moves, generateWhiteKingMoves(b)...)
	return moves
}

func generateAllBlackMoves(b *Board) []Move {
	target = Universe &^ b.BlackPieces //can't move to own pieces
	moves := generateBlackPawnMoves(b)
	moves = append(moves, generateKnightMoves(b, BlackKnight, b.BlackKnights)...)
	moves = append(moves, generateSlidingMoves(b, BlackRook, b.BlackRooks, RookMoves)...)
	moves = append(moves, generateSlidingMoves(b, BlackBishop, b.BlackBishops, BishopMoves)...)
	moves = append(moves, generateSlidingMoves(b, BlackQueen, b.BlackQueens, QueenMoves)...)
	moves = append(moves, generateBlackKingMoves(b)...)
	return moves
}

//generateBlackPawnMoves generates all moves for every black pawn currently on the board
//Includes special moves, enpassant and double movement
func generateBlackPawnMoves(b *Board) []Move {
	movearr := make([]Move, 0)
	move := Move(0)
	pawns := b.BlackPawns
	free := Universe &^ b.OccupiedSquares

	move.SetPiece(BlackPawn)
	for pawns != 0 { //for all pawns currently in play
		from := LSB(pawns) // get first pawn position (0-63)
		move.SetFrom(byte(from))
		moves := BlackPawnMoves[from] & free //normal moves
		if Ranks[from] == 7 && moves != 0 {  //on original square and not blocked
			moves |= BlackPawnDoubleMoves[from] & free //double moves
		}
		moves |= BlackPawnAttacks[from] & b.WhitePieces //captures
		for moves != 0 {                                //trim off all moves
			to := LSB(moves) // get first move
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
			if BlackPawnAttacks[from]&BitSet[b.Ep] != 0 {
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

//generateWhitePawnMoves generates all possible moves for the white pawns currently on the board
//Includes special rules: enpassant and double movement
func generateWhitePawnMoves(b *Board) []Move {
	movearr := make([]Move, 0)
	move := Move(0)
	pawns := b.WhitePawns
	free := Universe &^ b.OccupiedSquares

	move.SetPiece(WhitePawn)
	for pawns != 0 { //for all pawns currently in play
		from := LSB(pawns) // get first pawn position (0-63)
		move.SetFrom(byte(from))
		moves := WhitePawnMoves[from] & free //normal moves
		if Ranks[from] == 2 && moves != 0 {
			moves |= WhitePawnDoubleMoves[from] & free //double moves
		}
		moves |= WhitePawnAttacks[from] & b.BlackPieces //captures
		for moves != 0 {                                //trim off all moves
			to := LSB(moves) // get first move
			move.SetTo(byte(to))
			move.SetCapt(b.Squares[to])
			if Ranks[to] == 8 { // promotion
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
			if WhitePawnAttacks[from]&BitSet[b.Ep] != 0 {
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

//generateBlackKingMoves generates all possible moves for the black king
//TODO(gdth) Check for check and attacks around the king which prevent movement
func generateBlackKingMoves(b *Board) []Move {
	movearr := make([]Move, 0)
	move := Move(0)
	move.SetPiece(BlackKing)
	king := b.BlackKing

	for king != 0 {
		from := LSB(king)
		move.SetFrom(byte(from))
		moves := KingAttacks[from] & target
		for moves != 0 {
			to := LSB(moves)
			move.SetTo(byte(to))
			move.SetCapt(b.Squares[to])
			movearr = append(movearr, move)
			moves ^= BitSet[to]
		}

		//Black 00 castling
		if (b.BlackCastle & CanCastleOO) != 0 {
			if (MaskFG[BlackMove] & b.OccupiedSquares) == 0 {
				if !b.isAttacked(MaskEG[BlackMove], WhiteMove) {
					move = BlackCastleOO
					movearr = append(movearr, move)
				}
			}
		}
		//Black 000 castling
		if (b.BlackCastle & CanCastleOOO) != 0 {
			if (MaskBD[BlackMove] & b.OccupiedSquares) == 0 {
				if !b.isAttacked(MaskCE[BlackMove], WhiteMove) {
					move = BlackCastleOOO
					movearr = append(movearr, move)
				}
			}
		}
		king ^= BitSet[from]
		move.SetProm(Empty)
	}
	return movearr
}

//generateWhiteKingMoves generates all possible moves for the white king
//TODO(gdth) Check for check and attacks around the king which could impede movement
func generateWhiteKingMoves(b *Board) []Move {
	movearr := make([]Move, 0)
	move := Move(0)
	move.SetPiece(WhiteKing)
	king := b.WhiteKing

	for king != 0 {
		from := LSB(king)
		move.SetFrom(byte(from))
		moves := KingAttacks[from] & target
		for moves != 0 {
			to := LSB(moves)
			move.SetTo(byte(to))
			move.SetCapt(b.Squares[to])
			movearr = append(movearr, move)
			moves ^= BitSet[to]
		}

		//white 00 castling
		if (b.WhiteCastle & CanCastleOO) != 0 {
			if (MaskFG[WhiteMove] & b.OccupiedSquares) == 0 {
				if !b.isAttacked(MaskEG[WhiteMove], BlackMove) {
					move = WhiteCastleOO
					movearr = append(movearr, move)
				}
			}
		}
		//white 000 castling
		if (b.WhiteCastle & CanCastleOOO) != 0 {
			if (MaskBD[WhiteMove] & b.OccupiedSquares) == 0 {
				if !b.isAttacked(MaskCE[WhiteMove], BlackMove) {
					move = WhiteCastleOOO
					movearr = append(movearr, move)
				}
			}
		}
		king ^= BitSet[from]
		move.SetProm(Empty)
	}
	return movearr
}

//generateKnightMoves generates all possible moves for each knight of a particular colour on the board
//id is the unique id of a white knight, for the move object and pieces is the bitboard containing the
//positions of a particular colour's knights
func generateKnightMoves(b *Board, id byte, pieces uint64) []Move {
	movearr := make([]Move, 0)
	move := Move(0)
	move.SetPiece(id)
	knights := pieces

	for knights != 0 {
		from := LSB(knights)
		move.SetFrom(byte(from))
		moves := KnightAttacks[from] & target
		for moves != 0 {
			to := LSB(moves)
			move.SetTo(byte(to))
			move.SetCapt(b.Squares[to])
			movearr = append(movearr, move)
			moves ^= BitSet[to]
		}
		knights ^= BitSet[from]
	}
	return movearr
}

//generateSlidingMoves is a generic function for calculating all possible moves for any colour's rooks, bishops,
//or queens. It takes the id of the piece type along with the associated bitboard and a function to generate a move mask
//from. 
func generateSlidingMoves(b *Board, id byte, pieces uint64, moveGen func(*Board, uint) uint64) []Move {
	movearr := make([]Move, 0)
	move := Move(0)
	move.SetPiece(id)
	tempPieces := pieces

	for tempPieces != 0 {
		from := LSB(tempPieces)
		move.SetFrom(byte(from))
		moves := moveGen(b, from)
		for moves != 0 {
			to := LSB(moves)
			move.SetTo(byte(to))
			move.SetCapt(b.Squares[to])
			movearr = append(movearr, move)
			moves ^= BitSet[to]
		}
		tempPieces ^= BitSet[from]
	}
	return movearr
}

//rankMoves generates all moves along a particular rank on the board.
//It Uses pre-generated rank attack masks
func rankMoves(b *Board, from uint) uint64 {
	n := ((b.OccupiedSquares & RankMask[from]) >> uint(RankShift[from]))
	return RankAttacks[from][n] & target
}

//fileMoves generates all possible moves on a particular file of the board.
//It Uses pre-generated file attack masks and magic bitboards
func fileMoves(b *Board, from uint) uint64 {
	n := ((b.OccupiedSquares & FileMask[from]) * FileMagic[from]) >> 57
	return FileAttacks[from][n] & target
}

//diagA8H1Moves generates all possible moves on the A8H1 diagonal.
//It uses pre-generated diagonal magic bitboards
func diagA8H1Moves(b *Board, from uint) uint64 {
	n := ((b.OccupiedSquares & DiagA8H1Mask[from]) * DiagA8H1Magic[from]) >> 57
	return DiagA8H1Attacks[from][n] & target
}

//diagA1H8Moves generates all possible moves on the A1H8 diagonal.
//It uses pre-generated diagonal magic bitboards
func diagA1H8Moves(b *Board, from uint) uint64 {
	n := ((b.OccupiedSquares & DiagA1H8Mask[from]) * DiagA1H8Magic[from]) >> 57
	return DiagA1H8Attacks[from][n] & target
}

//bishopMoves generates all possible bishop moves from a square  on the board.
//It combines the moves generated from both the A1H8 diagonal and the A8H1 diagonal
func BishopMoves(b *Board, from uint) uint64 {
	return diagA1H8Moves(b, from) | diagA8H1Moves(b, from)
}

//rookMoves generates all possible rook moves from a square on the board. 
//It combines both the rank moves and the file moves 
func RookMoves(b *Board, from uint) uint64 {
	return rankMoves(b, from) | fileMoves(b, from)
}

//queenMoves generates all possible queen moves from a square on the board.
//It combines moves generated for a bishop and a rook on the same square
func QueenMoves(b *Board, from uint) uint64 {
	return BishopMoves(b, from) | RookMoves(b, from)
}

func findCheckBlocking(b *Board, moves []Move) []Move {
	blockers := make([]Move, 0)
	for _, move := range moves {
		b.MakeMove(&move)
		if !b.IsCheck() {
			blockers = append(blockers, move)
		}
		b.UnmakeMove(&move)
	}
	return blockers
}
