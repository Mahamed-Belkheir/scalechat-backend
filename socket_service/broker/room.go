package broker

import (
	"log"
	"sync"

	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
)

type room struct {
	users map[string]chan service.Message
}

func newRoom() *room {
	return &room{make(map[string]chan service.Message)}
}

func (r *room) register(user string, ch chan service.Message) {
	r.users[user] = ch
}

func (r *room) unregister(user string) {
	delete(r.users, user)
}

func (r *room) broadcast(msg service.Message) {
	for _, ch := range r.users {
		ch <- msg
	}
}

func (r *room) isEmpty() bool {
	if len(r.users) == 0 {
		return true
	}
	return false
}

type rooms struct {
	mut *sync.RWMutex
	rms map[string]*room
}

func newRooms() *rooms {
	return &rooms{
		mut: &sync.RWMutex{},
		rms: make(map[string]*room),
	}
}

func (r *rooms) register(user, roomName string, ch chan service.Message) bool {
	r.mut.Lock()
	defer r.mut.Unlock()
	newSub := false
	rm, ok := r.rms[roomName]
	if !ok {
		rm = newRoom()
		r.rms[roomName] = rm
		newSub = true
	}
	rm.register(user, ch)
	return newSub
}

func (r *rooms) unregister(user string, roomName string) bool {
	r.mut.Lock()
	defer r.mut.Unlock()
	lastSub := false
	rm, ok := r.rms[roomName]
	if !ok {
		log.Printf("error: attempted to unregister nonexisting room %v", roomName)
		return true
	}
	rm.unregister(user)
	if rm.isEmpty() {
		lastSub = true
		delete(r.rms, roomName)
	}
	return lastSub
}

func (r *rooms) broadcast(msg service.Message) {
	r.mut.RLock()
	defer r.mut.RUnlock()
	rm, ok := r.rms[msg.Room]
	if !ok {
		log.Printf("error: attempted to broadcast to nonexisting room %v", msg.Room)
		return
	}
	rm.broadcast(msg)
}
