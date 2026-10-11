package domain

import (
	"time"
)

type User struct {
	ID UUID
	Name,
	Email,
	PasswordHash string
	CreatedAt,
	UpdatedAt time.Time
}
