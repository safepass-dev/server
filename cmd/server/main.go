package main

import (
	"fmt"
	"net/http"

	"github.com/safepass/server/internal/api/handlers"
	"github.com/safepass/server/internal/api/middlewares"
	"github.com/safepass/server/internal/api/routes"
	"github.com/safepass/server/internal/config"
	"github.com/safepass/server/internal/database"
	"github.com/safepass/server/internal/logging"
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
	logger, err := logging.NewLogger(logging.INFO, "log.txt")

	userRepository := repositories.NewUserRepository(client)
	vaultRepository := repositories.NewVaultRepository(client, logger)
	passwordRepository := repositories.NewPasswordRepository(client, logger)

	userServices := services.NewUserServices(userRepository)
	vaultServices := services.NewVaultServices(vaultRepository, passwordRepository, &appConfig)
	authServices := services.NewAuthServices(userServices, vaultServices, &appConfig)

	authHandlers := handlers.NewAuthHandlers(*authServices)
	vaultHandlers := handlers.NewVaultHandlers(*vaultServices)

	if err != nil {
		panic(err)
	}

	logMiddleware := middlewares.NewLogMiddleware(logger)
	authMiddleware := middlewares.NewAuthMiddleware(logger, appConfig)

	router := routes.NewRouter(authMiddleware, authHandlers, vaultHandlers)
	mux := router.NewServer()

	loggedMux := logMiddleware.LogMiddlewareFunc(mux)

	err = http.ListenAndServe("0.0.0.0:5050", loggedMux)
	if err != nil {
		fmt.Println("Failed to start server")
	}
}
