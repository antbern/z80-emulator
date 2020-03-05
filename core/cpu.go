package core

import (
	"fmt"
	"log"

	"github.com/antbern/z80-emulator/io"
)

// a lookup table to be used when evaluating conditional call/jump op codes
var condTable = []Condition{NonZero, Zero, NoCarry, Carry, ParityOdd, ParityEven, SignPos, SignNeg}

// Z80 contains all internal registers and such for the Z80 processor
type Z80 struct {
	// the 16 bit registers
	AF, BC, DE, HL, IX, IY R16

	// their high and low parts
	A, F, B, C, D, E, H, L R8
	IXL, IXH, IYL, IYH     R8

	// and their alternatives
	AFa, BCa, DEa, HLa R16

	// the stack pointer and program counter
	SP, PC R16

	// memory and IO device
	Mem *RAM
	IO  io.Device

	// internal flags
	Halted, InterruptEnabled bool
}

// NewZ80 creates a new Z80 CPU instance with memory, and registers
func NewZ80() Z80 {
	z80 := Z80{Mem: NewRAM(), Halted: false}

	// set up all the registers
	z80.AF, z80.A, z80.F = NewR16()
	z80.BC, z80.B, z80.C = NewR16()
	z80.DE, z80.D, z80.E = NewR16()
	z80.HL, z80.H, z80.L = NewR16()
	z80.IX, z80.IXH, z80.IXL = NewR16()
	z80.IY, z80.IYH, z80.IYL = NewR16()

	// the alternatives
	z80.AFa = NewR16Single()
	z80.BCa = NewR16Single()
	z80.DEa = NewR16Single()
	z80.HLa = NewR16Single()

	// and the stack pointer and program counter
	z80.SP = NewR16Single()
	z80.PC = NewR16Single()
	return z80
}

// Step causes the CPU to handle the next instruction
func (z *Z80) Step() {
	// TODO: for now, just don't do anything upon halted. Later: let interrupt resume execution
	if z.Halted {
		return
	}

	// read next operand and move PC forward
	// opCode := uint8(0x58)
	opCode := z.Mem.read8Inc(z.PC)

	// TODO: check for prefixed multi-byte op codes
	if opCode == 0xCB { // bit manipulations and roll/shift
		op := parseOP(z.Mem.read8Inc(z.PC))
		reg := z.regTableR(op.z)
		switch op.x {
		case 0: // TODO: rot[y] r[z]
		case 1: // BIT y, r[z]: Z = NOT bit y in r[z]
			bit8(*reg, op.y, z.F)
		case 2: // RES y, r[z]
			*reg &^= (1 << op.y)
		case 3: // SET y, r[z]
			*reg |= (1 << op.y)
		}
		// don't continue parsing
		return
	}

	// normal op-code, parse operands
	op := parseOP(opCode)
	log.Printf("Operand: %#02x -> %+v", opCode, op)

	// handle the op-codes using a giant switch matrix
	switch op.x {
	case 0: // x
		switch op.z {
		case 0:
			switch op.y {
			case 0: // NOP
			case 1: // EX AF, AF'
				exchange16(z.AF, z.AFa)
			case 2: // DJNZ d
				// decrement B
				*z.B--
				if *z.B > 0 { // if B is not yet zero, jump
					disp := z.Mem.read8Inc(z.PC)
					*z.PC += uint16(int8(disp))
				} else {
					*z.PC++ // increment PC to skip the displacement byte (no jump performed)
				}
			case 3: // JR d
				// read displacement byte and add it to PC (note handling of signed/unsigned numbers!)
				disp := z.Mem.read8Inc(z.PC)
				*z.PC += uint16(int8(disp))
			case 4, 5, 6, 7: // JR cc[y-4], d
				if condTable[op.y-4].isTrue(z.F) {
					disp := z.Mem.read8Inc(z.PC)
					*z.PC += uint16(int8(disp))
				} else {
					*z.PC++ // increment PC to skip the displacement byte
				}
			}
		case 1:
			reg := z.regTableRP(op.p, false)
			if op.q == 0 { // LD rp[p], nn
				*reg = z.Mem.read16Inc(z.PC)
			} else if op.q == 1 { // ADD HL, rp[p]
				*z.HL += *reg
			}
		case 2:
		case 3:
			reg := z.regTableRP(op.p, false)
			if op.q == 0 { // INC rp[p]
				*reg++
			} else if op.q == 1 { // DEC rp[p]
				*reg--
			}
		case 4: // INC r[y]
			reg := z.regTableR(op.y)
			*reg++
		case 5: // DEC r[y]
			reg := z.regTableR(op.y)
			*reg--
		case 6: // LD r[y], n
			reg := z.regTableR(op.y)
			*reg = z.Mem.read8Inc(z.PC)
		case 7:
			// some accumulator operands
		}
	case 1: // x
		// z=6 AND y=6 -> HALT
		if op.z == 6 && op.y == 6 {
			// HALT!!
			log.Printf("HALT!")
			z.Halted = true
			break
		}
		// 	LD r[y], r[z]
		dst := z.regTableR(op.y)
		src := z.regTableR(op.z)
		*dst = *src
	case 2: // x: ALU operation alu[y] with argument r[z]
		reg := z.regTableR(op.z)
		switch op.y {
		case 0: // ADD A, reg
			*z.A = add8(*z.A, *reg, z.F, false)
		case 1: // ADC A, reg
			*z.A = add8(*z.A, *reg, z.F, true)
		case 2: // SUB A, reg
			*z.A = sub8(*z.A, *reg, z.F, false)
		case 3: // SBC A, reg
			*z.A = sub8(*z.A, *reg, z.F, true)
		case 4: // TODO: AND
		case 5: // TODO: XOR
		case 6: // TODO: OR
		case 7: // CP i.e, A-r
			sub8(*z.A, *reg, z.F, false)
		}
	case 3: // x
		switch op.z {
		case 0: // RET cc[y]
			if condTable[op.y].isTrue(z.F) {
				z.Mem.stackPop16(z.SP, z.PC)
			}
		case 1:
			if op.q == 0 { // POP rp2[p]
				reg := z.regTableRP(op.p, true)
				z.Mem.stackPop16(z.SP, reg)
			} else if op.q == 1 {
				switch op.p {
				case 0: // RET
					z.Mem.stackPop16(z.SP, z.PC)
				case 1: // EXX
					exchange16(z.BC, z.BCa)
					exchange16(z.DE, z.DEa)
					exchange16(z.HL, z.HLa)
				case 2: // JP HL / JP (HL)
					*z.PC = z.Mem.read16(*z.HL)
				case 3: // LD SP, HL
					*z.SP = *z.HL
				}
			}
		case 2: // JP cc[y], nn
			if condTable[op.y].isTrue(z.F) {
				*z.PC = z.Mem.read16(*z.PC)
			} else {
				*z.PC += 2 // increment PC to skip jump address
			}
		case 3:
			switch op.y {
			case 0: // JP nn
				*z.PC = z.Mem.read16(*z.PC)
			case 1: // CB prefix
			case 2: // OUT (n), A
				addr := z.Mem.read8Inc(z.PC)
				if z.IO != nil {
					z.IO.Write(addr, *z.A)
				}
			case 3: // IN A, (n)
				addr := z.Mem.read8Inc(z.PC)
				if z.IO != nil {
					*z.A = z.IO.Read(addr)
				}
			case 4: // EX (SP), HL
				tmp := *z.HL
				*z.HL = z.Mem.read16(*z.SP)
				z.Mem.put16(*z.SP, tmp)
			case 5: // EX DE, HL
				exchange16(z.DE, z.HL)
			case 6: // DI
				z.InterruptEnabled = false
			case 7: // EI
				z.InterruptEnabled = true
			}
		case 4: // CALL cc[y], nn
			if condTable[op.y].isTrue(z.F) {
				// read adress to call, push return pointer to the stack and move PC
				addr := z.Mem.read16Inc(z.PC)
				z.Mem.stackPush16(z.SP, z.PC)
				*z.PC = addr
			} else {
				*z.PC += 2 // increment PC to skip jump address
			}
		case 5:
			if op.q == 0 { // PUSH rp2[p]
				reg := z.regTableRP(op.p, true)
				z.Mem.stackPush16(z.SP, reg)
			} else if op.q == 1 && op.p == 0 { // CALL nn
				// read adress to call, push return pointer to the stack and move PC
				addr := z.Mem.read16Inc(z.PC)
				z.Mem.stackPush16(z.SP, z.PC)
				*z.PC = addr
				log.Printf("CALL to %#04X", addr)
			}
		case 6: // ALU[y] n
			n := z.Mem.read8Inc(z.PC)
			switch op.y {
			case 0: // ADD A, reg
				*z.A = add8(*z.A, n, z.F, false)
			case 1: // ADC A, reg
				*z.A = add8(*z.A, n, z.F, true)
			case 2: // SUB A, reg
				*z.A = sub8(*z.A, n, z.F, false)
			case 3: // SBC A, reg
				*z.A = sub8(*z.A, n, z.F, true)
			case 4: // TODO: AND
			case 5: // TODO: XOR
			case 6: // TODO: OR
			case 7: // CP i.e, A-n
				sub8(*z.A, n, z.F, false)
			}
		case 7: // RST y*8
			*z.PC = uint16(op.y) * 8
		}
	}

}

