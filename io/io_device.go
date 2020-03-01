package io

// Device defines the interface for communicating with an IO device using IN and OUT operations
type Device interface {
	// Init initializes the Device
	Init()

	// Write writes a single byte to the specified port
	Write(port, val uint8)

	// Read reads a single byte from the specified port
	Read(port uint8) uint8
}
