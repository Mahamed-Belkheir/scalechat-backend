package socket_service

type Message struct {
	ID        string
	UserID    string
	Room      string
	Body      string
	CreatedAt int64
}

type IMessageRepository interface {
	Add(*Message) error
	GetRoomMessages(room string) ([]Message, error)
	Delete(id, userId string) error
}
