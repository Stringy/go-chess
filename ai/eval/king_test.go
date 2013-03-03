package eval

import (
	"chess/ai/gen"
	"fmt"
	"testing"
)

func TestShielding(t *testing.T) {
	for i, mask := range WeakShieldBlack {
		fmt.Println(i, ":")
		gen.PrintBitboard(mask)
	}
}
