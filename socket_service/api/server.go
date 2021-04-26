package api

import (
	"net/http"

	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/broker"
	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/pool"
)

func StartWebServer(config service.Config, br *broker.MessageBroker, pool *pool.Pool) error {
	jwt := NewJWT(config)
	http.Handle("/ws", newWSHandler(br, pool, config, jwt))
	return http.ListenAndServe(config.Port, nil)
}
