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
		bit(y, z)
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

func rlc(v byte) byte {
	setCBool(v >= 0x80)
	v = v << 1
	setSBool(v >= 0x80)
	setZBool(v == 0)
	setShiftFlags(v)
	return v
}

func rrc(v byte) byte {
	setCBool((v & 0x01) != 0)
	v = v >> 1
	if getC() {
		v = v | 0x80
	}
	setSBool(v >= 0x80)
	setZBool(v == 0)
	setShiftFlags(v)
	return v
}

func rl(v byte) byte {
	c := (v >= 0x80)
	v = v << 1
	if getC() {
		v = v | 0x01
	}
	setCBool(c)
	setSBool(v >= 0x80)
	setZBool(v == 0)
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
	setSBool(v >= 0x80)
	setZBool(v == 0)
	setShiftFlags(v)
	return v
}

func sla(v byte) byte {
	setCBool((v & 0x01) != 0)
	v = v << 1
	resetC()
	setSBool(v >= 0x80)
	setZBool(v == 0)
	setShiftFlags(v)
	return v
}

func sra(v byte) byte {
	setCBool((v & 0x0001) != 0)
	if (v & 0x80) == 0 {
		v = v >> 1
		resetS()
	} else {
		v = (v >> 1) | 0x0080
		setS()
	}
	setZBool(v == 0)
	setShiftFlags(v)
	return v
}

func sll(v byte) byte {
	setCBool(v >= 0x80)
	v = (v << 1) | 0x01 // bug in the silicon
	resetZ()            // can never be zero
	setSBool(v >= 0x80)
	setShiftFlags(v)
	return v
}

func srl(v byte) byte {
	setCBool((v & 0x0001) != 0)
	v = v >> 1
	resetS()
	setZBool(v == 0)
	setShiftFlags(v)
	return v
}

func setShiftFlags(v byte) {
	resetH()
	resetN()
	setPVFromV(v)
	setUnusedFlagsFromV(v)
}

func bit(y, z byte) {
	v := load8r(z)
	resetS()
	set3Bool((v & flag_3) != 0)
	set5Bool((v & flag_5) != 0)
	v = v & setBitTable[y]
	store8r(v, z)
	setZ()
	setPV()
	resetN()
	setH()
}
