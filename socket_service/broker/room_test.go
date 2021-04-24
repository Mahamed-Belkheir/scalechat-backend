package broker

import (
	"reflect"
	"testing"
	"time"

	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
)

func assert(got, expected interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("\nexpected: %v \n got: %v", expected, got)
	}
}

func ok(err interface{}, t *testing.T) {
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}

func TestRoom(t *testing.T) {
	rm := newRoom()
	ch1, ch2 := make(chan socket_service.Message), make(chan socket_service.Message)
	assert(rm.isEmpty(), true, t)

	rm.register("1", ch1)
	rm.register("2", ch2)
	assert(rm.isEmpty(), false, t)

	msg := socket_service.Message{
		UserID:    "user",
		Room:      "room",
		Body:      "hello world",
		CreatedAt: time.Now().Unix(),
	}

	go rm.broadcast(msg)

	res1 := <-ch1
	res2 := <-ch2

	assert(res1, msg, t)
	assert(res2, msg, t)

	rm.unregister("1")
	rm.unregister("2")
	assert(rm.isEmpty(), true, t)
}

func TestRooms(t *testing.T) {
	rms := newRooms()
	ch1, ch2, ch3 := make(chan socket_service.Message), make(chan socket_service.Message), make(chan socket_service.Message)
	rms.register("1", "room1", ch1)
	rms.register("2", "room1", ch2)
	rms.register("3", "room2", ch3)

	msg1 := socket_service.Message{
		UserID:    "user",
		Room:      "room1",
		Body:      "hello world",
		CreatedAt: time.Now().Unix(),
	}
	msg2 := socket_service.Message{
		UserID:    "user",
		Room:      "room2",
		Body:      "hello world",
		CreatedAt: time.Now().Unix(),
	}

	go rms.broadcast(msg1)
	go rms.broadcast(msg2)

	res1, res2, res3 := <-ch1, <-ch2, <-ch3

	assert(res1, msg1, t)
	assert(res2, msg1, t)
	assert(res3, msg2, t)

	del := rms.unregister("1", "room1")
	assert(del, false, t)
	del = rms.unregister("2", "room1")
	assert(del, true, t)
	del = rms.unregister("3", "room2")
	assert(del, true, t)

}
