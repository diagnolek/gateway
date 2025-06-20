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
	Name        string `gorm:"column:name"`
	Lastname    string `gorm:"column:last_name"`
	Password    string `gorm:"column:password"`
	Email       string `gorm:"column:email"`
	AccessToken string `gorm:"column:access_token"`
}

type UserInfo struct {
	Name     string `json:"name"`
	Lastname string `json:"lastname"`
	Email    string `json:"email"`
}

type Message struct {
	gorm.Model
	From uint   `gorm:"column:from_user"`
	To   uint   `gorm:"column:to_user"`
	Date int    `gorm:"column:msg_date"`
	Text string `gorm:"column:msg_text"`
}

type Messages []Message

func (m Messages) Len() int {
	return len(m)
}

func (m Messages) Less(i, j int) bool {
	return m[i].Date < m[j].Date
}

func (m Messages) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
