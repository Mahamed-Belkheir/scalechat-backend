package db

import (
	"database/sql"
	"fmt"
	"log"

	service "github.com/Mahamed-Belkheir/scalechat-backend/chat_service"
)

type ChatRepository struct {
	conn *sql.DB

	fetchAll *sql.Stmt
	addChat  *sql.Stmt
	delChat  *sql.Stmt
}

var _ service.IChatRepository = ChatRepository{}

func NewChatRepository(conn *sql.DB) ChatRepository {
	repo := ChatRepository{
		conn: conn,
	}
	statement := prep(conn)
	repo.fetchAll = statement.prepare("SELECT id, name, user_id FROM chats")
	repo.addChat = statement.prepare("INSERT INTO chats(name, user_id) VALUES ($1, $2)")
	repo.delChat = statement.prepare("DELETE FROM chats WHERE ID = $1")
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

func (c ChatRepository) DelChat(id string) error {
	_, err := c.delChat.Exec(id)
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
		err = rows.Scan(&chat.ID, &chat.Name, &chat.UserID)
		if err != nil {
			return nil, fmt.Errorf("DB scanning chats: %w", err)
		}
		results = append(results, chat)
	}
	return results, nil
}
