package main

import (
	"chess/ai/gen"
	"chess/ai"
	"chess/ai/search"
	"strings"
	"os"
	"time"
	"errors"
	"fmt"
)

type Command interface {
	Exec()
}

/*
 Move Generation Command
*/
type MoveGenCmd struct {
	args []string
}

func (m MoveGenCmd) Exec() {
	//generate moves
	t0 := time.Now()	
	var board gen.Board
	board.Init()
	m.moves = gen.Generate(m.depth, &board)
	fmt.Println("Best move found:", move)
	fmt.Printf("Time taken to generate %d moves: %v", len(m.moves), time.Now().Sub(t0))
}

/*
 Print Command
*/
type MovePrinterCmd struct {
	args []string
	moves []gen.Move
}

func (m MovePrinterCmd) Exec(args...string) {
	//print moves
	for _, move := range m.moves {
		fmt.Print(move.Print())
	}
}

/*
 Evaluation Command
*/
type EvalCmd struct {
	args []string
}

func (e EvalCmd) Exec() { 
	t0 := time.Now()
	
	var board gen.Board
	board.Init()

	move := ai.GetBestMove(board, e.algo, e.depth)

	fmt.Printf("Time take to evaluate search space to depth %d: %v", depth, time.Now().Sub(t0))
	fmt.Println("Best Move found:", move)
}

/*
 Help Command
*/
type HelpCmd struct {
}

func (h HelpCmd) Exec() {
	fmt.Println("Command List\n")
}

/*
 Exit cmd
*/
type ExitCmd struct {} 

func (e ExitCmd) Exec() {
	//clear up 
	fmt.Println("Exitting")
	os.Exit(0)
}
func ParseCommandLine(s string) (Command, error) {
	var cmd Command

	parts := strings.FieldsFunc(cmd, func(r rune) bool {
		return r == ' ' || r == '\n' || r == ':'
	}
	
	switch parts[0] {
	case "movegen":
		return MoveGenCmd { parts[1:] }, nil
	case "eval":
		return EvalCmd { parts[1:] }, nil
	case "help":
		return HelpCmd {}, nil
	case "exit" || "quit":
		return ExitCmd{}, nil
	}
}