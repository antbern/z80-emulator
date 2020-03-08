package io

import "log"

// DebugDevice implements a Device that simply prints writes to stdin
type debugDevice struct{}

// NewDebugDevice returns a new IO device that can be used for debugging
func NewDebugDevice() debugDevice {
	return debugDevice{}
}

func (debugDevice) Write(port, val uint8) {
	log.Printf("IO: Write %#02x (%v, %v) to %#02x ", val, val, byte(val), port)
}

func (debugDevice) Read(port uint8) uint8 {
	log.Printf("IO: Read from %#02x ", port)
	return 0
}
