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
var count uint32

func Build(mem *hw.Memory, ports *hw.IO) error {
	if mem == nil {
		return errors.New("Memory not defined")
	}
	if ports == nil {
		return errors.New("I/O not defined")
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
		return false, errors.New("CPU not initialized")
	}
	return execute(), nil
}

// decode and execute one instruction
// return halt state
func execute() bool {
	inst := (*memory).Get(reg.pc)

	//if count > 0x2643BC00 {
	//	fmt.Printf("%06x ", count)
	//	fmt.Printf("addr: %04x ", reg.pc)
	//	fmt.Printf("inst: %02x ", inst)
	//	fmt.Printf("A:%02x%02x ", reg.a, reg.f)
	//	fmt.Printf("BC:")
	//	regAndMem(getBC())
	//	fmt.Printf("DE:")
	//	regAndMem(getDE())
	//	fmt.Printf("HL:")
	//	regAndMem(getHL())
	//	fmt.Printf("SP:")
	//	fmt.Printf("%04x ", reg.sp)
	//	fmt.Printf(" " + getFlags() + "\n")
	//	//		line(0x2c80)
	//}
	//
	//count++
	//
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

// utility

func getFlags() string {
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

func regAndMem(addr uint16) {
	fmt.Printf("%04x(%02x) ", addr, (*memory).Get(addr))
}

func line(addr uint16) {
	fmt.Printf("%04x ", addr)
	for i := uint16(0); i < 8; i++ {
		fmt.Printf("%02x ", (*memory).Get(addr+i))
	}
	fmt.Println()
}
