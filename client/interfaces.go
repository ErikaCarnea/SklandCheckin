package client

import "github.com/ErikaCarnea/Skland/models"

type SklandHTTPClient interface {
	GetSignHeaders(urlStr, method string, body []byte) (map[string]string, error)
	SetCred(cred string)
	SetSignToken(token string)
	ExecuteRequest(method, urlStr string, reqBody any, respTarget models.APIResponse, opts ...RequestOptions) error
}
