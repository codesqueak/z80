package internal

const OCT_DIGIT_Z byte = 0b00_000_111
const OCT_DIGIT_Y byte = 0b00_111_000
const OCT_DIGIT_X byte = 0b11_000_000
const MASK_P byte = 0b00000110
const MASK_Q byte = 0b00000001
const LSB uint16 = 0x00FF
const MSB uint16 = 0xFF00
const flagS byte = 0b1000_0000
const flagZ byte = 0b0100_0000
const flag5 byte = 0b0010_0000
const flagH byte = 0b0001_0000
const flag3 byte = 0b0000_1000
const flagPV byte = 0b0000_0100
const flagN byte = 0b0000_0010
const flagC byte = 0b0000_0001

// flag register bit positions for setting
const flag_S byte = 0x80
const flag_Z byte = 0x40
const flag_5 byte = 0x20
const flag_H byte = 0x10
const flag_3 byte = 0x08
const flag_PV byte = 0x04
const flag_N byte = 0x02
const flag_C byte = 0x01

// for resetting
const flag_S_N byte = flag_S ^ 0xFF
const flag_Z_N byte = flag_Z ^ 0xFF
const flag_5_N byte = flag_Z ^ 0xFF
const flag_H_N byte = flag_H ^ 0xFF
const flag_3_N byte = flag_3 ^ 0xFF
const flag_PV_N byte = flag_PV ^ 0xFF
const flag_N_N byte = flag_N ^ 0xFF
const flag_C_N byte = flag_C ^ 0xFF

var PARITY_TABLE [256]bool

func init() {
	PARITY_TABLE[0] = true // even PARITY_TABLE seed value
	position := 1          // table position
	for bit := 0; bit < 8; bit++ {
		for fill := 0; fill < position; fill++ {
			PARITY_TABLE[position+fill] = !PARITY_TABLE[fill]
		}
		position = position * 2
	}
}

func basicDecode(inst byte) (byte, byte, byte) {
	z := inst & OCT_DIGIT_Z
	y := (inst & OCT_DIGIT_Y) >> 3
	x := (inst & OCT_DIGIT_X) >> 6
	return x, y, z
}

func getPQ(y byte) (byte, byte) {
	return getP(y), getQ(y)
}

func getP(y byte) byte {
	return (y & MASK_P) >> 1
}

func getQ(y byte) byte {
	return y & MASK_Q
}

// relative jump using 8 bit signed offset
func relativeJump() {
	offset := uint16((*memory).Get(reg.pc))
	if offset > 0x7F {
		offset = offset | 0xFF00
	}
	reg.pc = reg.pc + 1 + offset
}

// condition evaluation on flag register
func cc(y byte) bool {
	switch y {
	case 0x0:
		return 0 == (reg.f & flagZ)
	case 0x1:
		return 0 != (reg.f & flagZ)
	case 0x2:
		return 0 == (reg.f & flagC)
	case 0x3:
		return 0 != (reg.f & flagC)
	case 0x4:
		return 0 == (reg.f & flagPV)
	case 0x5:
		return 0 != (reg.f & flagPV)
	case 0x6:
		return 0 == (reg.f & flagN)
	default:
		return 0 != (reg.f & flagN)
	}
}

// load 16 bit value lsb,msb
func load16FromRAM(addr uint16) uint16 {
	lsb := uint16((*memory).Get(addr))
	msb := uint16((*memory).Get(addr + 1))
	return (msb << 8) + lsb
}

// load 16 bit from PC address
func load16FromPC() uint16 {
	v := load16FromRAM(reg.pc)
	reg.pc = reg.pc + 2
	return v
}

// load 16 bit from SP address
func pop() uint16 {
	v := load16FromRAM(reg.sp)
	reg.sp = reg.sp + 2
	return v
}

// put 16 bit value into memory
func store16ToRAM(addr uint16, v uint16) {
	lsb := uint8(v & LSB)
	msb := uint8(v >> 8)
	(*memory).Put(addr, lsb)
	(*memory).Put(addr+1, msb)
}

