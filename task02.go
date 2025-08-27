package main

import (
	"fmt"
	"strconv"
)

//给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。
//回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。
//例如，121 是回文，而 123 不是。

func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	s := strconv.Itoa(x)

	n := len(s)
	fmt.Println(n)
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] { //如：s[2] != s[7-1-2]    对应值是1234321的：  3 和 3
			return false
		}
	}
	return true
}

func main() {
	fmt.Println(isPalindrome(1234321)) // true
	fmt.Println(isPalindrome(-121))    // false
	fmt.Println(isPalindrome(10))      // false
}
