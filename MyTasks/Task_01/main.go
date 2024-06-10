package main

import (
	"fmt"
)

// Функция twoSum находит индексы двух чисел, сумма которых равна target
func twoSum(nums []int, target int) []int {
	// Создаем карту для хранения чисел и их индексов
	numMap := make(map[int]int)

	// Проходим по массиву чисел
	for i, num := range nums {
		// Вычисляем разность между целевым числом и текущим числом
		complement := target - num
		// Проверяем, есть ли разность в карте
		if j, ok := numMap[complement]; ok {
			// Если есть, возвращаем индексы чисел
			return []int{j, i}
		}
		// Иначе сохраняем текущее число и его индекс в карту
		numMap[num] = i
	}
	// Если нет такой пары чисел, возвращаем пустой массив
	return []int{}
}

func main() {
	// Пример массива чисел
	nums := []int{2, 8, 11, 7}
	// Целевое число
	target := 9

	// Вызов функции twoSum и передача значений
	result := twoSum(nums, target)

	// Печать результата
	fmt.Println("Индексы чисел, сумма которых равна", target, ":", result)
}
