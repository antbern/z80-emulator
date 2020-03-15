package core

import (
	"fmt"
	"log"
)

// handleBDOS handles any CP/M BDOS calls (call 5)
// Basic usage in code:
// 	LD  DE,parameter
// 	LD  C,function
// 	CALL    5
// Information about the available functions can be found here:
// * http://seasip.info/Cpm/bdosfunc.html
// * http://www.gaby.de/cpm/manuals/archive/cpm22htm/ch5.htm#Section_5.6
// The most important ones are 2: Write Char and 9: Write String
// By implementing the above, the Zexdoc Z80 exerciser can be run on this system :D
func (z *Z80) handleBDOS() {
	log.Printf("BDOS CALL: C=%v, DE=%#04x ", *z.C, *z.DE)

	switch *z.C {
	case 2: // Write char, E = the ascii character to write
		fmt.Print(byte(*z.E))

	case 9: // Write string, DE = points to start of string, terminated with the $ character
		// a for loop to perform the printing
		addr := *z.DE
		for {
			c := rune(z.Mem.read8Inc(&addr))
			if c == '$' {
				break
			}
			fmt.Printf("%c", c)
		}

	}
}
