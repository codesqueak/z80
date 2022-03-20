package internal

func decodeX0IXIY(y, z byte) {
	switch z {
	case 0:
		relativeJumps0(y)
	case 1:
		loadAdd16Immediate1IXIY(y)
	case 2:
		indirectLoad2IXIY(y)
	case 3:
		incDec163IXIY(y)
	case 4:
		inc4IXIY(y)
	case 5:
		dec5IXIY(y)
	case 6:
		ld6IXIY(y)
	default:
		accFlagOps7(y)
	}
}

func loadAdd16Immediate1IXIY(y byte) {
	p, q := getPQ(y)
	if q == 0 { // LD rp[p], nn
		setRPIXIY(p, load16FromRAM(reg.pc))
		reg.pc = reg.pc + 2
	} else { // ADD IXIY, rp[p]
		ixiy := getIXIY()
		right := getRPIXIY(p)
		result := ixiy + right
		//
		resetN()
		temp := (ixiy & 0x0FFF) + (right & 0x0FFF) // upper 8 half carry
		if (temp & 0xF000) != 0 {
			setH()
		} else {
			resetH()
		}
		if (result & 0x0800) != 0 {
			set3()
		} else {
			reset3()
		}
		if (result & 0x2000) != 0 {
			set5()
		} else {
			reset5()
		}
		if result < ixiy { // overflow ?
			setC()
		} else {
			resetC()
		}
		setIXIY(result)
	}
}

func indirectLoad2IXIY(y byte) {
	p, q := getPQ(y)
	if q == 0 {
		switch p {
		case 0:
			(*memory).Put(getBC(), reg.a)
		case 1:
			(*memory).Put(getDE(), reg.a)
		case 2:
			store16ToRAM(load16FromPC(), getIXIY()) // LD (nn), IXIY
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
			setIXIY(load16FromRAM(load16FromPC())) // LD IXIY, (nn)
		default:
			reg.a = (*memory).Get(load16FromPC()) // LD A, (nn)
		}
	}
}

func incDec163IXIY(y byte) {
	p, q := getPQ(y)
	v := getRPIXIY(p)
	if q == 0 {
		// inc 16
		setRPIXIY(p, v+1)
	} else {
		// dec 16
		setRPIXIY(p, v-1)
	}
}

// INC r[y]
func inc4IXIY(y byte) {
	v, offset := load8rIXIY(y)
	setHalfCarryFlagAddValue(v, 1)
	setPVBool(v == 0x7F)
	v++
	store8rIXIY(v, y)
	setSFromV(v)
	setZFromV(v)
	resetN()
	setUnusedFlagsFromV(v)
	if offset {
		reg.pc++
	}
}

// DEC r[y]
func dec5IXIY(y byte) {
	v, offset := load8rIXIY(y)
	setHalfCarryFlagSubValue(v, 1)
	setPVBool(v == 0x80)
	v = v - 1
	store8rIXIY(v, y)
	setSFromV(v)
	setZFromV(v)
	setN()
	setUnusedFlagsFromV(v)
	if offset {
		reg.pc++
	}
}

// LD r[y], n
func ld6IXIY(y byte) {
	if y == 6 { // ld (ix+dd), nn
		v := (*memory).Get(reg.pc + 1)
		store8rIXIY(v, y)
		reg.pc = reg.pc + 2
	} else { // ld r, nn
		v := (*memory).Get(reg.pc)
		store8rIXIY(v, y)
		reg.pc++
	}
}
