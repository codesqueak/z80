package internal

// Block 2 instructions
func decodeX2(y, z byte) {
	v := load8r(z)
	switch y {
	case 0: // add
		alu8BitAdd(v)
	case 1: // adc
		alu8BitAdc(v)
	case 2: // sub
		alu8BitSub(v)
	case 3: // sbc
		alu8BitSbc(v)
	case 4: // and
		alu8BitAnd(v)
	case 5: // xor
		alu8BitXor(v)
	case 6: // or
		alu8BitOr(v)
	default: // cp
		alu8BitCp(v)
	}
}

/* 8 bit ADD */
func alu8BitAdd(v byte) {
	setHalfCarryFlagAdd(v)
	setOverflowFlagAdd(v, false)
	setCBool((uint16(reg.a) + uint16(v)) > 0x00FF)
	reg.a = reg.a + v
	resetN()
	aluSetStandardFlags()
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
	resetN()
	aluSetStandardFlags()
}

/* 8 bit SUB */
func alu8BitSub(v byte) {
	setHalfCarryFlagSub(v)
	setOverflowFlagSub(v, false)
	setCBool(v > reg.a)
	reg.a = reg.a - v
	setN()
	aluSetStandardFlags()
}

/* 8 bit SBC */
func alu8BitSbc(v byte) {
	var c byte
	if getC() {
		c = 1
	}
	setHalfCarryFlagSubCarry(reg.a, v, c)
	setOverflowFlagSub(v, c == 1)
	setCBool(uint16(v)+uint16(c) > uint16(reg.a))
	reg.a = reg.a - v - c
	setN()
	aluSetStandardFlags()
}

/* 8 bit AND  */
func alu8BitAnd(v byte) {
	reg.f = flag_H // set the H flag
	reg.a = reg.a & v
	setPVFromA()
	aluSetStandardFlags()
}

/* 8 bit OR  */
func alu8BitOr(v byte) {
	reg.f = 0
	reg.a = reg.a | v
	setPVFromA()
	aluSetStandardFlags()
}

/* 8 bit XOR  */
func alu8BitXor(v byte) {
	reg.f = 0
	reg.a = reg.a ^ v
	setPVFromA()
	aluSetStandardFlags()
}

/* 8 bit CP */
func alu8BitCp(v byte) {
	setHalfCarryFlagSub(v)
	setOverflowFlagSub(v, false)
	setCBool(v > reg.a)
	r := reg.a - v
	setN()
	setSFromV(r)
	setZFromV(r)
	setUnusedFlagsFromV(v)
}

func aluSetStandardFlags() {
	setSFromA()
	setZFromA()
	setUnusedFlagsFromA()
}
