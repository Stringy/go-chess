package pieces

type Knight struct {
	Col int
	Moved bool
}

func (n Knight) IsValidMove(x, y, xx, yy int, pieces []Piece) bool {
	var dest = pieces[yy * 8 + xx];

	valid := false

	if &n == nil {
		return valid
	} else if dest != nil && dest.GetCol() == n.GetCol() {
		return valid
	}

	//No blocking, just check coordinates

	if (y + 2) == yy || (y - 2) == yy { // +/- y
		if (x + 1) == xx || (x - 1) == xx {
			valid = true
		} else {
			valid = false
		}
	} else if (x + 2) == xx || (x - 2) == xx { // +/- x
		if (y + 1) == yy || (y - 1) == yy {
			valid = true
		} else {
			valid = false
		}
	} else {
		valid = false
	}
	return valid
}

func (n Knight) ToStr() string {
	if n.Col == BLACK {
		return " n "
	} else if n.Col == WHITE {
		return " N "
	}
	return " E "
}

func (n Knight) GetCol() int {
	return n.Col
}

func (n Knight) GetType() int {
	return KNIGHT
}

func (n Knight) HasMoved() bool {
	return n.Moved
}

func (n Knight) Move() {
	n.Moved = true
}