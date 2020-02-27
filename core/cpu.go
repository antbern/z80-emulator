package core

import (
	"fmt"
	"log"
)

/*
	"bufio"
	"log"
	"os"
*/

// Z80 contains all internal registers and such for the Z80 processor
type Z80 struct {
	A, F, B, C, D, E, H, L uint8
	IXL, IXH, IYL, IYH     uint8
	AF, BC, DE, HL, IX, IY reg16

	spL, spH uint8
	SP       reg16

	PC uint16

	Mem *RAM
}

func NewZ80() Z80 {
	z80 := Z80{Mem: NewRAM()}
	z80.AF = reg16{&z80.A, &z80.F}
	z80.BC = reg16{&z80.B, &z80.C}
	z80.DE = reg16{&z80.D, &z80.E}
	z80.HL = reg16{&z80.H, &z80.L}

	z80.IX = reg16{&z80.IXH, &z80.IXL}
	z80.IY = reg16{&z80.IYH, &z80.IYL}

	z80.SP = reg16{&z80.spH, &z80.spL}

	return z80
}

// Step causes the CPU to handle the next instruction
func (z *Z80) Step() {

	// infinite loop for procesing operands
	// read next operand and move PC forward

	// opCode := uint8(0x58) //z.Mem.read8Inc(&z.PC)
	opCode := z.Mem.read8Inc(&z.PC)

	//c.Reg.PC.Set(c.Reg.PC.Get() + 1)

	// TODO: check for prefixed multi-byte op codes
	if opCode == 0xCB { // bit manipulations and roll/shift
		op := ParseOP(z.Mem.read8Inc(&z.PC))
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
	op := ParseOP(opCode)
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
				reg.set(z.Mem.read16Inc(&z.PC))
			} else if op.q == 1 {
				// ADD HL, rp[p]
				z.HL.set(z.HL.get() + reg.get())
			}
		case 2:
		case 3:
			reg, _ := z.regTableRP(op.p, false)
			if op.q == 0 {
				reg.inc()
			} else if op.q == 1 {
				reg.dec()
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
			(*reg) = z.Mem.read8Inc(&z.PC)
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
				z.PC = z.Mem.read16(z.PC)
			case 2: // OUT (n), A
				log.Printf("OUT: %#02x -> %#02x", z.A, z.Mem.read8Inc(&z.PC))
			}
		case 4: // CALL cc[y], nn
		case 5:
		case 6: // ALU[y] n
		case 7: // RST y*8
			z.PC = uint16(op.y) * 8
		}
	}

}

// codeToRegister takes a bit-code and returns a pointer to the correct 8-bit register or memory address
func (z *Z80) regTableR(code uint8) (*uint8, string) {
	switch code {
	case 0:
		return &z.B, "B"
	case 1:
		return &z.C, "C"
	case 2:
		return &z.D, "D"
	case 3:
		return &z.E, "E"
	case 4:
		return &z.H, "H"
	case 5:
		return &z.L, "L"
	case 6:
		return z.Mem.ptr8(z.HL.get()), "(HL)"
	case 7:
		return &z.A, "A"
	}

	return nil, ""
}

func (z *Z80) regTableRP(code uint8, withAF bool) (*reg16, string) {
	switch code {
	case 0:
		return &z.BC, "BC"
	case 1:
		return &z.DE, "DE"
	case 2:
		return &z.HL, "HL"
	case 3:
		if withAF {
			return &z.AF, "AF"
		}
		return &z.SP, "SP"
	}
	return nil, ""
}

func (z *Z80) String() string {
	return fmt.Sprintf("PC: %#04x SP: %#04x\n A: %#02x F: %#02x B: %#02x C: %#02x D: %#02x E: %#02x H: %#02x L: %#02x\nAF: %#04x BC: %#04x DE: %#04x HL: %#04x IX: %#04x IY: %#04x",
		z.PC, z.SP.get(), z.A, z.F, z.B, z.C, z.D, z.E, z.H, z.L, z.AF.get(), z.BC.get(), z.DE.get(), z.HL.get(), z.IX.get(), z.IY.get())
}

// Operand describes a single operands actions
type Operand struct {
	name string
	// theese values are used to match the compiled operands
	opCodeMask, opCodeValue uint8
	handle                  func(*Z80)
}

var operands = []Operand{
	Operand{ // NOOP
		name:        "NOOP",
		opCodeMask:  0xff,
		opCodeValue: 0x00,
		handle:      func(*Z80) {},
	}, Operand{ // JR
		name:        "JR",
		opCodeMask:  0xff,
		opCodeValue: 0x18,
		handle: func(z *Z80) {
			// read offset and increment PC
			offset := int16(z.Mem.read8Inc(&z.PC))

			// perform the jump using ugly type conversions...
			z.PC = uint16(int16(z.PC) + offset)

			log.Printf("JP, offset: %d (%#x), new PC: %#x", offset, offset, z.PC)
		},
	}, Operand{
		name:        "CALL",
		opCodeMask:  0xff,
		opCodeValue: 0xCD,
		handle: func(z *Z80) {

			// read target adress
			target := z.Mem.read16Inc(&z.PC)

			// TODO: push PC onto the stack

			// move PC
			z.PC = target

			log.Printf("CALL, target: %#x", target)
		},
	},
}

// OP splits an op-code into its different parts according to the description in http://www.z80.info/decoding.htm
type OP struct {
	x, y, z uint8
	p, q    uint8
}

func ParseOP(op uint8) OP {
	return OP{
		x: op >> 6,
		y: (op >> 3) & 0x7,
		z: (op & 0x7),
		p: (op >> 4) & 0x3,
		q: (op >> 3) & 0x1,
	}
}
