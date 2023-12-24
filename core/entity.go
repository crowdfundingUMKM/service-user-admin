package core

import (
	"time"
)

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
	UpdateIdAdmin  string
	UpdateAtAdmin  time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type NotifAdmin struct {
	ID          int
	UserAdminId string
	Title       string
	Description string
	ToUser      string
	Document    string
	StatusNotif int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
