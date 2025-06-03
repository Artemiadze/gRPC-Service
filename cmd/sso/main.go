package main

import (
	"fmt"

	"github.com/Artemiadze/gRPC-Service/internal/config"
)

func main() {
	// инициализация обьекта config

	cfg := config.MustLoad()

	fmt.Println(cfg)
	// инициализация logger

	// инициализация приложения (app)

	// запуск gRPC-сервер приложения
}
