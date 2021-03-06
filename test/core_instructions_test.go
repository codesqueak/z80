package test

import (
	"fmt"
	"github.com/codesqueak/z80/processor/pkg"
	"github.com/codesqueak/z80/processor/pkg/hw"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInstructions(t *testing.T) {

	var mem hw.Memory = RAM{make([]byte, 65536)}
	var io hw.IO = PORTS{make([]byte, 256)}
	//
	for addr := uint16(0); addr < 0xFFFF; addr++ {
		mem.Put(addr, 0x76) // halt
	}
	//
	loadFile("testdata/NAS_Test.nas", &mem)
	loadTestCode(&mem)
	//
	err := pkg.Build(&mem, &io)
	if err != nil {
		fmt.Println("Failed to initialize processor. ", err)
		os.Exit(-1)
	}
	//
	addr := textWriter("Hello world!\n", 0x400, &mem)
	mem.Put(addr, 0x76)
	//
	pkg.SetStartAddress(0x1000)
	assert.Equal(t, uint16(0x1000), pkg.GetPC(), "Start address should have been set")
	assert.Equal(t, "        ", pkg.GetFlags(), "All flags shoudl be reset")
	pkg.AddressAndMem(0x0000)
	pkg.Line(0x1000)
	//
	for halt := false; !halt; {
		halt, err = pkg.RunOne()
		if err != nil {
			halt = true
			fmt.Println("CPU error")
		}

	}
	if err == nil {
		fmt.Println("CPU halted")
	}

}
