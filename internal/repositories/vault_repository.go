package repositories

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/safepass/server/internal/logging"
	"github.com/safepass/server/pkg/dtos/vault"
	"github.com/safepass/server/pkg/models"
	"github.com/supabase-community/supabase-go"
)

type VaultRepositoryMethods interface {
	GetVaults() ([]*models.Vault, *models.Error)
	GetVault(string) (*models.Vault, *models.Error)
	GetVaultByUserId(string) (*models.Vault, *models.Error)
	CreateVault(*vault.CreateVault) *models.Error
	UpdateVault(*vault.CreateVault) (*models.Vault, *models.Error)
}

type VaultRepository struct {
	client *supabase.Client
	logger *logging.Logger

	VaultRepositoryMethods
}

func NewVaultRepository(client *supabase.Client, logger *logging.Logger) *VaultRepository {
	return &VaultRepository{
		client: client,
		logger: logger,
	}
}

func (v *VaultRepository) GetVaults() ([]*models.Vault, *models.Error) {
	res, n, err := v.client.From("vaults").Select("*", "exact", false).Execute()
	if err != nil {
		description := "An error occurred while retrieving vaults."
		v.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	if n <= 0 {
		description := "No vaults found"
		return nil, models.NewError(404, "NotFound", description)
	}

	var vaults []*models.Vault
	err = json.Unmarshal(res, &vaults)
	if err != nil {
		description := "An error occurred while retrieving the vaults."
		v.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	return vaults, nil
}

func (v *VaultRepository) GetVault(id string) (*models.Vault, *models.Error) {
	res, n, err := v.client.From("vaults").Select("*", "1", false).Eq("id", id).Execute()
	if err != nil {
		description := fmt.Sprintf("An error occurred while retrieving the vault with id=%s.", id)
		v.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	if n <= 0 {
		description := "No vault found with id=" + id
		return nil, models.NewError(404, "NotFound", description)
	}

	var vault *models.Vault
	err = json.Unmarshal(res, &vault)
	if err != nil {
		description := fmt.Sprintf("An error occurred while retrieving the vault with id=%s.", id)
		v.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	return vault, nil
}

func (v *VaultRepository) GetVaultByUserId(id string) (*models.Vault, *models.Error) {
	res, _, err := v.client.From("vaults").Select("*, users (*)", "1", false).Eq("user_id", id).Execute()
	if err != nil {
		description := fmt.Sprintf("An error occurred while retrieving the vault with user_id=%s.", id)
		v.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	var vaults []*models.Vault
	err = json.Unmarshal(res, &vaults)
	if err != nil {
		description := fmt.Sprintf("An error occurred while retrieving the vault with user_id=%s.", id)
		v.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	return vaults[0], nil
}

func (v *VaultRepository) CreateVault(vault *vault.CreateVault) *models.Error {
	res, _, err := v.client.From("vaults").Insert(vault, false, "", "", "1").Execute()
	if err != nil {
		description := "An error occurred while creating the vault."
		statusCode := 500
		statusText := "InternalServerError"

		if strings.Contains(err.Error(), "duplicate") {
			description = "The vault alrady exists"
			statusCode = 409
			statusText = "Confilict"
		}

		v.logger.Error(err.Error())

		return models.NewError(statusCode, statusText, description)
	}

	var response []*models.Vault
	err = json.Unmarshal(res, &response)
	if err != nil {
		description := fmt.Sprintf("Error unmarshalling response: %s", err.Error())
		return models.NewError(500, "InternalError", description)
	}

	return nil
}

func (v *VaultRepository) UpdateVault(id string, vault *vault.CreateVault) (*models.Vault, *models.Error) {
	res, _, err := v.client.From("vaults").Update(vault, "", "1").Eq("id", id).Execute()
	if err != nil {
		description := "An error occurred while creating the vault."
		statusCode := 500
		statusText := "InternalServerError"

		if strings.Contains(err.Error(), "duplicate") {
			description = "The vault alrady exists"
			statusCode = 409
			statusText = "Confilict"
		}

		v.logger.Error(err.Error())

		return nil, models.NewError(statusCode, statusText, description)
	}

	var response []*models.Vault
	err = json.Unmarshal(res, &response)
	if err != nil {
		description := fmt.Sprintf("Error unmarshalling response: %s", err.Error())
		return nil, models.NewError(500, "InternalError", description)
	}

	return response[0], nil
}
