package juniour


func Sum(v []int) int {
	ans := 0
	for _, x := range v {
		ans += x
	}
	return ans
}