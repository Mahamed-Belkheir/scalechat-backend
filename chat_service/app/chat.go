package app

import (
	"fmt"

	service "github.com/Mahamed-Belkheir/scalechat-backend/chat_service"
)

type ChatControl struct {
	repo   service.IChatRepository
	events service.IChatEvents
}

func NewChatControl(repo service.IChatRepository, events service.IChatEvents) ChatControl {
	return ChatControl{
		repo:   repo,
		events: events,
	}
}

func (c ChatControl) AddChat(name, userId string) error {
	chat, err := c.repo.GetChatByName(name)
	if err != nil {
		return fmt.Errorf("Error fetching chat: %w", err)
	}
	if chat != nil {
		return fmt.Errorf("Chat name already in use")
	}
	err = c.repo.AddChat(name, userId)
	if err != nil {
		return fmt.Errorf("Error adding chat: %w", err)
	}
	return nil
}

func (c ChatControl) GetChats(id string) ([]service.Chat, error) {
	chats, err := c.repo.GetChats()
	if err != nil {
		return nil, fmt.Errorf("Error fetching chats: %w", err)
	}
	return chats, nil
}

func (c ChatControl) DeleteChat(userId, name string) error {
	chat, err := c.repo.GetChatByName(name)
	if err != nil {
		return fmt.Errorf("Error fetching chat: %w", err)
	}
	if chat == nil {
		return fmt.Errorf("chat not found")
	}
	if chat.UserID != userId {
		return fmt.Errorf("user does not own chat")
	}
	err = c.repo.DelChat(name)
	if err != nil {
		return fmt.Errorf("Error deleting chat: %w", err)
	}
	err = c.events.SendDelChat(name)
	if err != nil {
		return fmt.Errorf("Error sending delete event: %w", err)
	}
	return nil
}
