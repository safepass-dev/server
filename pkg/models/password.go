package models

type Password struct {
	ID                int    `json:"id"`
	VaultID           int    `json:"vault_id"`
	AppName           string `json:"app_name"`
	Uri               string `json:"uri"`
	Username          string `json:"username"`
	EncryptedPassword string `json:"encrypted_password"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}
