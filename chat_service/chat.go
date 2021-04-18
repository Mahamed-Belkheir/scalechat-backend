package chat_service

type Chat struct {
	Name   string
	UserID string
}

type IChatRepository interface {
	GetChats() ([]Chat, error)
	AddChat(name, userId string) error
	DelChat(id string) error
	GetChatByName(id string) (*Chat, error)
}

type IChatEvents interface {
	SendDelChat(name string) error
}
