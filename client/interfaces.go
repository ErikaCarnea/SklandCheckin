package client

import (
	"time"

	"github.com/ErikaCarnea/Skland/models"
)

type RequestOptions struct {
	NeedSign bool
	Headers  map[string]string
	Timeout  time.Duration
}

type SklandHTTPClient interface {
	GetSignHeaders(urlStr, method string, body any) (map[string]string, error)
	SetCred(cred string)
	SetSignToken(token string)
	ExecuteRequest(method, urlStr string, reqBody any, respTarget models.APIResponse, opts ...RequestOptions) error
}
