package main

import "fmt"

//题目：给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数

func TwoSum(nums []int, t int) []int {
	other := 0
	m := make(map[int]int)
	for i, _ := range nums {
		other = t - nums[i]

		if _, ok := m[other]; ok {
			return []int{m[other], i}
		}
		m[nums[i]] = i
	}
	return nil
}

func main() {

	arr := []int{1, 2, 3, 4, 5}
	target := 8
	k := TwoSum(arr, target)
	fmt.Println(k)
}
