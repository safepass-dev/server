package models

type Vault struct {
	ID                    int    `json:"id"`
	UserID                int    `json:"user_id"`
	ProtectedSymmetricKey string `json:"protected_symmetric_key"`
	IV                    string `json:"iv"`
	Algorithm             string `json:"algorithm"`
	CreatedAt             string `json:"created_at"`
	UpdatedAt             string `json:"updated_at"`

	User User `json:"-"`
}
