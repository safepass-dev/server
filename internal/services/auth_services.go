package services

import (
	"crypto/ecdsa"
	"encoding/base64"
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
	Login(userRequest *user.LoginRequest) (*models.TokenResponse, *models.Error)
	Register(userRequest *user.CreateUserRequest) (*models.TokenResponse, *models.Error)
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

func (a *AuthServices) Login(userRequest *user.LoginRequest) (*models.TokenResponse, *models.Error) {
	user, merr := a.userServices.GetUserByEmail(userRequest.Email)
	if merr != nil {
		return nil, merr
	}

	salt, err := base64.StdEncoding.DecodeString(user.Salt)
	if err != nil {
		description := "Error decoding salt"
		return nil, models.NewError(500, "InternalError", description)
	}

	masterPasswordHash, err := base64.StdEncoding.DecodeString(userRequest.MasterPasswordHash)
	if err != nil {
		description := "Password is not valid Base64"
		return nil, models.NewError(422, "UnprocessableContent", description)
	}
	newMasterPasswordHash := crypto.DeriveKeySha256([]byte(masterPasswordHash), salt, user.IterationCount, MASTER_PASSWORD_HASH_LENGTH)

	if user.MasterPasswordHash != base64.StdEncoding.EncodeToString(newMasterPasswordHash) {
		description := "Invalid master password or email"
		return nil, models.NewError(401, "Unauthorized", description)
	}

	var (
		key *ecdsa.PrivateKey
		t   *jwt.Token
		s   string
	)

	key, err = a.appConfig.GetJWTSecretKey()
	if err != nil {
		description := "Config internal error"
		return nil, models.NewError(500, "InternalError", description)
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
		description := "Error signing JWT token"
		return nil, models.NewError(500, "InternalError", description)
	}

	tokenResponse := &models.TokenResponse{
		UserID: user.ID,
		Token:  s,
	}

	return tokenResponse, nil
}

func (a *AuthServices) Register(userRequest *user.CreateUserRequest) []*models.Error {
	salt, err := crypto.CreateRandomSalt(32)
	if err != nil {
		var errors []*models.Error

		description := "Creating salt error"
		errModel := models.NewError(500, "InternalError", description)
		errors = append(errors, errModel)

		return errors
	}

	masterPasswordHash, err := base64.StdEncoding.DecodeString(userRequest.MasterPasswordHash)
	if err != nil {
		var errors []*models.Error
		description := "Password is not valid Base64"
		errModel := models.NewError(422, "UnprocessableContent", description)
		errors = append(errors, errModel)

		return errors
	}

	newMasterPasswordHash := crypto.DeriveKeySha256(masterPasswordHash, salt, MASTER_PASSWORD_HASH_ITERATION_COUNT, MASTER_PASSWORD_HASH_LENGTH)

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

	identityResult := a.userServices.CreateUser(user)
	if !identityResult.Succeeded {
		return identityResult.Errors
	}

	return nil
}
