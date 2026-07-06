package client

import (
	"context"

	"github.com/ErikaCarnea/Skland/internal/models"
)

// Signer 签名生成器接口 — 用于单元测试 mock。
type Signer interface {
	GetSignHeaders(ctx context.Context, urlStr, method string, body []byte) (map[string]string, error)
}

// CredentialManager 凭证管理接口。
type CredentialManager interface {
	SetCred(cred string)
	SetSignToken(token string)
}

// HTTPClient HTTP 执行器接口。
type HTTPClient interface {
	ExecuteRequest(ctx context.Context, method, urlStr string, reqBody any, respTarget models.APIResponse, opts ...RequestOption) error
}

// SklandHTTPClient 组合接口，方便依赖注入。
// 调用方可以只依赖子接口（Signer / HTTPClient）以实现更精准的 mock。
type SklandHTTPClient interface {
	Signer
	CredentialManager
	HTTPClient
}
