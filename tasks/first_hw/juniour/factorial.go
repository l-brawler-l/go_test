package juniour


func Factorial(x int) uint64 {
	ans := uint64(1)
	for i := uint64(2); i <= uint64(x); i++ {
		ans *= i
	}
	return ans
}