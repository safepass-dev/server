package models

type Error struct {
	Code        int    `json:"code"`
	CodeString  string `json:"code_string"`
	Description string `json:"description"`
}

func NewError(code int, codeString, description string) *Error {
	return &Error{
		Code:        code,
		CodeString:  codeString,
		Description: description,
	}
}
