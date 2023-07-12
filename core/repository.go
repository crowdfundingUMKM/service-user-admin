package core

import "gorm.io/gorm"

// KONTRAK
type Repository interface {
	Save(user User) (User, error)
	FindByUnixID(unix_id string) (User, error)
	FindByEmail(email string) (User, error)
	UpdateToken(user User) (User, error)
	FindByPhone(phone string) (User, error)
	Update(user User) (User, error)
	UpdateStatusAccount(user User) (User, error)
	DeleteUser(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *repository) FindByUnixID(unix_id string) (User, error) {
	var user User

	err := r.db.Where("unix_id = ?", unix_id).Find(&user).Error

	if err != nil {
		return user, err
	}
	return user, nil

}

// UpdateToken
func (r *repository) UpdateToken(user User) (User, error) {
	err := r.db.Model(&user).Update("token", user.Token).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error) {
	var user User

	err := r.db.Where("email = ?", email).Find(&user).Error

	if err != nil {
		return user, err
	}
	return user, nil

}

func (r *repository) FindByPhone(phone string) (User, error) {
	var user User

	err := r.db.Where("phone = ?", phone).Find(&user).Error

	if err != nil {
		return user, err
	}
	return user, nil

}

func (r *repository) Update(user User) (User, error) {
	err := r.db.Model(&user).Updates(User{Name: user.Name, Phone: user.Phone, Email: user.Email}).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) UpdateStatusAccount(user User) (User, error) {
	err := r.db.Model(&user).Update("status_account", user.StatusAccount).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

// delete user
func (r *repository) DeleteUser(user User) (User, error) {
	err := r.db.Delete(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}
