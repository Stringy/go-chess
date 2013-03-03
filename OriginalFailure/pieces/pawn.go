package pieces

import (

)

type Pawn struct {
	Col int
	Moved bool
}

func (p Pawn) IsValidMove(x, y, xx, yy int, pieces []Piece) bool {
	var dest = pieces[yy * 8 + xx]

	var valid = false

	if &p == nil {
		return valid
	}

	switch p.GetCol() {
	case WHITE:
		if (y - 1) == yy { //moving forward only
			if dest == nil {
				if x == xx {
					valid = true
				}
			} else if dest.GetCol() == BLACK {
				if (x + 1) == xx || (x - 1) == xx {
					valid = true
				}
			}
		} else if !p.HasMoved() && (y - 2) == yy { 
			if dest == nil {
				if x == xx {
					valid = true
				}
			}
		}
		break
	case BLACK:
		if (y + 1) == yy { //moving forward only
			if dest == nil {
				if x == xx {
					valid = true
				}
			} else if dest.GetCol() == WHITE {
				if (x + 1) == xx || (x - 1) == xx {
					valid = true
				}
			}
		} else if !p.HasMoved() && (y + 2) == yy {
			if dest == nil {
				if x == xx {
					valid = true
				}
			}
		}
		break
	default:
		valid = false
	}
	return valid
}

func (p Pawn) ToStr() string {
	if p.Col == BLACK {
		return " p "
	} else if p.Col == WHITE {
		return " P "
	}
	return " E "
}

func (p Pawn) GetCol() int {
	return p.Col
}

func (p Pawn) GetType() int {
	return PAWN
}

func (p Pawn) HasMoved() bool {
	return p.Moved
}

func (p Pawn) Move() {
	p.Moved = true
}