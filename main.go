package main

import (
	"errors"
	"gateway/pkg/adapter/http"
	"gateway/pkg/user"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
)

const jwtSecret = "secret"

func main() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load()
	if err != nil {
		log.Fatal("problem loading .env")
	}

	userInfra := user.DefaultUserInfra(db)
	userService := user.DefaultUserService(userInfra, jwtSecret)

	err = createAdminUser(userService)
	if err != nil {
		log.Fatal(err)
	}

	userSever := http.DefaultUserServer(userService)

	server := gin.Default()

	group := server.Group("/user/", userSever.Authorize)
	{
		group.GET("", userSever.GetInfo)
		group.POST("/user", userSever.CreateUser)
	}

	server.POST("/login", userSever.LoginUser)

	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}

}

func createAdminUser(u user.UserService) error {
	adminUser := user.User{}
	adminUser.Name = "admin"
	adminUser.Lastname = "admin"
	adminUser.Email = "admin@example.com"
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		return errors.New("admin password is empty")
	}
	adminUser.Password = adminPassword
	err := u.CreateUser(adminUser)
	if err != nil {
		if errors.Is(err, user.ErrUserExists) {
			return nil
		}
		return err
	}
	return nil
}
