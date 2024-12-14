package routes

import (
	"net/http"

	"github.com/safepass/server/internal/api/handlers"
	"github.com/safepass/server/internal/api/middlewares"
)

type Router struct {
	authMiddleware *middlewares.AuthMiddleware

	authHandlers  *handlers.AuthHandlers
	vaultHandlers *handlers.VaultHandlers
}

func NewRouter(
	autMiddleware *middlewares.AuthMiddleware,
	authHandlers *handlers.AuthHandlers,
	vaultHandlers *handlers.VaultHandlers,
) *Router {
	return &Router{
		authMiddleware: autMiddleware,
		authHandlers:   authHandlers,
		vaultHandlers:  vaultHandlers,
	}
}

func (r *Router) NewServer() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/auth/login", r.authHandlers.Login)
	mux.HandleFunc("/api/v1/auth/register", r.authHandlers.Register)

	mux.Handle("/api/v1/vault/@me", r.authMiddleware.AuthMiddlewareFunc(http.HandlerFunc(r.vaultHandlers.GetVault)))

	mux.Handle("/api/v1/vault/passwords", r.authMiddleware.AuthMiddlewareFunc(http.HandlerFunc(r.vaultHandlers.GetPasswords)))
	mux.Handle("/api/v1/vault/password", r.authMiddleware.AuthMiddlewareFunc(http.HandlerFunc(r.vaultHandlers.GetPassword)))
	mux.Handle("/api/v1/vault/password/create", r.authMiddleware.AuthMiddlewareFunc(http.HandlerFunc(r.vaultHandlers.CreatePassword)))
	mux.Handle("/api/v1/vault/password/update/", r.authMiddleware.AuthMiddlewareFunc(http.HandlerFunc(r.vaultHandlers.UpdatePassword)))
	mux.Handle("/api/v1/vault/password/delete/", r.authMiddleware.AuthMiddlewareFunc(http.HandlerFunc(r.vaultHandlers.DeletePassword)))

	return mux
}
