package be

import (
	"fmt"
)

/*
 Move Stores information in max 32 bits
 
 Information needed to be stored:
 - From index (0.. 63)
 - To Index (0.. 63)
 - Piece (0.. 15)
 - Caputured Piece (0.. 15)
 - Promotion (0.. 15)

 from = 6 bits
 to = 6 bits
 piece = 4 bits
 captured = 4 bits
 promotion = 4 bits

*/
type Move struct {
	Val uint
}

func (m *Move) Clear() {
	m.Val = 0
}

//set first 6 bits
func (m *Move) SetFrom(from byte) {
	m.Val &= 0xffffffc0; //clear first 6 bits
	m.Val |= uint(from & 0x0000003f) //set to from bits
}

//set second 6 bits 
func (m *Move) SetTo(to byte) {
	m.Val &= 0xfffff03f
	m.Val |= uint(to & 0x0000003f) << 6
}

// bits 12-15
func (m *Move) SetPiece(piece byte) {
	m.Val &= 0xffff0fff
	m.Val |= uint(piece & 0x0000000f) << 12
}

//bits 16-19
func (m *Move) SetCapt(capt byte) {
	m.Val &= uint(0xfff0ffff)
	m.Val |= uint(capt & 0x0000000f) << 16
}

//bits 20-23
func (m *Move) SetProm(prom byte) {
	m.Val &= uint(0xff0fffff)
	m.Val |= uint(prom & 0x0000000f) << 20
}


//Getters
func (m *Move) GetFrom() byte {
	return byte(m.Val & 0x0000003f)
}

func (m *Move) GetTo() byte {
	return byte((m.Val >> 6) & 0x0000003f)
}

func (m *Move) GetPiece() byte {
	return byte((m.Val >> 12) & 0x0000000f)
}

func (m *Move) GetCapt() byte {
	return byte((m.Val >> 16) & 0x0000000f)
}

func (m *Move) GetProm() byte {
	return byte((m.Val >> 20) & 0x0000000f)
}


/*
 Bool Move Checks
*/
func (m *Move) IsWhitemove() bool {   // piec is white: bit 15 must be 0
       return (0x00008000 &^ m.Val) == 0x00008000;
} 
 
func (m *Move) IsBlackmove() bool {   // piec is black: bit 15 must be 1
       return (m.Val & 0x00008000) == 0x00008000;
} 
 
func (m *Move) IsCapture() bool {   // capt is nonzero, bits 16 to 19 must be nonzero
       return (m.Val & 0x000f0000) != 0x00000000;
} 
 
func (m *Move) IsKingcaptured() bool {   // bits 17 to 19 must be 010
       return (m.Val & 0x00070000) == 0x00020000;
} 
 
func (m *Move) IsRookmove() bool {   // bits 13 to 15 must be 110
       return (m.Val & 0x00007000) == 0x00006000;
} 
 
func (m *Move) IsRookcaptured() bool {   // bits 17 to 19 must be 110
       return (m.Val & 0x00070000) == 0x00060000;
} 
 
func (m *Move) IsKingmove() bool {   // bits 13 to 15 must be 010
       return (m.Val & 0x00007000) == 0x00002000;
} 
 
func (m *Move) IsPawnmove() bool {   // bits 13 to 15 must be 001
       return (m.Val & 0x00007000) == 0x00001000;
} 
 
func (m *Move) IsPawnDoublemove() bool {
   // bits 13 to 15 must be 001 &
   //     bits 4 to 6 must be 001 (from rank 2) & bits 10 to 12 must be 011 (to rank 4)
   // OR: bits 4 to 6 must be 110 (from rank 7) & bits 10 to 12 must be 100 (to rank 5)
	
   return ((( m.Val & 0x00007000) == 0x00001000) &&
		(((( m.Val & 0x00000038) == 0x00000008) &&
		((( m.Val & 0x00000e00) == 0x00000600))) || 
      ((( m.Val & 0x00000038) == 0x00000030) && ((( m.Val & 0x00000e00) == 0x00000800)))));
} 
 
func (m *Move) IsEnpassant() bool {   // prom is a pawn, bits 21 to 23 must be 001
       return (m.Val & 0x00700000) == 0x00100000;
} 
 
func (m *Move) IsPromotion() bool {   // prom (with color bit removed), .xxx > 2 (not king or pawn)
       return (m.Val & 0x00700000) >  0x00200000;
} 
 
func (m *Move) IsCastle() bool {   // prom is a king, bits 21 to 23 must be 010
       return (m.Val & 0x00700000) == 0x00200000;
} 
 
func (m *Move) IsCastleOO() bool {   // prom is a king and tosq is on the g-file
       return (m.Val & 0x007001c0) == 0x00200180;
} 
 
func (m *Move) IsCastleOOO() bool {   // prom is a king and tosq is on the c-file
       return (m.Val & 0x007001c0) == 0x00200080;
} 

func (m *Move) Print() (ret string) {
	ret = fmt.Sprintf("To: %d From: %d Piece: %s Capt: %s\n", m.GetTo(), m.GetFrom(), PieceNames[m.GetPiece()], PieceNames[m.GetCapt()])
	return
}