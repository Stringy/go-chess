package gen

import ()

func (b *Board) MakeAllMoves(moves []Move) {
	for _, move := range moves {
		b.MakeMove(&move)
	}
}

//Board.MakeMove applies the move upon the board
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
	case WhitePawn:
		b.makeWhitePawnMove(move)
	case WhiteRook:
		b.makeWhiteRookMove(move)
	case WhiteKnight:
		b.makeGenericMove(move, &b.WhiteKnights, &b.WhitePieces)
	case WhiteBishop:
		b.makeGenericMove(move, &b.WhiteBishops, &b.WhitePieces)
	case WhiteQueen:
		b.makeGenericMove(move, &b.WhiteQueens, &b.WhitePieces)
	case WhiteKing:
		b.makeWhiteKingMove(move)
		//black moves
	case BlackPawn:
		b.makeBlackPawnMove(move)
	case BlackRook:
		b.makeBlackRookMove(move)
	case BlackKnight:
		b.makeGenericMove(move, &b.BlackKnights, &b.BlackPieces)
	case BlackBishop:
		b.makeGenericMove(move, &b.BlackBishops, &b.BlackPieces)
	case BlackQueen:
		b.makeGenericMove(move, &b.BlackQueens, &b.BlackPieces)
	case BlackKing:
		b.makeBlackKingMove(move)
	}
	next()
}

//Board.UnmakeMove undoes the move which has been previously applied to the board
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

//Board.makeCapture makes a capture of a piece
func (b *Board) makeCapture(capt, to byte) {
	movebb := BitSet[to]

	switch capt {
	case WhitePawn:
		b.WhitePawns ^= movebb
		b.WhitePieces ^= movebb
		b.Material -= PawnValue
		b.WMaterial -= PawnValue
		b.NumWPawns -= 1
	case WhiteRook:
		b.WhiteRooks ^= movebb
		b.WhitePieces ^= movebb
		b.Material -= RookValue
		b.WMaterial -= RookValue
		b.NumWRooks -= 1
		if to == SquareMap["A1"] {
			b.WhiteCastle = b.WhiteCastle &^ CanCastleOOO
		} else if to == SquareMap["H1"] {
			b.WhiteCastle = b.WhiteCastle &^ CanCastleOO
		}
	case WhiteKnight:
		b.WhiteKnights ^= movebb
		b.WhitePieces ^= movebb
		b.Material -= KnightValue
		b.WMaterial -= KnightValue
		b.NumWKnights -= 1
	case WhiteBishop:
		b.WhiteBishops ^= movebb
		b.WhitePieces ^= movebb
		b.Material -= BishopValue
		b.WMaterial -= BishopValue
		b.NumWBishops -= 1
	case WhiteQueen:
		b.WhiteQueens ^= movebb
		b.WhitePieces ^= movebb
		b.Material -= QueenValue
		b.WMaterial -= QueenValue
		b.NumWQueens -= 1
	case WhiteKing:
		b.WhiteKing ^= movebb
		b.WhitePieces ^= movebb

		//Black Captures
	case BlackPawn:
		b.BlackPawns ^= movebb
		b.BlackPieces ^= movebb
		b.Material += PawnValue
		b.BMaterial += PawnValue
		b.NumBPawns -= 1
	case BlackRook:
		b.BlackRooks ^= movebb
		b.BlackPieces ^= movebb
		b.Material += RookValue
		b.BMaterial += RookValue
		b.NumBRooks -= 1
		if to == SquareMap["A8"] {
			b.BlackCastle = b.BlackCastle &^ CanCastleOOO
		} else if to == SquareMap["H8"] {
			b.BlackCastle = b.BlackCastle &^ CanCastleOO
		}
	case BlackKnight:
		b.BlackKnights ^= movebb
		b.BlackPieces ^= movebb
		b.Material += KnightValue
		b.BMaterial += KnightValue
		b.NumBKnights -= 1
	case BlackBishop:
		b.BlackBishops ^= movebb
		b.BlackPieces ^= movebb
		b.Material += BishopValue
		b.BMaterial += BishopValue
		b.NumBBishops -= 1
	case BlackQueen:
		b.BlackQueens ^= movebb
		b.BlackPieces ^= movebb
		b.Material += QueenValue
		b.BMaterial += QueenValue
		b.NumBQueens -= 1
	case BlackKing:
		b.BlackKing ^= movebb
		b.BlackPieces ^= movebb
	}
	b.NumMoves = 0
}

