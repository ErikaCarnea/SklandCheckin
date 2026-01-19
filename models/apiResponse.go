package models

type APIResponse interface {
	GetCode() int
	GetMessage() string
}
