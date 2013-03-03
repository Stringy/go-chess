package main

import (
	"flag"
	"unsafe"
	"fmt"
	"strings"
	"bufio"
	"strconv"
	"os"
	"time"
	"runtime"
	"runtime/pprof"

//	"chess/ai"
	"chess/ai/gen"
)

var (
	timed = false

	gencmd movegen
	prntcmd moveprinter

	cpuprofile = flag.String("cpuprofile", "",  "write a cpu profile to file")
	memprofile = flag.String("memprofile", "", "write a memory profile to file")
)

type Command interface {
	exec(timed bool)
}

type moveprinter struct {
	start int
	end int
}

func (m *moveprinter) exec(timed bool) {
	t0 := time.Now()

	if m.start < len(gencmd.moves) && m.end <= len(gencmd.moves) {
		for i := m.start; i < m.end; i++ {
			fmt.Print(gencmd.moves[i].Print())
		}
	} else {
		fmt.Println("Invalid Start or End")
	}

	if timed {
		fmt.Println("Time taken to print the moves:", time.Now().Sub(t0))
		fmt.Println("Why on earth did you want to time this?")
	}
}

type movegen struct {
	moves []gen.Move
	plys []int
	depth int
}

func (m *movegen) exec(timed bool) {
	fmt.Println("Generating Moves...")
	t0 := time.Now()

	var board gen.Board
	board.Init()

	m.moves = gen.Generate(m.depth, &board)

	fmt.Println("Generated", len(m.moves), "moves")
	fmt.Printf("Size of generated data: %v bytes\n", int(len(m.moves) * int(unsafe.Sizeof(gen.Move { 0 }))))
	
	if timed {
		t1 := time.Now().Sub(t0)
		fmt.Printf("Time taken: %v\n", t1)
	}
}

func main() {
	fmt.Println("Chess Engine Interactive Test Bed")
	fmt.Println("type \"help\" for command list")

	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	
	reader := bufio.NewReader(os.Stdin)
	gen.InitAllData()
		
	var cmd string
	for {
		cmd = ""
		fmt.Print(">>> ")
		cmd, _= reader.ReadString('\n')

		parts := strings.FieldsFunc(cmd, func(r rune) bool {
			return r == ' ' || r == '\n' || r == ':'
		})
		
		if parts[0] == "time" {
			timed = true
			parts = parts[1:]
		} else {
			timed = false
		}

		switch parts[0] {
		case "movegen":
			if len(parts) == 1 { //movegen for single state
				fmt.Println("Move generation of a single state\n")
				gencmd = movegen { make([]gen.Move, 0), make([]int, 0), 0 }
				gencmd.exec(timed)
			} else if len(parts) == 2 { //movegen plus depth
				if depth, err := strconv.Atoi(parts[1]); err == nil {
					gencmd = movegen { make([]gen.Move, 0), make([]int, 0), depth }
					fmt.Println("Move Gen to depth", depth)
					gencmd.exec(timed)
				} else {
					fmt.Println("Invalid movegen depth")
				}
			} else {
				fmt.Println("Unknown number of args to movegen command")
			}
		case "print":
			switch len(parts) {
			case 1: 
				fmt.Println("printing all generated moves")
				prntcmd = moveprinter { 0, len(gencmd.moves) }
				prntcmd.exec(timed)
			case 2:
				if begin, err := strconv.Atoi(parts[1]); err == nil {
					prntcmd = moveprinter { begin, len(gencmd.moves) }
					prntcmd.exec(timed)
				} else {
					fmt.Println("Invalid Number Arg")
				}
			case 3:
				begin, err := strconv.Atoi(parts[1])
				end, err2 := strconv.Atoi(parts[2])
				if err != nil || err2 != nil {
					fmt.Println("Invalid Number Args") 
				} else {
					prntcmd = moveprinter { begin, end }
					prntcmd.exec(timed)
				}
			default: fmt.Println("Unknown number of args to print command")
			}
		case "dump":
			if len(parts) < 2 {
				dump("test.dump")
			} else {
				dump(parts[1])
			}
		case "clear": clear()
		case "help": help()
		case "exit": fallthrough  
		case "quit":
			clear()
			fmt.Println("Quitting")
			goto EXIT
		default: 
			fmt.Println("Unknown Command")
		}
	}	
	EXIT:
   if *memprofile != "" {
      f, err := os.Create(*memprofile)
      if err != nil {
         panic(err)
      }
      pprof.WriteHeapProfile(f)
      f.Close()
      return
   }
}

func help() {
	fmt.Println("Commands:\n")
	//add commands
	fmt.Println("   help")
	fmt.Println("   movegen [depth]\n\tGenerate moves to the specified depth")
	fmt.Println("   print [[begin] : [end]] or [[begin] :] \n\tPrint the generated moves")
	fmt.Println("   time [cmd [args]] \n\t time a command")
	fmt.Println("   dump [filename]\n\tWrite the moves to file")
}

func clear() {
	//empty command structures of memory
	gencmd = movegen { make([]gen.Move, 0), make([]int, 0), 0 }
	prntcmd = moveprinter { 0, 0 }
}

func dump(file string) {
	out, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	var numbytes = 0
	for _, move := range gencmd.moves {
		a := byte(move.Val & 0xff)
		b := byte(move.Val & 0x00ff)
		c := byte(move.Val & 0x0000ff)
		d := byte(move.Val & 0x000000ff)
		arr := []byte { a, b, c, d, }
		numbytes += 4
		out.Write(arr)
	}
	fmt.Println("Wrote", numbytes, "bytes to", file)
}