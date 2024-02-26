package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               int32     `json:"ID"`
	UserName         string    `json:"userName" validate:"required"`
	Email            string    `json:"email" validate:"required,email"`
	UnhashedPassword string    `json:"password" validate:"required"`
	Created          time.Time `json:"created"`
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (um *UserModel) Insert(user *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.UnhashedPassword), 12)
	if err != nil {
		return err
	}

	insert_statement := `
  INSERT INTO users (name, email, hashed_password)
  VALUES ($1, $2, $3)
  RETURNING id
  `
	args := []any{user.UserName, user.Email, string(hashedPassword)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := um.DB.Exec(ctx, insert_statement, args...); err != nil {
		return err
	}

	return nil
}