// echanges / swaps the values of two R16 registers
func exchange16(a, b R16) {
	*a, *b = *b, *a
}

// codeToRegister takes a bit-code and returns a pointer to the correct 8-bit register or memory address
func (z *Z80) regTableR(code uint8) R8 {
	switch code {
	case 0:
		return z.B
	case 1:
		return z.C
	case 2:
		return z.D
	case 3:
		return z.E
	case 4:
		return z.H
	case 5:
		return z.L
	case 6:
		return z.Mem.ptr8(*z.HL)
	case 7:
		return z.A
	}

	return nil
}

func (z *Z80) regTableRP(code uint8, withAF bool) R16 {
	switch code {
	case 0:
		return z.BC
	case 1:
		return z.DE
	case 2:
		return z.HL
	case 3:
		if withAF {
			return z.AF
		}
		return z.SP
	}
	return nil
}

func (z *Z80) String() string {
	return fmt.Sprintf("PC: %#04x SP: %#04x\n A: %#02x F: %#02x B: %#02x C: %#02x D: %#02x E: %#02x H: %#02x L: %#02x\nAF: %#04x BC: %#04x DE: %#04x HL: %#04x IX: %#04x IY: %#04x",
		*z.PC, *z.SP, *z.A, *z.F, *z.B, *z.C, *z.D, *z.E, *z.H, *z.L, *z.AF, *z.BC, *z.DE, *z.HL, *z.IX, *z.IY)
}

// OP splits an op-code into its different parts according to the description in http://www.z80.info/decoding.htm
type OP struct {
	x, y, z uint8
	p, q    uint8
}

func parseOP(opCode uint8) OP {
	return OP{
		x: opCode >> 6,
		y: (opCode >> 3) & 0x7,
		z: (opCode & 0x7),
		p: (opCode >> 4) & 0x3,
		q: (opCode >> 3) & 0x1,
	}
}
