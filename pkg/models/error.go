package models

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Method  string `json:"method"`
}
