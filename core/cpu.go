package core

import (
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

	PC, SP uint16

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

	return z80
}

// Step causes the CPU to handle the next instruction
func (z *Z80) Step() {

	// infinite loop for procesing operands
	// read next operand and move PC forward

	op := uint8(0x58) //z.Mem.read8Inc(&z.PC)

	//c.Reg.PC.Set(c.Reg.PC.Get() + 1)

	// TODO: check for prefixed multi-byte op codes

	// parse operands
	o := ParseOP(op)
	log.Printf("Operand: %#x -> %+v", op, o)

	// find operand that match
	/*
		matched := false
		for _, operand := range operands {
			if (op & operand.opCodeMask) == operand.opCodeValue {
				log.Printf("Matched OP: %s", operand.name)
				matched = true
				operand.handle(z)
			}
		}
		if !matched {
			log.Printf("No handler found for OP %#x", op)
		}
	*/

	// handle the op-codes using a giant switch matrix
	switch o.x {
	case 0:
	case 1:
		// z=6 AND y=6 -> HALT
		if o.z == 6 && o.y == 6 {
			// HALT!!
			log.Printf("HALT!")
		}

		// 	LD r[y], r[z]

		// lookup register
		targetReg, targetName := z.codeToRegister(o.y)
		sourceReg, sourceName := z.codeToRegister(o.z)

		// apply load
		*targetReg = *sourceReg

		// print debug message
		log.Printf("LD %s, %s", targetName, sourceName)

	case 2:
		// ALU operation alu[y] with argument r[z]
	case 3:
	}

}

// codeToRegister takes a bit-code and returns a pointer to the correct 8-bit register or memory address
func (z *Z80) codeToRegister(code uint8) (*uint8, string) {
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
