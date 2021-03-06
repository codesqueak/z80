package pkg

import (
	"github.com/codesqueak/z80/processor/internal"
	"github.com/codesqueak/z80/processor/pkg/hw"
)

func Build(mem *hw.Memory, ports *hw.IO) error {
	return internal.Build(mem, ports)
}

// RunOne execute one instruction
func RunOne() (bool, error) {
	return internal.RunOne()
}

// SetStartAddress Set the program counter
func SetStartAddress(addr uint16) {
	internal.SetStartAddress(addr)
}

// GetPC get the program counter
func GetPC() uint16 {
	return internal.GetPC()
}

func GetTStates() uint64 {
	return internal.GetTStates()
}

func ResetTStates() {
	internal.ResetTStates()
}

// GetFlags returns set flag values from SZ5H3PNC
func GetFlags() string {
	return internal.GetFlags()
}

// AddressAndMem outputs an address and the memory location it points to
func AddressAndMem(addr uint16) {
	internal.AddressAndMem(addr)
}

// Line outputs the 8 bytes in memory an address points to
func Line(addr uint16) {
	internal.Line(addr)
}
