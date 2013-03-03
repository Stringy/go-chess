package pieces

import (
	"math"
)

type Bishop struct {
	Col int 
	Moved bool
}

func (b Bishop) IsValidMove(x, y, xx, yy int, pieces []Piece) bool {
	dest := pieces[yy * 8 + xx];

	var valid = false

	if &b == nil {
		return valid
	} else if dest != nil && dest.GetCol() != b.GetCol() {
		return valid
	}

	if !IsDiagonalMove(x, y, xx, yy) {
		valid = false
	} else {
		if x < xx && y < yy { // + +
			for tx, ty := x+1, y+1; tx < xx && ty < yy; tx, ty = tx+1, ty+1 {
				if pieces[ty * 8 + tx] != nil {
					valid = false
					break
				}
			}
		} else if x > xx && y < yy { // - +
			for tx, ty := x-1, y+1; tx > xx && ty < yy; tx, ty = tx-1, ty+1 {
				if pieces[ty * 8 + tx] != nil {
					valid = false
					break
				}
			}
		} else if x > xx && y > yy { // - - 
			for tx, ty := x-1, y-1; tx > xx && ty > yy; tx, ty = tx-1, ty-1 {
				if pieces[ty * 8 + tx] != nil {
					valid = false
					break
				}
			}
		} else if x < xx && y > yy { // + -
			for tx, ty := x+1, y-1; tx < xx && ty > yy; tx, ty = tx+1, ty-1 {
				if pieces[ty * 8 + tx] != nil {
					valid = false
					break
				}
			}
		} else {
			valid = true
		}
	}
	return valid
}

func IsDiagonalMove(x, y, xx, yy int) bool {
	if math.Abs(float64(xx - x)) == math.Abs(float64(yy - y)) {
		return true;
	}
	return false
}

func (b Bishop) ToStr() string {
	if b.Col == BLACK {
		return " b "
	} else if b.Col == WHITE {
		return " B "
	}
	return " E "
}

func (b Bishop) GetCol() int {
	return b.Col
}

func (b Bishop) GetType() int {
	return BISHOP
}

func (b Bishop) HasMoved() bool {
	return b.Moved
}

func (b Bishop) Move() {
	b.Moved = true
}