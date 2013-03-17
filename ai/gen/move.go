package gen

import (
	"fmt"
	//	"strconv"
	"errors"
	"strings"
)

// Move Stores information in max 32 bits
// Information needed to be stored:
//    - From index (0.. 63)
//    - To Index (0.. 63)
//    - Piece (0.. 15)
//    - Caputured Piece (0.. 15)
//    - Promotion (0.. 15)
//
//    from = 6 bits
//    to = 6 bits
//    piece = 4 bits
//    captured = 4 bits
//    promotion = 4 bits
type Move uint

//IsLegalMove returns true if the generator has generated this move,
//false otherwise. This is a simple way of determining whether the move
//is allowed
func (m *Move) IsLegalMove(b *Board) bool {
	moves := GenerateAllMoves(b)
	for _, move := range moves {
		if move == *m {
			return true
		}
	}
	return false
}

func NewMove(cmd string, b *Board) (*Move, error) {
	//	ranks := "abcdefgh"
	var move Move
	if len(cmd) > 5 || len(cmd) < 4 {
		return nil, errors.New("Unimplemented or unknown move string")
	}
	move.Clear()
	cmd = strings.ToUpper(cmd)
	if to, ok := SquareMap[cmd[2:4]]; !ok {
		return nil, errors.New("Unknown square for to coordinate")
	} else {
		move.SetTo(to)
		if from, ok := SquareMap[cmd[:2]]; !ok {
			return nil, errors.New("Unknown square for from coordinate")
		} else {
			move.SetFrom(from)
		}
	}
	if len(cmd) == 5 {
		prom := cmd[4]
		move.SetProm(PieceMap[prom])
	}
	move.SetPiece(b.Squares[move.GetFrom()])
	move.SetCapt(b.Squares[move.GetTo()])
	//check capture / promotion etc
	return &move, nil
}

//Move.Clear returns the move to a zero state
func (m *Move) Clear() {
	*m = Move(0)
}

//*Move.SetFrom set first 6 bits, as the from index
func (m *Move) SetFrom(from byte) {
	*m &= Move(0xffffffc0)        //clear first 6 bits
	*m |= Move(from & 0x0000003f) //set to from bits
}

//*Move.SetTo set second 6 bits as the to index
func (m *Move) SetTo(to byte) {
	*m &= Move(0xfffff03f)
	*m |= Move(to&0x0000003f) << 6
}

//*Move.SetPiece bits 12-15, the piece being moved
func (m *Move) SetPiece(piece byte) {
	*m &= Move(0xffff0fff)
	*m |= Move(piece&0x0000000f) << 12
}

//*Move.SetCapt bits 16-19, the piece being captured
func (m *Move) SetCapt(capt byte) {
	*m &= Move(0xfff0ffff)
	*m |= Move(capt&0x0000000f) << 16
}

//*Move.SetProm bits 20-23, the piece type a pawn is being promoted to
//usually 0000
func (m *Move) SetProm(prom byte) {
	*m &= Move(0xff0fffff)
	*m |= Move(prom&0x0000000f) << 20
}

//*Move.GetFrom return the from index
func (m *Move) GetFrom() byte {
	return byte(*m & 0x0000003f)
}

//*Move.GetTo returns the to index
func (m *Move) GetTo() byte {
	return byte((*m >> 6) & 0x0000003f)
}

//*Move.GetPiece returns the piece type being moved
func (m *Move) GetPiece() byte {
	return byte((*m >> 12) & 0x0000000f)
}

//*Move.GetCapt returns the piece which has been captured
func (m *Move) GetCapt() byte {
	return byte((*m >> 16) & 0x0000000f)
}

//*Move.GetProm returns the promotion piece 
func (m *Move) GetProm() byte {
	return byte((*m >> 20) & 0x0000000f)
}

//*Move.IsWhitemove returns true if it is a white move, false otherwise
func (m *Move) IsWhitemove() bool { // piec is white: bit 15 must be 0
	return (0x00008000 &^ *m) == 0x00008000
}

