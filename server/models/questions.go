package models

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid"
)

type Question struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Question  string    `json:"question,omitempty"`
	Answer    string    `json:"answer,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	UserID    uuid.UUID `json:"user_id,omitempty"`
}

func (q *Question) Create(question Question) (*Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `INSERT INTO questions (question, answer, created_at, updated_at, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING *`

	_, err := db.ExecContext(
		ctx,
		query,
		question.Question,
		question.Answer,
		time.Now(),
		time.Now(),
		question.UserID,
	)

	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (q *Question) FindAll() ([]*Question, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `SELECT id, question, answer, created_at, updated_at, user_id FROM questions`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var questions []*Question
	for rows.Next() {
		var question Question
		err := rows.Scan(&question.ID, &question.Question, &question.Answer, &question.CreatedAt, &question.UpdatedAt, &question.UserID)
		if err != nil {
			return nil, err
		}

		questions = append(questions, &question)
	}

	if len(questions) == 0 {
		return nil, errors.New("No question found")
	}

	return questions, nil
}
