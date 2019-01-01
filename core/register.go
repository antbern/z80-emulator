package core

// State holds all the registers in the Z80
type State struct {
	// program counter, stack pointer and index registers
	PC, SP, IX, IY R16

	// interrupt register
	I R8

	// general purpose registers
	GP *GPReg

	// alternate general purpose register set
	GPAlt *GPReg
}

// NewState returns a new state with default values
func NewState() State {
	return State{
		PC:    NewR16Single(0),
		SP:    NewR16Single(0),
		IX:    NewR16Single(0),
		IY:    NewR16Single(0),
		I:     NewR8(0),
		GP:    NewGPReg(),
		GPAlt: NewGPReg(),
	}
}

// GPReg contains a set of general purpose registers
type GPReg struct {
	A, F, B, C, D, E, H, L R8
	AF, BC, DE, HL         R16
}

// NewGPReg constructs a new res of general purpose register
func NewGPReg() *GPReg {
	rs := GPReg{
		A: NewR8(0),
		F: NewR8(0),
		B: NewR8(0),
		C: NewR8(0),
		D: NewR8(0),
		E: NewR8(0),
		H: NewR8(0),
		L: NewR8(0),
	}

	rs.AF = NewR16Combined(&rs.A, &rs.F)
	rs.BC = NewR16Combined(&rs.B, &rs.C)
	rs.DE = NewR16Combined(&rs.D, &rs.E)
	rs.HL = NewR16Combined(&rs.H, &rs.L)

	return &rs
}

/////////////////////////////////////////////////////////////////

// R8 is the basic interface for an 8-bit register
type R8 interface {
	Get() uint8
	Set(val uint8)
}

// R16 is the basic type for a 16-bit register
type R16 interface {
	Get() uint16
	Set(uint16)
}

// internal representation of an 8-bit register
type r8 struct {
	value uint8
}

// NewR8 creates a new 8 bit register
func NewR8(val uint8) R8 {
	return &r8{value: val}
}

func (r *r8) Get() uint8 {
	return r.value
}

func (r *r8) Set(val uint8) {
	r.value = val
}

// r16combined represents a 16-bt register made up of two 8 bit registers
type r16combined struct {
	high, low *R8
}

// NewR16Combined creates a new 16-bit register made up of two 8-bit registers
func NewR16Combined(high, low *R8) R16 {
	return &r16combined{
		high: high,
		low:  low,
	}
}

func (r *r16combined) Get() uint16 {
	return uint16((*r.high).Get())<<8 | uint16((*r.low).Get())
}

func (r *r16combined) Set(val uint16) {
	(*r.high).Set(uint8(val >> 8))
	(*r.low).Set(uint8(val & 0xff))
}

// r16single represents a 16-bt register represented by a single 16 bit value
type r16single struct {
	value uint16
}

// NewR16Single creates a new 16-bt register represented by a single 16 bit value
func NewR16Single(val uint16) R16 {
	return r16single{value: val}
}

func (r r16single) Get() uint16 {
	return r.value
}

func (r r16single) Set(val uint16) {
	r.value = val
}
