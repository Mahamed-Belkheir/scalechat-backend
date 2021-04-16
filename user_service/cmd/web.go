package main

import (
	service "github.com/Mahamed-Belkheir/scalechat-backend/user_service"
	api "github.com/Mahamed-Belkheir/scalechat-backend/user_service/api"
	app "github.com/Mahamed-Belkheir/scalechat-backend/user_service/app"
	db "github.com/Mahamed-Belkheir/scalechat-backend/user_service/db"
)

func main() {
	config := service.GetConfig()
	dbConn := db.GetConnection(config)
	userRepo := db.NewUserRepo(dbConn)
	authApp := app.NewAuthentication(userRepo)
	api.StartWebServer(config, authApp)
}
