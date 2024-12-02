package user

import "time"

type UpdateUser struct {
	Username           string    `json:"username"`
	Email              string    `json:"email"`
	Name               string    `json:"name"`
	Surname            string    `json:"surname"`
	MasterPasswordHash string    `json:"master_password_hash"`
	Salt               string    `json:"salt"`
	IterationCount     int       `json:"iteration_count"`
	RoleId             int       `json:"role_id"`
	UpdatedAt          time.Time `json:"updated_at"`
}
