package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/safepass/server/pkg/dtos/user"
	"github.com/safepass/server/pkg/models"
	"github.com/supabase-community/supabase-go"
)

type UserRepositoryMethods interface {
	GetUsers() ([]*models.User, *models.Error)
	GetUserByID(id string) (*models.User, *models.Error)
	GetUserByUsername(username string) (*models.User, *models.Error)
	CreateUser(*user.CreateUser) (*models.User, *models.Error)
	UpdateUser(string *user.UpdateUser) (*models.User, *models.Error)
	DeleteUser(id string) (*models.User, *models.Error)
}

type UserRepository struct {
	client *supabase.Client

	UserRepositoryMethods
}

func NewUserRepository(client *supabase.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

func (u *UserRepository) GetUsers() ([]*models.User, *models.Error) {
	res, count, err := u.client.From("users").Select("*", "exact", false).Execute()
	if err != nil {
		return nil, &models.Error{
			Status:  500,
			Message: fmt.Sprintf("Error getting users: %s", err.Error()),
			Method:  "GetUsers",
		}
	}

	if count == 0 {
		return nil, &models.Error{
			Status:  404,
			Message: "No users found",
			Method:  "GetUsers",
		}
	}

	var users []*models.User
	err = json.Unmarshal(res, &users)
	if err != nil {
		return nil, &models.Error{
			Status:  500,
			Message: fmt.Sprintf("Error unmarshalling users: %s", err.Error()),
			Method:  "GetUsers",
		}
	}

	return users, nil
}

func (u *UserRepository) GetUserByID(id string) (*models.User, *models.Error) {
	res, _, err := u.client.From("users").Select("*", "1", false).Eq("id", id).Single().Execute()
	if err != nil {
		return nil, &models.Error{
			Status:  500,
			Message: fmt.Sprintf("Error getting users: %s", err.Error()),
			Method:  "GetUserByID",
		}
	}

	var user *models.User
	err = json.Unmarshal(res, &user)
	if err != nil {
		return nil, &models.Error{
			Status:  500,
			Message: fmt.Sprintf("Error unmarshalling users: %s", err.Error()),
			Method:  "GetUserByID",
		}
	}

	return user, nil
}

func (u *UserRepository) GetUserByUsername(username string) (*models.User, *models.Error) {
	res, _, err := u.client.From("users").Select("*", "1", false).Eq("username", username).Single().Execute()
	if err != nil {
		return nil, &models.Error{
			Status:  500,
			Message: fmt.Sprintf("Error getting users: %s", err.Error()),
			Method:  "GetUserByUsername",
		}
	}

	var user *models.User
	err = json.Unmarshal(res, &user)
	if err != nil {
		return nil, &models.Error{
			Status:  500,
			Message: fmt.Sprintf("Error unmarshalling users: %s", err.Error()),
			Method:  "GetUserByUsername",
		}
	}

	return user, nil
}

func (u *UserRepository) CreateUser(user *user.CreateUser) (*models.User, *models.Error) {
	res, _, err := u.client.From("users").Insert(user, false, "", "", "1").Execute()
	if err != nil {
		return nil, &models.Error{
			Status:  500,
			Message: fmt.Sprintf("Error creating user: %s", err.Error()),
			Method:  "CreateUser",
		}
	}

	var response []*models.User
	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, &models.Error{
			Status:  500,
			Message: fmt.Sprintf("Error unmarshalling response: %s", err.Error()),
			Method:  "CreateUser",
		}
	}

	return response[0], nil
}

func (u *UserRepository) UpdateUser(userId string, user *user.UpdateUser) (*models.User, *models.Error) {
	res, _, err := u.client.From("users").Update(user, "", "1").Eq("id", userId).Execute()
	if err != nil {
		return nil, &models.Error{
			Status:  500,
			Message: fmt.Sprintf("Error updating user: %s", err.Error()),
			Method:  "UpdateUser",
		}
	}

	var response []*models.User
	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, &models.Error{
			Status:  500,
			Message: fmt.Sprintf("Error unmarshalling response: %s", err.Error()),
			Method:  "UpdateUser",
		}
	}

	return response[0], nil
}

func (u *UserRepository) DeleteUser(id string) (*models.User, *models.Error) {
	res, _, err := u.client.From("users").Delete("", "1").Eq("id", id).Execute()
	if err != nil {
		return nil, &models.Error{
			Status:  500,
			Message: fmt.Sprintf("Error deleting user: %s", err.Error()),
			Method:  "DeleteUser",
		}
	}

	var response []*models.User
	err = json.Unmarshal(res, &response)
	if err != nil {
		return nil, &models.Error{
			Status:  500,
			Message: fmt.Sprintf("Error unmarshalling response: %s", err.Error()),
			Method:  "DeleteUser",
		}
	}

	if len(response) == 0 {
		return nil, &models.Error{
			Status:  404,
			Message: "No user found",
			Method:  "DeleteUser",
		}
	}

	return response[0], nil
}
