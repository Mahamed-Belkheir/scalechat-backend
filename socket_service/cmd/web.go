package main

import (
	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/api"
	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/broker"
	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/pool"
	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/pubsub"
)

func main() {
	config := service.GetConfig()
	pb := pubsub.NewNatsPubSub(config)
	br := broker.NewMessageBroker(pb)
	pool := pool.NewPool(10, 1000)
	go pool.Start()
	api.StartWebServer(config, br, pool)
}
