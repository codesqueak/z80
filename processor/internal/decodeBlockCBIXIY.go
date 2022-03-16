package internal

func decodeCBIXIY() {
	addr := getIndex() + getIXIY()
	inst := (*memory).Get(reg.pc)
	// fmt.Printf("inst %04x \n", inst)
	reg.pc++
	x, y, z := basicDecode(inst)
	switch x {
	case 0:
		ldRotIXIY(addr, y, z)
	case 1:
		bitIXIY(addr, y)
	case 2:
		ldResIXIY(addr, y, z)
	default:
		ldSetIXIY(addr, y, z)
	}
}

// very odd instructions
// RLC  (IX+nn), followed by LD   rr,(IX+nn), but not if rr = 6
func ldRotIXIY(addr uint16, y, z byte) {
	v := (*memory).Get(addr)
	switch y {
	case 0: // RLC
		v = rlc(v)
	case 1: // RRC
		v = rrc(v)
	case 2: // RL
		v = rl(v)
	case 3: //  RR
		v = rr(v)
	case 4: // SLA
		v = sla(v)
	case 5: // SRA
		v = sra(v)
	case 6: // SLL
		v = sll(v)
	default: // SRL
		v = srl(v)
	}
	(*memory).Put(addr, v)
	if z != 6 {
		store8r(v, z)
	}
}

// BIT y, (IX+d)
func bitIXIY(addr uint16, y byte) {
	v := (*memory).Get(addr)
	v = bitGeneric(v, y)
	(*memory).Put(addr, v)
}

// LD r[z], RES y, (IX+d)
func ldResIXIY(addr uint16, y, z byte) {
	v := (*memory).Get(addr)
	v = v & resetBitTable[y]
	(*memory).Put(addr, v)
	if z != 6 {
		store8r(v, z)
	}
}

// LD r[z], SET y, (IX+d)
func ldSetIXIY(addr uint16, y, z byte) {
	v := (*memory).Get(addr)
	v = v & setBitTable[y]
	(*memory).Put(addr, v)
	if z != 6 {
		store8r(v, z)
	}
}
