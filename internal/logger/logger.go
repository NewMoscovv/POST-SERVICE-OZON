package logger

import (
	"log"
	"os"
)

type Logger struct {
	Err  *log.Logger
	Info *log.Logger
}

// InitLogger - метод инициализирующий логирование
func InitLogger() *Logger {
	logger := Logger{}
	logger.Info = log.New(os.Stdout, "[ИНФО]\t", log.Ldate|log.Ltime)
	logger.Err = log.New(os.Stderr, "[ОШИБКА]\t", log.Ldate|log.Ltime)
	logger.Info.Print("Инициализация логгирования.")
	return &logger
}
