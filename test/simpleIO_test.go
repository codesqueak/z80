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
		fmt.Printf("\n")
	} else {
		if data < 128 {
			fmt.Print(string(data))
		} else {
			fmt.Print("?")
		}
	}
}