//Board.unmakeCapture unmakes the associated capture
func (b *Board) unmakeCapture(capt, to byte) {
	movebb := BitSet[to]

	switch capt {
	case WhitePawn:
		b.WhitePawns ^= movebb
		b.WhitePieces ^= movebb
		b.Material += PawnValue
		b.WMaterial += PawnValue
		b.NumWPawns++
		b.Squares[to] = WhitePawn
	case WhiteRook:
		b.WhiteRooks ^= movebb
		b.WhitePieces ^= movebb
		b.Material += RookValue
		b.WMaterial += RookValue
		b.NumBRooks++
		b.Squares[to] = WhiteRook
	case WhiteKnight:
		b.WhiteKnights ^= movebb
		b.WhitePieces ^= movebb
		b.Material += KnightValue
		b.WMaterial += KnightValue
		b.NumWKnights++
		b.Squares[to] = WhiteKnight
	case WhiteBishop:
		b.WhiteBishops ^= movebb
		b.WhitePieces ^= movebb
		b.Material += BishopValue
		b.WMaterial += BishopValue
		b.NumWBishops++
		b.Squares[to] = WhiteBishop
	case WhiteQueen:
		b.WhiteQueens ^= movebb
		b.WhitePieces ^= movebb
		b.Material += QueenValue
		b.Material += QueenValue
		b.NumWQueens++
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
		b.BMaterial -= PawnValue
		b.NumBPawns++
		b.Squares[to] = BlackPawn
	case BlackRook:
		b.BlackRooks ^= movebb
		b.BlackPieces ^= movebb
		b.Material -= RookValue
		b.BMaterial -= RookValue
		b.NumBRooks++
		b.Squares[to] = BlackRook
	case BlackKnight:
		b.BlackKnights ^= movebb
		b.BlackPieces ^= movebb
		b.Material -= KnightValue
		b.BMaterial -= KnightValue
		b.NumBKnights++
		b.Squares[to] = BlackKnight
	case BlackBishop:
		b.BlackBishops ^= movebb
		b.BlackPieces ^= movebb
		b.Material -= BishopValue
		b.BMaterial -= BishopValue
		b.NumBBishops++
		b.Squares[to] = BlackBishop
	case BlackQueen:
		b.BlackQueens ^= movebb
		b.BlackPieces ^= movebb
		b.Material -= QueenValue
		b.BMaterial -= QueenValue
		b.NumBQueens++
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

	*pieces = *pieces ^ movebb
	*colour = *colour ^ movebb
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
	*pieces = *pieces ^ movebb
	*colour = *colour ^ movebb
	b.Squares[from] = piece
	b.Squares[to] = Empty

	if capt != 0 {
		b.unmakeCapture(capt, to)
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}
}

//Board.makeWhitePawnMove makes a white pawn move on the board, including 
//special rules
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
			b.BlackPawns ^= BitSet[to-8]
			b.BlackPieces ^= BitSet[to-8]
			b.OccupiedSquares ^= movebb | BitSet[to-8]
			b.Squares[to-8] = Empty
			b.Material += PawnValue
			b.BMaterial += PawnValue
			b.NumBPawns -= 1
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

//Board.makeWhiteRookMove makes a rook move on the board, including castling rules
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

	if from == SquareMap["A1"] {
		b.WhiteCastle = b.WhiteCastle &^ CanCastleOOO
	}
	if from == SquareMap["H1"] {
		b.WhiteCastle = b.WhiteCastle &^ CanCastleOO
	}

	if capt != 0 {
		b.makeCapture(capt, to)
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}
}

//Board.makeWhiteKingMove makes a king move on the board including castling rules
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
	b.WKingSq = to
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
			b.WhiteRooks ^= BitSet[SquareMap["H1"]] | BitSet[SquareMap["F1"]]
			b.WhitePieces ^= BitSet[SquareMap["H1"]] | BitSet[SquareMap["F1"]]
			b.OccupiedSquares ^= BitSet[SquareMap["H1"]] | BitSet[SquareMap["F1"]]
			b.Squares[SquareMap["H1"]] = Empty
			b.Squares[SquareMap["F1"]] = WhiteRook
		} else {
			b.WhiteRooks ^= BitSet[SquareMap["A1"]] | BitSet[SquareMap["D1"]]
			b.WhitePieces ^= BitSet[SquareMap["A1"]] | BitSet[SquareMap["D1"]]
			b.OccupiedSquares ^= BitSet[SquareMap["A1"]] | BitSet[SquareMap["D1"]]
			b.Squares[SquareMap["A1"]] = Empty
			b.Squares[SquareMap["D1"]] = WhiteRook
		}
	}
}

//Board.unmakeWhitePawnMoves unmakes a pawn move, including undoing special rules
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
			b.BlackPawns ^= BitSet[to-8]
			b.BlackPieces ^= BitSet[to-8]
			b.OccupiedSquares ^= movebb | BitSet[to-8]
			b.Squares[to-8] = BlackPawn
			b.Material -= PawnValue
			b.BMaterial -= PawnValue
			b.NumBPawns++
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

