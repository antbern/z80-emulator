package core

import "testing"

func TestR8(t *testing.T) {
	r := NewR8(0x12)
	if r.Get() != 0x12 {
		t.Fail()
	}

	r.Set(0x55)
	if r.Get() != 0x55 {
		t.Fail()
	}
}

func TestR16(t *testing.T) {
	h, l := NewR8(0), NewR8(0)

	r := NewR16Combined(&h, &l)

	r.Set(0x1234)
	if r.Get() != 0x1234 {
		t.Fail()
	}

	if h.Get() != 0x12 {
		t.Fail()
	}

	if l.Get() != 0x34 {
		t.Fail()
	}

}
