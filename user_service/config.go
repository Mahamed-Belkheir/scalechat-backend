package user_service

import "os"

type Config struct {
	Port   string
	DBConn string
	DB     string
}

func GetConfig() Config {
	return Config{
		Port:   os.Getenv("PORT"),
		DBConn: os.Getenv("DB_CONN"),
		DB:     os.Getenv("DB"),
	}
}
