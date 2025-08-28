package main

import "fmt"

// 实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
func isSlice(s *[]int) {
	for i := range *s {
		(*s)[i] *= 2
	}
}

func main() {
	nums := []int{1, 2, 3, 4}
	fmt.Println("Before:", nums) // [1 2 3 4]
	isSlice(&nums)
	fmt.Println("After:", nums) // [2 4 6 8]
}
