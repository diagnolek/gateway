package user

import (
	"fmt"
	"gateway/pkg/helpers"
)

type UserService interface {
	CreateUser(user User) error
	Login(email string, password string) (string, error)
	GetUserInfo(id uint) (UserInfo, error)
	Authorize(accessToken string) (uint, error)
}

type userService struct {
	infra     UserInfrastructure
	jwtSecret string
}

func DefaultUserService(userInfra UserInfrastructure, jwtSecret string) UserService {
	return &userService{
		infra:     userInfra,
		jwtSecret: jwtSecret,
	}
}

func (u *userService) CreateUser(newUser User) error {
	if newUser.Email == "" {
		return ErrEmailIsEmpty
	}
	if newUser.Name == "" {
		return ErrNameIsEmpty
	}
	if newUser.Password == "" {
		return ErrPasswordIsEmpty
	}

	user, _ := u.infra.GetUserByEmail(newUser.Email)
	if newUser.Email == user.Email {
		return ErrUserExists
	}

	hashedPassword, err := helpers.HashPassword(newUser.Password)
	if err != nil {
		return ErrInternalDBError
	}
	newUser.Password = hashedPassword

	err = u.infra.CreateUser(newUser)
	if err != nil {
		return ErrInternalDBError
	}

	return nil
}

func (u *userService) Login(email string, password string) (string, error) {
	if email == "" {
		return "", ErrEmailIsEmpty
	}
	if password == "" {
		return "", ErrPasswordIsEmpty
	}
	user, err := u.infra.GetUserByEmail(email)
	if err != nil {
		return "", ErrPasswordOrEmailIsInvalid
	}

	if !helpers.ComparePasswords(password, user.Password) {
		return "", ErrPasswordOrEmailIsInvalid
	}

	token, err := helpers.GenerateJWTToken(user.ID, u.jwtSecret)
	if err != nil {
		return "", ErrorInternalServerError
	}

	err = u.infra.UpdateAccessToken(user.ID, token)
	if err != nil {
		fmt.Println(err)
		return "", ErrorInternalServerError
	}

	return token, nil
}

func (u *userService) GetUserInfo(id uint) (UserInfo, error) {
	user, err := u.infra.GetUser(id)
	if err != nil {
		return UserInfo{}, ErrorInternalServerError
	}
	userInfo := UserInfo{
		Name:     user.Name,
		Lastname: user.Lastname,
		Email:    user.Email,
	}
	return userInfo, nil
}

func (u *userService) Authorize(accessToken string) (uint, error) {
	_, err := helpers.ValidateJWTToken(accessToken, u.jwtSecret)
	if err != nil {
		return 0, ErrUnauthorized
	}
	user, err := u.infra.GetUserByToken(accessToken)
	if err != nil {
		return 0, ErrUnauthorized
	}
	return user.ID, nil

}
