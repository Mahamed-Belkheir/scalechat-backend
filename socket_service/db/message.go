package db

import (
	"fmt"

	"github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
	"github.com/gocql/gocql"
)

type MessageRepository struct {
	conn *gocql.Session
}

var _ socket_service.IMessageRepository = MessageRepository{}

func NewMessageRepository(conn *gocql.Session) MessageRepository {
	return MessageRepository{conn}
}

func (m MessageRepository) Add(msg *socket_service.Message) error {
	id, err := gocql.RandomUUID()
	if err != nil {
		return fmt.Errorf("error generating UUID: %w", err)
	}
	msg.ID = id.String()
	err = m.conn.Query(
		"INSERT INTO messages (ID, Room, Body, UserID, CreatedAt) VALUES ( ?, ?, ?, ?, ?)",
		id, msg.Room, msg.Body, msg.UserID, msg.CreatedAt,
	).Exec()
	if err != nil {
		return fmt.Errorf("error inserting message: %w", err)
	}
	return nil
}

func (m MessageRepository) Delete(id, userId string) error {
	err := m.conn.Query(
		"DELETE FROM messages WHERE id = ? AND userId = ?",
		id, userId,
	).Exec()
	if err != nil {
		return fmt.Errorf("error deleting message: %w", err)
	}
	return nil
}

func (m MessageRepository) GetRoomMessages(room string) ([]socket_service.Message, error) {
	scanner := m.conn.Query(
		"SELECT ID, Room, Body, UserID, CreatedAt FROM messages WHERE Room=? ALLOW FILTERING",
		room,
	).Iter().Scanner()
	result := []socket_service.Message{}
	for scanner.Next() {
		msg := socket_service.Message{}
		err := scanner.Scan(&msg.ID, &msg.Room, &msg.Body, &msg.UserID, &msg.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning message: %w", err)
		}
		result = append(result, msg)
	}
	return result, nil
}
