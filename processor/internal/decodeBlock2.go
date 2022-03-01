package internal

func decodeX2(y, z byte) {
	switch y {
	case 0: // add
		alu8BitAdd(load8r(z))
	case 1: // adc
		alu8BitAdc(load8r(z))
	case 2: // sub
		alu8BitSub(load8r(z))
	case 3: // sbc
		alu8BitSbc(load8r(z))
	case 4: // and
		alu8BitAnd(load8r(z))
	case 5: // xor
		alu8BitXor(load8r(z))
	case 6: // or
		alu8BitOr(load8r(z))
	default: // cp
		alu8BitCp(load8r(z))
	}
}

/* 8 bit ADD */
func alu8BitAdd(v byte) {
	setHalfCarryFlagAdd(v)
	setOverflowFlagAdd(v, false)
	setCBool(uint16(reg.a)+uint16(v) > 0x00FF)
	reg.a = reg.a + v
	setSFromA()
	setZFromA()
	resetN()
	setUnusedFlags()
}

/* 8 bit ADC */
func alu8BitAdc(v byte) {
	var c byte
	if getC() {
		c = 1
	}
	setHalfCarryFlagAddCarry(v, c)
	setOverflowFlagAdd(v, c == 1)
	setCBool(uint16(reg.a)+uint16(v)+uint16(c) > 0x00FF)
	reg.a = reg.a + v + c
	setSFromA()
	setZFromA()
	resetN()
	setUnusedFlags()
}

/* 8 bit SUB */
func alu8BitSub(v byte) {
	setHalfCarryFlagSub(v)
	setOverflowFlagSub(v, false)
	setCBool(v > reg.a)
	reg.a = reg.a - v
	setSFromA()
	setZFromA()
	setN()
	setUnusedFlags()
}

/* 8 bit SBC */
func alu8BitSbc(v byte) {
	var c byte
	if getC() {
		c = 1
	}
	setHalfCarryFlagSubCarry(reg.a, v, c)
	setOverflowFlagSub(v, c == 1)
	setCBool(v+c > reg.a)
	setSFromA()
	setZFromA()
	setN()
	setUnusedFlags()
}

/* 8 bit AND  */
func alu8BitAnd(v byte) {
	reg.f = flag_H // set the H flag
	reg.a = reg.a & v
	setSFromA()
	setZFromA()
	setPVFromA()
	setUnusedFlags()
}

/* 8 bit OR  */
func alu8BitOr(v byte) {
	reg.f = 0
	reg.a = reg.a | v
	setSFromA()
	setZFromA()
	setPVFromA()
	setUnusedFlags()
}

/* 8 bit XOR  */
func alu8BitXor(v byte) {
	reg.f = 0
	reg.a = reg.a ^ v
	setSFromA()
	setZFromA()
	setPVFromA()
	setUnusedFlags()
}

/* 8 bit CP */
func alu8BitCp(v byte) {
	a := reg.a
	reg.a = a - v
	reg.f = flag_N
	setCBool(v > a)
	setSFromA()
	set3Bool((v & flag_3) != 0)
	set5Bool((v & flag_5) != 0)
	setZFromA()
	setHBool((((a & 0x0f) - (v & 0x0f)) & flag_H) != 0)
	setPVBool(((a ^ v) & (a ^ reg.a) & 0x80) != 0)
}
