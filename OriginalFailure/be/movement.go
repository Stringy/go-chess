package be

func (b *Board) MakeMove(move *Move) {
	piece := move.GetPiece()
	
	next := func() {
		if b.NextMove == WhiteMove {
			b.NextMove = BlackMove
		} else {
			b.NextMove = WhiteMove
		}
	}
	
	switch piece {
		//white moves
	case  WhitePawn: b.makeWhitePawnMove(move)
	case WhiteRook: b.makeWhiteRookMove(move)
	case WhiteKnight: b.makeGenericMove(move, &b.WhiteKnights, &b.WhitePieces)
	case WhiteBishop: b.makeGenericMove(move, &b.WhiteBishops, &b.WhitePieces)
	case WhiteQueen: b.makeGenericMove(move, &b.WhiteQueens, &b.WhitePieces)
	case WhiteKing: b.makeWhiteKingMove(move)
		//black moves
	case BlackPawn: b.makeBlackPawnMove(move)
	case BlackRook: b.makeBlackRookMove(move)
	case BlackKnight: b.makeGenericMove(move, &b.BlackKnights, &b.BlackPieces)
	case BlackBishop: b.makeGenericMove(move, &b.BlackBishops, &b.BlackPieces)
	case BlackQueen: b.makeGenericMove(move, &b.BlackQueens, &b.BlackPieces)
	case BlackKing: b.makeBlackKingMove(move)
	}
	next()
}

func (b *Board) UnmakeMove(move *Move) {
	piece := move.GetPiece() 

	next := func() {
		if b.NextMove == WhiteMove {
			b.NextMove = BlackMove
		} else {
			b.NextMove = WhiteMove
		}		
	} 
	
	switch piece {
	case WhitePawn:
		b.unmakeWhitePawnMove(move)
	case WhiteRook:
		b.unmakeWhiteRookMove(move)
	case WhiteKnight:
		b.unmakeGenericMove(move, &b.WhiteKnights, &b.WhitePieces)
	case WhiteBishop:
		b.unmakeGenericMove(move, &b.WhiteBishops, &b.WhitePieces)
	case WhiteQueen:
		b.unmakeGenericMove(move, &b.WhiteQueens, &b.WhitePieces)
	case WhiteKing:
		b.unmakeWhiteKingMove(move)
		//Black unmake moves
	case BlackPawn: 
		b.unmakeBlackPawnMove(move)
	case BlackRook: 
		b.unmakeBlackRookMove(move)
	case BlackKnight:
		b.unmakeGenericMove(move, &b.BlackKnights, &b.BlackPieces)
	case BlackBishop:
		b.unmakeGenericMove(move, &b.BlackBishops, &b.BlackPieces)
	case BlackQueen: 
		b.unmakeGenericMove(move, &b.BlackQueens, &b.BlackPieces)
	case BlackKing: 
		b.unmakeBlackKingMove(move)
	}
	next()
}

func (b *Board) makeCapture(capt, to byte) {
	movebb := BitSet[to]
	switch capt {
	case WhitePawn:
		b.WhitePawns ^= movebb
		b.WhitePieces ^= movebb
		b.Material -= PawnValue
	case WhiteRook:
		b.WhiteRooks ^= movebb
		b.WhitePieces ^= movebb
		b.Material -= RookValue
		if to == A1 {
			b.WhiteCastle = b.WhiteCastle &^ CanCastleOOO
		} else if to == H1 {
			b.WhiteCastle = b.WhiteCastle &^ CanCastleOO
		}
	case WhiteKnight:
		b.WhiteKnights ^= movebb
		b.WhitePieces ^= movebb
		b.Material -= KnightValue
	case WhiteBishop:
		b.WhiteBishops ^= movebb
		b.WhitePieces ^= movebb
		b.Material -= BishopValue
	case WhiteQueen:
		b.WhiteQueens ^= movebb
		b.WhitePieces ^= movebb
		b.Material -= QueenValue
	case WhiteKing:
		b.WhiteKing ^= movebb
		b.WhitePieces ^= movebb

		//Black Captures
	case BlackPawn:
		b.BlackPawns ^= movebb
		b.BlackPieces ^= movebb
		b.Material += PawnValue
	case BlackRook:
		b.BlackRooks ^= movebb
		b.BlackPieces ^= movebb
		b.Material += RookValue
		if to == A8 {
			b.BlackCastle = b.BlackCastle &^ CanCastleOOO
		} else if to == H8 {
			b.BlackCastle = b.BlackCastle &^ CanCastleOO
		}
	case BlackKnight:
		b.BlackKnights ^= movebb
		b.BlackPieces ^= movebb
		b.Material += KnightValue
	case BlackBishop:
		b.BlackBishops ^= movebb
		b.BlackPieces ^= movebb
		b.Material += BishopValue
	case BlackQueen:
		b.BlackQueens ^= movebb
		b.BlackPieces ^= movebb
		b.Material += QueenValue
	case BlackKing:
		b.BlackKing ^= movebb
		b.BlackPieces ^= movebb
	}
	b.NumMoves = 0
}

