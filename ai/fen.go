package ai

import (
	"chess/ai/gen"
	"errors"
	"fmt"
	"os"
	"strings"
)

func ReadFenString(filename string, n int) (*gen.Board, error) {
	var (
		fenwhite,
		fenblack,
		fen,
		fencol,
		fencastle,
		fenep string

		fenhalf,
		fenfull int
	)
	if n <= 0 {
		return nil, errors.New("Invalid FEN string index")
	}
	if file, err := os.Open(filename); err != nil {
		return nil, err
	} else {
		defer file.Close()
		var word string
		for _, err := fmt.Fscanf(file, "%s", word); err == nil; {
			switch word {
			case "[White":
				fmt.Fscanf(file, "%s", fenwhite)
				fenwhite = fenwhite[1 : len(fenwhite)-2]
			case "[Black":
				fmt.Fscanf(file, "%s", fenblack)
				fenblack = fenblack[1 : len(fenblack)-2]
			case "[FEN":
				fmt.Fscanf(file, "%s", fen)
				fmt.Fscanf(file, "%s", fencol)
				fmt.Fscanf(file, "%s", fencastle)
				fmt.Fscanf(file, "%s", fenep)
				fmt.Fscanf(file, "%d", fenhalf)
				fmt.Fscanf(file, "%d", fenfull)
			}
		}
	}
	return SetupBoardFromFen(fen, fencol, fencastle, fenep, fenhalf, fenfull), nil
}

func SetupBoardFromFen(fen, fencol, fencastle, fenep string, fenhalf, fenfull int) *gen.Board {
	var board gen.Board
	squares := make([]byte, 64)
	for i := 0; i < 64; i++ {
		squares[i] = gen.Empty
	}

	file, rank := 1, 8
	for i, j := 0, 0; i < len(fen) && j < 64; i++ {
		if int(fen[i]) > 48 && int(fen[i]) < 57 {
			file += int(fen[i]) - 48
			j += int(fen[i]) - 48
		} else {
			switch fen[i] {
			case '/':
				rank -= 1
			case 'B':
				squares[j] = gen.WhiteBishop
				file++
				j++
			case 'R':
				squares[j] = gen.WhiteRook
				file++
				j++
			case 'Q':
				squares[j] = gen.WhiteQueen
				file++
				j++
			case 'K':
				squares[j] = gen.WhiteKing
				file++
				j++
			case 'p':
				squares[j] = gen.BlackPawn
				file++
				j++
			case 'n':
				squares[j] = gen.BlackKnight
				file++
				j++
			case 'b':
				squares[j] = gen.BlackBishop
				file++
				j++
			case 'r':
				squares[j] = gen.BlackRook
				file++
				j++
			case 'q':
				squares[j] = gen.BlackQueen
				file++
				j++
			case 'k':
				squares[j] = gen.BlackKing
				file++
				j++
			default:
				break
			}
		}
	}
	next := gen.WhiteMove
	if fencol[0] == 'b' {
		next = gen.BlackMove
	}
	wcastle := byte(0)
	bcastle := byte(0)
	if strings.Contains(fencastle, "K") {
		wcastle += gen.CanCastleOO
	}
	if strings.Contains(fencastle, "Q") {
		wcastle += gen.CanCastleOOO
	}
	if strings.Contains(fencastle, "k") {
		bcastle += gen.CanCastleOO
	}
	if strings.Contains(fencastle, "q") {
		bcastle += gen.CanCastleOOO
	}
	var ep int
	if strings.Contains(fenep, "-") {
		ep = 0
	} else {
		ep = (int(fenep[0]) - 96) + 8*(int(fenep[1])-48) - 9
	}
	board.InitFromGrid(squares, next, fenhalf, wcastle, bcastle, byte(ep))
	return &board
}
