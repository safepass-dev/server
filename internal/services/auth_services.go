package services

import (
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/safepass/server/internal/config"
	"github.com/safepass/server/internal/consts"
	"github.com/safepass/server/pkg/crypto"
	"github.com/safepass/server/pkg/dtos/user"
	"github.com/safepass/server/pkg/models"
)

const (
	MASTER_PASSWORD_HASH_ITERATION_COUNT = 600000
	MASTER_PASSWORD_HASH_LENGTH          = 32
)

type AuthServicesMethods interface {
	Login(userRequest *user.LoginRequest) (*models.TokenResponse, error)
	Register(userRequest *user.CreateUserRequest) (*models.TokenResponse, error)
}

type AuthServices struct {
	userServices *UserServices
	appConfig    *config.Config

	AuthServicesMethods
}

func NewAuthServices(userServices *UserServices, config *config.Config) *AuthServices {
	return &AuthServices{
		userServices: userServices,
		appConfig:    config,
	}
}

func (a *AuthServices) Login(userRequest *user.LoginRequest) (*models.TokenResponse, error) {
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
		return nil, fmt.Errorf("401: Invalid master password or email")
	}

	var (
		key *ecdsa.PrivateKey
		t   *jwt.Token
		s   string
	)

	key, err = a.appConfig.GetJWTSecretKey()
	if err != nil {
		return nil, err
	}

	t = jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": "safepass",
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Second * a.appConfig.JWT.ExpirationTime).Unix(),
		"aud": "safepass-mobile",
		"roles": []string{
			"user",
		},
		"username":    user.Username,
		"email":       user.Email,
		"auth_method": "master_password",
	})

	s, err = t.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("500: Error signing JWT token")
	}

	tokenResponse := &models.TokenResponse{
		UserID: user.ID,
		Token:  s,
	}

	return tokenResponse, nil
}

func (a *AuthServices) Register(userRequest *user.CreateUserRequest) (*models.TokenResponse, error) {
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

	var (
		key *ecdsa.PrivateKey
		t   *jwt.Token
		s   string
	)

	key, err = a.appConfig.GetJWTSecretKey()
	if err != nil {
		return nil, err
	}

	t = jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"iss": "safepass",
		"sub": createdUser.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Second * a.appConfig.JWT.ExpirationTime).Unix(),
		"aud": "safepass-mobile",
		"roles": []string{
			"user",
		},
		"username":    createdUser.Username,
		"email":       createdUser.Email,
		"auth_method": "master_password",
	})

	s, err = t.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("500: Error signing JWT token")
	}

	tokenResponse := &models.TokenResponse{
		UserID: createdUser.ID,
		Token:  s,
	}

	return tokenResponse, nil
}
