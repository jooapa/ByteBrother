package main

import (
	f "fmt"
)

func main() {
	a := []int{1, 2, 3, 4, 5}

	for i := 0; i < len(a); i++ {
		f.Println(a[i])
	}
}
