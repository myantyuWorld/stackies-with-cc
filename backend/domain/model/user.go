package model

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

// User はユーザーエンティティを表す
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Picture   string    `json:"picture"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser は新しいユーザーを作成する
func NewUser(id, email, name, picture string) (*User, error) {
	user := &User{
		ID:        id,
		Email:     email,
		Name:      name,
		Picture:   picture,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if !user.IsValid() {
		return nil, errors.New("invalid user data")
	}

	return user, nil
}

// IsValid はユーザーデータが有効かどうかを検証する
func (u *User) IsValid() bool {
	if strings.TrimSpace(u.ID) == "" {
		return false
	}
	if strings.TrimSpace(u.Email) == "" {
		return false
	}
	if strings.TrimSpace(u.Name) == "" {
		return false
	}
	if !u.isValidEmail(u.Email) {
		return false
	}
	return true
}

// UpdateProfile はユーザープロフィールを更新する
func (u *User) UpdateProfile(name, picture string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name cannot be empty")
	}
	
	u.Name = name
	u.Picture = picture
	u.UpdatedAt = time.Now()
	
	return nil
}

// isValidEmail はメールアドレスの形式が有効かどうかを検証する
func (u *User) isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
