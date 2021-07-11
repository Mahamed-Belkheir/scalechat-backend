package socket_service

import (
	"os"
	"strings"
)

type Config struct {
	Port       string
	Clusters   []string
	Secret     string
	PubSubConn string
}

func GetConfig() Config {

	return Config{
		Port:       os.Getenv("PORT"),
		Clusters:   strings.Split(os.Getenv("CLUSTERS"), ","),
		Secret:     os.Getenv("SECRET"),
		PubSubConn: os.Getenv("PUBSUB_CONN"),
	}
}
