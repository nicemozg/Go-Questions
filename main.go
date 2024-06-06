package main

import (
	"fmt"
	"time"
)

func main() {
	// Создаем два канала
	ch1 := make(chan string)
	ch2 := make(chan string)

	// Горутина отправляет "Hello" в канал ch1 через 2 секунды
	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "Hello"
	}()

	// Горутина отправляет "World" в канал ch2 через 1 секунду
	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- "World"
	}()

	// Используем оператор select для ожидания и получения значений из каналов
	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received from ch1:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received from ch2:", msg2)
		}
	}
}
