package ai

import (
	"chess/ai/gen"
	"testing"
	"time"
)

func TestPonderUntilInput(t *testing.T) {
	stop := make(chan struct{})
	var board gen.Board
	board.Init()

	go PonderUntilInput(&board, stop)
	for {
		select {
		case <-time.After(2 * time.Minute):
			close(stop)
			PrintDebug()
		default:
		}
	}
}
