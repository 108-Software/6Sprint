package main

import (
	"log"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {

	logger := log.New(os.Stdout, "APP: ", log.LstdFlags|log.Lshortfile)

	srv := server.New(logger)

	logger.Println("Запуск сервера на порту 8080")
	if err := srv.Start(); err != nil {
		logger.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
