package admin

import (
	"errors"
	"os"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	DeactivateAccountUser(input DeactiveUserInput) (bool, error)
	ActivateAccountUser(input DeactiveUserInput) (bool, error)
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	SaveToken(UnixID string, Token string) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	IsPhoneAvailable(input CheckPhoneInput) (bool, error)
	GetUserByUnixID(UnixID string) (User, error)
	UpdateUserByUnixID(UnixID string, input UpdateUserInput) (User, error)
	DeleteToken(UnixID string) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) DeactivateAccountUser(input DeactiveUserInput) (bool, error) {
	user, err := s.repository.FindByUnixID(input.UnixID)
	user.StatusAccount = "Deactive"
	_, err = s.repository.UpdateStatusAccount(user)

	if err != nil {
		return false, err
	}

	if user.UnixID == "" {
		return true, nil
	}
	return true, nil
}

func (s *service) ActivateAccountUser(input DeactiveUserInput) (bool, error) {
	user, err := s.repository.FindByUnixID(input.UnixID)
	user.StatusAccount = "Active"
	_, err = s.repository.UpdateStatusAccount(user)

	if err != nil {
		return false, err
	}

	if user.UnixID == "" {
		return true, nil
	}
	return true, nil
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.UnixID = uuid.New().String()[:12]
	user.Name = input.Name
	user.Email = input.Email
	user.Phone = input.Phone

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)
	// convert data os env to string
	user.StatusAccount = string(os.Getenv("STATUS_ACCOUNT"))

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("No user found on that email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return user, err
	}

	return user, nil
}

// save token to database
func (s *service) SaveToken(UnixID string, Token string) (User, error) {
	user, err := s.repository.FindByUnixID(UnixID)
	user.Token = Token
	_, err = s.repository.UpdateToken(user)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}

	return false, nil
}

func (s *service) IsPhoneAvailable(input CheckPhoneInput) (bool, error) {
	phone := input.Phone

	user, err := s.repository.FindByPhone(phone)
	if err != nil {
		return false, err
	}

	if user.UnixID == "" {
		return true, nil
	}

	return false, nil
}

func (s *service) GetUserByUnixID(UnixID string) (User, error) {
	user, err := s.repository.FindByUnixID(UnixID)
	if err != nil {
		return user, err
	}

	if user.UnixID == "" {
		return user, errors.New("No user found on with that ID")
	}

	return user, nil
}

func (s *service) UpdateUserByUnixID(UnixID string, input UpdateUserInput) (User, error) {
	user, err := s.repository.FindByUnixID(UnixID)
	if err != nil {
		return user, err
	}

	if user.UnixID == "" {
		return user, errors.New("No user found on with that ID")
	}

	user.Name = input.Name
	user.Email = input.Email
	user.Phone = input.Phone

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

// logout
func (s *service) DeleteToken(UnixID string) (User, error) {
	user, err := s.repository.FindByUnixID(UnixID)
	if err != nil {
		return user, err
	}

	if user.UnixID == "" {
		return user, errors.New("No user found on with that ID")
	}

	user.Token = ""

	updatedUser, err := s.repository.UpdateToken(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}