// put 16 bit value onto teh stack
func push(v uint16) {
	lsb := uint8(v & LSB)
	msb := uint8(v >> 8)
	reg.sp--
	(*memory).Put(reg.sp, msb)
	reg.sp--
	(*memory).Put(reg.sp, lsb)
}

func getHL() uint16 {
	return uint16(reg.h)<<8 + uint16(reg.l)
}

func getDE() uint16 {
	return uint16(reg.d)<<8 + uint16(reg.e)
}

func getBC() uint16 {
	return uint16(reg.b)<<8 + uint16(reg.c)
}

func setHL(hl uint16) {
	reg.h = byte(hl >> 8)
	reg.l = byte(hl & LSB)
}

func setDE(de uint16) {
	reg.d = byte(de >> 8)
	reg.e = byte(de & LSB)
}

func setBC(bc uint16) {
	reg.d = byte(bc >> 8)
	reg.e = byte(bc & LSB)
}

func setAF(af uint16) {
	reg.a = byte(af >> 8)
	reg.f = byte(af & LSB)
}

func setRP(p byte, v uint16) {
	lsb := byte(v & LSB)
	msb := byte(v >> 8)
	switch p {
	case 0:
		reg.b = msb
		reg.c = lsb
	case 1:
		reg.d = msb
		reg.e = lsb
	case 2:
		reg.h = msb
		reg.l = lsb
	default:
		reg.sp = v
	}
}

func getRP(p byte) uint16 {
	switch p {
	case 0:
		return uint16(reg.b)<<8 + uint16(reg.c)
	case 1:
		return uint16(reg.d)<<8 + uint16(reg.e)
	case 2:
		return uint16(reg.h)<<8 + uint16(reg.l)
	default:
		return reg.sp
	}
}

func setRP2(p byte, v uint16) {
	lsb := byte(v & LSB)
	msb := byte(v >> 8)
	switch p {
	case 0:
		reg.b = msb
		reg.c = lsb
	case 1:
		reg.d = msb
		reg.e = lsb
	case 2:
		reg.h = msb
		reg.l = lsb
	default:
		reg.a = msb
		reg.f = lsb
	}
}

func getRP2(p byte) uint16 {
	switch p {
	case 0:
		return uint16(reg.b)<<8 + uint16(reg.c)
	case 1:
		return uint16(reg.d)<<8 + uint16(reg.e)
	case 2:
		return uint16(reg.h)<<8 + uint16(reg.l)
	default:
		return uint16(reg.a)<<8 + uint16(reg.f)
	}
}

func make16(lsb, msb byte) uint16 {
	return (uint16(msb) << 8) + uint16(lsb)
}

// 8 bit load
func load8r(y byte) byte {
	switch y {
	case 0:
		return reg.b
	case 1:
		return reg.c
	case 2:
		return reg.e
	case 3:
		return reg.e
	case 4:
		return reg.h
	case 5:
		return reg.l
	case 6:
		return (*memory).Get(make16(reg.h, reg.l))
	default:
		return reg.a

	}
}

// 8 bit store
func store8r(v byte, y byte) {
	switch y {
	case 0:
		reg.b = v
	case 1:
		reg.c = v
	case 2:
		reg.d = v
	case 3:
		reg.e = v
	case 4:
		reg.h = v
	case 5:
		reg.l = v
	case 6:
		(*memory).Put(make16(reg.h, reg.l), v)
	case 7:
		reg.a = v
	}
}

// flag manipulation

// set
func setN() {
	reg.f = reg.f | flag_N
}

func setH() {
	reg.f = reg.f | flag_H
}

func set3() {
	reg.f = reg.f | flag_3
}

func set5() {
	reg.f = reg.f | flag_5
}

func setC() {
	reg.f = reg.f | flag_C
}

func setPV() {
	reg.f = reg.f | flag_PV
}

