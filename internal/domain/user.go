// Package domain
package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID
	Name,
	Email,
	PasswordHash string
	CreatedAt,
	UpdatedAt time.Time
}
