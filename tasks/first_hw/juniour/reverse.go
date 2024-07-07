package juniour


func Reverse(s string) string {
	ans := ""
	for _, c := range s {
		ans = string(c) + ans
	}
	return ans
}