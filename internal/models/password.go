package models

type Password struct {
	Website   string `json:"website"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
