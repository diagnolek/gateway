package user

import (
	"errors"
	"gorm.io/gorm"
)

type UserInfrastructure interface {
	CreateUser(user User) error
	GetUser(id uint) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUserByToken(accessToken string) (User, error)
	UpdateAccessToken(id uint, newAccessToken string) error
	//Chat
	CreateMessage(message Message) error
	GetMessageHistory(from uint, to uint) (Messages, error)
}

type userInfra struct {
	db *gorm.DB
}

func DefaultUserInfra(db *gorm.DB) UserInfrastructure {
	return &userInfra{
		db,
	}
}

func (u *userInfra) CreateUser(user User) error {
	result := u.db.Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *userInfra) GetUser(id uint) (User, error) {
	user := User{}
	result := u.db.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, nil
		}
		return User{}, result.Error
	}
	return user, nil
}

func (u *userInfra) GetUserByEmail(email string) (User, error) {
	user := User{}
	result := u.db.First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, nil
		}
		return User{}, result.Error
	}
	return user, nil
}

func (u *userInfra) GetUserByToken(accessToken string) (User, error) {
	user := User{}
	result := u.db.First(&user, "access_token = ?", accessToken)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return User{}, nil
		}
		return User{}, result.Error
	}
	return user, nil
}

func (u *userInfra) UpdateAccessToken(id uint, newAccessToken string) error {
	user := User{}
	user.ID = id

	result := u.db.Model(&user).Update("access_token", newAccessToken)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		return result.Error
	}
	return nil

}

func (u *userInfra) CreateMessage(message Message) error {
	result := u.db.Create(&message)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *userInfra) GetMessageHistory(from uint, to uint) (Messages, error) {
	messages := []Message{}
	result := u.db.Where("from_user = ? AND to_user = ?", from, to).Find(&messages)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []Message{}, nil
		}
		return []Message{}, result.Error
	}
	return messages, nil
}