func (b *Board) unmakeCapture(capt, to byte) {
	movebb := BitSet[to]
	switch capt {
	case WhitePawn:
		b.WhitePawns ^= movebb
		b.WhitePieces ^= movebb
		b.Material += PawnValue
		b.Squares[to] = WhitePawn
	case WhiteRook:
		b.WhiteRooks ^= movebb
		b.WhitePieces ^= movebb
		b.Material += RookValue
		b.Squares[to] = WhiteRook
	case WhiteKnight:
		b.WhiteKnights ^= movebb
		b.WhitePieces ^= movebb
		b.Material += KnightValue
		b.Squares[to] = WhiteKnight
	case WhiteBishop:
		b.WhiteBishops ^= movebb
		b.WhitePieces ^= movebb
		b.Material += BishopValue
		b.Squares[to] = WhiteBishop
	case WhiteQueen:
		b.WhiteQueens ^= movebb
		b.WhitePieces ^= movebb
		b.Material += QueenValue
		b.Squares[to] = WhiteQueen
	case WhiteKing:
		b.WhiteKing ^= movebb
		b.WhitePieces ^= movebb
		b.Squares[to] = WhiteKing

		//Black Captures
	case BlackPawn:
		b.BlackPawns ^= movebb
		b.BlackPieces ^= movebb
		b.Material -= PawnValue
		b.Squares[to] = BlackPawn
	case BlackRook:
		b.BlackRooks ^= movebb
		b.BlackPieces ^= movebb
		b.Material -= RookValue
		b.Squares[to] = BlackRook
	case BlackKnight:
		b.BlackKnights ^= movebb
		b.BlackPieces ^= movebb
		b.Material -= KnightValue
		b.Squares[to] = BlackKnight
	case BlackBishop:
		b.BlackBishops ^= movebb
		b.BlackPieces ^= movebb
		b.Material -= BishopValue
		b.Squares[to] = BlackBishop
	case BlackQueen:
		b.BlackQueens ^= movebb
		b.BlackPieces ^= movebb
		b.Material -= QueenValue
		b.Squares[to] = BlackQueen
	case BlackKing:
		b.BlackKing ^= movebb
		b.BlackPieces ^= movebb
		b.Squares[to] = BlackKing
	}
}


/*
 Can Be used for non-special movement types
 e.g. Queens, Knights and Bishops

 They have no special movment rules
 This is also colour independant
*/
func (b *Board) makeGenericMove(move *Move, pieces *uint64, colour *uint64) {
	from := move.GetFrom()
	to := move.GetTo() 
	piece := move.GetPiece()
	capt := move.GetCapt()
	
	movebb := BitSet[from] | BitSet[to]
	
	*pieces ^= movebb
	*colour ^= movebb
	b.Squares[from] = Empty
	b.Squares[to] = piece
	b.Ep = 0
	b.NumMoves++
	
	if capt != 0 {
		b.makeCapture(capt, to)
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}
}

//Unmake generic move,
//same rules as above
func (b *Board) unmakeGenericMove(move *Move, pieces *uint64, colour *uint64) {
	from := move.GetFrom()
	to := move.GetTo()
	piece := move.GetPiece()
	capt := move.GetCapt()
	
	movebb := BitSet[from] | BitSet[to]
	*pieces ^= movebb
	*colour ^= movebb
	b.Squares[from] = piece
	b.Squares[to] = Empty
	
	if capt != 0 {
		b.unmakeCapture(capt, to)
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}
}

/***************************************************/
/*         White Piece Movement / Unmove           */
/*             Pawns, Rooks and King               */
/***************************************************/
func (b *Board) makeWhitePawnMove(move *Move) {
	from := move.GetFrom()
	to := move.GetTo()
	piece := move.GetPiece()
	capt := move.GetCapt()

	frombb := BitSet[from]
	movebb := frombb | BitSet[to]

	b.WhitePawns ^= movebb
	b.WhitePieces ^= movebb
	b.Squares[from] = Empty
	b.Squares[to] = piece
	b.Ep = 0
	b.NumMoves = 0

	if Ranks[from] == 2 && Ranks[to] == 4 {
		b.Ep = int(from) + 8
	}

	if capt != 0 {
		if move.IsEnpassant() {
			b.BlackPawns ^= BitSet[to - 8]
			b.BlackPieces ^= BitSet[to - 8]
			b.OccupiedSquares ^= movebb | BitSet[to - 8]
			b.Squares[to - 8] = Empty
			b.Material += PawnValue
		} else {
			b.makeCapture(capt, to)
		}
	} else {
		b.OccupiedSquares ^= movebb
	}

	if move.IsPromotion() {
		b.makeWhiteProm(move.GetProm(), to)
		b.Squares[to] = move.GetProm()
	}
}

