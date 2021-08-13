package app

import (
	"fmt"

	service "github.com/Mahamed-Belkheir/scalechat-backend/user_service"
)

type Authentication struct {
	repository service.IUserRepository
}

func NewAuthentication(repo service.IUserRepository) Authentication {
	return Authentication{repo}
}

func (a Authentication) Login(username, password string) (*service.User, error) {
	user, err := a.repository.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("error fetching user %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	err = user.ComparePassword(password)
	if err != nil {
		return nil, fmt.Errorf("Incorrect password %w", err)
	}
	return user, nil
}

func (a Authentication) Register(username, password string) (*service.User, error) {
	user, err := a.repository.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database %w", err)
	}
	if user != nil {
		return nil, fmt.Errorf("username already used")
	}
	newUser := service.User{
		Username: username,
		Password: password,
	}
	err = newUser.Hash()
	if err != nil {
		return nil, fmt.Errorf("error hashing user password %w", err)
	}
	err = a.repository.AddUser(newUser)
	if err != nil {
		return nil, fmt.Errorf("error saving user to database %w", err)
	}
	user, err = a.repository.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database %w", err)
	}
	return user, nil
}
