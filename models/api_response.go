package models

type ApiResponse interface {
	GetCode() int
	GetMessage() string
}
