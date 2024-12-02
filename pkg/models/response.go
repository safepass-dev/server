package models

type Response struct {
	Status     int         `json:"status"`
	StatusText string      `json:"statusText"`
	Data       interface{} `json:"data"`
}
