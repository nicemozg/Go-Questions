package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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

	// Установка обработчика для корневого пути
	http.HandleFunc("/", handler)

	// Запуск сервера на порту 8080
	fmt.Println("Server started at port 8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Создаем контекст с таймаутом в 1 минуту
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Minute)
	defer cancel()

	// Запускаем горутину для выполнения долгой операции
	go func() {
		// Выполняем долгую операцию
		time.Sleep(2 * time.Minute)

		// После выполнения операции проверяем, не закончился ли контекст
		select {
		case <-ctx.Done():
			// Если контекст истек, завершаем выполнение
			return
		default:
			// Если операция завершилась успешно, возвращаем результат
			fmt.Fprintln(w, "Long operation completed successfully")
		}
	}()

	// Возвращаем ответ, чтобы показать, что запрос был принят и обрабатывается
	fmt.Fprintln(w, "Request accepted and being processed...")
}
