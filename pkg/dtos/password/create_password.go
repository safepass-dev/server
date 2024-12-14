package password

type CreatePassword struct {
	VaultID           int    `json:"vault_id" validate:"required"`
	AppName           string `json:"app_name,omitempty"`
	Uri               string `json:"uri,omitempty"`
	Username          string `json:"username,omitempty"`
	EncryptedPassword string `json:"encrypted_password" validate:"required"`
}
