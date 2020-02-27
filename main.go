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
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/antbern/z80-emulator/core"
)

func main() {
	// read the contents of the binary into a byte slice
	data, err := readBinary("input/count.bin")

	if err != nil {
		log.Println("Error leading file: ", err)
		return
	}

	log.Printf("Data: [%# x...", data[:16])

	mainLoop(data)
}

func readBinary(filename string) ([]byte, error) {
	log.Println("Loading file", filename)

	data, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func mainLoop(code []byte) {
	testReg := core.NewR8(0)
	testReg.Get()
	//var r core.r8 = core.NewR8(0)

	//r6 := core.NewR16(&testReg, &testReg)

	cpu := core.NewZ80()

	//data := []uint8{0x01, 0x02, 0x30}
	cpu.Mem.Write(0x0000, &code)

	// start with PC at 0x0000
	cpu.PC = 0x0000

	//log.Printf("%s", cpu.Mem.Dump(0x0000, 0x1000))
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
		}

		println(cpu.String())
	}
outside:
}
