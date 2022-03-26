package hw

// Interface definitions

type Memory interface {
	// Get byte from memory
	Get(addr uint16) byte

	// Put byte into memory
	Put(addr uint16, data byte)

	// Load block data into memory
	Load(addr uint16, block []byte)
}

type IO interface {
	// Get a byte from an I/O port
	Get(addr byte) byte

	// Put a byte to an I/O port
	Put(addr byte, data byte)
}
