package core

import "unsafe"

// R16 represents a 16-bit register
type R16 *uint16

// R8 represents a 8-bit register
type R8 *uint8

// NewR16 creates a new R16 register and returns a pointer to it as well as the high and low R8 parts
func NewR16() (R16, R8, R8) {
	val := new(uint16)
	high := (R8)(unsafe.Pointer(uintptr(unsafe.Pointer(val)) + 1*unsafe.Sizeof(uint8(0))))
	low := (R8)(unsafe.Pointer(uintptr(unsafe.Pointer(val)) + 0*unsafe.Sizeof(uint8(0))))
	return val, high, low
}

type regs struct {
	// the 16 bit registers
	AF, BC, DE, HL, IX, IY R16

	// and their high and low part
	A, F, B, C, D, E, H, L R8
	IXL, IXH, IYL, IYH     R8
}

// newRegs returns a pointer to a new initialized regs struct
func newRegs() *regs {
	r := regs{}
	r.AF, r.A, r.F = NewR16()
	r.BC, r.B, r.C = NewR16()
	r.DE, r.D, r.E = NewR16()
	r.HL, r.H, r.L = NewR16()
	r.IX, r.IXH, r.IXL = NewR16()
	r.IY, r.IYH, r.IYL = NewR16()
	return &r
}
