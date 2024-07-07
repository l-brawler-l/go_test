package juniour

func Max(x, y, z int) int {
	switch {
	case x >= y && x >= z:
		return x
	case y >= x && y >= z:
		return y
	default:
		return z
	}
}