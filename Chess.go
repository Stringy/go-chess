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

	Players     = make([]ai.Player, 2)
	GameHistory = make([]gen.Move, 500)
	MoveCount   = 0
	currentTurn = 0

	playerOne ai.Player
	playerTwo ai.Player

	evaluator = eval.BasicEvaluator
)

//init is run before the main func
//initialises players and depth from command line args
func init() {
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
		playerOne = &HumanPlayer{}
		playerTwo = &HumanPlayer{}
		Players = []ai.Player{playerOne, playerTwo}
	} else {
		playerOne = &HumanPlayer{}
		playerTwo = ai.NewAIPlayer()
		Players = []ai.Player{playerOne, playerTwo}
	}
}

func main() {
	defer cleanUp()

	fmt.Println("Aperture Science Simulation of Heuristic Analysis Techniques")
	fmt.Println("A.S.S.H.A.T")
	var board = gen.NewBoard()
	var move *gen.Move
	var ponderMoves = []*gen.Move{nil, nil}

	for {
		stop := make(chan struct{})
		board.PrintBoard()

		//start pondering in new goroutine
		go Players[(currentTurn+1)%2].PonderUntilInput(board, stop, ponderMoves[(currentTurn+1)%2])

		//if current player discovered solution while pondering last turn
		//make that move
		if ponderMoves[currentTurn] != nil {
			board.MakeMove(ponderMoves[currentTurn])
			GameHistory[MoveCount] = *ponderMoves[currentTurn] //record game history
		} else { //otherwise get new move
			move = Players[currentTurn].GetBestMove(board, depth)
			board.MakeMove(move)
			GameHistory[MoveCount] = *move
		}
		ponderMoves[currentTurn] = nil      //reinitialise ponderMove to nil
		currentTurn = (currentTurn + 1) % 2 //change turn
		MoveCount++
		close(stop)
	}
}

//cleanUp simply writes a memory profile and/or stops CPU profiling, provided 
//they have been requested
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

type HumanPlayer struct{}

// ai.Player interface conformation
// HumanPlayers don't need to ponder programmatically
func (h *HumanPlayer) PonderUntilInput(_ *gen.Board, _ chan struct{}, _ *gen.Move) {
	//humans ponder in their heads
}

//ai.Player interface conformation
//Human players don't need to debug themselves
func (h *HumanPlayer) Debug() {
	//humans know what they're thinking
}

//HumanPlayer.GetBestMove implementing ai.Player.GetBestMove. Handles commands
//and returns a move parsed from the command line. Will loop until a valid move 
//is given 
func (h *HumanPlayer) GetBestMove(b *gen.Board, _ int) *gen.Move {
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
				fmt.Print(m, ", ")
			}
			fmt.Println()
		case "undo":
			if MoveCount == 0 {
				break
			}
			MoveCount -= 1
			b.UnmakeMove(&GameHistory[MoveCount])
		case "debug":
			if *twoPlayer {
				continue
			} else {
				Players[1].Debug()
			}
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
			if err != nil {
				fmt.Println(err)
			} else {
				if move.IsLegalMove(b) {
					return move
				} else {
					fmt.Println("Invalid move for this position")
				}
			}
		}
	}
	return &move
}

//help prints all of the commands and their respective details
func help() {
	fmt.Println("help     \tThis help text")
	fmt.Println("[move]   \tMake a move on the board")
	fmt.Println("generate \tGenerate all moves for the current position")
	fmt.Println("undo     \tUndo the last played move")
	fmt.Println("debug    \tPrint debug information for the last search")
	fmt.Println("eval     \tPrint evaluation information for the current position")
	fmt.Println("print    \tPrint the board")
	fmt.Println("exit/quit\tQuit the game")
}
