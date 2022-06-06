package internal

import (
	"errors"
	"fmt"
	"github.com/codesqueak/z80/processor/pkg/hw"
)

var memory *hw.Memory
var io *hw.IO
var reg Registers
var initialized = false

// var count uint32

func Build(mem *hw.Memory, ports *hw.IO) error {
	if mem == nil {
		return errors.New("memory not defined")
	}
	if ports == nil {
		return errors.New("i/o not defined")
	}
	memory = mem
	io = ports
	reg = Registers{}
	initialized = true
	return nil
}

// execute one instruction
func RunOne() (bool, error) {
	if !initialized {
		return false, errors.New("the CPU is not initialized")
	}
	return execute(), nil
}

// decode and execute one instruction
// return halt state
func execute() bool {
	inst := (*memory).Get(reg.pc)
	reg.tStates += opcode8TStates[inst]
	reg.pc++
	if inst == 0x76 { // halt
		return true
	}
	if inst == 0 { // nop
		return false
	}

	// Initial instruction decode
	x, y, z := basicDecode(inst)
	switch x {
	case 0:
		decodeX0(y, z) // various
	case 1:
		store8r(load8r(z), y) //  LD r[y], r[z]
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

func GetPC() uint16 {
	return reg.pc
}

func GetTStates() uint64 {
	return reg.tStates
}

func ResetTStates() {
	reg.tStates = 0
}

// utility

func GetFlags() string {
	flagChar := "SZ5H3PNC"
	flags := ""
	mask := byte(0x80)
	for mask > 0 {
		if reg.f&mask != 0 {
			flags = flags + string(flagChar[0])
		} else {
			flags = flags + " "
		}
		flagChar = flagChar[1:]
		mask = mask >> 1
	}
	return flags
}

func AddressAndMem(addr uint16) {
	fmt.Printf("%04x(%02x) ", addr, (*memory).Get(addr))
}

func Line(addr uint16) {
	fmt.Printf("%04x ", addr)
	for i := uint16(0); i < 8; i++ {
		fmt.Printf("%02x ", (*memory).Get(addr+i))
	}
	fmt.Println()
}
