package internal

type Registers struct {
	// 8 bit
	a, b, c, d, e, h, l byte
	// alt 8 bit
	a_, b_, c_, d_, e_, h_, l_ byte
	// index
	ix, iy, pc, sp uint16
	// control
	f, f_ byte
	// internal
	i, r           byte
	iff1, iff2     bool
	nmi_ff         bool
	interruptMode  byte
	ddMode, fdMode bool
}