func setS() {
	reg.f = reg.f | flag_S
}

func setZ() {
	reg.f = reg.f | flag_Z
}

// reset
func resetN() {
	reg.f = reg.f & flag_N_N
}

func resetH() {
	reg.f = reg.f & flag_H_N
}

func reset3() {
	reg.f = reg.f & flag_3_N
}

func reset5() {
	reg.f = reg.f & flag_5_N
}

func resetC() {
	reg.f = reg.f & flag_C_N
}

func resetPV() {
	reg.f = reg.f & flag_PV_N
}

func resetS() {
	reg.f = reg.f & flag_S_N
}

func resetZ() {
	reg.f = reg.f & flag_Z_N
}

// set to boolean

func setPVFromA() {
	setPVFromV(reg.a)
}

func setPVFromV(v byte) {
	setPVBool(PARITY_TABLE[v])
}

func setPVBool(b bool) {
	if b {
		setH()
	} else {
		resetH()
	}
}

func setHBool(b bool) {
	if b {
		setH()
	} else {
		resetH()
	}
}

func setSFromA() {
	if reg.a >= 0x80 {
		setS()
	} else {
		resetS()
	}
}

func setZFromA() {
	if reg.a == 0 {
		setZ()
	} else {
		resetZ()
	}
}

func setCBool(b bool) {
	if b {
		setC()
	} else {
		resetC()
	}
}

func setSBool(b bool) {
	if b {
		setS()
	} else {
		resetS()
	}
}

func set3Bool(b bool) {
	if b {
		set3()
	} else {
		reset3()
	}
}

func set5Bool(b bool) {
	if b {
		set5()
	} else {
		reset5()
	}
}

// check flag status

func getC() bool {
	return (reg.f & flagC) != 0
}

func getH() bool {
	return (reg.f & flagH) != 0
}

func getN() bool {
	return (reg.f & flagN) != 0
}

// set unused flags (3 & 5) to bits in result
// Z80 ALU weirdness
func setUnusedFlags() {
	v := reg.a & 0x28
	reg.f = reg.f & 0xD7
	reg.f = reg.f | v
}

/* half carry flag control */
func setHalfCarryFlagAddCarry(right, carry byte) {
	left := reg.a & 0x0f
	right = right & 0x0f
	if (right + left + carry) > 0x0f {
		setH()
	} else {
		resetH()
	}
}

/* half carry flag control */

func setHalfCarryFlagAdd(right byte) {
	setHalfCarryFlagAddValue(reg.a, right)
}

func setHalfCarryFlagAddValue(left, right byte) {
	left = left & 0x0F
	right = right & 0x0F
	if (right + left) > 0x0f {
		setH()
	} else {
		resetH()
	}
}

func setHalfCarryFlagSub(right byte) {
	setHalfCarryFlagSubValue(reg.a, right)
}

func setHalfCarryFlagSubValue(left, right byte) {
	left = left & 0x0F
	right = right & 0x0F
	if left < right {
		setH()
	} else {
		resetH()
	}
}

/* half carry flag control */
func setHalfCarryFlagSubCarry(left, right, carry byte) {
	left = left & 0x0F
	right = right & 0x0F
	if left < (right + carry) {
		setH()
	} else {
		resetH()
	}
}

/* 2's compliment overflow flag control */
//
// Find better way to do these ops !!!
//
func setOverflowFlagAdd(rv byte, c bool) {
	left := int16(reg.a)
	right := int16(rv)
	if left > 127 {
		left = left - 256
	}
	if right > 127 {
		right = right - 256
	}
	left = left + right
	if c {
		left++
	}
	setPVBool((left < -128) || (left > 127))
}

func setOverflowFlagSub(rv byte, c bool) {
	left := int16(reg.a)
	right := int16(rv)
	if left > 127 {
		left = left - 256
	}
	if right > 127 {
		right = right - 256
	}
	left = left - right
	if c {
		left--
	}
	setPVBool((left < -128) || (left > 127))
}
