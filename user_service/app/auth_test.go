package app

import (
	"errors"
	"reflect"
	"testing"

	service "github.com/Mahamed-Belkheir/scalechat-backend/user_service"
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

type mockUserRepository struct {
	users map[string]service.User
}

func newMockUserRepo() mockUserRepository {
	return mockUserRepository{
		make(map[string]service.User),
	}
}

func (m mockUserRepository) AddUser(user service.User) error {
	_, ok := m.users[user.Username]
	if ok {
		return errors.New("unique contraint violation: username")
	}
	m.users[user.Username] = user
	return nil
}

func (m mockUserRepository) GetUserByUsername(username string) (*service.User, error) {
	user, ok := m.users[username]
	if !ok {
		return nil, nil
	}
	return &user, nil
}

func (m mockUserRepository) GetUsers() ([]service.User, error) {
	users := make([]service.User, len(m.users))
	i := 0
	for _, user := range m.users {
		users[i] = user
		i++
	}
	return users, nil
}

func TestAuth(t *testing.T) {
	authApp := NewAuthentication(newMockUserRepo())

	_, err := authApp.Login("username", "password")
	assert(err, errors.New("user not found"), t)

	_, err = authApp.Register("username", "password")
	ok(err, t)

	_, err = authApp.Register("username", "password")
	assert(err, errors.New("username already used"), t)

	_, err = authApp.Login("username", "password")
	ok(err, t)

}
