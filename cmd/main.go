package main

import (
	"my_app/internal/config"
	"my_app/internal/logger"
	"os"
)

func main() {
	lgr := logger.InitLogger()

	// проверка на наличие файлов конфига
	// если более двух то используется последний
	envfile := ".env"
	if len(os.Args) >= 2 {
		envfile = os.Args[1]
	}

	lgr.Info.Printf("Инициализация конфигурации.\nЧтение файла %s", envfile)
	if err := config.Init(envfile); err != nil {
		lgr.Err.Fatal(err.Error())
	}
}
