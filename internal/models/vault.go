package models

type Vault struct {
	ID                    int        `json:"id"`
	UserID                int        `json:"user_id"`
	ProtectedSymmetricKey string     `json:"protected_symmetric_key"`
	Passwords             []Password `json:"passwords"`

	User User `json:"-"`
}
