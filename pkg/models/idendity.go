package models

type IdentityResult struct {
	Errors    []*Error
	Succeeded bool
	Message   string
}
