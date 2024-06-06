package main

import (
	"errors"
	"io"
	"log"
	"os"
)

func main() {
	// Открытие файла в режиме добавления или создание нового файла, если он не существует
	file, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Cannot open log file:", err)
	}
	defer file.Close()

	// Установка многопоточного вывода в файл и в консоль
	multiWriter := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multiWriter)

	// Логирование сообщений
	log.Println("Программа запущена")
	log.Printf("Значение переменной: %d", 42)

	// Логирование ошибки
	err = someFunction()
	if err != nil {
		log.Printf("Ошибка в someFunction: %v", err)
	}
}

// Пример функции, которая может возвращать ошибку
func someFunction() error {
	// Имитация ошибки
	return errors.New("произошла ошибка в someFunction")
}
