package pkg

import (
	"z80/processor/internal"
	"z80/processor/pkg/hw"
)

func Build(mem *hw.Memory, ports *hw.IO) error {
	return internal.Build(mem, ports)
}

// execute one instruction
func RunOne() (bool, error) {
	return internal.RunOne()
}

func SetStartAddress(addr uint16) {
	internal.SetStartAddress(addr)
}
