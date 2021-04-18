package app

import (
	"fmt"
	"reflect"
	"testing"

	service "github.com/Mahamed-Belkheir/scalechat-backend/chat_service"
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

type MockChatRepo struct {
	chats map[string]service.Chat
}

func (m MockChatRepo) AddChat(name, userId string) error {
	if _, ok := m.chats[name]; ok {
		return fmt.Errorf("chat name already in use")
	}
	m.chats[name] = service.Chat{Name: name, UserID: userId}
	return nil
}

func (m MockChatRepo) DelChat(name string) error {
	if _, ok := m.chats[name]; !ok {
		return fmt.Errorf("chat not found")
	}
	delete(m.chats, name)
	return nil
}

func (m MockChatRepo) GetChatByName(name string) (*service.Chat, error) {
	if chat, ok := m.chats[name]; !ok {
		return nil, nil
	} else {
		return &chat, nil
	}
}

func (m MockChatRepo) GetChats() ([]service.Chat, error) {
	res := make([]service.Chat, len(m.chats))
	i := 0
	for _, chat := range m.chats {
		res[i] = chat
		i++
	}
	return res, nil
}

func newMockChatRepo() MockChatRepo {
	return MockChatRepo{
		chats: make(map[string]service.Chat, 0),
	}
}

type MockChatEvents struct {
	deletedChats *[]string
}

func (m MockChatEvents) SendDelChat(name string) error {
	fmt.Println("got here", name)
	*m.deletedChats = append(*(m.deletedChats), name)
	return nil
}

func newMockChatEvents() MockChatEvents {
	return MockChatEvents{
		&[]string{},
	}
}

func TestChatController(t *testing.T) {
	events := newMockChatEvents()
	chatCont := NewChatControl(newMockChatRepo(), events)

	chats, _ := chatCont.GetChats()
	assert(chats, []service.Chat{}, t)

	err := chatCont.AddChat("room1", "userman")
	ok(err, t)

	err = chatCont.AddChat("room1", "userman")
	assert(err, fmt.Errorf("chat name already in use"), t)

	chats, _ = chatCont.GetChats()
	assert(chats, []service.Chat{
		{
			Name:   "room1",
			UserID: "userman",
		},
	}, t)

	err = chatCont.DeleteChat("userman", "room2")
	assert(err, fmt.Errorf("chat not found"), t)

	err = chatCont.DeleteChat("fakeUser", "room1")
	assert(err, fmt.Errorf("user does not own chat"), t)

	err = chatCont.DeleteChat("userman", "room1")
	ok(err, t)
	assert(events.deletedChats, &[]string{"room1"}, t)

}