//Board.unmakeRookMove unmakes a rook move, including capture
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

//Board.unmakeWhiteKingMove unmakes a king move including castling
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
	b.WKingSq = from

	if capt != 0 {
		b.unmakeCapture(capt, to)
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}

	if move.IsCastle() {
		if move.IsCastleOO() {
			b.WhiteRooks ^= BitSet[SquareMap["H1"]] |
				BitSet[SquareMap["F1"]]
			b.WhitePieces ^= BitSet[SquareMap["H1"]] |
				BitSet[SquareMap["F1"]]
			b.OccupiedSquares ^= BitSet[SquareMap["H1"]] |
				BitSet[SquareMap["F1"]]
			b.Squares[SquareMap["H1"]] = WhiteRook
			b.Squares[SquareMap["F1"]] = Empty
		} else {
			b.WhiteRooks ^= BitSet[SquareMap["A1"]] | BitSet[SquareMap["D1"]]
			b.WhitePieces ^= BitSet[SquareMap["A1"]] | BitSet[SquareMap["D1"]]
			b.OccupiedSquares ^= BitSet[SquareMap["A1"]] | BitSet[SquareMap["D1"]]
			b.Squares[SquareMap["A1"]] = WhiteRook
			b.Squares[SquareMap["D1"]] = Empty
		}
	}
}

//Board.makeWhiteProm makes a white promotion
func (b *Board) makeWhiteProm(piece, to byte) {
	movebb := BitSet[to]

	b.WhitePawns ^= movebb //clear pawn bit
	b.Material -= PawnValue
	b.WMaterial -= PawnValue

	switch piece {
	case WhiteQueen:
		b.WhiteQueens ^= movebb
		b.Material += QueenValue
		b.WMaterial += QueenValue
		b.NumWQueens++
	case WhiteRook:
		b.WhiteRooks ^= movebb
		b.Material += RookValue
		b.WMaterial += RookValue
		b.NumWRooks++
	case WhiteBishop:
		b.WhiteBishops ^= movebb
		b.Material += BishopValue
		b.WMaterial += BishopValue
		b.NumWBishops++
	case WhiteKnight:
		b.WhiteKnights ^= movebb
		b.Material += KnightValue
		b.WMaterial += KnightValue
		b.NumWKnights++
	}
}

//Board.unmakeWhiteProm unmakes a white promotion
func (b *Board) unmakeWhiteProm(piece, to byte) {
	movebb := BitSet[to]

	b.WhitePawns ^= movebb //clear pawn bit
	b.Material += PawnValue
	b.WMaterial += PawnValue

	switch piece {
	case WhiteQueen:
		b.WhiteQueens ^= movebb
		b.Material -= QueenValue
		b.WMaterial -= QueenValue
		b.NumWQueens -= 1
	case WhiteRook:
		b.WhiteRooks ^= movebb
		b.Material -= RookValue
		b.WMaterial -= RookValue
		b.NumWRooks -= 1
	case WhiteBishop:
		b.WhiteBishops ^= movebb
		b.Material -= BishopValue
		b.WMaterial -= BishopValue
		b.NumWBishops -= 1
	case WhiteKnight:
		b.WhiteKnights ^= movebb
		b.Material -= KnightValue
		b.WMaterial -= KnightValue
		b.NumWKnights -= 1
	}
}

//Board.makeBlackPawnMove makes a black pawn move on the board including special rules
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
			b.BlackPawns ^= BitSet[to+8]
			b.BlackPieces ^= BitSet[to+8]
			b.OccupiedSquares ^= movebb | BitSet[to+8]
			b.Squares[to+8] = Empty
			b.Material -= PawnValue
			b.WMaterial -= PawnValue
			b.NumWPawns -= 1
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

//Board.makeBlackRookMove makes black rook move including castling moves
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

	if from == SquareMap["A8"] {
		b.BlackCastle = b.BlackCastle &^ CanCastleOOO
	}
	if from == SquareMap["H8"] {
		b.BlackCastle = b.BlackCastle &^ CanCastleOO
	}

	if capt != 0 {
		b.makeCapture(capt, to)
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}
}