func (b *Board) makeWhiteRookMove(move *Move) {
	from := move.GetFrom()
	to := move.GetTo()
	piece := move.GetPiece()
	capt := move.GetCapt()

	frombb := BitSet[from]
	movebb := frombb | BitSet[to]

	b.WhiteRooks ^= movebb
	b.WhitePieces ^= movebb
	b.Squares[from] = Empty
	b.Squares[to] = piece
	b.Ep = 0
	b.NumMoves++
	
	if from == A1 {
		b.WhiteCastle = b.WhiteCastle &^ CanCastleOOO
	}
	if from == H1 {
		b.WhiteCastle = b.WhiteCastle &^ CanCastleOO
	}

	if capt != 0 {
		b.makeCapture(capt, to) 
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}
}

func (b *Board) makeWhiteKingMove(move *Move) {
	from := move.GetFrom()
	to := move.GetTo()
	piece := move.GetPiece()
	capt := move.GetCapt()

	frombb := BitSet[from]
	movebb := frombb | BitSet[to]

	b.WhiteKing ^= movebb
	b.WhitePieces ^= movebb
	b.Squares[from] = Empty
	b.Squares[to] = piece
	b.Ep = 0
	b.NumMoves++
	
	if capt != 0 {
		b.makeCapture(capt, to)
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}
	
	if move.IsCastle() {
		if move.IsCastleOO() {
			b.WhiteRooks ^= BitSet[H1] | BitSet[F1]
			b.WhitePieces ^= BitSet[H1] | BitSet[F1]
			b.OccupiedSquares ^= BitSet[H1] | BitSet[F1]
			b.Squares[H1] = Empty
			b.Squares[F1] = WhiteRook
		} else {
			b.WhiteRooks ^= BitSet[A1] | BitSet[D1]
			b.WhitePieces ^= BitSet[A1] | BitSet[D1]
			b.OccupiedSquares ^= BitSet[A1] | BitSet[D1]
			b.Squares[A1] = Empty
			b.Squares[D1] = WhiteRook
		}
	}
}

func (b *Board) unmakeWhitePawnMove(move *Move) {
	to := move.GetTo()
	from := move.GetFrom()
	capt := move.GetCapt()

	movebb := BitSet[from] | BitSet[to]

	b.WhitePawns ^= movebb
	b.WhitePieces ^= movebb
	b.Squares[from] = WhitePawn
	b.Squares[to] = Empty

	if capt != 0 {
		if move.IsEnpassant() {
			b.BlackPawns ^= BitSet[to - 8]
			b.BlackPieces ^= BitSet[to - 8]
			b.OccupiedSquares ^= movebb | BitSet[to - 8]
			b.Squares[to - 8] = BlackPawn
			b.Material -= PawnValue
		} else {
			b.unmakeCapture(capt, to)
			b.OccupiedSquares ^= movebb
		}
	} else {
		b.OccupiedSquares ^= movebb
	}
	if move.IsPromotion() {
		b.unmakeWhiteProm(move.GetProm(), to)
	}
}

func (b *Board) unmakeWhiteRookMove(move *Move) {
	from := move.GetFrom()
	to := move.GetTo()
	piece := move.GetPiece()
	capt := move.GetCapt()

	frombb := BitSet[from]
	movebb := frombb | BitSet[to]

	b.WhiteRooks ^= movebb
	b.WhitePieces ^= movebb
	b.Squares[from] = piece
	b.Squares[to] = Empty
	
	if capt != 0 {
		b.unmakeCapture(capt, to) 
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}
}

