package main

import (
	"chess/ai"
	"chess/ai/eval"
	"chess/ai/gen"
	// "chess/ai/search"
	"flag"
	"fmt"
	"log"
	"os"
	//	"runtime"
	"runtime/pprof"
)

var (
	//cmd line flags
	cpu       = flag.String("cpuprofile", "", "write cpu profile to file")
	mem       = flag.String("memprofile", "", "write a memory profile")
	depthFlag = flag.Int("depth", 4, "Define Search Depth")
	twoPlayer = flag.Bool("twoplayer", false, "Two humans fight to the death")
	depth     = 4

	Players     = make([]Player, 2)
	GameHistory = make([]gen.Move, 500)
	MoveCount   = 0
	currentTurn = 0

	playerOne Player
	playerTwo Player

	evaluator = eval.BasicEvaluator
)

const ()

type Player interface {
	GetMove(b *gen.Board) gen.Move
}

type HumanPlayer struct{}

type AIPlayer struct{}

func main() {
	Initialise()
	defer cleanUp()

	fmt.Println("Aperture Science Simulation of Heuristic Analysis Techniques")
	fmt.Println("A.S.S.H.A.T")
	var board gen.Board
	var move gen.Move

	board.Init()

	for {
		stop := make(chan struct{})
		board.PrintBoard()
		if currentTurn == 0 { //human
			go ai.PonderUntilInput(&board, stop)
		}
		move = Players[currentTurn].GetMove(&board)
		board.MakeMove(&move)
		currentTurn = (currentTurn + 1) % 2
		GameHistory[MoveCount] = move
		MoveCount++
		close(stop)
	}
}

func cleanUp() {
	if *mem != "" {
		f, err := os.Create(*mem)
		if err != nil {
			os.Exit(1)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
	pprof.StopCPUProfile()
}

func Initialise() {

	flag.Parse()
	if *cpu != "" {
		f, err := os.Create(*cpu)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}

	if *depthFlag != 4 {
		depth = *depthFlag
	}

	if *twoPlayer {
		playerOne = HumanPlayer{}
		playerTwo = HumanPlayer{}
		Players = []Player{playerOne, playerTwo}
	} else {
		playerOne = HumanPlayer{}
		playerTwo = AIPlayer{}
		Players = []Player{playerOne, playerTwo}
	}
}

func (h HumanPlayer) GetMove(b *gen.Board) gen.Move {
	move := gen.Move(0)
	var cmd string
	for {
		fmt.Print(">>> ")
		fmt.Scanf("%s", &cmd)
		switch cmd {
		case "generate":
			moves := gen.GenerateAllMoves(b)
			for i, m := range moves {
				if i%5 == 0 {
					fmt.Println()
				}
				fmt.Print(m.String(), ", ")
			}
			fmt.Println()
		case "undo":
			if MoveCount == 0 {
				break
			}
			MoveCount -= 1
			b.UnmakeMove(&GameHistory[MoveCount])
		case "debug":
			ai.PrintDebug()
		case "eval":
			score := evaluator.Eval(b)
			evaluator.Debug()
			fmt.Println("EVALUATED SCORE:", score)
		case "print":
			b.PrintBoard()
		case "help":
			help()
		case "exit":
			fallthrough
		case "quit":
			cleanUp()
			os.Exit(0)
		default: //assume it's a move
			if len(cmd) > 5 {
				continue
			}
			move, err := gen.NewMove(cmd, b)
			if err == nil && move.IsLegalMove(b) {
				return *move
			} else {
				fmt.Println(err)
			}
		}
	}
	return move
}

func (a AIPlayer) GetMove(b *gen.Board) gen.Move {
	return *ai.GetBestMove(*b, depth)
}

func help() {
	fmt.Println("help     \tThis help text")
	fmt.Println("generate \tgenerate all moves for the current position")
	fmt.Println("undo     \tundo the last played move")
	fmt.Println("debug    \tprint debug information for the last search")
	fmt.Println("eval     \tprint evaluation information for the current position")
	fmt.Println("print    \tprint the board")
	fmt.Println("exit/quit\tquit the game")
}
