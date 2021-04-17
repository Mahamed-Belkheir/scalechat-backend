package chat_service

import "os"

type Config struct {
	Port   string
	DBConn string
	DB     string
	Secret string
}

func GetConfig() Config {
	return Config{
		Port:   os.Getenv("PORT"),
		DBConn: os.Getenv("DB_CONN"),
		DB:     os.Getenv("DB"),
		Secret: os.Getenv("SECRET"),
	}
}
