package core

import (
	"encoding/hex"
	"log"
)

// RAM represents the RAM in the Z80
type RAM struct {
	data []uint8
}

const ramSize = 0x10000

// NewRAM makes a new RAM object with size 64k and populates with the initial data supplied
func NewRAM() *RAM {
	ram := RAM{
		data: make([]uint8, ramSize),
	}
	return &ram
}

func (ram *RAM) Write(addr uint16, data *[]uint8) {
	if int(addr)+len(*data) >= ramSize {
		log.Panic("[RAM] Tried to write outside RAM")
	}
	copy(ram.data[addr:], *data)
}

func (ram *RAM) read8(addr uint16) uint8 {
	return ram.data[addr]
}

func (ram *RAM) read8Inc(addr *uint16) uint8 {
	(*addr)++
	return ram.data[(*addr)-1]
}

func (ram *RAM) put8(addr uint16, data uint8) {
	ram.data[addr] = data
}

func (ram *RAM) read16(addr uint16) uint16 {
	return uint16(ram.data[addr+1])<<8 | uint16(ram.data[addr])
}

func (ram *RAM) read16Inc(addr *uint16) uint16 {
	(*addr) += 2
	return uint16(ram.data[*addr-1])<<8 | uint16(ram.data[*addr-2])
}

func (ram *RAM) put16(addr uint16, data uint16) {
	ram.data[addr] = uint8(data & 0xff)
	ram.data[addr+1] = uint8((data >> 8) & 0xff)
}

func (ram *RAM) ptr8(addr uint16) *uint8 {
	return &ram.data[addr]
}

// Dump prints the RAM contents to the provided writer
func (ram *RAM) Dump(start, length uint16) string {
	return hex.Dump(ram.data[start : start+length])
}
