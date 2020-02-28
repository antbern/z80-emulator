package core

import (
	"fmt"
	"log"
)

// Z80 contains all internal registers and such for the Z80 processor
type Z80 struct {
	R, Ralt *regs
	SP      R16
	PC      R16
	Mem     *RAM
}

// NewZ80 creates a new Z80 CPU instance with memory, and registers
func NewZ80() Z80 {
	z80 := Z80{Mem: NewRAM(), R: newRegs(), Ralt: newRegs()}
	z80.SP, _, _ = NewR16()
	z80.PC, _, _ = NewR16()
	return z80
}

// Step causes the CPU to handle the next instruction
func (z *Z80) Step() {

	// read next operand and move PC forward
	// opCode := uint8(0x58)
	opCode := z.Mem.read8Inc(z.PC)

	//c.Reg.PC.Set(c.Reg.PC.Get() + 1)

	// TODO: check for prefixed multi-byte op codes
	if opCode == 0xCB { // bit manipulations and roll/shift
		op := parseOP(z.Mem.read8Inc(z.PC))
		reg, _ := z.regTableR(op.z)
		switch op.x {
		case 0: // rot[y] r[z]
		case 1: // BIT y, r[z]: Z = NOT bit y in r[z]
		case 2: // RES y, r[z]
			*reg &^= (1 << op.y)
		case 3: // SET y, r[z]
			*reg |= (1 << op.y)
		}

	}

	// parse operands
	op := parseOP(opCode)
	log.Printf("Operand: %#x -> %+v", opCode, op)

	// handle the op-codes using a giant switch matrix
	switch op.x {
	case 0:
		switch op.z {
		case 0:
		case 1:
			reg, _ := z.regTableRP(op.p, false)
			if op.q == 0 {
				// LD rp[p], nn
				*reg = z.Mem.read16Inc(z.PC)
			} else if op.q == 1 {
				// ADD HL, rp[p]
				*z.R.HL += *reg
			}
		case 2:
		case 3:
			reg, _ := z.regTableRP(op.p, false)
			if op.q == 0 {
				*reg++
			} else if op.q == 1 {
				*reg--
			}
		case 4:
			// INC r[y]
			reg, _ := z.regTableR(op.y)
			(*reg)++
		case 5:
			// DEC r[y]
			reg, _ := z.regTableR(op.y)
			(*reg)--
		case 6:
			// LD r[y], n
			reg, _ := z.regTableR(op.y)
			*reg = z.Mem.read8Inc(z.PC)
		case 7:
			// some accumulator operands
		}
	case 1:
		// z=6 AND y=6 -> HALT
		if op.z == 6 && op.y == 6 {
			// HALT!!
			log.Printf("HALT!")
		}

		// 	LD r[y], r[z]

		// lookup register
		targetReg, targetName := z.regTableR(op.y)
		sourceReg, sourceName := z.regTableR(op.z)

		// apply load
		*targetReg = *sourceReg

		// print debug message
		log.Printf("LD %s, %s", targetName, sourceName)

	case 2:
		// ALU operation alu[y] with argument r[z]
	case 3:
		switch op.z {
		case 0: // RET cc[y]
		case 1:
		case 2: // JP cc[y], nn
		case 3:
			switch op.y {
			case 0: // JP nn
				*z.PC = z.Mem.read16(*z.PC)
			case 2: // OUT (n), A
				log.Printf("OUT: %#02x -> %#02x", *z.R.A, z.Mem.read8Inc(z.PC))
			}
		case 4: // CALL cc[y], nn
		case 5:
		case 6: // ALU[y] n
		case 7: // RST y*8
			*z.PC = uint16(op.y) * 8
		}
	}

}

// codeToRegister takes a bit-code and returns a pointer to the correct 8-bit register or memory address
func (z *Z80) regTableR(code uint8) (R8, string) {
	switch code {
	case 0:
		return z.R.B, "B"
	case 1:
		return z.R.C, "C"
	case 2:
		return z.R.D, "D"
	case 3:
		return z.R.E, "E"
	case 4:
		return z.R.H, "H"
	case 5:
		return z.R.L, "L"
	case 6:
		return z.Mem.ptr8(*z.R.HL), "(HL)"
	case 7:
		return z.R.A, "A"
	}

	return nil, ""
}

func (z *Z80) regTableRP(code uint8, withAF bool) (R16, string) {
	switch code {
	case 0:
		return z.R.BC, "BC"
	case 1:
		return z.R.DE, "DE"
	case 2:
		return z.R.HL, "HL"
	case 3:
		if withAF {
			return z.R.AF, "AF"
		}
		return z.SP, "SP"
	}
	return nil, ""
}

func (z *Z80) String() string {
	return fmt.Sprintf("PC: %#04x SP: %#04x\n A: %#02x F: %#02x B: %#02x C: %#02x D: %#02x E: %#02x H: %#02x L: %#02x\nAF: %#04x BC: %#04x DE: %#04x HL: %#04x IX: %#04x IY: %#04x",
		*z.PC, *z.SP, *z.R.A, *z.R.F, *z.R.B, *z.R.C, *z.R.D, *z.R.E, *z.R.H, *z.R.L, *z.R.AF, *z.R.BC, *z.R.DE, *z.R.HL, *z.R.IX, *z.R.IY)
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
