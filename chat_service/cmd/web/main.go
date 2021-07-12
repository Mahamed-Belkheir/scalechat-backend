package main

import (
	service "github.com/Mahamed-Belkheir/scalechat-backend/chat_service"
	api "github.com/Mahamed-Belkheir/scalechat-backend/chat_service/api"
	app "github.com/Mahamed-Belkheir/scalechat-backend/chat_service/app"
	db "github.com/Mahamed-Belkheir/scalechat-backend/chat_service/db"
	events "github.com/Mahamed-Belkheir/scalechat-backend/chat_service/events"
)

func main() {
	config := service.GetConfig()
	conn := db.GetConnection(config)
	chatRepo := db.NewChatRepository(conn)
	chatEvents := events.NewMockChatEvents()
	chatApp := app.NewChatControl(chatRepo, chatEvents)
	api.StartWebServer(config, chatApp)
}
