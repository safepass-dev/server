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
	GetUserByEmail(email string) (*models.User, *models.Error)
	CreateUser(*user.CreateUser) *models.IdentityResult
	UpdateUser(string *user.UpdateUser) (*models.User, *models.IdentityResult)
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
		description := "No users found"
		errModel := models.NewError(404, "NotFound", description)

		return nil, errModel
	}

	if count <= 0 {
		description := "No users found"
		errModel := models.NewError(404, "NotFound", description)

		return nil, errModel
	}

	var users []*models.User
	err = json.Unmarshal(res, &users)
	if err != nil {
		description := fmt.Sprintf("Error unmarshalling users: %s", err.Error())
		errModel := models.NewError(500, "InternalError", description)

		return nil, errModel
	}

	return users, nil
}

func (u *UserRepository) GetUserByID(id string) (*models.User, *models.Error) {
	res, _, err := u.client.From("users").Select("*", "1", false).Eq("id", id).Single().Execute()
	if err != nil {
		description := "No user found"
		errModel := models.NewError(404, "NotFound", description)

		return nil, errModel
	}

	var user *models.User
	err = json.Unmarshal(res, &user)
	if err != nil {
		description := fmt.Sprintf("Error unmarshalling users: %s", err.Error())
		errModel := models.NewError(500, "InternalError", description)

		return nil, errModel
	}

	return user, nil
}

func (u *UserRepository) GetUserByEmail(email string) (*models.User, *models.Error) {
	res, _, err := u.client.From("users").Select("*", "1", false).Eq("email", email).Single().Execute()
	if err != nil {
		description := "No user found"
		errModel := models.NewError(404, "NotFound", description)

		return nil, errModel
	}

	var user *models.User
	err = json.Unmarshal(res, &user)
	if err != nil {
		description := fmt.Sprintf("Error unmarshalling users: %s", err.Error())
		errModel := models.NewError(500, "InternalError", description)

		return nil, errModel
	}

	return user, nil
}

func (u *UserRepository) CreateUser(user *user.CreateUser) *models.IdentityResult {
	res, _, err := u.client.From("users").Insert(user, false, "", "", "1").Execute()
	if err != nil {
		var errors []*models.Error

		if err.Error() == "(23505) duplicate key value violates unique constraint \"users_username_key\"" {
			description := "Username already exists"
			errModel := models.NewError(409, "Confilict", description)

			errors = append(errors, errModel)
		} else if err.Error() == "(23505) duplicate key value violates unique constraint \"users_email_key\"" {
			description := "Email already exists"
			errModel := models.NewError(409, "Confilict", description)

			errors = append(errors, errModel)
		} else {
			description := "Content is not valid user"
			errModel := models.NewError(422, "UnprocessableContent", description)

			errors = append(errors, errModel)
		}

		return &models.IdentityResult{
			Errors:    errors,
			Succeeded: false,
			Message:   "Registration error",
		}
	}

	var response []*models.User
	err = json.Unmarshal(res, &response)
	if err != nil {
		var errors []*models.Error

		description := fmt.Sprintf("Error unmarshalling response: %s", err.Error())
		errModel := models.NewError(500, "InternalError", description)

		errors = append(errors, errModel)

		return &models.IdentityResult{
			Errors:    errors,
			Succeeded: false,
			Message:   "Registration error",
		}
	}

	return &models.IdentityResult{
		Errors:    nil,
		Succeeded: true,
		Message:   response,
	}
}

func (u *UserRepository) UpdateUser(userId string, user *user.UpdateUser) (*models.User, *models.IdentityResult) {
	res, _, err := u.client.From("users").Update(user, "", "1").Eq("id", userId).Execute()
	if err != nil {
		var errors []*models.Error

		if err.Error() == "(23505) duplicate key value violates unique constraint \"users_username_key\"" {
			description := "Username already exists"
			errModel := models.NewError(409, "Confilict", description)

			errors = append(errors, errModel)
		} else if err.Error() == "(23505) duplicate key value violates unique constraint \"users_email_key\"" {
			description := "Email already exists"
			errModel := models.NewError(409, "Confilict", description)

			errors = append(errors, errModel)
		} else {
			description := "Content is not valid user"
			errModel := models.NewError(422, "UnprocessableContent", description)

			errors = append(errors, errModel)
		}

		return nil, &models.IdentityResult{
			Errors:    errors,
			Succeeded: false,
			Message:   "Update error",
		}
	}

	var response []*models.User
	err = json.Unmarshal(res, &response)
	if err != nil {
		var errors []*models.Error

		description := fmt.Sprintf("Error unmarshalling response: %s", err.Error())
		errModel := models.NewError(500, "InternalError", description)

		errors = append(errors, errModel)

		return nil, &models.IdentityResult{
			Errors:    errors,
			Succeeded: false,
			Message:   "Update error",
		}
	}

	return response[0], &models.IdentityResult{
		Errors:    nil,
		Succeeded: true,
		Message:   "Update successful",
	}
}

func (u *UserRepository) DeleteUser(id string) (*models.User, *models.Error) {
	res, _, err := u.client.From("users").Delete("", "1").Eq("id", id).Execute()
	if err != nil {
		description := fmt.Sprintf("Error deleting user: %s", err.Error())
		errModel := models.NewError(500, "InternalError", description)

		return nil, errModel
	}

	var response []*models.User
	err = json.Unmarshal(res, &response)
	if err != nil {
		description := fmt.Sprintf("Error unmarshalling response: %s", err.Error())
		errModel := models.NewError(500, "InternalError", description)

		return nil, errModel
	}

	if len(response) == 0 {
		description := "No user found"
		errModel := models.NewError(404, "NotFound", description)

		return nil, errModel
	}

	return response[0], nil
}
