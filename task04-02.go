package main

import "fmt"

//题目：查找字符串数组中的最长公共前缀

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	// 找出最短字符串长度
	minLen := len(strs[0])
	for _, s := range strs {
		if len(s) < minLen {
			minLen = len(s)
		}
	}

	// 纵向扫描
	for i := 0; i < minLen; i++ {
		char := strs[0][i]
		for _, s := range strs {
			if s[i] != char {
				return strs[0][:i]
			}
		}
	}
	return strs[0][:minLen]
}

func main() {
	testCases := [][]string{
		{"flower", "flow", "flight"},
		{"dog", "racecar", "car"},
		{"interspecies", "interstellar", "interstate"},
		{},
	}
	for _, tc := range testCases {
		fmt.Printf("Input: %v -> Output: '%s'\n", tc, longestCommonPrefix(tc))
	}
}
