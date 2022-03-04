package internal

func decodeED() {
	inst := (*memory).Get(reg.pc)
	reg.pc++
	x, y, z := basicDecode(inst)
	switch x {
	case 0: // NONI
		return
	case 1:
		variousED(y, z)
	case 2:
		return
	default: // NONI
		return
	}
}

func variousED(y, z byte) {
	switch z {
	case 0: // IN r[y], (C) / IN (C)
		inC(y)
	case 1: // OUT (C), r[y] / OUT (C), 0
		outC(y)
	case 2: // SBC HL, rp[p] / ADC HL, rp[p]
		return
	case 3:
		return
	case 4:
		return
	case 5:
		return
	case 6:
		return
	default:
		return
	}
}

func inC(y byte) {
	v := (*io).Get(reg.c)
	if y != 6 { // special case - only changes flags
		store8r(v, y)
	}
	setSBool(v&0x80 != 0)
	setZBool(v == 0)
	setPVFromV(v)
	resetPV()
	resetN()
	resetH()
}

func outC(y byte) {
	if y == 6 {
		(*io).Put(reg.c, 0)
	} else {
		(*io).Put(reg.c, load8r(y))
	}
}

func sbcadchl(y byte) {
	p, q := getPQ(y)
	if q == 0 { // SBC HL, rp[p]
		hl := getHL()
		v := getRP(p)
		var c uint32 = 0
		if getC() {
			c = 1
		}
		ans32 := uint32(hl) - uint32(v) - c
		ans16 := uint16(ans32 & 0xFFFF)
		setSBool((ans16 & 0x8000) != 0)
		set3Bool((ans16 & 0x0800) != 0)
		set5Bool((ans16 & 0x2000) != 0)
		setZBool(ans16 == 0)
		setCBool(ans32 < 0)
		setOverflowFlagSub16(hl, v, c)
		setHBool((((hl & 0x0fff) - (v & 0x0fff) - uint16(c)) & 0x1000) != 0)
		setN()
		setHL(ans16)
	} else { // ADC HL, rp[p]
		hl := getHL()
		v := getRP(p)
		var c uint32 = 0
		if getC() {
			c = 1
		}
		ans32 := uint32(hl) + uint32(v) + c
		ans16 := uint16(ans32 & 0xFFFF)
		setSBool((ans16 & 0x8000) != 0)
		set3Bool((ans16 & 0x0800) != 0)
		set5Bool((ans16 & 0x2000) != 0)
		setZBool(ans16 == 0)
		setCBool(ans32 > 0xFFFF)
		setOverflowFlagAdd16(hl, v, c)
		setHBool((((hl & 0x0fff) + (v & 0x0fff) + uint16(c)) & 0x1000) != 0)
		setN()
		setHL(ans16)
	}
}

/* 2's compliment overflow flag control */
func setOverflowFlagSub16(rr, nn uint16, cc uint32) {
	left := int32(int16(rr))
	right := int32(int16(nn))
	carry := int32(cc)
	r := left - right - carry
	setPVBool((r < -32768) || (r > 32767))
}

/* 2's compliment overflow flag control */
func setOverflowFlagAdd16(rr, nn uint16, cc uint32) {
	left := int32(int16(rr))
	right := int32(int16(nn))
	carry := int32(cc)
	r := left + right + carry
	setPVBool((r < -32768) || (r > 32767))
}

func block(y, z byte) {

}
