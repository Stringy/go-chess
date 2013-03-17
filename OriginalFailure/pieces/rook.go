package pieces

type Rook struct {
	Col int
	Moved bool
}

func (r Rook) IsValidMove(x, y, xx, yy int, pieces []Piece) bool {
	var dest = pieces[yy * 8 + xx];
	
	var valid = false;

	if &r == nil {
		return valid
	} else if dest != nil && dest.GetCol() != r.GetCol() {
		return valid
	}

	if IsVerticalMove(x, y, xx, yy) {
		if y > yy { //backwards movement
			for ty := y - 1; ty > yy; ty-- {
				if pieces[ty * 8 + x] != nil {
					valid = false
					break
				}
			}
		} else if y < yy { //forward movement
			for ty := y + 1; ty < yy; ty++ {
				if pieces[ty * 8 + x] != nil {
					valid = false
					break
				}
			}
		}
	} else if IsHorizontalMove(x, y, xx, yy) {
		if x > xx { //left movement
			for tx := x - 1; tx > xx; tx-- {
				if pieces[y * 8 + tx] != nil {
					valid = false
					break
				}
			}
		} else if x < xx { //right movement
			for tx := x + 1; tx < xx; tx++ {
				if pieces[y * 8 + tx] != nil {
					valid = false
					break
				}
			}
		}
	}
	return valid
} 

func IsVerticalMove(x, y, xx, yy int) bool {
	if x != xx {
		return false
	} 
	return true
}

func IsHorizontalMove(x, y, xx, yy int) bool {
	if y != yy {
		return false
	}
	return true
}

func (r Rook) ToStr() string {
	if r.Col == BLACK {
		return " r "
	} else if r.Col == WHITE {
		return " R "
	}
	return " E "
}

func (r Rook) GetCol() int {
	return r.Col
}

func (r Rook) GetType() int {
	return ROOK
}

func (r Rook) HasMoved() bool {
	return r.Moved
}

func (r Rook) Move() {
	r.Moved = true
}