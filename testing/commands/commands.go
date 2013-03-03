package commands

import (
	"chess/ai"
	"chess/ai/eval"
	"chess/ai/gen"
	"chess/ai/search"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Command interface {
	Exec(args ...string) error
}

/*
 Move Generation Command
*/
type MoveGenCmd struct{}

func (m MoveGenCmd) Exec(args ...string) error {
	//generate moves
	t0 := time.Now()
	var board gen.Board
	board.Init()

	var moves uint64

	switch len(args) {
	case 0:
		moves = uint64(len(gen.GenerateAllMoves(&board)))
	case 1:
		if depth, err := strconv.Atoi(args[0]); err == nil {
			moves = gen.Perft(&board, depth)
		} else {
			return err
		}
	default:
		return errors.New("Unknown args passed to movegen")
	}
	fmt.Printf("Time taken to generate %d moves: %v\n", moves, time.Now().Sub(t0))
	return nil
}

/*
 Print Command
*/
type MovePrinterCmd struct {
	moves []gen.Move
}

func (m MovePrinterCmd) Exec(args ...string) error {
	//print moves
	for _, move := range m.moves {
		fmt.Print(move.Print())
	}
	return nil
}

/*
 Evaluation Command
*/
type EvalCmd struct {
}

func (e EvalCmd) Exec(args ...string) error {
	t0 := time.Now()

	var board gen.Board
	board.Init()

	var move *gen.Move
	heur := new(eval.StaticHeuristic)

	if len(args) == 0 {
		move = ai.GetBestMove(board, search.AlphaBeta{heur}, 1)
	} else {
		if depth, err := strconv.Atoi(args[0]); err == nil {
			move = ai.GetBestMove(board, search.AlphaBeta{heur}, depth)
		} else {
			return err
		}
	}

	fmt.Printf("Time take to evaluate search space: %v\n", time.Now().Sub(t0))
	fmt.Println("Best Move found:", move)
	return nil
}

/*
 Help Command
*/
type HelpCmd struct {
}

func (h HelpCmd) Exec(args ...string) error {
	fmt.Println("Command List\n")
	fmt.Println("\n\tmovegen [depth]")
	fmt.Println("\teval [depth]")
	return nil
}

/*
 Exit cmd
*/
type ExitCmd struct{}

func (e ExitCmd) Exec(args ...string) error {
	//clear up 
	fmt.Println("Exitting")
	os.Exit(0)
	return nil
}
func GetCommand(s string) (Command, error) {
	var cmd Command

	switch s {
	case "movegen":
		cmd = MoveGenCmd{}
	case "eval":
		cmd = EvalCmd{}
	case "help":
		cmd = HelpCmd{}
	case "exit":
		fallthrough
	case "quit":
		cmd = ExitCmd{}
	default:
		return nil, errors.New("Invalid Command")
	}
	return cmd, nil
}
