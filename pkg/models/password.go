package models

type Password struct {
	ID                int    `json:"id"`
	VaultID           int    `json:"vault_id"`
	AppName           string `json:"app_name"`
	WebsiteUrl        string `json:"website_url"`
	Username          string `json:"username"`
	EncryptedPassword string `json:"encrypted_password"`
	Nonce             string `json:"nonce"`
	Mac               string `json:"mac"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}
