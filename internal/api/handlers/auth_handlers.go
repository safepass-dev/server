package handlers

import (
	"encoding/json"
	"net/http"

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

	jwtResponse, merr := a.authServices.Login(loginRequest)
	if merr != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(merr.Code)

		response := models.Response{
			Status:     merr.Code,
			StatusText: http.StatusText(merr.Code),
			Data: map[string]string{
				"message": merr.Description,
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

	errors := a.authServices.Register(createUser)
	if errors != nil {
		merr := errors[0]

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(merr.Code)

		response := models.Response{
			Status:     merr.Code,
			StatusText: http.StatusText(merr.Code),
			Data: map[string]string{
				"message": merr.Description,
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
		Data: map[string]string{
			"message": "Registration successful",
		},
	}

	json.NewEncoder(w).Encode(response)
}
