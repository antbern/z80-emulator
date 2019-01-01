package core

import (
	"log"
)

/*
	"bufio"
	"log"
	"os"
*/

type CPU struct {
	Mem *RAM
	Reg *State
}

// NewCPU created a new CPU structure
func NewCPU() CPU {
	return CPU{
		Mem: NewRAM(),
		Reg: &State{},
	}
}

// Step causes the CPU to handle the next instruction
func (c *CPU) Step() {

	// infinite loop for procesing operands
	// read next operand and move PC forward

	op := c.Mem.read8(c.Reg.PC.Get())

	c.Reg.PC.Set(c.Reg.PC.Get() + 1)

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
				operand.handle(c.Reg, c.Mem)
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
		//*c.Reg.codeToRegister(o.y) = *c.Reg.codeToRegister(o.z)

	case 2:
		// ALU operation alu[y] with argument r[z]
	case 3:
	}

}

/*
// codeToRegister takes a bit-code and returns a pointer to the correct 8-bit register
func (r *State) codeToRegister(code uint8) *uint8 {
	switch code {
	case 0:
		return &r.RegSet.B
	case 1:
		return &r.RegSet.C
	case 2:
		return &r.RegSet.D
	case 3:
		return &r.RegSet.E
	case 4:
		return &r.RegSet.H
	case 5:
		return &r.RegSet.L
	case 6:
		return nil //&r.RegSet.HL
	case 7:
		return &r.RegSet.A
	}

	return nil
}
*/

// Operand describes a single operands actions
type Operand struct {
	name string
	// theese values are used to match the compiled operands
	opCodeMask, opCodeValue uint8
	handle                  func(*State, *RAM)
}

/*
var operands = []Operand{
	Operand{ // NOOP
		name:        "NOOP",
		opCodeMask:  0xff,
		opCodeValue: 0x00,
		handle:      func(*State, *RAM) {},
	}, Operand{ // JR
		name:        "JR",
		opCodeMask:  0xff,
		opCodeValue: 0x18,
		handle: func(s *State, r *RAM) {
			// read offset and increment PC
			offset := int16(r.read8Inc(&s.PC))

			// perform the jump using ugly type conversions...
			s.PC = uint16(int16(s.PC) + offset)

			log.Printf("JP, offset: %d (%#x), new PC: %#x", offset, offset, s.PC)
		},
	}, Operand{
		name:        "CALL",
		opCodeMask:  0xff,
		opCodeValue: 0xCD,
		handle: func(s *State, r *RAM) {

			// read target adress
			target := r.read16Inc(&s.PC)

			// TODO: push PC onto the stack

			// move PC
			s.PC = target

			log.Printf("CALL, target: %#x", target)
		},
	},
}
*/

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
