package main

import "fmt"

//题目：给定一个表示 大整数 的整数数组 digits，其中 digits[i] 是整数的第 i 位数字。这些数字按从左到右，从最高位到最低位排列。
//这个大整数不包含任何前导

func plusOne(digits []int) []int {
	n := len(digits)
	for i := n - 1; i >= 0; i-- {
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}

	return append([]int{1}, digits...)
}

func main() {
	testCases := [][]int{
		{1, 2, 3},
		{4, 3, 2, 1},
		{9, 9, 9},
		{0},
	}
	for _, tc := range testCases {
		fmt.Printf("Input: %v -> Output: %v\n", tc, plusOne(tc))
	}
}
