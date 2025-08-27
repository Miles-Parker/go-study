package main

import "fmt"

func main() {
	var nums = [5]int{22, 33, 22, 33, 44}
	m := make(map[int]int)

	for _, num := range nums {
		m[num]++ //统计数据 数组的元素写入map key进行自增，map value记录重复数量。
	}
	for k, v := range m {
		if v == 1 {
			fmt.Printf("只出现一次的元素是：%d\n", k)
		}
	}
}
