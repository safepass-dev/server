package models

type Vault struct {
	ID                    int    `json:"id"`
	ProtectedSymmetricKey string `json:"protected_symmetric_key"`
	Mac                   string `json:"mac"`
	Algorithm             string `json:"algorithm"`
	CreatedAt             string `json:"created_at"`
	UpdatedAt             string `json:"updated_at"`

	UserID int  `json:"user_id"`
	User   User `json:"users"`
}
