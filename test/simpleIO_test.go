package test

import "fmt"

type PORTS struct {
	Ports []byte
}

func (io PORTS) Get(addr byte) byte {
	return io.Ports[addr]
}

func (io PORTS) Put(addr byte, data byte) {
	io.Ports[addr] = data

	if data < 32 {
		fmt.Println()
	} else {
		fmt.Print(string(data))
	}
}