//*Move.IsBlackmove returns true if it is a black move, false otherwise
func (m *Move) IsBlackmove() bool { // piec is black: bit 15 must be 1
	return (*m & 0x00008000) == 0x00008000
}

//*Move.IsCapture returns true if capture bits are non-zero, false otherwise
func (m *Move) IsCapture() bool { // capt is nonzero, bits 16 to 19 must be nonzero
	return (*m & 0x000f0000) != 0x00000000
}

//*Move.IsKingcaptured returns true if capture bits are king, false otherwise
func (m *Move) IsKingcaptured() bool { // bits 17 to 19 must be 010
	return (*m & 0x00070000) == 0x00020000
}

//*Move.IsRookmove returns true if piece bits are rook, false otherwise 
func (m *Move) IsRookmove() bool { // bits 13 to 15 must be 110
	return (*m & 0x00007000) == 0x00006000
}

//*Move.IsRookcaptured returns true if capture bits are rook, false otherwise 
func (m *Move) IsRookcaptured() bool { // bits 17 to 19 must be 110
	return (*m & 0x00070000) == 0x00060000
}

//*Move.IsKingmove returns true if piece bits are king, false otherwise
func (m *Move) IsKingmove() bool { // bits 13 to 15 must be 010
	return (*m & 0x00007000) == 0x00002000
}

//*Move.IsPawnmove returns true if piece bits are pawn, false otherwise
func (m *Move) IsPawnmove() bool { // bits 13 to 15 must be 001
	return (*m & 0x00007000) == 0x00001000
}

//*Move.IsPawnDoublemove returns true if pawn double move, false otherwise
//     bits 13 to 15 must be 001 &
//     bits 4 to 6 must be 001 (from rank 2) & bits 10 to 12 must be 011 (to rank 4)
//     OR: bits 4 to 6 must be 110 (from rank 7) & bits 10 to 12 must be 100 (to rank 5)

func (m *Move) IsPawnDoublemove() bool {
	return (((*m & 0x00007000) == 0x00001000) &&
		((((*m & 0x00000038) == 0x00000008) &&
			((*m & 0x00000e00) == 0x00000600)) ||
			(((*m & 0x00000038) == 0x00000030) && ((*m & 0x00000e00) == 0x00000800))))
}

//*Move.IsEnpassant returns true if enpassant move, false otherwise 
func (m *Move) IsEnpassant() bool { // prom is a pawn, bits 21 to 23 must be 001
	return (*m & 0x00700000) == 0x00100000
}

//*Move.IsPromotion returns true if promotion bits are non-zero, false otherwise
func (m *Move) IsPromotion() bool { // prom (with color bit removed), .xxx > 2 (not king or pawn)
	return (*m & 0x00700000) > 0x00200000
}

//*Move.IsCastle returns  true if castling move, false otherwise
func (m *Move) IsCastle() bool { // prom is a king, bits 21 to 23 must be 010
	return (*m & 0x00700000) == 0x00200000
}

//*Move.IsCastleOO returns true if two square castle move, false otherwise 
func (m *Move) IsCastleOO() bool { // prom is a king and tosq is on the g-file
	return (*m & 0x007001c0) == 0x00200180
}

//*Move.IsCastleOOO returns true if three square castle move, false otherwise
func (m *Move) IsCastleOOO() bool { // prom is a king and tosq is on the c-file
	return (*m & 0x007001c0) == 0x00200080
}

//*Move.Print returns a string representation of the move
func (m *Move) Print() (ret string) {
	to := m.GetTo()
	from := m.GetFrom()

	ret = fmt.Sprintf("From: %s%d To: %s%d Piece: %s Capt: %s Prom: %s\n",
		string(from%8+'A'), (from/8 + 1),
		string(to%8+'A'), (to/8 + 1),
		PieceNames[m.GetPiece()],
		PieceNames[m.GetCapt()],
		PieceNames[m.GetProm()])
	return
}

func (m *Move) String() string {
	from := m.GetFrom()
	to := m.GetTo()
	return fmt.Sprintf("%s%d%s%d",
		string(from%8+'a'), (from/8 + 1),
		string(to%8+'a'), (to/8 + 1))
}
