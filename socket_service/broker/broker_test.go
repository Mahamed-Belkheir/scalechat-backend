package broker

import (
	"log"
	"testing"
	"time"

	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
)

type MockPubSub struct {
	listeners map[string]func(service.Message)
}

func (m *MockPubSub) PublishMessage(msg service.Message) {
	callback, ok := m.listeners[msg.Room]
	if !ok {
		log.Fatalf("panic: tried to publish to unregistered channel: %v", msg.Room)
	}
	callback(msg)
}

func (m *MockPubSub) StartListening(roomName string, callback func(service.Message)) {
	m.listeners[roomName] = callback
}

func (m *MockPubSub) StopListening(roomName string) {
	_, ok := m.listeners[roomName]
	if !ok {
		log.Fatalf("panic: tried to unregister an unregistered channel: %v", roomName)
	}
	delete(m.listeners, roomName)
}

var _ service.IPubSub = &MockPubSub{}

func TestBroker(t *testing.T) {
	mockPubSub := &MockPubSub{make(map[string]func(service.Message))}
	br := NewMessageBroker(mockPubSub)

	ch1, ch2 := make(chan service.Message), make(chan service.Message)

	br.Register("user1", "room1", ch1)
	br.Register("user2", "room1", ch2)

	br.rms.mut.Lock()
	assert(br.rms.rms, map[string]*room{"room1": {users: map[chan service.Message]string{ch1: "user1", ch2: "user2"}}}, t)
	br.rms.mut.Unlock()

	msg1 := service.Message{
		UserID:    "user1",
		Room:      "room1",
		Body:      "Hello World!",
		CreatedAt: time.Now().Unix(),
	}

	go br.SendMessage(msg1)

	res1, res2 := <-ch1, <-ch2

	assert(res1, msg1, t)
	assert(res2, msg1, t)
	br.Unregister("user1", "room1", ch1)
	br.Unregister("user2", "room1", ch2)

	br.rms.mut.Lock()
	assert(br.rms.rms, map[string]*room{}, t)
	br.rms.mut.Unlock()
}
