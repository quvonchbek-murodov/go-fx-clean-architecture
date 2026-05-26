package user

import (
	"time"

	"golang-project-structure/internal/dto"
)

type UserDTO struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListOutput struct {
	Items []*UserDTO  `json:"items"`
	Meta  dto.MetaDTO `json:"meta"`
}
