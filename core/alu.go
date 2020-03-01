package core

// This file contains implementations for all the ALU (Arithmetic Logic Unit) operations
// and handles correct setting and testing of flag register flags

// constants for the flag register F
const (
	// FlagS Sign flag (bit 7)
	FlagS = 1 << 7

	// FlagZ Zero flag (bit 6)
	FlagZ = 1 << 6

	// FlagH Half Carry flag (bit 4)
	FlagH = 1 << 4

	// FlagP Parity/Overflow flag (bit 2) (set if the result byte has an even number of bits set)
	FlagP = 1 << 2

	// FlagN Add/Subtract flag (bit 1)
	FlagN = 1 << 1

	// FlagC Carry flag (bit 0)
	FlagC = 1 << 0
)

// Condition represents the possible conditions that can be tested for by the ALU
type Condition uint16

// Constants for all possible conditions that can be tested for. High byte is mask, and low byte is result required for true evaluation
const (
	NonZero    Condition = (FlagZ<<8 | 0)
	Zero                 = (FlagZ<<8 | FlagZ)
	NoCarry              = (FlagC<<8 | 0)
	Carry                = (FlagC<<8 | FlagC)
	ParityEven           = (FlagP<<8 | FlagP)
	ParityOdd            = (FlagP<<8 | 0)
	SignPos              = (FlagS<<8 | 0)
	SignNeg              = (FlagS<<8 | FlagS)
)

// isTrue tests whether the condition is true based on the provided Flag register
func (c Condition) isTrue(F R8) bool {
	return (*F & uint8(c>>8)) == uint8(c&0xff)
}

func add8(dst, val R8) uint8 {
	return *dst + *val
}

// // alu defines the basics for the alu
// type alu struct {
// 	// reference to the F register where all the flags are stored
// 	F R8
// }
