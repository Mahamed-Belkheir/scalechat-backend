package events

import "fmt"

type MockChatEvents struct {
}

func (m MockChatEvents) SendDelChat(name string) error {
	fmt.Println("chat deleted: ", name)
	return nil
}

func NewMockChatEvents() MockChatEvents {
	return MockChatEvents{}
}
