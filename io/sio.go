package io

import (
	"io"
	"log"
)

// SIO implements the IODevice interface and represents the Z80 Serial Input/Output device
type SIO struct {
	reader *io.Reader
	writer *io.Writer
}

// constants defining the address used
const (
	SioBase  = 0x20
	SioAData = SioBase + 0 + 0
	SioACtrl = SioBase + 0 + 2
	SioBData = SioBase + 1 + 0
	SioBCtrl = SioBase + 1 + 2
)

// NewSIO returns a new SIO/2 serial input output device tied to the console
func NewSIO(r *io.Reader, w *io.Writer) *SIO {
	return &SIO{reader: r, writer: w}
}

func (s *SIO) Write(port, val uint8) {
	switch port {
	case SioAData:
		// TODO: Implement real writing
		log.Printf("SIO: Write %#02x (%v, '%v') ", val, val, string(val))
	case SioACtrl:
		// TODO: Handle Control bits
	default:
		log.Printf("SIO: Port B write not implemented yet!")
	}
}

func (s *SIO) Read(port uint8) uint8 {
	switch port {
	case SioAData:
		// TODO implement serial console read here
		return 0
	case SioACtrl:
		return (1 << 2) | (1 << 0) // bit 2 = tx ready, bit 0 = 1 if buffer has data
	default:
		log.Printf("SIO: Port B read not implemented yet!")
	}
	return 0
}
