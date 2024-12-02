package models

type TokenResponse struct {
	UserID int    `json:"user_id"`
	Token  string `json:"token"`
}
