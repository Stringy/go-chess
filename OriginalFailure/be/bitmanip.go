package be

import (
	"fmt"
)



const (
	Emp = uint64(0)
	Universe = uint64(0xffffffffffffffff)
)
/*
 Calculates the bitcount of n
 */
func BitCount(n uint64) uint {
	//Constants for MIT Hakmem algorithm
	M1 := uint64(0x5555555555555555)
	M2 := uint64(0x3333333333333333)
	M4 := uint64(0x0f0f0f0f0f0f0f0f)
	M8 := uint64(0x00ff00ff00ff00ff)
	M16 := uint64(0x0000ffff0000ffff)
	M32 := uint64(0x00000000ffffffff)

	n = (n & M1) + ((n >> 1) & M1)
	n = (n & M2) + ((n >> 2) & M2)
	n = (n & M4) + ((n >> 4) & M4)
	n = (n & M8) + ((n >> 8) & M8)
	n = (n & M16) + ((n >> 16) & M16)
	n = (n & M32) + ((n >> 32) & M32)
	
	return uint(n)
}

func FirstOne(n uint64) uint {
	Index64 := []int {
		63,  0, 58,  1, 59, 47, 53,  2,
      60, 39, 48, 27, 54, 33, 42,  3,
      61, 51, 37, 40, 49, 18, 28, 20,
      55, 30, 34, 11, 43, 14, 22,  4,
      62, 57, 46, 52, 38, 26, 32, 41,
      50, 36, 17, 19, 29, 10, 13, 21,
      56, 45, 25, 31, 35, 16,  9, 12,
      44, 24, 15,  8, 23,  7,  6,  5,
	}
	
	Debruijn := uint64(0x07edd5e59a4e28c2)
	return uint(Index64[((n & -n) * Debruijn) >> 58])
}

func LastOne(n uint64) (result uint) {
	result = 0
	if n > 0xffffffff {
		n >>= 32
		result = 32
	}
	if n > 0xffff {
		n >>= 16
		result += 16
	}
	if n > 0xff {
		n >>= 8
		result += 8
	}
	result += uint(MsbTable[n])
	return
}

func PrintBitboard(n uint64) {
	grid := make([]byte, 64)
	for i := 0; i < 64; i++ {
		if n & BitSet[i] != 0 {
			grid[i] = '1'
		} else {
			grid[i] = '.'
		}
	}
	
	for rank := 8; rank > 0; rank-- {
		fmt.Print(" ", rank, " ")
		for file := 1; file < 9; file++ {
			fmt.Print(string(grid[Index[file][rank]]))
		}
		fmt.Println()
	}
	fmt.Println("   abcdefgh")
}