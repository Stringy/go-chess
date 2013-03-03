package main 

import (
	"runtime"
	"bufio"
	"fmt"
	"os"
	"strings"
	"chess/testing/commands"
)

func main() {
	fmt.Println("Chess Engine Interactive Test Bed")
	fmt.Println("type \"help\" for command list")

	runtime.GOMAXPROCS(runtime.NumCPU())
	
	reader := bufio.NewReader(os.Stdin)
	var line string
	for {
		line = ""
		fmt.Print(">>> ")
		line, _ = reader.ReadString('\n')
		
		parts := strings.FieldsFunc(line, func(r rune) bool {
			return r == ' ' || r == '\n' || r == ':'
		})

		if len(parts) == 0 {
			fmt.Println("Enter a Command")
		}

		cmd, err := commands.GetCommand(parts[0])
		if err != nil {
			fmt.Println(err)
		} else {
			cmd.Exec(parts[1:]...)
		}
	}
}