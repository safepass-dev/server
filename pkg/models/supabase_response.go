package models

type SupabaseResponse struct {
	Status     int    `json:"status"`
	StatusText string `json:"statusText"`
}
