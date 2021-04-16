package db

import (
	"database/sql"
	"errors"
	"log"

	service "github.com/Mahamed-Belkheir/scalechat-backend/user_service"
)

type UserRepository struct {
	conn *sql.DB

	fetchAll        *sql.Stmt
	fetchByUsername *sql.Stmt
	addUser         *sql.Stmt
}

func NewUserRepo(conn *sql.DB) UserRepository {
	var user UserRepository
	statement := prep(conn)
	user.fetchAll = statement.prepare("SELECT id, username, password FROM users;")
	user.fetchByUsername = statement.prepare("SELECT id, username, password from users WHERE (username=$1);")
	user.addUser = statement.prepare("INSERT INTO users(username, password) VALUES ($1, $2)")
	if statement.err != nil {
		log.Fatalf("error preparing statements: %v", statement.err)
	}
	return user
}

func (u UserRepository) GetUsers() ([]service.User, error) {
	result := make([]service.User, 0)
	rows, err := u.fetchAll.Query()
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user service.User
		err = rows.Scan(&user.ID, &user.Username, &user.Password)
		if err != nil {
			return nil, err
		}
		result = append(result, user)
	}
	return result, nil
}

func (u UserRepository) GetUserByUsername(username string) (*service.User, error) {
	var user service.User
	row := u.fetchByUsername.QueryRow(username)
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (u UserRepository) AddUser(user service.User) error {
	_, err := u.addUser.Exec(user.Username, user.Password)
	return err
}
