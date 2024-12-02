package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/safepass/server/internal/services"
	"github.com/safepass/server/pkg/dtos/user"
	"github.com/safepass/server/pkg/models"
)

type AuthHandlers interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
}

type authHandlers struct {
	authServices services.AuthServices

	AuthHandlers
}

func NewAuthHandlers(authServices services.AuthServices) *authHandlers {
	return &authHandlers{
		authServices: authServices,
	}
}

func (a *authHandlers) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		response := models.Response{
			Status:     http.StatusMethodNotAllowed,
			StatusText: http.StatusText(http.StatusMethodNotAllowed),
			Data:       nil,
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	var loginRequest *user.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		response := models.Response{
			Status:     http.StatusBadRequest,
			StatusText: http.StatusText(http.StatusBadRequest),
			Data:       nil,
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	validate := validator.New()
	err = validate.Struct(loginRequest)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		response := models.Response{
			Status:     http.StatusBadRequest,
			StatusText: http.StatusText(http.StatusBadRequest),
			Data:       nil,
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	jwtResponse, err := a.authServices.Login(loginRequest)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)

		response := models.Response{
			Status:     http.StatusUnauthorized,
			StatusText: http.StatusText(http.StatusUnauthorized),
			Data: map[string]string{
				"message": "Email or password is incorrect",
			},
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.Response{
		Status:     http.StatusOK,
		StatusText: http.StatusText(http.StatusOK),
		Data:       jwtResponse,
	}

	json.NewEncoder(w).Encode(response)
}

func (a *authHandlers) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		response := models.Response{
			Status:     http.StatusMethodNotAllowed,
			StatusText: http.StatusText(http.StatusMethodNotAllowed),
			Data:       nil,
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	var createUser *user.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&createUser)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		response := models.Response{
			Status:     http.StatusBadRequest,
			StatusText: http.StatusText(http.StatusBadRequest),
			Data:       nil,
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	validate := validator.New()
	err = validate.Struct(createUser)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		response := models.Response{
			Status:     http.StatusBadRequest,
			StatusText: http.StatusText(http.StatusBadRequest),
			Data:       nil,
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	jwtResponse, err := a.authServices.Register(createUser)
	if err != nil {
		var statusCode int
		var statusText string

		parts := strings.Split(err.Error(), ":")
		if len(parts) != 2 {
			statusCode = http.StatusInternalServerError
			statusText = http.StatusText(http.StatusInternalServerError)
		} else {
			statusCode, err = strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				statusCode = http.StatusInternalServerError
				statusText = http.StatusText(http.StatusInternalServerError)
			} else if statusCode < 400 || statusCode >= 511 {
				statusCode = http.StatusInternalServerError
				statusText = http.StatusText(http.StatusInternalServerError)
			} else {
				statusText = strings.TrimSpace(parts[1])
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		response := models.Response{
			Status:     statusCode,
			StatusText: http.StatusText(statusCode),
			Data: map[string]string{
				"message": statusText,
			},
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.Response{
		Status:     http.StatusOK,
		StatusText: http.StatusText(http.StatusOK),
		Data:       jwtResponse,
	}

	json.NewEncoder(w).Encode(response)
}
