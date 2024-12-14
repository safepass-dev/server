package password

type CreatePasswordRequest struct {
	AppName           string `json:"app_name,omitempty"`
	Uri               string `json:"uri,omitempty"`
	Username          string `json:"username,omitempty"`
	EncryptedPassword string `json:"encrypted_password" validate:"required"`
}
