package internal

import "github.com/codesqueak/z80/processor/pkg/hw"

// DD // FD Prefixed Index Instructions
func decodeDDFD(dd bool) {
	if dd {
		reg.ddMode = true
		reg.fdMode = false
	} else {
		reg.ddMode = false
		reg.fdMode = true
	}
	inst := (*memory).Get(reg.pc)
	reg.tStates = +opcodeDDFDStates[inst]
	reg.pc++
	x, y, z := basicDecode(inst)
	switch x {
	case 0:
		decodeX0IXIY(y, z) // various
	case 1:
		ldryrzIXIY(y, z) // LD r[y], r[z]
	case 2:
		decodeX2IXIY(y, z) // alu[y] r[z]
	default:
		decodeX3IXIY(y, z) // various
	}
}

func decodeX2IXIY(y, z byte) {
	v, offset := load8rIXIY(z)
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
	if offset {
		reg.pc++
	}
}

func decodeX3IXIY(y, z byte) {
	switch z {
	case 0: // RET cc[z]
		ret(y)
	case 1: // POP rp2[p] / RET / EXX
		popRetExxIXIY(y)
	case 2: // JP cc[z], nn
		jpcc(y)
	case 3: // various
		various3_3IXIY(y, io)
	case 4: // CALL cc[z], nn
		callcc(y)
	case 5: // various
		various3_5IXIY(y)
	case 6: //
		aluImmediate(y)
	default: // RST z*8
		rst(y)
	}
}

func popRetExxIXIY(y byte) {
	p, q := getPQ(y)
	if q == 0 { // POP rp2ixiy[p]
		v := pop()
		setRP2IXIY(p, v)
	} else {
		switch p {
		case 0: // RET
			reg.pc = pop()
		case 1: // EXX
			t := reg.b
			reg.b = reg.b_
			reg.b_ = t
			t = reg.c
			reg.c = reg.c_
			reg.c_ = t
			t = reg.d
			reg.d = reg.d_
			reg.d_ = t
			t = reg.e
			reg.e = reg.e_
			reg.e_ = t
			t = reg.h
			reg.h = reg.h_
			reg.h_ = t
			t = reg.l
			reg.l = reg.l_
			reg.l_ = t
		case 2: // JP (IX)
			reg.pc = getIXIY()
		default: // LD SP, HL
			reg.sp = getHL()
		}
	}
}

// various
func various3_3IXIY(y byte, io *hw.IO) {
	switch y {
	case 0: // JP nn
		reg.pc = load16FromPC()
	case 1: // (CB prefix)
		decodeCBIXIY()
	case 2: // OUT (n), A
		port := (*memory).Get(reg.pc)
		(*io).Put(port, reg.a)
		reg.pc++
	case 3: // IN A, (n)
		port := (*memory).Get(reg.pc)
		reg.a = (*io).Get(port)
		reg.pc++
	case 4: // EX (SP), IXIY
		t := load16FromRAM(reg.sp)
		store16ToRAM(reg.sp, getIXIY())
		setHL(t)
	case 5: // EX DE, HL - really is HL, not IX or IY - weird
		t := reg.d
		reg.d = reg.h
		reg.h = t
		t = reg.e
		reg.e = reg.l
		reg.l = t
	case 6: // DI
		reg.iff1 = false
	default: // EI
		reg.iff1 = true
	}
}

// various
func various3_5IXIY(y byte) {
	p, q := getPQ(y)
	if q == 0 { // PUSH rp2ixiy[p]
		push(getRP2IXIY(p))
	} else {
		switch p {
		case 0: // CALL nn
			addr := load16FromPC()
			push(reg.pc)
			reg.pc = addr
		case 1:
			decodeDDFD(true) // really !
		case 2:
			decodeED() // DD / FD gets ignored
		default:
			decodeDDFD(false) // really !
		}
	}
}
