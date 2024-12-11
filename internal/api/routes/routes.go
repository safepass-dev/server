package routes

import (
	"net/http"

	"github.com/safepass/server/internal/api/handlers"
)

type Router struct {
	authHandlers handlers.AuthHandlers
}

func NewRouter(authHandlers handlers.AuthHandlers) *Router {
	return &Router{
		authHandlers: authHandlers,
	}
}

func NewServer(router *Router) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/auth/login", router.authHandlers.Login)
	mux.HandleFunc("/api/v1/auth/register", router.authHandlers.Register)

	return mux
}
