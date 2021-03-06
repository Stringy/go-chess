package gen

import (
	"fmt"
	"testing"
	"time"
)

type sq struct {
	i     string
	piece byte
}

func populateBoard(squares ...sq) Board {
	var board Board
	grid := make([]byte, 64)
	for _, s := range squares {
		grid[SquareMap[s.i]] = s.piece
	}
	board.InitFromGrid(grid,
		WhiteMove,
		0,
		CanCastleOO+CanCastleOOO,
		CanCastleOO+CanCastleOOO,
		0)
	return board
}

func TestRookMoveGeneration(t *testing.T) {
	fmt.Println("\n=======================ROOK TEST==========================\n")
	fmt.Println("The purpose of this test is to check an apparent bug in rook move generation")
	fmt.Println("where it is generating moves to take own pieces")

	board := populateBoard(
		sq{"A1", WhiteRook},
		//sq{"B1", WhiteKnight},
		sq{"A2", WhitePawn},
		sq{"C1", WhiteBishop},
	)
	board.PrintBoard()
	moves := GenerateAllMoves(&board)
	for _, move := range moves {
		board.MakeMove(&move)
		fmt.Println(move.String())
	}
	board.PrintBoard()
}

func TestQueenMoveGeneration(t *testing.T) {
	fmt.Println("\n======================QUEEN TEST===========================\n")
	fmt.Println("Populating Board with White Queen piece on D4 and a blocking Black Bishop on F6")
	fmt.Println("G4, D6, and B6. Note Bishops used purely for blocking.\n")
	fmt.Println("Expected no moves past these squares along the direction of travel of the queen")
	fmt.Println("Inclusion of these squares is expected\n")
	fmt.Println("Queen move generation also tests rook / bishop generation due to using same mechanisms\n")
	board := populateBoard(
		sq{"D4", WhiteQueen},
		sq{"F6", BlackBishop},
		sq{"G4", BlackBishop},
		sq{"D6", BlackBishop},
		sq{"B6", BlackBishop},
	)

	//	board.PrintBoard()
	fmt.Println("Generating Moves and Making them on the board")
	moves := GenerateAllMoves(&board)

	if len(moves) == 0 {
		fmt.Println("ERROR: No moves generated")
		t.Fail()
	}

	fmt.Println("Making All Moves")
	for _, move := range moves {
		board.MakeMove(&move)
	}

	if board.Squares[SquareMap["G7"]] != Empty {
		fmt.Println("ERROR: Blocking Bishop on F6 (A1H8 Diag) is Ineffective")
		t.Fail()
	}
	if board.Squares[SquareMap["H4"]] != Empty {
		fmt.Println("ERROR: Blocking Bishop on G4 (Rank Attack) is Ineffective")
		t.Fail()
	}
	if board.Squares[SquareMap["D7"]] != Empty {
		fmt.Println("ERROR: Blocking Bishop on D6 (File Attack) is Ineffective")
		t.Fail()
	}
	if board.Squares[SquareMap["A7"]] != Empty {
		fmt.Println("ERROR: Blocking Bishop on B6 (A8H1 Diag) is Ineffective")
		t.Fail()
	}
	//	board.PrintBoard()
	fmt.Println("Passed")
}

func TestPawnMoveGeneration(t *testing.T) {
	fmt.Println("\n========================PAWN TEST===========================\n")
	fmt.Println("Populating board with White Pawn on D2 along with two Black Pawns along")
	fmt.Println("the White Pawn's attack lines (C3, E3)\n")
	fmt.Println("Expecting both attacks to occour along with the initial double move and single move\n")
	board := populateBoard(
		sq{"D2", WhitePawn},
		sq{"E3", BlackPawn},
		sq{"C3", BlackPawn},
	)
	fmt.Println("Generating all moves for the board")
	moves := GenerateAllMoves(&board)
	fmt.Println("Making all moves on the board")
	for _, move := range moves {
		board.MakeMove(&move)
	}
	//	board.PrintBoard()
	if board.Squares[SquareMap["E3"]] != WhitePawn {
		fmt.Println("ERROR: No attack to BlackPawn on E3")
		t.Fail()
	}
	if board.Squares[SquareMap["C3"]] != WhitePawn {
		fmt.Println("ERROR: No attack to BlackPawn on C3")
		t.Fail()
	}
	if board.Squares[SquareMap["D3"]] != WhitePawn {
		fmt.Println("ERROR: No single square move from D2 to D3")
		t.Fail()
	}
	if board.Squares[SquareMap["D4"]] != WhitePawn {
		fmt.Println("ERROR: No double square move from D2 to D4")
		t.Fail()
	}
}

