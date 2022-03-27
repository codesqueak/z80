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
