package internal

func decodeFD() {
	inst := (*memory).Get(reg.pc)
	reg.pc++
	_, _, z := basicDecode(inst)
	switch z {
	case 0:
		return
	case 1:
		return
	case 2:
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
