# z80-emulator
A Z80 emulator written in Go. Under slow development :)

**Goal:** To be able to run and develop my [z80-monitor](https://github.com/antbern/z80-monitor) as a console application on my own computer

Here are some thoughts on the project so far:

Features that are implemented
* Loading of binary code
* Control flow operations such as CALL, JP and RET, including conditional jumps
* Operations for loading registers with values
* Single stepping or bulk stepping through the instructions

Features that still need to be implementated
* Loading of Intel HEX files
* All arithmetic operations, including correct manipulation of the flag bits
* Block I/O instructions 
* Prefixed instructions such as for using IX/IY, IX+d/IY+d and so on
* Interrupts
* SIO console
* Interactive mode and its triggers


## Interactive Mode
Setting special triggers will enter so called _interactive mode_. Without specifying any triggers, the program will just execute and the console will be used as SIO input and output. Here is a list of considered triggers (that should be specifiable on the command line):

* Directly
* After N clock cycles/instructions
* Upon execution of specific OP-code?
* When PC reaches a specific address
* When any register contains a specific value
* Halt was executed? (maybe not since execution should continue once and interrupt is recieved)

Upon entering _interactive mode_:
* Display current register values in console
* Allow certain commands to be executed, for example:
    * Continue execution
    * Single stepping of instructions
    * Examine and change memory contents
    * (Set register contents)
    * Set up / remove breakpoints, either register or PC value
    * Quit