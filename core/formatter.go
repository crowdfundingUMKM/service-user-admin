package core

type UserAdminFormatter struct {
	ID            int    `json:"id"`
	UnixID        string `json:"unix_id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	Token         string `json:"token"`
	StatusAccount string `json:"status_account"`
	RefAdmin      string `json:"ref_admin"`
}

func FormatterUser(user User, token string) UserAdminFormatter {
	formatter := UserAdminFormatter{
		ID:            user.ID,
		UnixID:        user.UnixID,
		Name:          user.Name,
		Phone:         user.Phone,
		Email:         user.Email,
		Token:         token,
		StatusAccount: user.StatusAccount,
		RefAdmin:      user.RefAdmin,
	}
	return formatter
}

type UserDetailFormatter struct {
	ID            int    `json:"id"`
	UnixID        string `json:"unix_id"`
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	StatusAccount string `json:"status_account"`
}

func FormatterUserDetail(user User, updatedUser User) UserDetailFormatter {
	formatter := UserDetailFormatter{
		ID:            user.ID,
		UnixID:        user.UnixID,
		Name:          user.Name,
		Phone:         user.Phone,
		Email:         user.Email,
		StatusAccount: user.StatusAccount,
	}
	// read data before update if null use old data
	if updatedUser.Name != "" {
		formatter.Name = updatedUser.Name
	}
	if updatedUser.Phone != "" {
		formatter.Phone = updatedUser.Phone
	}
	if updatedUser.Email != "" {
		formatter.Email = updatedUser.Email
	}
	if updatedUser.StatusAccount != "" {
		formatter.StatusAccount = updatedUser.StatusAccount
	}
	return formatter
}

// for api to other service
type UserAdmin struct {
	UnixAdmin          string `json:"unix_admin"`
	StatusAccountAdmin string `json:"status_account_admin"`
}

// get user admin status
func FormatterUserAdminID(user User) UserAdmin {
	formatter := UserAdmin{
		UnixAdmin:          user.UnixID,
		StatusAccountAdmin: user.StatusAccount,
	}
	return formatter
}

// formater all user
// type AllUserFormatter struct {
// 	ID            int    `json:"id"`
// 	UnixID        string `json:"unix_id"`
// 	Name          string `json:"name"`
// 	Phone         string `json:"phone"`
// 	Email         string `json:"email"`
// 	StatusAccount string `json:"status_account"`
// 	RefAdmin      string `json:"ref_admin"`
// 	UpdateIdAdmin string `json:"update_id_admin"`
// 	UpdateAtAdmin string `json:"update_at_admin"`
// 	CreatedAt     string `json:"created_at"`
// 	UpdatedAt     string `json:"updated_at"`
// }

// func FormatterAllUsers(users []User) []AllUserFormatter {
// 	var formatters []AllUserFormatter

// 	for _, user := range users {
// 		formatter := AllUserFormatter{
// 			ID:            user.ID,
// 			UnixID:        user.UnixID,
// 			Name:          user.Name,
// 			Phone:         user.Phone,
// 			Email:         user.Email,
// 			StatusAccount: user.StatusAccount,
// 			RefAdmin:      user.RefAdmin,
// 			UpdateIdAdmin: user.UpdateIdAdmin,
// 			UpdateAtAdmin: user.UpdateAtAdmin.String(),
// 			CreatedAt:     user.CreatedAt.String(),
// 			UpdatedAt:     user.UpdatedAt.String(),
// 		}
// 		formatters = append(formatters, formatter)
// 	}

// 	return formatters
// }
