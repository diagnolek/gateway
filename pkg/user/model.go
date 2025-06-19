package user

import (
	"errors"
	"gorm.io/gorm"
)

var (
	ErrEmailIsEmpty             = errors.New("email-is-empty")
	ErrNameIsEmpty              = errors.New("name-is-empty")
	ErrLastNameIsEmpty          = errors.New("lastname-is-empty")
	ErrPasswordIsEmpty          = errors.New("password-is-empty")
	ErrUserNotFound             = errors.New("user-not-found")
	ErrUserExists               = errors.New("user-exists")
	ErrPasswordOrEmailIsInvalid = errors.New("password-or-email-is-invalid")
	ErrInternalDBError          = errors.New("internal-db-error")
	ErrorInternalServerError    = errors.New("internal-server-error")
	ErrUnauthorized             = errors.New("unauthorized")
)

type User struct {
	gorm.Model
	Name        string `gorm:"name"`
	Lastname    string `gorm:"last_name"`
	Password    string `gorm:"password"`
	Email       string `gorm:"email"`
	AccessToken string `gorm:"access_token"`
}

type UserInfo struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Email    string `json:"email"`
}
