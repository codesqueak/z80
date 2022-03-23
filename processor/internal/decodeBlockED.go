package internal

func decodeED() {
	inst := (*memory).Get(reg.pc)
	//	fmt.Printf("inst ED : %x\n", inst)
	reg.pc++
	x, y, z := basicDecode(inst)
	switch x {
	case 0: // NONI
		return
	case 1:
		variousED(y, z)
	case 2:
		block(y, z)
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
		sbcadchl(y)
	case 3: // 	LD (nn), rp[p] / LD rp[p], (nn)
		ld16Indirect(y)
	case 4: // NEG
		neg()
	case 5: // RETN / RETI
		retIN(y)
	case 6: // IM im[y]
		im(y)
	default: // LD / RRD / RLD
		ldrrdrld(y)
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
		rp := getRP(p)
		var c uint32 = 0
		if getC() {
			c = 1
		}
		ans32 := uint32(hl) - uint32(rp) - c
		ans16 := uint16(ans32)
		setSBool((ans16 & 0x8000) != 0)
		set3Bool((ans16 & 0x0800) != 0)
		set5Bool((ans16 & 0x2000) != 0)
		setZBool(ans16 == 0)
		setCBool(ans32 > 0xFFFF)
		setOverflowFlagSub16(hl, rp, c)
		setHBool((hl&0x0fff)-(rp&0x0fff)-uint16(c) >= 0x1000)
		setN()
		setHL(ans16)
	} else { // ADC HL, rp[p]
		hl := getHL()
		rp := getRP(p)
		var c uint32 = 0
		if getC() {
			c = 1
		}
		ans32 := uint32(hl) + uint32(rp) + c
		ans16 := uint16(ans32)
		setSBool((ans16 & 0x8000) != 0)
		set3Bool((ans16 & 0x0800) != 0)
		set5Bool((ans16 & 0x2000) != 0)
		setZBool(ans16 == 0)
		setCBool(ans32 > 0xFFFF)
		setOverflowFlagAdd16(hl, rp, c)
		setHBool((hl&0x0fff)+(rp&0x0fff)+uint16(c) >= 0x1000)
		resetN()
		setHL(ans16)
	}
}

func ld16Indirect(y byte) {
	p, q := getPQ(y)
	if q == 0 { // LD (nn), rp[p]
		addr := load16FromPC()
		v := getRP(p)
		store16ToRAM(addr, v)
	} else { // LD rp[p], (nn)
		addr := load16FromPC()
		v := load16FromRAM(addr)
		setRP(p, v)
	}
}

func neg() {
	v := reg.a
	setHBool((v & 0x0f) != 0x00)
	setPVBool(v == 0x80)
	setCBool(v != 0)
	reg.a = 0 - v
	setZFromA()
	setSFromA()
	setN()
	setUnusedFlagsFromA()
}

func retIN(y byte) {
	reg.pc = pop()
	if y != 1 { // RETN
		reg.iff1 = reg.iff2
	}
}

func im(y byte) {
	switch y {
	case 0, 4:
		reg.interruptMode = 0
	case 1, 2, 5, 6:
		reg.interruptMode = 1
	case 3, 7:
		reg.interruptMode = 2
	}
}

func ldrrdrld(y byte) {
	switch y {
	case 0: // LD I, A
		reg.i = reg.a
	case 1: // LD R, A
		reg.r = reg.a
	case 2: // LD A, I
		ldai()
	case 3: // LD A, R
		ldar()
	case 4: // RRD
		rrd()
	case 5: // RLD
		rld()
	case 6: // NOP
		return
	default: // NOP
		return
	}
}

func ldai() {
	reg.a = reg.i
	setSFromA()
	setZFromA()
	resetH()
	resetN()
	setPVBool(reg.iff2)
	setUnusedFlagsFromA()
}

func ldar() {
	reg.a = reg.r & 0x7F
	resetS()
	setZFromA()
	resetH()
	resetN()
	setPVBool(reg.iff2)
	setUnusedFlagsFromA()
}

func rrd() {
	temp := (*memory).Get(getHL())
	nibble1 := (reg.a & 0x00F0) >> 4
	nibble2 := reg.a & 0x000F
	nibble3 := (temp & 0x00F0) >> 4
	nibble4 := temp & 0x000F
	//
	reg.a = (nibble1 << 4) | nibble4
	temp = (nibble2 << 4) | nibble3
	(*memory).Put(getHL(), temp)
	//
	setSFromA()
	setZFromA()
	resetH()
	setPVFromA()
	resetN()
	setUnusedFlagsFromA()
}

