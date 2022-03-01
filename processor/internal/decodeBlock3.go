package internal

import (
	"z80/processor/pkg/hw"
)

func decodeX3(y, z byte) {
	switch z {
	case 0: // RET cc[z]
		ret(y)
	case 1: // POP rp2[p] / RET / EXX
		popRetExx(y)
	case 2: // JP cc[z], nn
		jpcc(y)
	case 3: // various
		various3_3(y, io)
	case 4: // CALL cc[z], nn
		callcc(y)
	case 5: // various
		various3_5(y)
	case 6: //
		aluImmediate(y)
	default: // RST z*8
		rst(y)
	}
}

// RET cc[y]
func ret(y byte) {
	if cc(y) {
		reg.pc = pop()
	}
}

func popRetExx(y byte) {
	p, q := getPQ(y)
	if q == 0 { // POP rp2[p]
		v := pop()
		setRP2(p, v)
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
		case 2: // JP (HL)
			reg.pc = load16FromRAM(getHL())
		default: // LD SP, HL
			reg.sp = getHL()
		}
	}
}

// JP cc[y], nn
func jpcc(y byte) {
	if cc(y) {
		reg.pc = load16FromPC()
	} else {
		reg.pc = reg.pc + 2
	}
}

// various
func various3_3(y byte, io *hw.IO) {
	switch y {
	case 0: // JP nn
		reg.pc = load16FromPC()
	case 1: // (CB prefix)
		decodeCB()
	case 2: // OUT (n), A
		port := (*memory).Get(reg.pc)
		(*io).Put(port, reg.a)
		reg.pc++
	case 3: // IN A, (n)
		port := (*memory).Get(reg.pc)
		reg.a = (*io).Get(port)
		reg.pc++
	case 4: // EX (SP), HL
		t := load16FromRAM(reg.sp)
		store16ToRAM(reg.sp, getHL())
		setHL(t)
	case 5: // EX DE, HL
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

// CALL cc[y], nn
func callcc(y byte) {
	if cc(y) {
		addr := load16FromPC()
		reg.pc = reg.pc + 2
		push(reg.pc)
		reg.pc = addr
	} else {
		reg.pc = reg.pc + 2
	}
}

// various
func various3_5(y byte) {
	p, q := getPQ(y)
	if q == 0 { // POP rp2[p]
		v := pop()
		setRP2(p, v)
	} else {
		switch p {
		case 0: // CALL nn
			addr := load16FromPC()
			reg.pc = reg.pc + 2
			push(reg.pc)
			reg.pc = addr
		case 1:
			decodeDD()
		case 2:
			decodeED()
		default:
			decodeFD()
		}
	}
}

func aluImmediate(y byte) {
	v := (*memory).Get(reg.pc)
	reg.pc++
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

func rst(y byte) {
	push(reg.pc)
	reg.pc = uint16(y) * 8
}
