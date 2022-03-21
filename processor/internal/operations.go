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
const flag_5_N byte = flag_5 ^ 0xFF
const flag_H_N byte = flag_H ^ 0xFF
const flag_3_N byte = flag_3 ^ 0xFF
const flag_PV_N byte = flag_PV ^ 0xFF
const flag_N_N byte = flag_N ^ 0xFF
const flag_C_N byte = flag_C ^ 0xFF

var PARITY_TABLE [256]bool
var setBitTable = [8]byte{0x01, 0x02, 0x04, 0x08, 0x10, 0x20, 0x40, 0x80}
var resetBitTable = [8]byte{0xFE, 0xFD, 0xFB, 0xF7, 0xEF, 0xDF, 0xBF, 0x7F}

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
	case 0x0: // NZ
		return 0 == (reg.f & flagZ)
	case 0x1: // Z
		return 0 != (reg.f & flagZ)
	case 0x2: // NC
		return 0 == (reg.f & flagC)
	case 0x3: // C
		return 0 != (reg.f & flagC)
	case 0x4: // PO
		return 0 == (reg.f & flagPV)
	case 0x5: // PE
		return 0 != (reg.f & flagPV)
	case 0x6: // P
		return 0 == (reg.f & flagS)
	default: // M
		return 0 != (reg.f & flagS)
	}
}

// load 16 bit value lsb,msb
func load16FromRAM(addr uint16) uint16 {
	lsb := uint16((*memory).Get(addr))
	msb := uint16((*memory).Get(addr + 1))
	return (msb << 8) | lsb
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

// put 16 bit value onto the stack
func push(v uint16) {
	lsb := uint8(v & LSB)
	msb := uint8(v >> 8)
	reg.sp--
	(*memory).Put(reg.sp, msb)
	reg.sp--
	(*memory).Put(reg.sp, lsb)
}

func getHL() uint16 {
	return uint16(reg.h)<<8 | uint16(reg.l)
}

func getIXIY() uint16 {
	if reg.ddMode {
		return reg.ix
	} else {
		return reg.iy
	}
}

func getDE() uint16 {
	return uint16(reg.d)<<8 + uint16(reg.e)
}

func getBC() uint16 {
	return uint16(reg.b)<<8 + uint16(reg.c)
}

func getAF() uint16 {
	return uint16(reg.a)<<8 + uint16(reg.f)
}

func setHL(hl uint16) {
	reg.h = byte(hl >> 8)
	reg.l = byte(hl & LSB)
}

func setIXIY(ixiy uint16) {
	if reg.ddMode {
		reg.ix = ixiy
	} else {
		reg.iy = ixiy
	}
}

func setDE(de uint16) {
	reg.d = byte(de >> 8)
	reg.e = byte(de & LSB)
}

func setBC(bc uint16) {
	reg.b = byte(bc >> 8)
	reg.c = byte(bc & LSB)
}

func setAF(af uint16) {
	reg.a = byte(af >> 8)
	reg.f = byte(af & LSB)
}

func setRP(p byte, v uint16) {
	switch p {
	case 0:
		setBC(v)
	case 1:
		setDE(v)
	case 2:
		setHL(v)
	default:
		reg.sp = v
	}
}

func setRPIXIY(p byte, v uint16) {
	switch p {
	case 0:
		setBC(v)
	case 1:
		setDE(v)
	case 2:
		setIXIY(v)
	default:
		reg.sp = v
	}
}

func getRP(p byte) uint16 {
	switch p {
	case 0:
		return getBC()
	case 1:
		return getDE()
	case 2:
		return getHL()
	default:
		return reg.sp
	}
}

func getRPIXIY(p byte) uint16 {
	switch p {
	case 0:
		return getBC()
	case 1:
		return getDE()
	case 2:
		return getIXIY()
	default:
		return reg.sp
	}
}

func setRP2(p byte, v uint16) {
	switch p {
	case 0:
		setBC(v)
	case 1:
		setDE(v)
	case 2:
		setHL(v)
	default:
		setAF(v)
	}
}

func setRP2IXIY(p byte, v uint16) {
	switch p {
	case 0:
		setBC(v)
	case 1:
		setDE(v)
	case 2:
		setIXIY(v)
	default:
		setAF(v)
	}
}

func getRP2(p byte) uint16 {
	switch p {
	case 0:
		return getBC()
	case 1:
		return getDE()
	case 2:
		return getHL()
	default:
		return getAF()
	}
}

func getRP2IXIY(p byte) uint16 {
	switch p {
	case 0:
		return getBC()
	case 1:
		return getDE()
	case 2:
		return getIXIY()
	default:
		return getAF()
	}
}

func make16(msb, lsb byte) uint16 {
	return (uint16(msb) << 8) + uint16(lsb)
}

// get the +dd part of an index address
func getIndex() uint16 {
	dd := uint16((*memory).Get(reg.pc))
	if dd > 0x7F {
		dd = dd | 0xFF00
	}
	return dd
}

// 8 bit load
func load8r(y byte) byte {
	switch y {
	case 0:
		return reg.b
	case 1:
		return reg.c
	case 2:
		return reg.d
	case 3:
		return reg.e
	case 4:
		return reg.h
	case 5:
		return reg.l
	case 6:
		return (*memory).Get(getHL())
	default:
		return reg.a
	}
}

// 8 bit load
func load8rIXIY(y byte) (byte, bool) {
	switch y {
	case 0:
		return reg.b, false
	case 1:
		return reg.c, false
	case 2:
		return reg.d, false
	case 3:
		return reg.e, false
	case 4:
		if reg.ddMode {
			return byte((reg.ix & 0xFF00) >> 8), false
		} else {
			return byte((reg.iy & 0xFF00) >> 8), false
		}
	case 5:
		if reg.ddMode {
			return byte(reg.ix), false
		} else {
			return byte(reg.iy), false
		}
	case 6:
		addr := indexedAddr()
		return (*memory).Get(addr), true
	default:
		return reg.a, false
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
		(*memory).Put(getHL(), v)
	case 7:
		reg.a = v
	}
}

// 8 bit store
func store8rIXIY(v byte, y byte) bool {
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
		if reg.ddMode {
			reg.ix = (reg.ix & 0x00FF) | (uint16(v) << 8)
		} else {
			reg.iy = (reg.iy & 0x00FF) | (uint16(v) << 8)
		}
	case 5:
		if reg.ddMode {
			reg.ix = (reg.ix & 0xFF00) | uint16(v)
		} else {
			reg.iy = (reg.iy & 0xFF00) | uint16(v)
		}
	case 6:
		addr := indexedAddr()
		(*memory).Put(addr, v)
		return true
	case 7:
		reg.a = v
	}
	return false
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
		setPV()
	} else {
		resetPV()
	}
}

