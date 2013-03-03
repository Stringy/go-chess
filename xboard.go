package main

import (
	//	"bufio"
	"chess/ai"
	"chess/ai/gen"
	"fmt"
	//	"io"
	"os"
	"strconv"
	"strings"
)

const (
	White   = gen.WhiteMove
	Black   = gen.BlackMove
	None    = 0
	Analyse = 3

	Invalid = 666

	MaxMoves = 500
	MaxPly   = 60

	Off = 0
	On  = 1

	DefaultFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
)

var (
	GameHistory = make([]gen.Move, MaxMoves)
	MoveCount   = 0
	board       gen.Board
	depth       int
	force       bool
)

func main() {
	var cmd string
	//	in := bufio.NewReader(os.Stdin)
	board.Init()
	for {
		os.Stdout.Sync()
		bytes := make([]byte, 100)
		n, err := os.Stdin.Read(bytes)
		if err != nil {
			fmt.Println(err)
		}
		cmd = string(bytes)
		fmt.Printf(cmd)
		HandleCommand(cmd[:n])
		if board.NextMove == Black {
			move := ai.GetBestMove(board, 4)
			board.MakeMove(move)
			os.Stdout.Write([]byte("move " + move.String()))
		}
	}
}

func HandleCommand(cmd string) {
	fmt.Println("ECHO:", cmd)

	if cmd == "go" {
		move := ai.GetBestMove(board, 4)
		board.MakeMove(move)
		os.Stdout.Write([]byte("move " + move.String()))
	}

	if cmd == "protover" {
		os.Stdout.Write([]byte("feature ping=1 setboard=1 colors=0 usermove=1 memory=1 debug=1"))
		os.Stdout.Write([]byte("feature option=\"Resign -check 0\""))
		os.Stdout.Write([]byte("feature option=\"Contempt -spin -200 200\""))
		os.Stdout.Write([]byte("featrue done=1"))
		return
	}

	if isMoveCommand(cmd) {
		if move, err := gen.NewMove(cmd, &board); err == nil {
			if move.IsLegalMove(&board) {
				board.MakeMove(move)
				GameHistory[MoveCount] = *move
				MoveCount++
			} else {
				os.Stdout.Write([]byte("Illegal Move: " + cmd))
			}
		}
		return
	}

	if len(cmd) > 8 && cmd[:8] == "setboard" {
		//set fen strings
		fen := cmd[8:]
		words := strings.Split(fen, " ")
		var half, full int
		if num, err := strconv.Atoi(words[4]); err == nil {
			half = num
		}
		if num, err := strconv.Atoi(words[5]); err == nil {
			full = num
		}
		board = *ai.SetupBoardFromFen(
			words[0],
			words[1],
			words[2],
			words[3],
			half,
			full)
		return
	}
	if cmd == "undo" {
		MoveCount -= 1
		board.UnmakeMove(&GameHistory[MoveCount])
		return
	}
	if cmd == "remove" {
		MoveCount -= 1
		board.UnmakeMove(&GameHistory[MoveCount])
		MoveCount -= 1
		board.UnmakeMove(&GameHistory[MoveCount])
		return
	}
	if len(cmd) > 4 && cmd[:4] == "ping" {
		num, err := strconv.Atoi(cmd[5:len(cmd)])
		if err != nil {

		}
		pong := fmt.Sprintf("pong %d", num)
		os.Stdout.Write([]byte(pong))
		return
	}

	if cmd == "xboard" {
		//do nothing, no xboard mode necessary
		return
	}

	if cmd == "new" {
		board.Init()
		return
	}

	if cmd == "force" {
		force = true
		return
	}

	if cmd == "quit" {
		os.Exit(0)
	}
}

func isMoveCommand(cmd string) bool {
	if len(cmd) > 5 || len(cmd) < 4 {
		return false
	}
	if cmd == "quit" || cmd == "ping" {
		return false
	}
	return true
}
