package models

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"time"

	"github.com/gofrs/uuid"
	"server/helpers"
)

type User struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name,omitempty" validate:"required"`
	Email     string    `json:"email,omitempty" validate:"required,email"`
	Password  string    `json:"password,omitempty" validate:"required"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (u *User) Create(user User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	// Check if user already exists
	dupe, err := u.FindByEmail(user.Email)

	if err != nil {
		log.Error().Err(err).Msg("Error finding user")
		return nil, errors.New("Error finding user")
	}

	if dupe != nil {
		return nil, errors.New("User already exists")
	}

	// Create user
	hasedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		log.Error().Err(err).Msg("Error hashing password")
		return nil, err
	}

	user.Password = hasedPassword

	query := `INSERT INTO users (name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING *`

	_, err = db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		log.Error().Err(err).Msg("Error creating user")
		return nil, err
	}

	return &user, nil
}

func (u *User) FindAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `SELECT id, name, email, password, created_at, updated_at FROM users`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("Error finding users")
		return nil, err
	}

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Error().Err(err).Msg("Error scanning users")
			return nil, err
		}

		users = append(users, &user)
	}

	if len(users) == 0 {
		return nil, errors.New("No user found")
	}

	return users, nil
}

func (u *User) FindByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1`

	rows, err := db.QueryContext(ctx, query, email)
	if err != nil {
		log.Error().Err(err).Msg("Error finding user")
		return nil, err
	}

	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Error().Err(err).Msg("Error scanning user")
			return nil, err
		}

		users = append(users, &user)
	}

	if len(users) == 0 {
		return nil, errors.New("No user found")
	}

	u.ID = users[0].ID
	u.Name = users[0].Name
	u.Email = users[0].Email
	u.Password = users[0].Password
	u.CreatedAt = users[0].CreatedAt
	u.UpdatedAt = users[0].UpdatedAt

	return u, nil
}

func (u *User) UpdateByEmail(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	query := `UPDATE users SET name = $1, email = $2, password = $3, updated_at = $4 WHERE email = $5`

	_, err := db.ExecContext(ctx, query, user.Name, user.Email, user.Password, time.Now(), user.Email)
	if err != nil {
		log.Error().Err(err).Msg("Error updating user")
		return err
	}

	return nil
}
