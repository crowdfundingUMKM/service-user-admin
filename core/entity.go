package core

import (
	"time"
)

type User struct {
	ID             int       `json:"id"`
	UnixID         string    `json:"unix_id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	PasswordHash   string    `json:"password_hash"`
	AvatarFileName string    `json:"avatar_file_name"`
	StatusAccount  string    `json:"status_account"`
	Token          string    `json:"token"`
	RefAdmin       string    `json:"ref_admin"`
	UpdateIdAdmin  string    `column:"update_id_admin"`
	UpdateAtAdmin  time.Time `column:"update_at_admin"`
	CreatedAt      time.Time `column:"created_at"`
	UpdatedAt      time.Time `column:"updated_at"`
}

type NotifAdmin struct {
	ID          int       `json:"id"`
	UserAdminId string    `json:"user_admin_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ToUser      string    `json:"to_user"`
	Document    string    `json:"document"`
	StatusNotif int       `json:"status_notif"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
