package core

// constants for the flag register F

// FlagS Sign flag (bit 7)
const FlagS = 1 << 7

// FlagZ Zero flag (bit 6)
const FlagZ = 1 << 6

// FlagH Half Carry flag (bit 4)
const FlagH = 1 << 4

// FlagP Parity/Overflow flag (bit 2)
const FlagP = 1 << 2

// FlagN Add/Subtract flag (bit 1)
const FlagN = 1 << 1

// FlagC Carry flag (bit 0)
const FlagC = 1 << 0
