package vault

type CreateVault struct {
	UserID                int    `json:"user_id"`
	ProtectedSymmetricKey string `json:"protected_symmetric_key"`
	Mac                   string `json:"mac"`
	Algorithm             string `json:"algorithm"`
}
