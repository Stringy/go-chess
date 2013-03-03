package pieces

import (

)

const (
	BLACK int = 100
	WHITE int = 200

	PAWN int = 10
	ROOK int = 20
	KNIGHT int = 30
	BISHOP int = 40
	QUEEN int = 50
	KING int = 60
)

type Piece interface {
	IsValidMove(x, y, xx, yy int, pieces []Piece) bool
	GetCol() int
	GetType() int
	HasMoved() bool
	ToStr() string
	Move()
}