func setHBool(b bool) {
	if b {
		setH()
	} else {
		resetH()
	}
}

func setSFromV(v byte) {
	if (v & 0x80) != 0 {
		setS()
	} else {
		resetS()
	}
}

func setSFromA() {
	setSFromV(reg.a)
}

func setZFromV(v byte) {
	if v == 0 {
		setZ()
	} else {
		resetZ()
	}
}

func setZFromA() {
	setZFromV(reg.a)
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

func setZBool(b bool) {
	if b {
		setZ()
	} else {
		resetZ()
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

func getZ() bool {
	return (reg.f & flagZ) != 0
}

func getH() bool {
	return (reg.f & flagH) != 0
}

func getN() bool {
	return (reg.f & flagN) != 0
}

// set unused flags (3 & 5) to bits in result
// Z80 ALU weirdness
func setUnusedFlagsFromA() {
	setUnusedFlagsFromV(reg.a)
}

func setUnusedFlagsFromV(v byte) {
	v = v & 0x28
	reg.f = reg.f & 0xD7
	reg.f = reg.f | v
}

/* half carry flag control */
func setHalfCarryFlagAddCarry(right, carry byte) {
	left := reg.a & 0x0f
	right = right & 0x0f
	setHBool((right+left+carry)&0xF0 != 0)
}

/* half carry flag control */

func setHalfCarryFlagAdd(right byte) {
	setHalfCarryFlagAddValue(reg.a, right)
}

func setHalfCarryFlagAddValue(left, right byte) {
	left = left & 0x0F
	right = right & 0x0F
	setHBool((right+left)&0xF0 != 0)
}

func setHalfCarryFlagSub(right byte) {
	setHalfCarryFlagSubValue(reg.a, right)
}

func setHalfCarryFlagSubValue(left, right byte) {
	left = left & 0x0F
	right = right & 0x0F
	setHBool((left-right)&0xF0 != 0)
}

/* half carry flag control */
func setHalfCarryFlagSubCarry(left, right, carry byte) {
	left = left & 0x0F
	right = right & 0x0F
	setHBool((left-right-carry)&0xF0 != 0)
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

/* 2's compliment overflow flag control */
func setOverflowFlagSub16(rr, nn uint16, cc uint32) {
	left := int32(rr)
	right := int32(nn)
	carry := int32(cc)
	if left > 32767 {
		left = left - 65536
	}
	if right > 32767 {
		right = right - 65536
	}
	r := left - right - carry
	setPVBool((r < -32768) || (r > 32767))
}

/* 2's compliment overflow flag control */
func setOverflowFlagAdd16(rr, nn uint16, cc uint32) {
	left := int32(rr)
	right := int32(nn)
	if left > 32767 {
		left = left - 65536
	}
	if right > 32767 {
		right = right - 65536
	}
	carry := int32(cc)
	r := left + right + carry
	setPVBool((r < -32768) || (r > 32767))
}
