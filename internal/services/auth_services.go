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
	Login(userRequest *user.LoginRequest) ([]byte, error)
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

func (a *AuthServices) Login(userRequest *user.LoginRequest) ([]byte, error) {
	user, err := a.userServices.GetUserByEmail(userRequest.Email)
	if err != nil {
		return nil, err
	}

	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		return nil, fmt.Errorf("500: Error decoding salt")
	}

	masterPasswordHash := userRequest.MasterPasswordHash
	newMasterPasswordHash := crypto.DeriveKeySha256([]byte(masterPasswordHash), salt, user.IterationCount, MASTER_PASSWORD_HASH_LENGTH)

	if user.MasterPasswordHash != base64.StdEncoding.EncodeToString(newMasterPasswordHash) {
		return nil, fmt.Errorf("401: Invalid master password")
	}

	return []byte(base64.StdEncoding.EncodeToString(newMasterPasswordHash)), nil
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
