package internal

// // LD r[y], r[z]
func ldryrzIXIY(y, z byte) {
	if z == 6 { // Use h,l not ixh, ixl as we have an (ix)
		addr := indexedAddr()
		v := (*memory).Get(addr)
		store8r(v, y) // ld r[y], (ix+dd)
	} else {
		if y == 6 { // Use h,l not ixh, ixl as we have an (ix)
			addr := indexedAddr()
			v := load8r(z)
			(*memory).Put(addr, v) // ld (ix+dd),r[z]
		} else {
			store8rIXIY(load8rIXIY(z), y)
		}
	}
}

func indexedAddr() uint16 {
	addr := getIndex()
	if reg.ddMode {
		return addr + reg.ix
	} else {
		return addr + reg.iy
	}
}
