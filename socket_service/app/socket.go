package app

import (
	"log"

	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service/broker"
)

type SocketApplication struct {
	br           *broker.MessageBroker
	messagesRepo socket_service.IMessageRepository
}

func NewSocketApplication(br *broker.MessageBroker, msg socket_service.IMessageRepository) SocketApplication {
	return SocketApplication{
		br, msg,
	}
}

func (s SocketApplication) Register(userId, roomName string, ch chan socket_service.Message) {
	s.br.Register(userId, roomName, ch)
}

func (s SocketApplication) Unregister(userId, roomName string, ch chan socket_service.Message) {
	s.br.Unregister(userId, roomName, ch)
}

func (s SocketApplication) SendMessage(msg socket_service.Message) {
	err := s.messagesRepo.Add(&msg)
	if err != nil {
		log.Printf("error: error adding message to db: %v", err)
		return
	}
	s.br.SendMessage(msg)
}
