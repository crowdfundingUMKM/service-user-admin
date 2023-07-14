package core

import "time"

type User struct {
	ID             int
	UnixID         string
	Name           string
	Email          string
	Phone          string
	PasswordHash   string
	AvatarFileName string
	StatusAccount  string
	Token          string
	RefAdmin       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
