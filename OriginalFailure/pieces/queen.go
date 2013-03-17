package pieces

type Queen struct {
	Col int
	Moved bool
}

func (q Queen) IsValidMove(x, y, xx, yy int, pieces []Piece) bool {
	var rook = Rook { Col : q.GetCol(), Moved : q.HasMoved() }
	var bishop = Bishop { Col : q.GetCol(), Moved : q.HasMoved() }

	return rook.IsValidMove(x, y, xx, yy, pieces) || bishop.IsValidMove(x, y, xx, yy, pieces);
}

func (q Queen) ToStr() string {
	if q.Col == BLACK {
		return " q "
	} else if q.Col == WHITE {
		return " Q "
	}
	return " E "
}

func (q Queen) GetCol() int {
	return q.Col
}

func (q Queen) GetType() int {
	return QUEEN
}

func (q Queen) HasMoved() bool {
	return q.Moved
}

func (q Queen) Move() {
	q.Moved = true
}