package test

import "fmt"

type RAM struct {
	Storage []byte
}

func init() {
	fmt.Println("--init--")
}

func (r RAM) Get(addr uint16) byte {
	return r.Storage[addr]
}

func (r RAM) Put(addr uint16, data byte) {
	r.Storage[addr] = data
}

func (r RAM) Load(addr uint16, block []byte) {
	length := uint16(len(block))
	for i := uint16(16); i < length; i++ {
		r.Storage[addr+i] = block[i]
	}
}
