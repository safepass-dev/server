package services

import (
	"strings"

	"github.com/safepass/server/internal/config"
	"github.com/safepass/server/internal/repositories"
	"github.com/safepass/server/pkg/dtos/vault"
	"github.com/safepass/server/pkg/models"
)

type VaultServicesMethods interface {
	GetVaultByUserID(string) (*models.Vault, *models.Error)
	CreateVault(int, string) *models.Error
}

type VaultServices struct {
	vaultRepository *repositories.VaultRepository
	appConfig       *config.Config

	VaultServicesMethods
}

func NewVaultServices(vaultRepository *repositories.VaultRepository, config *config.Config) *VaultServices {
	return &VaultServices{
		vaultRepository: vaultRepository,
		appConfig:       config,
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
