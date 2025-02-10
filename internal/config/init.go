package config

import "github.com/joho/godotenv"

func Init() error {
	err := godotenv.Load()

	return err
}
