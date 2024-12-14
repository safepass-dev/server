package repositories

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/safepass/server/internal/logging"
	"github.com/safepass/server/pkg/dtos/password"
	"github.com/safepass/server/pkg/models"
	"github.com/supabase-community/supabase-go"
)

type PasswordRepositoryMethods interface {
	GetPasswords() ([]*models.Password, *models.Error)
	GetPassword(string) (*models.Password, *models.Error)
	GetPasswordsByVaultID(string) ([]*models.Password, *models.Error)
	CreatePassword(*password.CreatePassword) (*models.Password, *models.Error)
	UpdatePassword(string, *password.CreatePassword) (*models.Password, *models.Error)
	DeletePassword(string, string) (*models.Password, *models.Error)
}

type PasswordRepository struct {
	client *supabase.Client
	logger *logging.Logger

	PasswordRepositoryMethods
}

func NewPasswordRepository(client *supabase.Client, logger *logging.Logger) *PasswordRepository {
	return &PasswordRepository{
		client: client,
		logger: logger,
	}
}

func (p *PasswordRepository) GetPasswords() ([]*models.Password, *models.Error) {
	res, n, err := p.client.From("passwords").Select("*", "exact", false).Execute()
	if err != nil {
		description := "An error occurred while retrieving passwords."
		p.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	if n <= 0 {
		return []*models.Password{}, nil
	}

	var passwords []*models.Password
	err = json.Unmarshal(res, &passwords)
	if err != nil {
		description := "An error occurred while retrieving passwords."
		p.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	return passwords, nil
}

func (p *PasswordRepository) GetPassword(id string) (*models.Password, *models.Error) {
	res, _, err := p.client.From("passwords").Select("*", "1", false).Eq("id", id).Execute()
	if err != nil {
		description := "An error occurred while retrieving password."
		p.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	if string(res) == "[]" {
		description := "Password not found with id=" + id
		return nil, models.NewError(404, "NotFound", description)
	}

	var password []*models.Password
	err = json.Unmarshal(res, &password)
	if err != nil {
		description := fmt.Sprintf("An error occurred while retrieving the vault with id=%s.", id)
		p.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	return password[0], nil
}

func (p *PasswordRepository) GetPasswordsByVaultID(vaultID string) ([]*models.Password, *models.Error) {
	res, n, err := p.client.From("passwords").Select("*", "exact", false).Eq("vault_id", vaultID).Execute()
	if err != nil {
		description := "An error occurred while retrieving passwords."
		p.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	if n <= 0 {
		return []*models.Password{}, nil
	}

	var passwords []*models.Password
	err = json.Unmarshal(res, &passwords)
	if err != nil {
		description := "An error occurred while retrieving passwords."
		p.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	return passwords, nil
}

func (p *PasswordRepository) CreatePassword(createPassword *password.CreatePassword) (*models.Password, *models.Error) {
	res, _, err := p.client.From("passwords").Insert(createPassword, false, "", "", "1").Execute()
	if err != nil {
		description := "An error occurred while retrieving passwords."
		p.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	var passwords []*models.Password
	err = json.Unmarshal(res, &passwords)
	if err != nil {
		description := "An error occurred while retrieving passwords."
		p.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	return passwords[0], nil
}

func (p *PasswordRepository) UpdatePassword(passwordID string, createPassword *password.CreatePassword) (*models.Password, *models.Error) {
	res, _, err := p.client.From("passwords").Update(createPassword, "", "1").Eq("id", passwordID).Eq("vault_id", strconv.Itoa(createPassword.VaultID)).Execute()
	if err != nil {
		description := "An error occurred while retrieving passwords."
		p.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	var passwords []*models.Password
	err = json.Unmarshal(res, &passwords)
	if err != nil {
		description := "An error occurred while retrieving passwords."
		p.logger.Error(err.Error())

		return nil, models.NewError(500, "InternalServerError", description)
	}

	return passwords[0], nil
}

func (v *PasswordRepository) DeletePassword(passwordID string, vaultID string) (*models.Password, *models.Error) {
	res, _, err := v.client.From("passwords").Delete("", "1").Eq("id", passwordID).Eq("vault_id", vaultID).Execute()
	if err != nil {
		description := fmt.Sprintf("Error deleting password: %s", err.Error())
		errModel := models.NewError(500, "InternalError", description)

		return nil, errModel
	}

	var response []*models.Password
	err = json.Unmarshal(res, &response)
	if err != nil {
		description := fmt.Sprintf("Error unmarshalling response: %s", err.Error())
		errModel := models.NewError(500, "InternalError", description)

		return nil, errModel
	}

	if len(response) == 0 {
		description := "No password found"
		errModel := models.NewError(404, "NotFound", description)

		return nil, errModel
	}

	return response[0], nil
}
