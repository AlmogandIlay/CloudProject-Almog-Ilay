package main

import (
	"fmt"
)

func incrementer() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

func main() {
	increment := incrementer()
	fmt.Printf("%d\n", increment())
	fmt.Printf("%d", increment())
}
