package services

import (
	"strconv"
	"strings"

	"github.com/safepass/server/internal/config"
	"github.com/safepass/server/internal/repositories"
	"github.com/safepass/server/pkg/dtos/password"
	"github.com/safepass/server/pkg/dtos/vault"
	"github.com/safepass/server/pkg/models"
)

type VaultServicesMethods interface {
	GetVaultByUserID(string) (*models.Vault, *models.Error)
	CreateVault(int, string) *models.Error

	GetPasswords(vaultID string) ([]*models.Password, *models.Error)
	GetPassword(passwordID string, vaultID int) (*models.Password, *models.Error)
	CreatePassword(vaultID int, passwordRequest *password.CreatePasswordRequest) (*models.Password, *models.Error)
	UpdatePassword(passwordID int, vaultID int, passwordRequest *password.CreatePasswordRequest) (*models.Password, *models.Error)
	DeletePassword(id int, vaultID int) (*models.Password, *models.Error)
}

type VaultServices struct {
	vaultRepository    *repositories.VaultRepository
	passwordRepository *repositories.PasswordRepository
	appConfig          *config.Config

	VaultServicesMethods
}

func NewVaultServices(vaultRepository *repositories.VaultRepository, passwordRepository *repositories.PasswordRepository, config *config.Config) *VaultServices {
	return &VaultServices{
		vaultRepository:    vaultRepository,
		passwordRepository: passwordRepository,
		appConfig:          config,
	}
}

func (v *VaultServices) GetVaultByUserID(userID string) (*models.Vault, *models.Error) {
	vault, err := v.vaultRepository.GetVaultByUserId(userID)

	return vault, err
}

func (v *VaultServices) CreateVault(userID int, protectedSymmetricKey string) *models.Error {
	parts := strings.Split(protectedSymmetricKey, ":")
	if len(parts) != 2 {
		return models.NewError(422, "Unprocessable Content", "Protected symmetric key is not valid.")
	}

	mac := parts[0]
	symmetricKey := parts[1]

	vault := &vault.CreateVault{
		UserID:                userID,
		ProtectedSymmetricKey: symmetricKey,
		Mac:                   mac,
		Algorithm:             "AESCBCPKCS5Padding",
	}

	err := v.vaultRepository.CreateVault(vault)
	return err
}

func (v *VaultServices) GetPasswords(vaultID string) ([]*models.Password, *models.Error) {
	passwords, err := v.passwordRepository.GetPasswordsByVaultID(vaultID)

	return passwords, err
}

func (v *VaultServices) GetPassword(passwordID string, vaultID int) (*models.Password, *models.Error) {
	password, err := v.passwordRepository.GetPassword(passwordID)
	if err != nil {
		return nil, err
	}

	if password.VaultID != vaultID {
		return nil, models.NewError(401, "Unauthorized", "")
	}

	return password, nil
}

func (v *VaultServices) CreatePassword(vaultID int, passwordRequest *password.CreatePasswordRequest) (*models.Password, *models.Error) {
	pw := &password.CreatePassword{
		VaultID:           vaultID,
		AppName:           passwordRequest.AppName,
		Uri:               passwordRequest.Uri,
		Username:          passwordRequest.Username,
		EncryptedPassword: passwordRequest.EncryptedPassword,
	}

	newPw, merr := v.passwordRepository.CreatePassword(pw)
	return newPw, merr
}

func (v *VaultServices) UpdatePassword(passwordID int, vaultID int, passwordRequest *password.CreatePasswordRequest) (*models.Password, *models.Error) {
	pw := &password.CreatePassword{
		VaultID:           vaultID,
		AppName:           passwordRequest.AppName,
		Uri:               passwordRequest.Uri,
		Username:          passwordRequest.Username,
		EncryptedPassword: passwordRequest.EncryptedPassword,
	}

	newPw, merr := v.passwordRepository.UpdatePassword(strconv.Itoa(passwordID), pw)
	return newPw, merr
}

func (v *VaultServices) DeletePassword(id int, vaultID int) (*models.Password, *models.Error) {
	newPw, merr := v.passwordRepository.DeletePassword(strconv.Itoa(id), strconv.Itoa(vaultID))

	return newPw, merr
}
