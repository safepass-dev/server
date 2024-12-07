package user

type UpdateUserRequest struct {
	Username           string `json:"username,omitempty"`
	Email              string `json:"email,omitempty"`
	Name               string `json:"name,omitempty"`
	Surname            string `json:"surname,omitempty"`
	MasterPasswordHash string `json:"master_password_hash,omitempty"`
}
