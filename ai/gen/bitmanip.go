package gen

import (
	"fmt"
)

var (
	Index64 = [64]uint{
		63, 0, 58, 1, 59, 47, 53, 2,
		60, 39, 48, 27, 54, 33, 42, 3,
		61, 51, 37, 40, 49, 18, 28, 20,
		55, 30, 34, 11, 43, 14, 22, 4,
		62, 57, 46, 52, 38, 26, 32, 41,
		50, 36, 17, 19, 29, 10, 13, 21,
		56, 45, 25, 31, 35, 16, 9, 12,
		44, 24, 15, 8, 23, 7, 6, 5,
	}
)

const (
	//Empty bitboard for various calculations
	Emp = uint64(0)
	//universal bitboard for calculating unoccupied squares
	Universe = uint64(0xffffffffffffffff)
)

/*
BitCount calculates the bitcount of n
*/
func BitCount(n uint64) uint {
	//Constants for MIT Hakmem algorithm
	const M1 = uint64(0x5555555555555555)
	const M2 = uint64(0x3333333333333333)
	const M4 = uint64(0x0f0f0f0f0f0f0f0f)
	const M8 = uint64(0x00ff00ff00ff00ff)
	const M16 = uint64(0x0000ffff0000ffff)
	const M32 = uint64(0x00000000ffffffff)

	n = (n & M1) + ((n >> 1) & M1)
	n = (n & M2) + ((n >> 2) & M2)
	n = (n & M4) + ((n >> 4) & M4)
	n = (n & M8) + ((n >> 8) & M8)
	n = (n & M16) + ((n >> 16) & M16)
	n = (n & M32) + ((n >> 32) & M32)

	return uint(n)
}

//LSB returns the index of the first 1 bit in the bitboard
func LSB(n uint64) uint {
	const Debruijn = uint64(0x07edd5e59a4e28c2)
	return Index64[((n&-n)*Debruijn)>>58]
}

//MSB returns the index of the last 1 bit in the bitboard
func MSB(n uint64) (result uint) {
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

//PrintBitboard prints a formatted bitstring of the bitboard
func BitboardToStr(n uint64) string {
	grid := make([]byte, 64)
	str := ""
	for i := 0; i < 64; i++ {
		if n&BitSet[i] != 0 {
			grid[i] = '1'
		} else {
			grid[i] = '.'
		}
	}

	for rank := 8; rank > 0; rank-- {
		str += fmt.Sprintf(" %d ", rank)
		for file := 1; file < 9; file++ {
			str += string(grid[Index[file][rank]])
		}
		str += "\n"
	}
	str += "   abcdefgh\n"
	return str
}
