package services

import (
	"encoding/base64"
	"fmt"

	"github.com/safepass/server/internal/consts"
	"github.com/safepass/server/pkg/crypto"
	"github.com/safepass/server/pkg/dtos/user"
)

const (
	MASTER_PASSWORD_HASH_ITERATION_COUNT = 600000
	MASTER_PASSWORD_HASH_LENGTH          = 32
)

type AuthServicesMethods interface {
	Register(userRequest *user.CreateUserRequest) ([]byte, error)
}

type AuthServices struct {
	userServices *UserServices

	AuthServicesMethods
}

func NewAuthServices(userServices *UserServices) *AuthServices {
	return &AuthServices{
		userServices: userServices,
	}
}

func (a *AuthServices) Register(userRequest *user.CreateUserRequest) ([]byte, error) {
	salt, err := crypto.CreateRandomSalt(32)
	if err != nil {
		return nil, err
	}

	masterPasswordHash := userRequest.MasterPasswordHash
	newMasterPasswordHash := crypto.DeriveKeySha256([]byte(masterPasswordHash), salt, MASTER_PASSWORD_HASH_ITERATION_COUNT, MASTER_PASSWORD_HASH_LENGTH)

	user := &user.CreateUser{
		Username:           userRequest.Username,
		Email:              userRequest.Email,
		Name:               userRequest.Name,
		Surname:            userRequest.Surname,
		MasterPasswordHash: base64.StdEncoding.EncodeToString(newMasterPasswordHash),
		Salt:               base64.StdEncoding.EncodeToString(salt),
		IterationCount:     MASTER_PASSWORD_HASH_ITERATION_COUNT,
		RoleId:             consts.Roles.USER,
	}

	createdUser, err := a.userServices.CreateUser(user)
	if err != nil {
		return nil, err
	}

	fmt.Println(createdUser.Username)

	return []byte(base64.StdEncoding.EncodeToString(newMasterPasswordHash)), nil
}
