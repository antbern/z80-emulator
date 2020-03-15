package main

/* Useful Links
https://github.com/remogatto/z80/blob/master/z80.go
https://github.com/floooh/chips/blob/master/systems/cpc.h
https://floooh.github.io/2016/07/12/z80-rust-ms1.html
Instruction table: http://clrhome.org/table/
Decoding instructions: http://www.z80.info/decoding.htm
*/

import (
	"bufio"
	"encoding/hex"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/antbern/z80-emulator/core"
	"github.com/antbern/z80-emulator/io"
)

func main() {

	fileName := flag.String("i", "input/monitor.bin", "The binary input file to load")
	origin := flag.String("o", "0x0000", "The origin/base address of the code. Decides where the loaded file will be placed in memory")
	flag.Parse()

	// parse the origin base address
	baseAddr, err := strconv.ParseUint(*origin, 0, 16)

	if err != nil {
		log.Printf("Error parsing origin argument %v: %v\n", *origin, err)
		return
	}

	// read the contents of the binary into a byte slice
	log.Println("Loading file", *fileName)
	data, err := ioutil.ReadFile(*fileName)
	if err != nil {
		log.Println("Error loading file: ", err)
		return
	}

	log.Printf("\n%s", hex.Dump(data[:64]))

	mainLoop(data, uint16(baseAddr))
}

func mainLoop(code []byte, origin uint16) {

	cpu := core.NewZ80()
	cpu.IO = io.NewSIO(nil, nil)

	log.Printf("Writing loaded file (%v bytes) into memory at base address %#04x", len(code), origin)
	cpu.Mem.Write(origin, &code)

	// for the CP/M to know where the stack can start (used for zexdoc exerciser)
	// cpu.Mem.Write(0x0006, &[]byte{0xff, 0x00})
	// cpu.EnableBDOS = true

	// start with PC at the origin for now since the rest is just zeroes
	*cpu.PC = origin

	// infinite loop for procesing operands
	reader := bufio.NewReader(os.Stdin)
	for {
		// read a single character from stdin that decides what to do
		print(">")
		text, _ := reader.ReadBytes('\n')
		switch strings.TrimSpace(string(text)) {
		case "", "n":
			cpu.Step()
		case "nt":
			for i := 1; i <= 100; i++ {
				cpu.Step()
			}
		case "q":
			goto outside
		case "o":
			for {
				cpu.Step()
				if *cpu.PC == 0x04EC {
					break
				}
			}
		}

		println(cpu.String())
	}
outside:
}
