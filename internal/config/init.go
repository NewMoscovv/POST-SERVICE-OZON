package config

import "github.com/joho/godotenv"

func Init(file string) error {
	err := godotenv.Load(file)

	return err
}
