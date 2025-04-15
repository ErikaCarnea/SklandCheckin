package client

import (
	"Skland/models"
	"net/http"
)

type HTTPClient interface {
	DoRequest(method, url string, body any, headers map[string]string) (*http.Response, error)
	GetSignHeaders(urlStr, method string, body any) (map[string]string, error)
	SetCred(cred string)
	SetSignToken(token string)
	ExecuteRequest(method, urlStr string, reqBody any, respTarget models.APIResponse) error
}
