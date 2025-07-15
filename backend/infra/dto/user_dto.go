package dto

import (
	"stackies-backend/domain/model"
	"time"
)

// UserDTO はデータベース用のユーザー構造体を表す
type UserDTO struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Name      string    `db:"name"`
	Picture   string    `db:"picture"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// ToDomain はDTOからドメインモデルに変換する
func (dto *UserDTO) ToDomain() *model.User {
	return &model.User{
		ID:        dto.ID,
		Email:     dto.Email,
		Name:      dto.Name,
		Picture:   dto.Picture,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
	}
}

// FromDomain はドメインモデルからDTOに変換する
func (dto *UserDTO) FromDomain(user *model.User) {
	dto.ID = user.ID
	dto.Email = user.Email
	dto.Name = user.Name
	dto.Picture = user.Picture
	dto.CreatedAt = user.CreatedAt
	dto.UpdatedAt = user.UpdatedAt
}
