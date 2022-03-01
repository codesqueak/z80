package test

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"z80/processor/pkg/hw"
)

// load  a .nas format file into memory
func loadFile(location string, mem *hw.Memory) {
	file, err := os.Open(location)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			// can't close file - bah!
		}
	}(file)
	var baseAddr uint16
	baseFound := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 10 {
			if !baseFound {
				baseFound = true
				base := line[0:4]
				addr, err := strconv.ParseUint(base, 16, 16)
				if err != nil {
					log.Fatal(err, "Bad file format in NAS_Test.nas")
				}
				baseAddr = uint16(addr)
			}
			line := line[5:28]
			values := strings.Split(line, " ")
			for _, b := range values {
				v, err := strconv.ParseUint(b, 16, 8)
				if err != nil {
					log.Fatal(err, "Bad file format in NAS_Test.nas")
				}
				(*mem).Put(baseAddr, byte(v))
				baseAddr++
			}
		}
	}
}

// various bits of utility code to help debug and test
func loadTestCode(mem *hw.Memory) {

	// A very simple I/O routine to simulate NAS-SYS character output call
	(*mem).Put(0x30, 0xD3) // out (00), a
	(*mem).Put(0x31, 0x00) //
	(*mem).Put(0x32, 0xC9) // ret
	//
	// A very simple I/O routine to simulate CP/M BDOS string output calls
	var addr uint16 = 0x05
	(*mem).Put(addr, 0xC3) // jp
	addr++
	(*mem).Put(addr, 0x45) //
	addr++
	(*mem).Put(addr, 0x00) //

	addr = 0x45            // get it out of the way of RST addresses
	(*mem).Put(addr, 0x79) // ld a,c
	addr++
	(*mem).Put(addr, 0xFE) // cp a, 09
	addr++
	(*mem).Put(addr, 0x09) //
	addr++
	(*mem).Put(addr, 0xCA) // jp z
	addr++
	(*mem).Put(addr, 0x4F) //
	addr++
	(*mem).Put(addr, 0x00) //
	addr++
	// output single char BDOS 2
	(*mem).Put(addr, 0x7B) // ld a,e
	addr++
	(*mem).Put(addr, 0xD3) // out (00), a
	addr++
	(*mem).Put(addr, 0x00) //
	addr++
	(*mem).Put(addr, 0xC9) // ret
	addr++
	// Output string BDOS 6
	// addr 000F
	(*mem).Put(addr, 0x1A) // ld a,(de)
	addr++
	(*mem).Put(addr, 0xFE) // cp a, '$'
	addr++
	(*mem).Put(addr, 0x24) //
	addr++
	(*mem).Put(addr, 0xC8) // ret z
	addr++
	(*mem).Put(addr, 0xD3) // out (00), a
	addr++
	(*mem).Put(addr, 0x00) //
	addr++
	(*mem).Put(addr, 0x13) // inc de
	addr++
	(*mem).Put(addr, 0xC3) // jp 0
	addr++
	(*mem).Put(addr, 0x4F) //
	addr++
	(*mem).Put(addr, 0x00) //
	addr++
	//
}

// generate a bit of code to write a string out to a port
func textWriter(text string, addr uint16, m *hw.Memory) uint16 {
	for _, c := range text {
		(*m).Put(addr, 0x3E)
		addr++
		(*m).Put(addr, byte(c))
		addr++
		(*m).Put(addr, 0xF7) // use rst 30 defined above
		addr++
	}
	return addr
}
