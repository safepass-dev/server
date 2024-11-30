package models

type User struct {
	ID                 int    `json:"id"`
	Username           string `json:"username"`
	Email              string `json:"email"`
	Name               string `json:"name"`
	Surname            string `json:"surname"`
	MasterPasswordHash string `json:"master_password_hash"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}
