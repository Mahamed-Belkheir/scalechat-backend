package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	service "github.com/Mahamed-Belkheir/scalechat-backend/chat_service"
)

type ChatRepository struct {
	conn *sql.DB

	fetchAll *sql.Stmt
	fetchOne *sql.Stmt
	addChat  *sql.Stmt
	delChat  *sql.Stmt
}

var _ service.IChatRepository = ChatRepository{}

func NewChatRepository(conn *sql.DB) ChatRepository {
	repo := ChatRepository{
		conn: conn,
	}
	statement := prep(conn)
	repo.fetchAll = statement.prepare("SELECT name, user_id FROM chats")
	repo.fetchOne = statement.prepare("SELECT name, user_id from chats WHERE name = $1")
	repo.addChat = statement.prepare("INSERT INTO chats(name, user_id) VALUES ($1, $2)")
	repo.delChat = statement.prepare("DELETE FROM chats WHERE name = $1")
	if statement.err != nil {
		log.Fatalf("error preparing statements: %v", statement.err)
	}

	return repo
}

func (c ChatRepository) AddChat(name, userId string) error {
	_, err := c.addChat.Exec(name, userId)
	if err != nil {
		return fmt.Errorf("DB Error adding chat: %w", err)
	}
	return nil
}

func (c ChatRepository) DelChat(name string) error {
	_, err := c.delChat.Exec(name)
	if err != nil {
		return fmt.Errorf("DB Error deleting chat: %w", err)
	}
	return nil
}

func (c ChatRepository) GetChats() ([]service.Chat, error) {
	results := make([]service.Chat, 0)

	rows, err := c.fetchAll.Query()
	if err != nil {
		return nil, fmt.Errorf("DB error fetching chats: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var chat service.Chat
		err = rows.Scan(&chat.Name, &chat.UserID)
		if err != nil {
			return nil, fmt.Errorf("DB error scanning chats: %w", err)
		}
		results = append(results, chat)
	}
	return results, nil
}

func (c ChatRepository) GetChatByName(name string) (*service.Chat, error) {
	var chat service.Chat
	row := c.fetchOne.QueryRow(name)
	err := row.Scan(&chat.Name, &chat.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("DB error fetching chat: %w", err)
	}
	return &chat, nil
}
