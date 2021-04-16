package user_service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string
	Username string
	Password string
}

func (u *User) Hash() error {
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	u.Password = string(password)
	return nil
}

func (u *User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("failed password comparsion %w", err)
	}
	return nil
}