func (b *Board) unmakeWhiteKingMove(move *Move) {
	from := move.GetFrom()
	to := move.GetTo()
	piece := move.GetPiece()
	capt := move.GetCapt()

	frombb := BitSet[from]
	movebb := frombb | BitSet[to]

	b.WhiteKing ^= movebb
	b.WhitePieces ^= movebb
	b.Squares[from] = piece
	b.Squares[to] = Empty
	
	if capt != 0 {
		b.unmakeCapture(capt, to)
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}
	
	if move.IsCastle() {
		if move.IsCastleOO() {
			b.WhiteRooks ^= BitSet[H1] | BitSet[F1]
			b.WhitePieces ^= BitSet[H1] | BitSet[F1]
			b.OccupiedSquares ^= BitSet[H1] | BitSet[F1]
			b.Squares[H1] = WhiteRook
			b.Squares[F1] = Empty
		} else {
			b.WhiteRooks ^= BitSet[A1] | BitSet[D1]
			b.WhitePieces ^= BitSet[A1] | BitSet[D1]
			b.OccupiedSquares ^= BitSet[A1] | BitSet[D1]
			b.Squares[A1] = WhiteRook
			b.Squares[D1] = Empty
		}
	}	
}

func (b *Board) makeWhiteProm(piece, to byte) {
	movebb := BitSet[to]
	
	b.WhitePawns ^= movebb //clear pawn bit
	b.Material -= PawnValue

	switch piece {
	case WhiteQueen: 
		b.WhiteQueens ^= movebb
		b.Material += QueenValue
	case WhiteRook:
		b.WhiteRooks ^= movebb
		b.Material += RookValue
	case WhiteBishop:
		b.WhiteBishops ^= movebb
		b.Material += BishopValue
	case WhiteKnight:
		b.WhiteKnights ^= movebb
		b.Material += KnightValue
	}
}

func (b *Board) unmakeWhiteProm(piece, to byte) {
	movebb := BitSet[to]
	
	b.WhitePawns ^= movebb //clear pawn bit
	b.Material += PawnValue

	switch piece {
	case WhiteQueen: 
		b.WhiteQueens ^= movebb
		b.Material -= QueenValue
	case WhiteRook:
		b.WhiteRooks ^= movebb
		b.Material -= RookValue
	case WhiteBishop:
		b.WhiteBishops ^= movebb
		b.Material -= BishopValue
	case WhiteKnight:
		b.WhiteKnights ^= movebb
		b.Material -= KnightValue
	}	
}

/***************************************************/
/*             Black Piece Movement / Unmove       */
/***************************************************/
func (b *Board) makeBlackPawnMove(move *Move) {
	from := move.GetFrom()
	to := move.GetTo()
	piece := move.GetPiece()
	capt := move.GetCapt()

	frombb := BitSet[from]
	movebb := frombb | BitSet[to]

	b.BlackPawns ^= movebb
	b.BlackPieces ^= movebb
	b.Squares[from] = Empty
	b.Squares[to] = piece
	b.Ep = 0
	b.NumMoves = 0

	if Ranks[from] == 7 && Ranks[to] == 5 {
		b.Ep = int(from) - 8
	}

	if capt != 0 {
		if move.IsEnpassant() {
			b.BlackPawns ^= BitSet[to + 8]
			b.BlackPieces ^= BitSet[to + 8]
			b.OccupiedSquares ^= movebb | BitSet[to + 8]
			b.Squares[to + 8] = Empty
			b.Material -= PawnValue
		} else {
			b.makeCapture(capt, to)
		}
	} else {
		b.OccupiedSquares ^= movebb
	}

	if move.IsPromotion() {
		b.makeBlackProm(move.GetProm(), to)
		b.Squares[to] = move.GetProm()
	}
}

func (b *Board) makeBlackRookMove(move *Move) {
	from := move.GetFrom()
	to := move.GetTo()
	piece := move.GetPiece()
	capt := move.GetCapt()

	frombb := BitSet[from]
	movebb := frombb | BitSet[to]

	b.BlackRooks ^= movebb
	b.BlackPieces ^= movebb
	b.Squares[from] = Empty
	b.Squares[to] = piece
	b.Ep = 0
	b.NumMoves++
	
	if from == A8 {
		b.BlackCastle = b.BlackCastle &^ CanCastleOOO
	}
	if from == H8 {
		b.BlackCastle = b.BlackCastle &^ CanCastleOO
	}

	if capt != 0 {
		b.makeCapture(capt, to) 
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}
}

