package models

import "time"

type User struct {
	ID                 int       `json:"id"`
	Username           string    `json:"username"`
	Email              string    `json:"email"`
	Name               string    `json:"name"`
	Surname            string    `json:"surname"`
	MasterPasswordHash string    `json:"master_password_hash"`
	Salt               string    `json:"salt"`
	IterationCount     int       `json:"iteration_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	RoleId             int       `json:"role_id"`
}
