package main

import "fmt"

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func longestCommonSubsequence(s, t string) int {
	m, n := len(s), len(t)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s[i-1] == t[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp[m][n]
}

func main() {
	s := "abcde"
	t := "ac"
	result := longestCommonSubsequence(s, t)
	fmt.Println("Length of Longest Common Subsequence:", result)
}
