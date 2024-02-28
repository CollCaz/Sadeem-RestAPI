package models

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var defaultPFP = os.Getenv("DEFAULT_PROFILE_PICTURE")

type User struct {
	ID               int       `json:"ID"`
	UserName         string    `json:"userName" validate:"required"`
	Email            string    `json:"email" validate:"required,email"`
	UnhashedPassword string    `json:"password" validate:"required"`
	Created          time.Time `json:"created"`
	PicturePath      string
	IsAdmin          bool
}

// Custom marshaling function so we only show information we want to show
func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID       int    `json:"ID"`
		UserName string `json:"userName"`
		Email    string `json:"email"`
		IsAdmin  bool   `json:"isAdmin"`
	}{
		ID:       u.ID,
		UserName: u.UserName,
		Email:    u.Email,
		IsAdmin:  u.IsAdmin,
	})
}

type UserModel struct {
	DB *pgxpool.Pool
}

func (um *UserModel) Insert(user *User) error {
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

func (um *UserModel) DeleteUserByName(userName string) error {
	statement := `
  DELETE FROM users WHERE name = ($1)
  `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := um.DB.Exec(ctx, statement, userName)
	if err != nil {
		return err
	}

	return nil
}

func (um *UserModel) GetUserByName(userName string) (*User, error) {
	user := new(User)

	statement := `
  SELECT id, name, email, created, profile_picture_path FROM users
  WHERE name = ($1)
  `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRow(ctx, statement, userName).Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.Created,
		&user.PicturePath,
	)
	if err != nil {
		return nil, err
	}
	um.SetUserRole(user)

	return user, nil
}

// Sets the role of the user (admin or not)
func (um *UserModel) SetUserRole(user *User) {
	selectStatement := `
  SELECT name FROM users
  JOIN admin_users 
  ON admin_users.user_id = users.id
  WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRow(ctx, selectStatement, user.Email).Scan(&user.UserName)
	if err != nil {
		user.IsAdmin = false
		return
	}

	user.IsAdmin = true
}

func (um *UserModel) ResetPicture(userName string) error {
	statement := `
  UPDATE users
  SET profile_picture_path = $1
  WHERE name = $2
  `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := um.DB.Exec(ctx, statement, userName, defaultPFP)
	if err != nil {
		return err
	}

	return nil
}

func (um *UserModel) ValidateLogin(user *User) error {
	var hashedPassword []byte
	selectStatement := `
  SELECT hashed_password FROM users
  WHERE email = ($1)
  `
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRow(ctx, selectStatement, &user.Email).Scan(&hashedPassword)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(user.UnhashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func (um *UserModel) UpdatePicture(user *User) error {
	statement := `
  UPDATE users
  SET profile_picture_path = ($1)
  WHERE name = ($2)
  RETURNING name
  `
	args := []any{&user.PicturePath, &user.UserName}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRow(ctx, statement, args...).Scan(&user.UserName)
	if err != nil {
		return err
	}
	if user.UserName == "" {
		return errors.New("User does not exist")
	}

	return nil
}