func TestKnightMoveGeneration(t *testing.T) {
	fmt.Println("\n======================KNIGHT TEST===========================\n")
	fmt.Println("Populating board with a knight on F3 to test all possible move positions")
	fmt.Println("and another on the edge of the board at A5 to check boundary moves\n")
	fmt.Println("Expecting moves to the following positions:")
	fmt.Println("\tF3-H2")
	fmt.Println("\tF3-H4")
	fmt.Println("\tF3-G5")
	fmt.Println("\tF3-E5")
	fmt.Println("\tF3-G1")
	fmt.Println("\tF3-E1")
	fmt.Println("\tF3-D2")
	fmt.Println("\tF3-D4")
	fmt.Println("\tA5-B3")
	fmt.Println("\tA5-C4")
	fmt.Println("\tA5-C6")
	fmt.Println("\tA5-B7")

	board := populateBoard(
		sq{"F3", WhiteKnight},
		sq{"A5", WhiteKnight},
	)
	moves := GenerateAllMoves(&board)
	for _, move := range moves {
		board.MakeMove(&move)
	}
	//	board.PrintBoard()
	//bitcount of occupied squares is expected to be 14 because it includes the
	//two original knights
	if BitCount(board.OccupiedSquares) != 14 || len(moves) != 12 {
		fmt.Println("Unexpected number of moves generated", BitCount(board.OccupiedSquares))
		board.PrintBoard()
		t.FailNow()
	}
	expectedSquares := []string{
		"H2", "H4", "G5", "E5", "G1", "E1",
		"D2", "D4", "B3", "C4", "C6", "B7",
	}
	fmt.Println("Checking all expected moves")
	for _, expected := range expectedSquares {
		if board.Squares[SquareMap[expected]] != WhiteKnight {
			fmt.Println("ERROR: No move on square", expected)
			t.Fail()
		}
	}
}

func TestKingMoveGeneration(t *testing.T) {
	fmt.Println("\n========================KING TEST===========================\n")
	fmt.Println("Populating the board with a WhiteKing on its starting position E1")
	fmt.Println("along with a two rooks, each also on its starting position H1 and A1 respectively\n")
	fmt.Println("Expecting castling moves and normal moves to be generated")

	board := populateBoard(
		sq{"E1", WhiteKing},
		sq{"H1", WhiteRook},
		sq{"A1", WhiteRook},
	)
	moves := GenerateAllMoves(&board)
	for _, move := range moves {
		board.MakeMove(&move)
	}
	board.PrintBoard()

	expectedKingSquares := []string{
		"G1", "C1", "F2", "E2", "D2",
	}
	for _, expected := range expectedKingSquares {
		if board.Squares[SquareMap[expected]] != WhiteKing {
			fmt.Println("ERROR: No move to", expected)
			t.Fail()
		}
	}
}

func TestCheckMateDetection(t *testing.T) {
	board := populateBoard(
		sq{"H8", WhiteKing},
		sq{"A8", BlackRook},
		sq{"H1", BlackRook},
		sq{"A1", BlackQueen})
	board.PrintBoard()
	board.NextMove = WhiteMove

	board.BlackCastle = 0
	board.WhiteCastle = 0

	if board.IsCheckMate() {
		fmt.Println("WIN")
	} else {
		t.Fail()
	}
}

func TestCompleteGeneration(t *testing.T) {
	for depth := 0; depth < 11; depth++ {
		b := NewBoard()
		t0 := time.Now()
		nodes := Perft(b, depth)
		fmt.Println(depth, ":", nodes, "nodes/second : ", float64(nodes)/time.Since(t0).Seconds())
	}
}
