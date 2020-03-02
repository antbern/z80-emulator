package core

import "testing"

func TestExchange(t *testing.T) {
	// and the stack pointer and program counter
	A, _, _ := NewR16()
	B, _, _ := NewR16()

	*A = 0x55AA
	*B = 0x4488

	exchange16(A, B)

	// make sure the values were right after the swap
	if *A != 0x4488 || *B != 0x55AA {
		t.Errorf("Expected A = 0x4488 and B = 0x55AA, got %#04x and %#04x", *A, *B)
	}

}
