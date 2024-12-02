package main

import (
	"fmt"

	"github.com/safepass/server/internal/config"
	"github.com/safepass/server/internal/database"
	"github.com/safepass/server/internal/repositories"
	"github.com/safepass/server/internal/services"
	"github.com/safepass/server/pkg/dotenv"
	"github.com/safepass/server/pkg/dtos/user"
)

func main() {
	dotenv.LoadEnv(".env")

	var appConfig config.Config
	config.LoadConfig(&appConfig)

	context, err := database.NewAppContextDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	client := context.GetSupabaseClient()

	userRepository := repositories.NewUserRepository(client)

	userServices := services.NewUserServices(userRepository)
	authServices := services.NewAuthServices(userServices, &appConfig)

	userRequest := &user.CreateUserRequest{
		Email:              "test588@gmail.com",
		MasterPasswordHash: "test1",
	}

	response, err := authServices.Register(userRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(response.Token)
}
