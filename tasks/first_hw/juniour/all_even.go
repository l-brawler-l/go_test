package juniour

import "fmt"

func AllEven(x int) {
	primes := make([]bool, x+1)
	for i := 2; i <= x; i++ {
		if !primes[i] {
			fmt.Println(i)
			for j := i*i; j <= x; j += i {
				primes[j] = true
			}
		}
	}
}