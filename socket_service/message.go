package socket_service

type Message struct {
	UserID    string
	Room      string
	Body      string
	CreatedAt int
}
