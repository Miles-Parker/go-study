package main

//题目：给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效

import (
	"fmt"
)

func isValid(s string) bool {
	stack := make([]rune, 0)
	mapping := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	for _, char := range s {
		if closing, exists := mapping[char]; exists {
			if len(stack) == 0 || stack[len(stack)-1] != closing {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, char)
		}
	}
	return len(stack) == 0
}

func main() {
	testCases := []string{"()[]{}", "([)]", "{[]}"}
	for _, s := range testCases {
		fmt.Printf("%s: %v\n", s, isValid(s))
	}
}
