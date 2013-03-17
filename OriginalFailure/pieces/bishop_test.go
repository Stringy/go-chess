package pieces

import (
	"testing"
	"fmt"
)

/*
 Tests Bishop Move validation
 
 Various situations are created and validation testing is conducted over them
*/
func TestBishopMoveValidation(t *testing.T) {
	var grid = make([]Piece, 64)

	/*
	  | t - t - t - - - |
	  | - - - - - - - - |
	  | t - B - t - - - |
	  | - - - - - - - - |
	  | t - t - t - - - |
	  | - - - - - Q - - |
	  | - - - - - - t - |
	  | - - - - - - - - |

	  B - Bishop
	  Q - Blocking Queen
	  t - tested move
	 */

	var x, y = 2, 2
	grid[y * 8 + x] = Bishop { WHITE, false }
	grid[5 * 8 + 5] = Queen { WHITE, false }

	var bish = grid[y * 8 + x]
	
	//first move 0, 0 expected valid
	var xx, yy = 0, 0
	if !bish.IsValidMove(x, y, xx, yy, grid) {
		fmt.Println("FAIL: Expected Valid Move for", xx, yy)
		t.Fail()
	}

	//second 2, 0 expected false
	xx, yy = 2, 0
	if bish.IsValidMove(x, y, xx, yy, grid) {
		fmt.Println("FAIL: Expected Invalid Move for", xx, yy)
		t.Fail()
	}

	//third 4, 0 expected true
	xx, yy = 4, 0
	if !bish.IsValidMove(x, y, xx, yy, grid) {
		fmt.Println("FAIL: Expected Valid Move for", xx, yy)
		t.Fail()
	}

	//fourth 0, 2 expected false
	xx, yy = 0, 2 
	if bish.IsValidMove(x, y, xx, yy, grid) {
		fmt.Println("FAIL: Expected Invalid Move for", xx, yy)
		t.Fail()
	}

	//fifth 0, 4 expected false
	xx, yy = 0, 4
	if bish.IsValidMove(x, y, xx, yy, grid) {
		fmt.Println("FAIL: Expected Invalid Move for", xx, yy)
		t.Fail()
	}

	//sixth 0, 6 expected true
	xx, yy = 0, 6
	if !bish.IsValidMove(x, y, xx, yy, grid) {
		fmt.Println("FAIL: Expected Valid Move for", xx, yy)
		t.Fail()
	}

	//seventh 2, 6 expected false
	xx, yy = 2, 6
	if bish.IsValidMove(x, y, xx, yy, grid) {
		fmt.Println("FAIL: Expected Invalid Move for", xx, yy)
		t.Fail()
	}

	//eighth 4, 6 expected true
	xx, yy = 4, 6 
	if !bish.IsValidMove(x, y, xx, yy, grid) {
		fmt.Println("FAIL: Expected Valid Move for", xx, yy)
		t.Fail()
	}

	//nineth 6, 6 expected false (blocked)
	xx, yy = 6, 6 
	if bish.IsValidMove(x, y, xx, yy, grid) {
		fmt.Println("FAIL: Expected Invalid Move for", xx, yy)
		t.Fail()
	}
}

/*
 Tests whether a set of coordinates are identified as diagonal or not
*/
func TestIsDiagonalMove(t *testing.T) {
	x := 0
	y := 0
	xx := 4
	yy := -4

	if !IsDiagonalMove(x, y, xx, yy) { //Expected True
		fmt.Println("Diagonal Move not recognised")
		t.Fail()
	}

	xx = 2

	if IsDiagonalMove(x, y, xx, yy) { //Expected false
		fmt.Println("Diagonal Move Incorrectly Identified")
		t.Fail()
	}
}

