package main

import (
	"fmt"
	"unicode"
)

// isPalindrome проверяет, является ли строка палиндромом.
func isPalindrome(s string) bool {
	left, right := 0, len(s)-1

	for left < right {
		// Игнорируем неалфавитные и нецифровые символы слева
		for left < right && !unicode.IsLetter(rune(s[left])) && !unicode.IsDigit(rune(s[left])) {
			left++
		}
		// Игнорируем неалфавитные и нецифровые символы справа
		for left < right && !unicode.IsLetter(rune(s[right])) && !unicode.IsDigit(rune(s[right])) {
			right--
		}
		// Сравниваем символы без учета регистра
		if left < right && unicode.ToLower(rune(s[left])) != unicode.ToLower(rune(s[right])) {
			return false
		}
		left++
		right--
	}
	return true
}

func main() {
	// Примеры использования функции isPalindrome
	tests := []struct {
		input  string
		output bool
	}{
		{"A man, a plan, a canal: Panama", true},
		{"race a car", false},
		{"No lemon, no melon", true},
		{"Was it a car or a cat I saw?", true},
	}

	for _, test := range tests {
		result := isPalindrome(test.input)
		fmt.Printf("Input: %q, Output: %v, Expected: %v\n", test.input, result, test.output)
	}
}
