package internal

func decodeX0(y, z byte) {
	switch z {
	case 0:
		relativeJumps0(y)
	case 1:
		loadAdd16Immediate1(y)
	case 2:
		indirectLoad2(y)
	case 3:
		incDec163(y)
	case 4:
		inc4(y)
	case 5:
		dec5(y)
	case 6:
		ld6(y)
	default:
		accFlagOps7(y)
	}
}

func relativeJumps0(y byte) {
	switch y {
	case 0: // nop
		return
	case 1: // ex af,af'
		t := reg.a
		reg.a = reg.a_
		reg.a_ = t
		t = reg.f
		reg.f = reg.f_
		reg.f_ = t
	case 2: // djnz dd
		reg.b = reg.b - 1
		if reg.b != 0 {
			relativeJump()
		} else {
			reg.pc++
		}
		return
	case 3: // jr dd
		relativeJump()
	default: // jr cc[y-4] dd
		if cc(y - 4) {
			relativeJump()
		} else {
			reg.pc++
		}
	}
}

func loadAdd16Immediate1(y byte) {
	p, q := getPQ(y)
	if q == 0 { // LD rp[p], nn
		setRP(p, load16FromRAM(reg.pc))
		reg.pc = reg.pc + 2
	} else { // ADD HL, rp[p]
		hl := getHL()
		right := getRP(p)
		result := hl + right
		//
		resetN()
		temp := (hl & 0x0FFF) + (right & 0x0FFF) // upper 8 half carry
		setHBool(temp&0xF000 != 0)
		set3Bool((result & 0x0800) != 0)
		set5Bool((result & 0x2000) != 0)
		setCBool(uint32(hl)+uint32(right) > 0xFFFF)
		setHL(result)
	}
}

func indirectLoad2(y byte) {
	p, q := getPQ(y)
	if q == 0 {
		switch p {
		case 0:
			(*memory).Put(getBC(), reg.a)
		case 1:
			(*memory).Put(getDE(), reg.a)
		case 2:
			store16ToRAM(load16FromPC(), getHL()) // LD (nn), HL
		default:
			(*memory).Put(load16FromPC(), reg.a) // LD (nn), A
		}
	} else {
		switch p {
		case 0:
			reg.a = (*memory).Get(getBC())
		case 1:
			reg.a = (*memory).Get(getDE())
		case 2:
			setHL(load16FromRAM(load16FromPC())) // LD HL, (nn)
		default:
			reg.a = (*memory).Get(load16FromPC()) // LD A, (nn)
		}
	}
}

func incDec163(y byte) {
	p, q := getPQ(y)
	v := getRP(p)
	if q == 0 {
		// inc 16
		setRP(p, v+1)
	} else {
		// dec 16
		setRP(p, v-1)
	}
}

// INC r[y]
func inc4(y byte) {
	v := load8r(y)
	setHalfCarryFlagAddValue(v, 1)
	setPVBool(v == 0x7F)
	v++
	store8r(v, y)
	setSFromV(v)
	setZFromV(v)
	resetN()
	setUnusedFlagsFromV(v)
}

// DEC r[y]
func dec5(y byte) {
	v := load8r(y)
	setHalfCarryFlagSubValue(v, 1)
	setPVBool(v == 0x80)
	v = v - 1
	store8r(v, y)
	setSFromV(v)
	setZFromV(v)
	setN()
	setUnusedFlagsFromV(v)
}

// LD r[y], n
func ld6(y byte) {
	v := (*memory).Get(reg.pc)
	reg.pc++
	store8r(v, y)
}

func accFlagOps7(y byte) {
	switch y {
	case 0: // RLCA
		carry := reg.a >= 0x80
		reg.a = reg.a << 1
		if carry {
			setC()
			reg.a = reg.a | 0x01
		} else {
			resetC()
		}
		resetH()
		resetN()
		setUnusedFlagsFromA()
	case 1: // RRCA
		carry := (reg.a & 0x01) != 0
		reg.a = reg.a >> 1
		if carry {
			setC()
			reg.a = reg.a | 0x80
		} else {
			resetC()
		}
		resetH()
		resetN()
		setUnusedFlagsFromA()
	case 2: // RLA
		carry := reg.a >= 0x80
		reg.a = reg.a << 1
		if getC() {
			reg.a = reg.a | 0x01
		}
		if carry {
			setC()
		} else {
			resetC()
		}
		resetH()
		resetN()
		setUnusedFlagsFromA()
	case 3: // RRA
		carry := (reg.a & 0x01) != 0
		reg.a = reg.a >> 1
		if getC() {
			reg.a = reg.a | 0x80
		}
		if carry {
			setC()
		} else {
			resetC()
		}
		resetH()
		resetN()
		setUnusedFlagsFromA()
	case 4: // DAA is weird, can't find Zilog algorithm so using +0110 if Nibble>9 algorithm.
		ans := reg.a
		var incr byte = 0
		carry := getC()
		if getH() || ((ans & 0x0f) > 0x09) {
			incr = 0x06
		}
		if carry || (ans > 0x9f) || ((ans > 0x8f) && ((ans & 0x0f) > 0x09)) {
			incr |= 0x60
		}
		if ans > 0x99 {
			carry = true
		}
		if getN() {
			alu8BitSub(incr) // sub_a(incr)
		} else {
			alu8BitAdd(incr) // add_a(incr)
		}
		if carry {
			setC()
		} else {
			resetC()
		}
		setPVFromA()
	case 5: // CPL
		reg.a = reg.a ^ 0xFF
		setH()
		setN()
		setUnusedFlagsFromA()
	case 6: // SCF
		setC()
		resetH()
		resetN()
		setUnusedFlagsFromA()
	default: // CCF
		if getC() {
			setH()
			resetC()
		} else {
			resetH()
			setC()
		}
		resetN()
		setUnusedFlagsFromA()
	}
}
