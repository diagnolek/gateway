package main

import (
	"errors"
	grpcServer "gateway/pkg/adapter/grpc"
	"gateway/pkg/adapter/http"
	"gateway/pkg/helpers"
	"gateway/pkg/user"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
)

var (
	jwtSecret   = helpers.GetEnv("JWT_SECRET", "secret")
	grpcPort    = helpers.GetEnv("GRPC_PORT", "8081")
	httpPort    = helpers.GetEnv("HTTP_PORT", "8080")
	postgresdsn = helpers.GetEnv("POSTGRES_DSN", "")
)

func main() {
	var db *gorm.DB
	var err error
	if postgresdsn == "" {
		db, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	} else {
		db, err = gorm.Open(postgres.Open(postgresdsn), &gorm.Config{})
	}

	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&user.User{}, &user.Message{})
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

	go startHTTPServer(userService)
	startGRPCServer(userService)
}

func startGRPCServer(userService user.UserService) {
	listener, err := net.Listen("tcp", "0.0.0.0:"+grpcPort)
	if err != nil {
		log.Fatal("Failed to listen: 8081")
	}

	server := grpc.NewServer()
	userServer := grpcServer.DefaultGRPCUserServer(userService)

	grpcServer.RegisterUserServiceServer(server, userServer)

	err = server.Serve(listener)
	if err != nil {
		log.Fatal("Failed to serve: GrpcServer")
	}

	log.Println("Server GRPC started")
}

func startHTTPServer(userService user.UserService) {

	server := gin.Default()
	userSever := http.DefaultUserServer(userService)

	group := server.Group("/user/", userSever.Authorize)
	{
		group.GET("", userSever.GetInfo)
		group.POST("/user", userSever.CreateUser)
	}

	server.POST("/login", userSever.LoginUser)

	err := server.Run("0.0.0.0:" + httpPort)
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
