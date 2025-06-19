package http

import (
	"errors"
	userService "gateway/pkg/user"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserServer interface {
	CreateUser(c *gin.Context)
	LoginUser(c *gin.Context)
	Authorize(c *gin.Context)
	GetInfo(c *gin.Context)
}

func DefaultUserServer(service userService.UserService) UserServer {
	return &userServer{
		service: service,
	}
}

type userServer struct {
	service userService.UserService
}

func (u *userServer) CreateUser(c *gin.Context) {
	user := CreateUserRequest{}
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	createUser := userService.User{}
	createUser.Name = user.Name
	createUser.Email = user.Email
	createUser.Lastname = user.LastName
	createUser.Password = user.Password

	err = u.service.CreateUser(createUser)
	if err != nil {
		if errors.Is(err, userService.ErrInternalDBError) || errors.Is(err, userService.ErrorInternalServerError) {
			c.JSON(http.StatusInternalServerError, nil)
			return
		}
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (u *userServer) LoginUser(c *gin.Context) {
	loginUser := LoginRequest{}
	err := c.BindJSON(&loginUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	accessToken, err := u.service.Login(loginUser.Email, loginUser.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorOccurredModel{Message: err.Error()})
	}
	response := LoginResponse{
		AccessToken: accessToken,
	}
	c.JSON(http.StatusOK, response)
}

func (u *userServer) Authorize(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		c.JSON(http.StatusUnauthorized, ErrorOccurredModel{Message: "Authorization is empty"})
		c.Abort()
		return
	}
	userId, err := u.service.Authorize(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, nil)
		c.Abort()
		return
	}
	c.Request.Header.Add("user_id", strconv.Itoa(int(userId)))
	c.Next()
}

func (u *userServer) GetInfo(c *gin.Context) {
	userId := c.GetHeader("user_id")
	if userId == "" {
		c.JSON(http.StatusUnauthorized, ErrorOccurredModel{Message: "User id is empty"})
		return
	}
	id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
		c.Abort()
	}
	userInfo, err := u.service.GetUserInfo(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	response := GetUserResponse{
		Name:     userInfo.Name,
		LastName: userInfo.Lastname,
		Email:    userInfo.Email,
	}
	c.JSON(http.StatusOK, response)
}