func (b *Board) makeBlackKingMove(move *Move) {
	from := move.GetFrom()
	to := move.GetTo()
	piece := move.GetPiece()
	capt := move.GetCapt()

	frombb := BitSet[from]
	movebb := frombb | BitSet[to]

	b.BlackKing ^= movebb
	b.BlackPieces ^= movebb
	b.Squares[from] = Empty
	b.Squares[to] = piece
	b.Ep = 0
	b.NumMoves++
	
	if capt != 0 {
		b.makeCapture(capt, to)
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}
	
	if move.IsCastle() {
		if move.IsCastleOO() {
			b.BlackRooks ^= BitSet[H8] | BitSet[F8]
			b.BlackPieces ^= BitSet[H8] | BitSet[F8]
			b.OccupiedSquares ^= BitSet[H8] | BitSet[F8]
			b.Squares[H8] = Empty
			b.Squares[F1] = BlackRook
		} else {
			b.BlackRooks ^= BitSet[A8] | BitSet[D8]
			b.BlackPieces ^= BitSet[A8] | BitSet[D8]
			b.OccupiedSquares ^= BitSet[A8] | BitSet[D8]
			b.Squares[A1] = Empty
			b.Squares[D1] = BlackRook
		}
	}
}

func (b *Board) unmakeBlackPawnMove(move *Move) {
	from := move.GetFrom()
	to := move.GetTo()
	piece := move.GetPiece()
	capt := move.GetCapt()

	frombb := BitSet[from]
	movebb := frombb | BitSet[to]

	b.BlackPawns ^= movebb
	b.BlackPieces ^= movebb
	b.Squares[from] = piece
	b.Squares[to] = Empty

	if capt != 0 {
		if move.IsEnpassant() {
			b.BlackPawns ^= BitSet[to + 8]
			b.BlackPieces ^= BitSet[to + 8]
			b.OccupiedSquares ^= movebb | BitSet[to + 8]
			b.Squares[to + 8] = WhitePawn
			b.Material += PawnValue
		} else {
			b.unmakeCapture(capt, to)
			b.OccupiedSquares = movebb
		}
	} else {
		b.OccupiedSquares ^= movebb
	}

	if move.IsPromotion() {
		b.unmakeBlackProm(move.GetProm(), to)
	}
}

func (b *Board) unmakeBlackRookMove(move *Move) {
	from := move.GetFrom()
	to := move.GetTo()
	piece := move.GetPiece()
	capt := move.GetCapt()

	frombb := BitSet[from]
	movebb := frombb | BitSet[to]

	b.BlackRooks ^= movebb
	b.BlackPieces ^= movebb
	b.Squares[from] = piece
	b.Squares[to] = Empty
	
	if capt != 0 {
		b.unmakeCapture(capt, to) 
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}
}

func (b *Board) unmakeBlackKingMove(move *Move) {
	from := move.GetFrom()
	to := move.GetTo()
	piece := move.GetPiece()
	capt := move.GetCapt()

	frombb := BitSet[from]
	movebb := frombb | BitSet[to]

	b.BlackKing ^= movebb
	b.BlackPieces ^= movebb
	b.Squares[from] = piece
	b.Squares[to] = Empty
	
	if capt != 0 {
		b.unmakeCapture(capt, to)
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}
	
	if move.IsCastle() {
		if move.IsCastleOO() {
			b.BlackRooks ^= BitSet[H8] | BitSet[F8]
			b.BlackPieces ^= BitSet[H8] | BitSet[F8]
			b.OccupiedSquares ^= BitSet[H8] | BitSet[F8]
			b.Squares[H8] = BlackRook
			b.Squares[F1] = Empty
		} else {
			b.BlackRooks ^= BitSet[A8] | BitSet[D8]
			b.BlackPieces ^= BitSet[A8] | BitSet[D8]
			b.OccupiedSquares ^= BitSet[A8] | BitSet[D8]
			b.Squares[A1] = BlackRook
			b.Squares[D1] = Empty
		}
	}
}

func (b *Board) makeBlackProm(prom, to byte) {
	movebb := BitSet[to]
	
	b.BlackPawns ^= movebb //clear pawn bit
	b.Material += PawnValue

	switch prom {
	case BlackQueen: 
		b.BlackQueens ^= movebb
		b.Material -= QueenValue
	case BlackRook:
		b.BlackRooks ^= movebb
		b.Material -= RookValue
	case BlackBishop:
		b.BlackBishops ^= movebb
		b.Material -= BishopValue
	case BlackKnight:
		b.BlackKnights ^= movebb
		b.Material -= KnightValue
	}	
}

func (b *Board) unmakeBlackProm(prom, to byte) {
	movebb := BitSet[to]
	
	b.BlackPawns ^= movebb //clear pawn bit
	b.Material -= PawnValue

	switch prom {
	case BlackQueen: 
		b.BlackQueens ^= movebb
		b.Material += QueenValue
	case BlackRook:
		b.BlackRooks ^= movebb
		b.Material += RookValue
	case BlackBishop:
		b.BlackBishops ^= movebb
		b.Material += BishopValue
	case BlackKnight:
		b.BlackKnights ^= movebb
		b.Material += KnightValue
	}	
}
