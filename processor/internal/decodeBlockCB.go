package internal

// CB-PREFIXED OPCODES
func decodeCB() {
	inst := (*memory).Get(reg.pc)
	reg.pc++
	x, y, z := basicDecode(inst)
	switch x {
	case 0: // rot[y] r[z]
		rotate(y, z)
	case 1: // BIT y, r[z]
		if z == 6 {
			bitInMemory(y, getHL())
		} else {
			bit(y, z)
		}
	case 2: // RES y, r[z]
		store8r(load8r(z)&resetBitTable[y], z)
	default: // SET y, r[z]
		store8r(load8r(z)|setBitTable[y], z)
	}
}

func rotate(y, z byte) {
	switch y {
	case 0: // RLC
		store8r(rlc(load8r(z)), z)
	case 1: // RRC
		store8r(rrc(load8r(z)), z)
	case 2: // RL
		store8r(rl(load8r(z)), z)
	case 3: //  RR
		store8r(rr(load8r(z)), z)
	case 4: // SLA
		store8r(sla(load8r(z)), z)
	case 5: // SRA
		store8r(sra(load8r(z)), z)
	case 6: // SLL
		store8r(sll(load8r(z)), z)
	default: // SRL
		store8r(srl(load8r(z)), z)
	}
}

func setShiftFlags(v byte) {
	resetH()
	resetN()
	setPVFromV(v)
	setUnusedFlagsFromV(v)
}

func rlc(v byte) byte {
	setCBool((v & 0x80) != 0)
	v = v << 1
	if getC() {
		v = v | 0x01
	}
	setSFromV(v)
	setZFromV(v)
	setShiftFlags(v)
	return v
}

func rrc(v byte) byte {
	setCBool((v & 0x01) != 0)
	v = v >> 1
	if getC() {
		v = v | 0x80
	}
	setSFromV(v)
	setZBool(v == 0)
	setShiftFlags(v)
	return v
}

func rl(v byte) byte {
	c := (v & 0x80) != 0
	v = v << 1
	if getC() {
		v = v | 0x01
	}
	setCBool(c)
	setSFromV(v)
	setZFromV(v)
	setShiftFlags(v)
	return v
}

func rr(v byte) byte {
	c := getC()
	setCBool((v & 0x01) != 0)
	v = v >> 1
	if c {
		v = v | 0x80
	}
	setSFromV(v)
	setZFromV(v)
	setShiftFlags(v)
	return v
}

func sla(v byte) byte {
	setCBool((v & 0x80) != 0)
	v = v << 1
	setSFromV(v)
	setZFromV(v)
	setShiftFlags(v)
	return v
}

func sra(v byte) byte {
	setCBool((v & 0x01) != 0)
	if (v & 0x80) == 0 {
		v = v >> 1
		resetS()
	} else {
		v = (v >> 1) | 0x80
		setS()
	}
	setZFromV(v)
	setShiftFlags(v)
	return v
}

func sll(v byte) byte {
	setCBool((v & 0x80) != 0)
	v = (v << 1) | 0x01 // bug in the silicon
	resetZ()            // can never be zero
	setSFromV(v)
	setShiftFlags(v)
	return v
}

func srl(v byte) byte {
	setCBool((v & 0x01) != 0)
	v = v >> 1
	resetS()
	setZFromV(v)
	setShiftFlags(v)
	return v
}

func bit(y, z byte) {
	v := load8r(z)
	setUnusedFlagsFromV(v) // very odd, only for direct reg access, not (rr)
	bitGeneric(v, y)
}

func bitInMemory(y byte, addr uint16) {
	v := (*memory).Get(addr)
	bitGeneric(v, y)
}

func bitGeneric(v, y byte) {
	v = v & setBitTable[y]
	setSBool((y == 7) && (v != 0))
	setH()
	setZFromV(v)
	setPVBool(v == 0)
	resetN()
}
