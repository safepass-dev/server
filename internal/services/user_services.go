package services

import (
	"github.com/safepass/server/internal/consts"
	"github.com/safepass/server/internal/repositories"
	"github.com/safepass/server/pkg/dtos/user"
	"github.com/safepass/server/pkg/models"
)

type UserServicesMethods interface {
	GetUsers() ([]*models.User, *models.Error)
	GetUserByID(id string) (*models.User, *models.Error)
	GetUserByEmail(email string) (*models.User, *models.Error)
	CreateUser(*user.CreateUser) *models.IdentityResult
	UpdateUser(id string, user *user.UpdateUserRequest) (*models.User, *models.IdentityResult)
	DeleteUser(id string) (*models.User, *models.Error)
}

type UserServices struct {
	userRepository *repositories.UserRepository

	UserServicesMethods
}

func NewUserServices(userRepository *repositories.UserRepository) *UserServices {
	return &UserServices{
		userRepository: userRepository,
	}
}

func (u *UserServices) GetUsers() ([]*models.User, *models.Error) {
	users, err := u.userRepository.GetUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserServices) GetUserByID(id string) (*models.User, *models.Error) {
	user, err := u.userRepository.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserServices) GetUserByEmail(email string) (*models.User, *models.Error) {
	user, err := u.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserServices) CreateUser(userRequest *user.CreateUser) *models.IdentityResult {
	identityResult := u.userRepository.CreateUser(userRequest)

	return identityResult
}

func (u *UserServices) UpdateUser(id string, userRequest *user.UpdateUserRequest) (*models.User, *models.IdentityResult) {
	newUser := &user.UpdateUser{
		Username: userRequest.Username,
		Email:    userRequest.Email,
		RoleId:   consts.Roles.USER,
	}

	res, identityResult := u.userRepository.UpdateUser(id, newUser)
	return res, identityResult
}

func (u *UserServices) DeleteUser(id string) (*models.User, *models.Error) {
	res, err := u.userRepository.DeleteUser(id)

	return res, err
}
