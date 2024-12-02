package user

type CreateUser struct {
	Username           string `json:"username"`
	Email              string `json:"email"`
	Name               string `json:"name"`
	Surname            string `json:"surname"`
	MasterPasswordHash string `json:"master_password_hash"`
	Salt               string `json:"salt"`
	IterationCount     int    `json:"iteration_count"`
	RoleId             int    `json:"role_id"`
}
