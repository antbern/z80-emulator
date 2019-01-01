package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"github.com/antbern/z80-emulator/core"
)

func main() {
	// read the contents of the binary into a byte slice
	data, err := readBinary("input/main.bin")

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

	cpu := core.NewCPU()

	//data := []uint8{0x01, 0x02, 0x30}
	cpu.Mem.Write(0x0000, &code)

	// start with PC at 0x0000
	cpu.Reg.PC.Set(0x0000)

	//log.Printf("%s", cpu.Mem.Dump(0x0000, 0x1000))
	// infinite loop for procesing operands
	for {
		cpu.Step()
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

}
