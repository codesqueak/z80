package internal

import (
	"errors"
	"fmt"
	"z80/processor/pkg/hw"
)

var memory *hw.Memory
var io *hw.IO
var reg Registers
var initialized = false

func Build(mem *hw.Memory, ports *hw.IO) error {
	if mem == nil {
		return errors.New("Memory not defined")
	}
	if ports == nil {
		return errors.New("I/O not defined")
	}
	memory = mem
	io = ports
	reg = Registers{a: 0xaa, f: 0x55}
	initialized = true
	return nil
}

// execute one instruction
func RunOne() (bool, error) {
	if !initialized {
		return false, errors.New("CPU not initialized")
	}
	//
	halt := execute()
	//
	return halt, nil
}

// decode and execute one instruction
func execute() bool {
	inst := (*memory).Get(reg.pc)
	reg.pc++
	if inst == 0x76 { // halt
		return true
	}
	if inst == 0 { // nop
		return false
	}

	x, y, z := basicDecode(inst)
	switch x {
	case 0:
		decodeX0(y, z) // various
	case 1:
		store8r(load8r(z), y) // LD r[y], r[z]
	case 2:
		decodeX2(y, z) // alu[y] r[z]
	default:
		decodeX3(y, z) // various
	}
	return false
}

func SetStartAddress(addr uint16) {
	reg.pc = addr
}

// utility

func dumpRegs() {
	fmt.Printf("reg A: %d \n", reg.a)
	fmt.Printf("reg F: %d \n", reg.f)
}

func dump() {

}
