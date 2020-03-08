package core

import "math/bits"

// This file contains implementations for all the ALU (Arithmetic Logic Unit) operations
// and handles correct setting and testing of flag register flags

// constants for the flag register F
const (
	// FlagS Sign flag (bit 7)
	FlagSshift = 7
	FlagS      = 1 << FlagSshift

	// FlagZ Zero flag (bit 6)
	FlagZshift = 6
	FlagZ      = 1 << FlagZshift

	// FlagH Half Carry flag (bit 4)
	FlagHshift = 4
	FlagH      = 1 << FlagHshift

	// FlagP Parity (bit 2) (set if the result byte has an even number of bits set)
	FlagPshift = 2
	FlagP      = 1 << FlagPshift

	// FlagV Overflow (bit 2)
	FlagVshift = FlagPshift
	FlagV      = 1 << FlagVshift

	// FlagN Add/Subtract flag (bit 1)
	FlagNshift = 1
	FlagN      = 1 << FlagNshift

	// FlagC Carry flag (bit 0)
	FlagCshift = 0
	FlagC      = 1 << FlagCshift
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

// isZero returns 1 if value is zero, and 0 otherwise
func isZero(val uint8) uint8 {
	if val == 0 {
		return 1
	}
	return 0
}

// add8 simply adds two uint8 together with the carry flag (if addCarry is true) and affects the flags in the correct way
// inspiration for overflow and half carry calculation taken from
// https://stackoverflow.com/questions/8034566/overflow-and-carry-flags-on-z80#8037485
// It affects the following flags (from the Z80 User Manual)
// * S is set if result is negative; otherwise, it is reset.
// * Z is set if result is 0; otherwise, it is reset.
// * H is set if carry from bit 3: otherwise, it is reset.
// * P/V is set if overflow; otherwise, it is reset.
// * N is reset.
// * C is set if carry from bit 7; otherwise, it is reset.
func add8(a, b uint8, F R8, addCarry bool) uint8 {
	// perform the addition using 16 bit numbers + the carry flag if wanted
	res16 := uint16(a) + uint16(b)
	if addCarry {
		res16 += uint16((*F & FlagC) >> FlagCshift)
	}

	res := uint8(res16)

	// the carry that was created during this operation
	carryOut := uint8((res16 >> 8) & 0x1)

	// Calculate all carry-ins.
	// Remembering that each bit of the sum = addend a's bit XOR addend b's bit XOR carry-in,
	// we can work out all carry-ins from a, b and their sum.
	carryIns := res ^ a ^ b

	// Calculate the overflow using the carry-out and most significant carry-in.
	overflow := (carryIns >> 7) ^ carryOut

	// calculate the half carry too
	halfCarryOut := (carryIns >> 4) & 1

	// fmt.Printf("%#02x + %#02x = %#02x C=%v V=%v \n", a, b, res, carryOut, overflow)
	// fmt.Printf("%4v + %4v = %4v C=%v V=%v \n", a, b, res, carryOut, overflow)
	// fmt.Printf("%3v(%4v) + %3v(%4v) + %v = %3v(%4v) CY=%v OV=%v \n", a, int8(a), b, int8(b), (*F&FlagC)>>FlagCshift, res, int8(res), carryOut, overflow)
	// fmt.Printf("%#02x, %#02x, %#02x\n", res, overflow, (isZero(res) << FlagZshift))

	// copy the sign flag from res
	*F = (*F &^ FlagS) | (((res & (1 << 7)) >> 7) << FlagSshift)

	// carry flag
	*F = (*F &^ FlagC) | (carryOut << FlagCshift)

	// half-carry flag
	*F = (*F &^ FlagH) | (halfCarryOut << FlagHshift)

	// overflow flag
	*F = (*F &^ FlagV) | (overflow << FlagVshift)

	// set the zero flag correctly
	*F = (*F &^ FlagZ) | (isZero(res) << FlagZshift)

	// clear the add/subtract flag
	*F &^= FlagN

	return res

}

// sub8 performs a - b - carry (if subCarry is true) and returns the result while manipulating the flags
// based on the same answer as for addc8 above
func sub8(a, b uint8, F R8, subCarry bool) uint8 {
	// a - b - c = a + ~b + 1 - c = a + ~b + !c
	*F ^= FlagC
	res := add8(a, ^b, F, subCarry)
	*F ^= FlagC
	*F ^= FlagH // should probably toggle the half carry flag too
	*F |= FlagN // set add/subtract flag for subtract operation
	return res
}

// bit8 sets the Z flag to the inverse of bit 'bit' of 'a'
// also sets the H and resets the N flag, but does not affect C flag
// modification of S and P/V is undefined
func bit8(a, bit uint8, F R8) {
	// set Z to 1 if bit is 0, otherwise to 0 (= inverse of bit)
	*F = (*F &^ FlagZ) | (((^a >> bit) & 1) << FlagZshift)

	// set half-carry flag
	*F |= FlagH

	// reset the add/sub flag
	*F &^= FlagN
}

// and8 performs a logical AND between a and b and sets the flags accordingly
func and8(a, b uint8, F R8) uint8 {
	return _bitwiseFlagSet(a&b, F)
}

// or8 performs a logical OR between a and b and sets the flags accordingly
func or8(a, b uint8, F R8) uint8 {
	return _bitwiseFlagSet(a|b, F)
}

// xor8 performs a logical XOR between a and b and sets the flags accordingly
func xor8(a, b uint8, F R8) uint8 {
	return _bitwiseFlagSet(a^b, F)
}

func _bitwiseFlagSet(result uint8, F R8) uint8 {
	// copy the sign flag from res
	*F = (*F &^ FlagS) | (((result & (1 << 7)) >> 7) << FlagSshift)

	// set the zero flag correctly
	*F = (*F &^ FlagZ) | (isZero(result) << FlagZshift)

	// reset N and C, and set H
	*F = (*F &^ (FlagN | FlagC)) | FlagH

	// set the parity flag depending on result
	*F = (*F &^ FlagP) | ((^(uint8(bits.OnesCount8(result)) % 2) & 1) << FlagPshift)

	return result
}

/*
func inc8(val *uint8, F R8) {
	// do increment
	*val++

	// copy the sign flag (works since FlagS is (1<<7))
	*F = (*F &^ FlagS) | (*val & FlagS)

	// set the zero flag correctly
	*F = (*F &^ FlagZ) | (FlagS & isZero(*val))

	// copy the half-carry flag (works since FlagH is (1<<4))
	*F = (*F &^ FlagH) | (*val & FlagH)

	// half-carry flag from https://stackoverflow.com/questions/8868396/game-boy-what-constitutes-a-half-carry#8874607
	//((a & 0xf) + (value & 0xf)) & 0x10

}
*/
