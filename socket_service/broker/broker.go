package broker

import (
	"log"

	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
)

type MessageBroker struct {
	rms    *rooms
	pubsub service.IPubSub
}

func NewMessageBroker(pubsub service.IPubSub) *MessageBroker {
	return &MessageBroker{
		newRooms(),
		pubsub,
	}
}

func (m *MessageBroker) Register(userId, roomName string, ch chan service.Message) {
	log.Printf("debug: registering new user in broker, user: %s, in room: %s", userId, roomName)
	sub := m.rms.register(userId, roomName, ch)
	if sub {
		log.Printf("first room usage, subbing to %s", roomName)
		m.pubsub.StartListening(roomName, func(msg service.Message) {
			m.rms.broadcast(msg)
		})
	}
}

func (m *MessageBroker) Unregister(userId, roomName string, ch chan service.Message) {
	log.Printf("debug: unregistering user in broker, user: %s, in room %s", userId, roomName)
	unsub := m.rms.unregister(ch, roomName)
	if unsub {
		log.Printf("debug: last room usage, unsubbing from %s", roomName)
		m.pubsub.StopListening(roomName)
	}
}

func (m *MessageBroker) SendMessage(msg service.Message) {
	m.pubsub.PublishMessage(msg)
}
