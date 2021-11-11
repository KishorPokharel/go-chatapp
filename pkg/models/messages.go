package models

import (
	"context"
	"database/sql"
	"time"
)

type Message struct {
	ID        int64
	CreatedAt time.Time
	UserID    int64
	Body      string
	Username  string
}

type MessageModel struct {
	DB *sql.DB
}

func (m *MessageModel) Insert(message *Message) error {
	query := `
			INSERT INTO messages (user_id, body)
			VALUES ($1, $2)
			RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(
		ctx,
		query,
		message.UserID,
		message.Body,
	).Scan(&message.ID, &message.CreatedAt)
	return err
}

func (m *MessageModel) GetAll() ([]*Message, error) {
	query := `
			SELECT 
				messages.id, messages.created_at, messages.user_id, messages.body,
				users.username
			FROM messages
			INNER JOIN users on users.id = messages.user_id
			ORDER BY messages.created_at ASC
			LIMIT 100`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []*Message{}
	for rows.Next() {
		var message Message
		err := rows.Scan(
			&message.ID,
			&message.CreatedAt,
			&message.UserID,
			&message.Body,
			&message.Username,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}
