package models

import (
	"context"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               int       `json:"ID"`
	UserName         string    `json:"userName" validate:"required"`
	Email            string    `json:"email" validate:"required,email"`
	UnhashedPassword string    `json:"password" validate:"required"`
	Created          time.Time `json:"created"`
	PicturePath      string
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (um *UserModel) Insert(user *User) error {
	defaultPFP := os.Getenv("DEFAULT_PROFILE_PICTURE")
	if defaultPFP == "" {
		defaultPFP = "pics/default_pfp.png"
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.UnhashedPassword), 12)
	if err != nil {
		return err
	}

	insertStatement := `
  INSERT INTO users (name, email, hashed_password, profile_picture_path)
  VALUES ($1, $2, $3, $4)
  RETURNING id
  `
	args := []any{user.UserName, user.Email, string(hashedPassword), defaultPFP}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = um.DB.QueryRow(ctx, insertStatement, args...).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (um *UserModel) UpdatePicture(user *User) error {
	statement := `
  UPDATE users
  SET profile_picture_path = ($1)
  WHERE id = ($2)
  RETURNING id
  `

	args := []any{&user.PicturePath, &user.ID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRow(ctx, statement, args...).Scan(&user.ID)
	if err != nil {
		return err
	}

	return nil
}
