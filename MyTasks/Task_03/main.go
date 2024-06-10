package main

import "fmt"

func fibonacciIterative(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
		fmt.Printf("a = %d, b = %d\n", a, b)
	}
	return b
}

func main() {
	// Примеры использования функции fibonacci
	fmt.Println(fibonacciIterative(10)) // 55
}
