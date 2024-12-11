package main

import (
	"fmt"
	"net/http"

	"github.com/safepass/server/internal/api/handlers"
	"github.com/safepass/server/internal/api/routes"
	"github.com/safepass/server/internal/config"
	"github.com/safepass/server/internal/database"
	"github.com/safepass/server/internal/repositories"
	"github.com/safepass/server/internal/services"
	"github.com/safepass/server/pkg/dotenv"
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

	authHandlers := handlers.NewAuthHandlers(*authServices)

	router := routes.NewRouter(authHandlers)
	mux := routes.NewServer(router)

	err = http.ListenAndServe("0.0.0.0:5050", mux)
	if err != nil {
		fmt.Println("Failed to start server")
	}
}
