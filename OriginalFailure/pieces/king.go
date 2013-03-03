package pieces

import (
	"math"
)

type King struct {
	Col int
	Moved bool
}

func (k King) IsValidMove(x, y, xx, yy int, pieces []Piece) bool {
	dest := pieces[yy * 8 + xx]

	var valid = false

	if &k == nil {
		return valid
	} else if dest != nil && dest.GetCol() == k.GetCol() {
		return valid
	}

	if !(math.Abs(float64(xx - x)) == 0 || math.Abs(float64(xx - x)) == 1) { //x too high 
		valid = false
	} else if !(math.Abs(float64(yy - y)) == 0 || math.Abs(float64(yy - y)) == 1) { // y too high
		valid = false
	} else {
		var rook = Rook { Col : k.GetCol(), Moved : k.HasMoved() }
		var bish = Bishop { Col : k.GetCol(), Moved : k.HasMoved() }
		valid = rook.IsValidMove(x, y, xx, yy, pieces)
		valid = valid || bish.IsValidMove(x, y, xx, yy, pieces)
	}
	return valid
}

func (k King) ToStr() string {
	if k.Col == BLACK {
		return " k "
	} else if k.Col == WHITE {
		return " K "
	}
	return " E "
}

func (k King) GetCol() int {
	return k.Col
}

func (k King) GetType() int {
	return KING
}

func (k King) HasMoved() bool {
	return k.Moved
}

func (k King) Move() {
	k.Moved = true
}