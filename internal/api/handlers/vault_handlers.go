package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/safepass/server/internal/services"
	"github.com/safepass/server/pkg/dtos/password"
	"github.com/safepass/server/pkg/models"
)

type VaultHandlersFuncs interface {
	GetVault(w http.ResponseWriter, r *http.Request)
}

type VaultHandlers struct {
	vaultServices services.VaultServices

	VaultHandlersFuncs
}

func NewVaultHandlers(vaultServices services.VaultServices) *VaultHandlers {
	return &VaultHandlers{
		vaultServices: vaultServices,
	}
}

func (v *VaultHandlers) GetVault(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpError(w, http.StatusMethodNotAllowed, nil)
		return
	}

	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		data := map[string]string{"message": "An error occurred during token decryption."}
		httpError(w, http.StatusInternalServerError, data)
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		data := map[string]string{"message": "Token does not contain a userID."}
		httpError(w, http.StatusBadRequest, data)
		return
	}
	vault, merr := v.vaultServices.GetVaultByUserID(strconv.Itoa(int(userID)))
	if merr != nil {
		data := map[string]string{"message": merr.Description}
		httpError(w, merr.Code, data)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.Response{
		Status:     http.StatusOK,
		StatusText: http.StatusText(http.StatusOK),
		Data:       vault,
	}

	json.NewEncoder(w).Encode(response)
}

func (v *VaultHandlers) GetPasswords(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpError(w, http.StatusMethodNotAllowed, nil)
		return
	}

	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		data := map[string]string{"message": "An error occurred during token decryption."}
		httpError(w, http.StatusInternalServerError, data)
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		data := map[string]string{"message": "Token does not contain a userID."}
		httpError(w, http.StatusBadRequest, data)
		return
	}

	vault, merr := v.vaultServices.GetVaultByUserID(strconv.Itoa(int(userID)))
	if merr != nil {
		data := map[string]string{"message": merr.Description}
		httpError(w, merr.Code, data)
		return
	}

	passwords, merr := v.vaultServices.GetPasswords(strconv.Itoa(vault.ID))
	if merr != nil {
		data := map[string]string{"message": merr.Description}
		httpError(w, merr.Code, data)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.Response{
		Status:     http.StatusOK,
		StatusText: http.StatusText(http.StatusOK),
		Data:       passwords,
	}

	json.NewEncoder(w).Encode(response)
}

func (v *VaultHandlers) GetPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		httpError(w, http.StatusMethodNotAllowed, nil)
		return
	}

	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		data := map[string]string{"message": "An error occurred during token decryption."}
		httpError(w, http.StatusInternalServerError, data)
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		data := map[string]string{"message": "Token does not contain a userID."}
		httpError(w, http.StatusBadRequest, data)
		return
	}

	vault, merr := v.vaultServices.GetVaultByUserID(strconv.Itoa(int(userID)))
	if merr != nil {
		data := map[string]string{"message": merr.Description}
		httpError(w, merr.Code, data)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	passwords, merr := v.vaultServices.GetPassword(id, vault.ID)
	if merr != nil {
		data := map[string]string{"message": merr.Description}
		httpError(w, merr.Code, data)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.Response{
		Status:     http.StatusOK,
		StatusText: http.StatusText(http.StatusOK),
		Data:       passwords,
	}

	json.NewEncoder(w).Encode(response)
}

func (v *VaultHandlers) CreatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpError(w, http.StatusMethodNotAllowed, nil)
		return
	}

	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		data := map[string]string{"message": "An error occurred during token decryption."}
		httpError(w, http.StatusInternalServerError, data)
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		data := map[string]string{"message": "Token does not contain a userID."}
		httpError(w, http.StatusBadRequest, data)
		return
	}

	vault, merr := v.vaultServices.GetVaultByUserID(strconv.Itoa(int(userID)))
	if merr != nil {
		data := map[string]string{"message": merr.Description}
		httpError(w, merr.Code, data)
		return
	}

	var passwordRequest *password.CreatePasswordRequest
	err := json.NewDecoder(r.Body).Decode(&passwordRequest)
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
	err = validate.Struct(passwordRequest)
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

	password, merr := v.vaultServices.CreatePassword(vault.ID, passwordRequest)
	if merr != nil {
		data := map[string]string{"message": merr.Description}
		httpError(w, merr.Code, data)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.Response{
		Status:     http.StatusOK,
		StatusText: http.StatusText(http.StatusOK),
		Data:       password,
	}

	json.NewEncoder(w).Encode(response)
}

func (v *VaultHandlers) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpError(w, http.StatusMethodNotAllowed, nil)
		return
	}

	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		data := map[string]string{"message": "An error occurred during token decryption."}
		httpError(w, http.StatusInternalServerError, data)
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		data := map[string]string{"message": "Token does not contain a userID."}
		httpError(w, http.StatusBadRequest, data)
		return
	}

	vault, merr := v.vaultServices.GetVaultByUserID(strconv.Itoa(int(userID)))
	if merr != nil {
		data := map[string]string{"message": merr.Description}
		httpError(w, merr.Code, data)
		return
	}

	var passwordRequest *password.CreatePasswordRequest
	err := json.NewDecoder(r.Body).Decode(&passwordRequest)
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
	err = validate.Struct(passwordRequest)
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

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/vault/password/update/")
	if path == "" || strings.Contains(path, "/") {
		httpError(w, http.StatusBadRequest, nil)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		httpError(w, http.StatusBadRequest, nil)
		return
	}

	password, merr := v.vaultServices.UpdatePassword(id, vault.ID, passwordRequest)
	if merr != nil {
		data := map[string]string{"message": merr.Description}
		httpError(w, merr.Code, data)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.Response{
		Status:     http.StatusOK,
		StatusText: http.StatusText(http.StatusOK),
		Data:       password,
	}

	json.NewEncoder(w).Encode(response)
}

func (v *VaultHandlers) DeletePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		httpError(w, http.StatusMethodNotAllowed, nil)
		return
	}

	claims, ok := r.Context().Value("claims").(jwt.MapClaims)
	if !ok {
		data := map[string]string{"message": "An error occurred during token decryption."}
		httpError(w, http.StatusInternalServerError, data)
		return
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		data := map[string]string{"message": "Token does not contain a userID."}
		httpError(w, http.StatusBadRequest, data)
		return
	}

	vault, merr := v.vaultServices.GetVaultByUserID(strconv.Itoa(int(userID)))
	if merr != nil {
		data := map[string]string{"message": merr.Description}
		httpError(w, merr.Code, data)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/vault/password/delete/")
	if path == "" || strings.Contains(path, "/") {
		httpError(w, http.StatusBadRequest, nil)
		return
	}

	id, err := strconv.Atoi(path)
	if err != nil {
		httpError(w, http.StatusBadRequest, nil)
		return
	}

	_, merr = v.vaultServices.DeletePassword(id, vault.ID)
	if merr != nil {
		data := map[string]string{"message": merr.Description}
		httpError(w, merr.Code, data)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.Response{
		Status:     http.StatusOK,
		StatusText: http.StatusText(http.StatusOK),
		Data:       map[string]any{"id": id, "succeeded": true, "operation": "delete"},
	}

	json.NewEncoder(w).Encode(response)
}

func httpError(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	response := models.Response{
		Status:     code,
		StatusText: http.StatusText(code),
		Data:       data,
	}

	json.NewEncoder(w).Encode(response)
}
