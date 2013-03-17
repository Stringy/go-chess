package main

import (
	"chess/ai/gen"
	"fmt"
	"os"
)

func main() {
	file, _ := os.Create("diagA8H1Magics.txt")
	for i, magic := range gen.GeneralSlidingAttacks {
		hex := fmt.Sprintf("%d: 0x%x\n", i, magic)
		file.Write([]byte(hex))
		file.Write([]byte(gen.BitboardToStr(magic)))
	}
}
