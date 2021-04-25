package broker

import (
	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
)

type MessageBroker struct {
	rms    *rooms
	pubsub service.IPubSub
}

func (m *MessageBroker) Register(userId, roomName string, ch chan service.Message) {
	sub := m.rms.register(userId, roomName, ch)
	if sub {
		m.pubsub.StartListening(roomName, func(msg service.Message) {
			m.rms.broadcast(msg)
		})
	}
}

func (m *MessageBroker) Unregister(userId, roomName string) {
	unsub := m.rms.unregister(userId, roomName)
	if unsub {
		m.pubsub.StopListening(roomName)
	}
}

func (m *MessageBroker) SendMessage(msg service.Message) {
	m.pubsub.PublishMessage(msg)
}
