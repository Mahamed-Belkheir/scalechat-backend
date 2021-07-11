package api

import (
	"log"
	"net/http"

	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/app"
	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/pool"
)

func StartWebServer(config service.Config, sockApp app.SocketApplication, pool *pool.Pool) {
	jwt := NewJWT(config)
	http.Handle("/ws", newWSHandler(sockApp, pool, config, jwt))
	log.Printf("info: server listening at %v", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, nil))
}
