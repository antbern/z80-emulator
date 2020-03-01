package io

import "log"

// DebugDevice implements a Device that simply prints writes to stdin
type DebugDevice struct{}

func (DebugDevice) Init() {}

func (DebugDevice) Write(port, val uint8) {
	log.Printf("IO: Write %#02x to %#02x ", val, port)
}

func (DebugDevice) Read(port uint8) uint8 {
	log.Printf("IO: Read from %#02x ", port)
	return 0
}