func rld() {
	temp := (*memory).Get(getHL())
	nibble1 := (reg.a & 0x00F0) >> 4
	nibble2 := reg.a & 0x000F
	nibble3 := (temp & 0x00F0) >> 4
	nibble4 := temp & 0x000F
	//
	reg.a = (nibble1 << 4) | nibble3
	temp = (nibble4 << 4) | nibble2
	//
	(*memory).Put(getHL(), temp)
	//
	setSFromA()
	setZFromA()
	resetH()
	setPVFromA()
	resetN()
	setUnusedFlagsFromA()
}

// block moves
func block(y, z byte) {
	if (z <= 3) && (y >= 4) {
		switch z {
		case 0:
			blockLD(y)
		case 1:
			blockCP(y)
		case 2:
			blockIN(y)
		default:
			blockOUT(y)
		}
	} else {
		return // NOP
	}
}

func blockLD(y byte) {
	switch y {
	case 4:
		LDI()
	case 5:
		LDD()
	case 6:
		LDIR()
	default:
		LDDR()
	}
}

func LDI() {
	v := (*memory).Get(getHL())
	(*memory).Put(getDE(), v)
	setHL(getHL() + 1)
	setDE(getDE() + 1)
	setBC(getBC() - 1)
	resetH()
	resetN()
	setPVBool(getBC() != 0)
	temp := v + reg.a
	set5Bool((temp & 0x02) != 0)
	set3Bool((temp & 0x08) != 0)
}

func LDIR() {
	LDI()
	if getBC() != 0 {
		reg.pc = reg.pc - 2
	}
}

func LDD() {
	v := (*memory).Get(getHL())
	(*memory).Put(getDE(), v)
	setHL(getHL() - 1)
	setDE(getDE() - 1)
	setBC(getBC() - 1)
	resetH()
	resetN()
	setPVBool(getBC() != 0)
	temp := v + reg.a
	set5Bool((temp & 0x02) != 0)
	set3Bool((temp & 0x08) != 0)
}

func LDDR() {
	LDD()
	if getBC() != 0 {
		reg.pc = reg.pc - 2
	}
}

func blockCP(y byte) {
	switch y {
	case 4:
		CPI()
	case 5:
		CPD()
	case 6:
		CPIR()
	default:
		CPDR()
	}
}

func CPI() {
	v := (*memory).Get(getHL())
	result := reg.a - v
	setHL(getHL() + 1)
	setBC(getBC() - 1)
	setSBool((result & 0x80) != 0)
	setZBool(result == 0)
	setHalfCarryFlagSubValue(reg.a, v)
	setPVBool(getBC() != 0)
	setN()
	if getH() {
		result--
	}
	set5Bool((result & 0x02) != 0)
	set3Bool((result & 0x08) != 0)
}

func CPIR() {
	CPI()
	if !getZ() && (getBC() != 0) {
		reg.pc = reg.pc - 2
	}
}

func CPD() {
	v := (*memory).Get(getHL())
	result := reg.a - v
	setHL(getHL() - 1)
	setBC(getBC() - 1)
	setSBool((result & 0x80) != 0)
	setZBool(result == 0)
	setHalfCarryFlagSubValue(reg.a, v)
	setPVBool(getBC() != 0)
	setN()
	if getH() {
		result--
	}
	set5Bool((result & 0x02) != 0)
	set3Bool((result & 0x08) != 0)
}

func CPDR() {
	CPD()
	if !getZ() && (getBC() != 0) {
		reg.pc = reg.pc - 2
	}
}

func blockIN(y byte) {
	switch y {
	case 4:
		INI()
	case 5:
		IND()
	case 6:
		INIR()
	default:
		INDR()
	}
}

func INI() {
	(*memory).Put(getHL(), (*io).Get(reg.c))
	reg.b--
	setHL(getHL() + 1)
	setZBool(reg.b == 0)
	setN()
}

func INIR() {
	INI()
	if !getZ() {
		reg.pc = reg.pc - 2
	}
}

func IND() {
	(*memory).Put(getHL(), (*io).Get(reg.c))
	reg.b--
	setHL(getHL() - 1)
	setZBool(reg.b == 0)
	setN()
}

func INDR() {
	IND()
	if !getZ() {
		reg.pc = reg.pc - 2
	}
}

func blockOUT(y byte) {
	switch y {
	case 4: // OUTI
		OUTI()
	case 5: // OUTD
		OUTD()
	case 6: // OTIR
		OTIR()
	default: // OTDR
		OTDR()
	}
}

func OUTI() {
	(*io).Put((*memory).Get(getHL()), reg.c)
	reg.b--
	setHL(getHL() + 1)
	setZBool(reg.b == 0)
	setN()
}

func OTIR() {
	OUTI()
	if !getZ() {
		reg.pc = reg.pc - 2
	}
}

func OUTD() {
	(*io).Put((*memory).Get(getHL()), reg.c)
	reg.b--
	setHL(getHL() - 1)
	setZBool(reg.b == 0)
	setN()
}

func OTDR() {
	OUTD()
	if !getZ() {
		reg.pc = reg.pc - 2
	}
}