//Board.makeBlackKingMove makes a black king move including castling rules
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
	b.BKingSq = to
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
			b.BlackRooks ^= BitSet[SquareMap["H8"]] | BitSet[SquareMap["F8"]]
			b.BlackPieces ^= BitSet[SquareMap["H8"]] | BitSet[SquareMap["F8"]]
			b.OccupiedSquares ^= BitSet[SquareMap["H8"]] | BitSet[SquareMap["F8"]]
			b.Squares[SquareMap["H8"]] = Empty
			b.Squares[SquareMap["F1"]] = BlackRook
		} else {
			b.BlackRooks ^= BitSet[SquareMap["A8"]] | BitSet[SquareMap["D8"]]
			b.BlackPieces ^= BitSet[SquareMap["A8"]] | BitSet[SquareMap["D8"]]
			b.OccupiedSquares ^= BitSet[SquareMap["A8"]] | BitSet[SquareMap["D8"]]
			b.Squares[SquareMap["A1"]] = Empty
			b.Squares[SquareMap["D1"]] = BlackRook
		}
	}
}

//Board.unmakeBlackPawnMove unmakes a black pawn move including special rules
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
			b.BlackPawns ^= BitSet[to+8]
			b.BlackPieces ^= BitSet[to+8]
			b.OccupiedSquares ^= movebb | BitSet[to+8]
			b.Squares[to+8] = WhitePawn
			b.Material += PawnValue
			b.WMaterial += PawnValue
			b.NumWPawns++
		} else {
			b.unmakeCapture(capt, to)
			b.OccupiedSquares ^= movebb
		}
	} else {
		b.OccupiedSquares ^= movebb
	}

	if move.IsPromotion() {
		b.unmakeBlackProm(move.GetProm(), to)
	}
}

//Board.unmakeBlackRookMove unmakes a rook move including castling
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

//Board.unmakeBlackKingMove unmakes a king move including castling 
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
	b.BKingSq = from

	if capt != 0 {
		b.unmakeCapture(capt, to)
		b.OccupiedSquares ^= movebb
	} else {
		b.OccupiedSquares ^= movebb
	}

	if move.IsCastle() {
		if move.IsCastleOO() {
			b.BlackRooks ^= BitSet[SquareMap["H8"]] | BitSet[SquareMap["F8"]]
			b.BlackPieces ^= BitSet[SquareMap["H8"]] | BitSet[SquareMap["F8"]]
			b.OccupiedSquares ^= BitSet[SquareMap["H8"]] | BitSet[SquareMap["F8"]]
			b.Squares[SquareMap["H8"]] = BlackRook
			b.Squares[SquareMap["F1"]] = Empty
		} else {
			b.BlackRooks ^= BitSet[SquareMap["A8"]] | BitSet[SquareMap["D8"]]
			b.BlackPieces ^= BitSet[SquareMap["A8"]] | BitSet[SquareMap["D8"]]
			b.OccupiedSquares ^= BitSet[SquareMap["A8"]] | BitSet[SquareMap["D8"]]
			b.Squares[SquareMap["A1"]] = BlackRook
			b.Squares[SquareMap["D1"]] = Empty
		}
	}
}

//Board.makeBlackProm makes a black promotion move
func (b *Board) makeBlackProm(prom, to byte) {
	movebb := BitSet[to]

	b.BlackPawns ^= movebb //clear pawn bit
	b.Material += PawnValue
	b.BMaterial += PawnValue

	switch prom {
	case BlackQueen:
		b.BlackQueens ^= movebb
		b.Material -= QueenValue
		b.BMaterial -= QueenValue
		b.NumBQueens++
	case BlackRook:
		b.BlackRooks ^= movebb
		b.Material -= RookValue
		b.BMaterial -= RookValue
		b.NumBRooks++
	case BlackBishop:
		b.BlackBishops ^= movebb
		b.Material -= BishopValue
		b.BMaterial -= BishopValue
		b.NumBBishops++
	case BlackKnight:
		b.BlackKnights ^= movebb
		b.Material -= KnightValue
		b.BMaterial -= KnightValue
		b.NumBKnights++
	}
}

//Board.unmakeBlackProm unmakes a black promotion move
func (b *Board) unmakeBlackProm(prom, to byte) {
	movebb := BitSet[to]

	b.BlackPawns ^= movebb //clear pawn bit
	b.Material -= PawnValue
	b.BMaterial -= PawnValue

	switch prom {
	case BlackQueen:
		b.BlackQueens ^= movebb
		b.Material += QueenValue
		b.BMaterial += QueenValue
		b.NumBQueens -= 1
	case BlackRook:
		b.BlackRooks ^= movebb
		b.Material += RookValue
		b.BMaterial += RookValue
		b.NumBRooks -= 1
	case BlackBishop:
		b.BlackBishops ^= movebb
		b.Material += BishopValue
		b.BMaterial += BishopValue
		b.NumBBishops -= 1
	case BlackKnight:
		b.BlackKnights ^= movebb
		b.Material += KnightValue
		b.BMaterial += KnightValue
		b.NumBKnights -= 1
	}
}
