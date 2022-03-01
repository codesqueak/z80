package hw

type Memory interface {
	Get(addr uint16) byte

	Put(addr uint16, data byte)

	Load(addr uint16, block []byte)
}

type IO interface {
	Get(addr byte) byte

	Put(addr byte, data byte)
}
