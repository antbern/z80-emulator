package core

import "testing"

func TestRegs(t *testing.T) {
	// create register
	r, h, l := NewR16()

	// test read
	*r = 0x5577
	if *h != 0x55 {
		t.Errorf("High read failed: got %#02x, wanted %#02x", *h, 0x55)
	}
	if *l != 0x77 {
		t.Errorf("Low read failed: got %#02x, wanted %#02x", *l, 0x77)
	}

	// test write
	*r = 0x5577
	*h = 0xaa
	if *r != 0xaa77 {
		t.Errorf("High write failed: got %#02x, wanted %#02x", *r, 0xaa77)
	}
	*l = 0x88
	if *r != 0xaa88 {
		t.Errorf("Low write failed: got %#02x, wanted %#02x", *r, 0xaa88)
	}

	// test for correct increments of high and low parts
	*r = 0xffff
	*r++
	if *r != 0 {
		t.Errorf("Increment failed: got %#04x, wanted %#04x", *r, 0)
	}

	*r = 0xffff
	*h++
	if *r != 0x00ff {
		t.Errorf("High increment failed: got %#04x, wanted %#04x", *r, 0x00ff)
	}

	*l++
	if *r != 0x0000 {
		t.Errorf("Low increment failed: got %#04x, wanted %#04x", *r, 0)
	}

	*r = 0x00ff
	*r++
	if *r != 0x0100 {
		t.Errorf("Low increment failed: got %#04x, wanted %#04x", *r, 0x100)
	}

}
