package user

type CreateUserRequest struct {
	Username           string `json:"username" validate:"required"`
	Email              string `json:"email" validate:"required,email"`
	Name               string `json:"name,omitempty"`
	Surname            string `json:"surname,omitempty"`
	MasterPasswordHash string `json:"master_password_hash" validate:"required"`
}
