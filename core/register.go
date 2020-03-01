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
