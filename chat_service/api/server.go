package api

import (
	"log"
	"net/http"

	service "github.com/Mahamed-Belkheir/scalechat-backend/chat_service"
	app "github.com/Mahamed-Belkheir/scalechat-backend/chat_service/app"
)

func StartWebServer(config service.Config, chatController app.ChatControl) {
	chatHandler := ChatHandler{NewJWT(config), chatController}
	log.Printf("server listening at %v", config.Port)
	log.Fatal(http.ListenAndServe(config.Port, chatHandler))
}
