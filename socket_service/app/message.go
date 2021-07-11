package app

import "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"

type MessagesApplication struct {
	messagesRepo socket_service.IMessageRepository
}

func NewMessagesApplication(messagesRepo socket_service.IMessageRepository) MessagesApplication {
	return MessagesApplication{
		messagesRepo,
	}
}

func (m MessagesApplication) GetRoomMessages(room string) ([]socket_service.Message, error) {
	return m.messagesRepo.GetRoomMessages(room)
}

func (m MessagesApplication) DeleteUserMessage(id, userId string) error {
	return m.messagesRepo.Delete(id, userId)
}
