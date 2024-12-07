package user

type LoginRequest struct {
	Email              string `json:"email" validate:"required,email"`
	MasterPasswordHash string `json:"master_password_hash" validate:"required"`
}
