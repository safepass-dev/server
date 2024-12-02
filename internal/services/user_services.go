package services

import (
	"fmt"

	"github.com/safepass/server/internal/consts"
	"github.com/safepass/server/internal/repositories"
	"github.com/safepass/server/pkg/dtos/user"
	"github.com/safepass/server/pkg/models"
)

type UserServicesMethods interface {
	GetUsers() ([]*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	CreateUser(*user.CreateUser) (*models.User, error)
	UpdateUser(id string, user *user.UpdateUserRequest) (*models.User, error)
	DeleteUser(id string) (*models.User, error)
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

func (u *UserServices) GetUsers() ([]*models.User, error) {
	users, err := u.userRepository.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err.Method, err.Message)
	}

	return users, nil
}

func (u *UserServices) CreateUser(userRequest *user.CreateUser) (*models.User, error) {
	res, err := u.userRepository.CreateUser(userRequest)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err.Method, err.Message)
	}

	return res, nil
}

func (u *UserServices) GetUserByUsername(username string) (*models.User, error) {
	user, err := u.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err.Method, err.Message)
	}

	return user, nil
}

func (u *UserServices) GetUserByID(id string) (*models.User, error) {
	user, err := u.userRepository.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err.Method, err.Message)
	}

	return user, nil
}

func (u *UserServices) UpdateUser(id string, userRequest *user.UpdateUserRequest) (*models.User, error) {
	newUser := &user.UpdateUser{
		Username: userRequest.Username,
		Email:    userRequest.Email,
		RoleId:   consts.Roles.USER,
	}

	res, err := u.userRepository.UpdateUser(id, newUser)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err.Method, err.Message)
	}

	return res, nil
}

func (u *UserServices) DeleteUser(id string) (*models.User, error) {
	res, err := u.userRepository.DeleteUser(id)
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err.Method, err.Message)
	}

	return res, nil
}
