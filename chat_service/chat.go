package chat_service

type Chat struct {
	ID     string
	Name   string
	UserID string
}

type IChatRepository interface {
	GetChats() ([]Chat, error)
	AddChat(name, userId string) error
	DelChat(id string) error
}