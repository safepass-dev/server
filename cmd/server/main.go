package main

import (
	"fmt"

	"github.com/safepass/server/internal/database"
	"github.com/safepass/server/internal/repositories"
	"github.com/safepass/server/internal/services"
	"github.com/safepass/server/pkg/dtos/user"
)

func main() {
	context, err := database.NewAppContextDB()
	if err != nil {
		panic(err)
	}
	client := context.GetSupabaseClient()

	userRepository := repositories.NewUserRepository(client)

	userServices := services.NewUserServices(userRepository)
	authServices := services.NewAuthServices(userServices)

	newUser := &user.CreateUserRequest{
		Username:           "test58",
		Email:              "test58@gmail.com",
		MasterPasswordHash: "test",
	}

	p, err := authServices.Register(newUser)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(p))

	users, err := userServices.GetUsers()
	if err != nil {
		panic(err)
	}

	for _, user := range users {
		fmt.Println(user.Username, user.ID)
	}
}
