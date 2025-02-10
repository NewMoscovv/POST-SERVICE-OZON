package main

import (
	"my_app/internal/config"
	"my_app/internal/logger"
)

func main() {
	lgr := logger.InitLogger()

	lgr.Info.Print("Инициализация конфигурации.")

	if err := config.Init(); err != nil {
		lgr.Err.Fatal(err.Error())
	}
}
