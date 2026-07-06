package client

import (
	"context"
	"net/http"
)

// Signer 签名生成器接口。
type Signer interface {
	GetSignHeaders(ctx context.Context, urlStr, method string, body []byte) (map[string]string, error)
}

// CredentialManager 凭证管理接口。
type CredentialManager interface {
	SetCred(cred string)
	SetSignToken(token string)
}

// HTTPDoer 原始 HTTP 执行接口 — 单一职责，便于 mock。
type HTTPDoer interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// HeaderProvider 默认请求头提供者。
type HeaderProvider interface {
	DefaultHeaders() map[string]string
}

// HTTPClient 组合接口，供 ExecuteRequest 泛型函数使用。
type HTTPClient interface {
	Signer
	HTTPDoer
	HeaderProvider
}

// SklandHTTPClient 全功能客户端接口，组合签名、凭证、HTTP 能力。
// 供 API 层注入使用。
type SklandHTTPClient interface {
	Signer
	CredentialManager
	HTTPDoer
	HeaderProvider
}
