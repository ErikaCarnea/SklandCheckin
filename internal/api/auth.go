package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ErikaCarnea/Skland/internal/client"
	"github.com/ErikaCarnea/Skland/internal/config"
	"github.com/ErikaCarnea/Skland/internal/models"
	"github.com/rs/zerolog/log"
)

const (
	AppCode               = "4ca99fa6b56cc2ba"
	PhonePasswordURL      = "https://as.hypergryph.com/user/auth/v1/token_by_phone_password"
	SendPhoneCodeURL      = "https://as.hypergryph.com/general/v1/send_phone_code"
	TokenByPhoneCodeURL   = "https://as.hypergryph.com/user/auth/v2/token_by_phone_code"
	GenerateCredURL       = "https://zonai.skland.com/api/v1/user/auth/generate_cred_by_code"
	GrantURL              = "https://as.hypergryph.com/user/oauth2/v2/grant"
)

type AuthAPI struct {
	client client.SklandHTTPClient
}

func NewAuthAPI(client client.SklandHTTPClient) *AuthAPI {
	return &AuthAPI{client: client}
}

// LoginByPassword 手机号+密码登录。
func (a *AuthAPI) LoginByPassword(ctx context.Context, phone, password string) (*models.CredResult, error) {
	log.Info().Str("action", "password_login").Msg("开始密码登录流程")

	reqBody := map[string]string{
		"phone":    phone,
		"password": password,
	}
	loginResult, err := client.ExecuteRequest[models.LoginResult](ctx, a.client, http.MethodPost, PhonePasswordURL, reqBody, client.WithoutSign())
	if err != nil {
		return nil, err
	}

	if _, err := config.SaveToken(loginResult.Data.Token); err != nil {
		log.Error().Err(err).Msg("保存Token失败")
	}
	return a.GetCredByToken(ctx, loginResult.Data.Token)
}

// SendPhoneCode 发送手机验证码。
func (a *AuthAPI) SendPhoneCode(ctx context.Context, phone string) error {
	sendCodeResult, err := client.ExecuteRequest[models.SendCodeResult](ctx, a.client, http.MethodPost, SendPhoneCodeURL, map[string]any{"phone": phone, "type": 2}, client.WithoutSign())
	if err != nil {
		return fmt.Errorf("发送验证码失败: %w", err)
	}

	if sendCodeResult.GetCode() != 0 {
		return fmt.Errorf("API错误: %s", sendCodeResult.GetMessage())
	}
	return nil
}

// LoginByPhoneCode 验证码登录。
func (a *AuthAPI) LoginByPhoneCode(ctx context.Context, phone, code string) (*models.CredResult, error) {
	log.Info().Str("action", "phone_code_login").Msg("开始手机验证码登录流程")

	loginResult, err := client.ExecuteRequest[models.LoginResult](ctx, a.client, http.MethodPost, TokenByPhoneCodeURL, map[string]string{"phone": phone, "code": code}, client.WithoutSign())
	if err != nil {
		return nil, err
	}

	if _, err := config.SaveToken(loginResult.Data.Token); err != nil {
		log.Error().Err(err).Msg("保存Token失败")
	}

	return a.GetCredByToken(ctx, loginResult.Data.Token)
}

// LoginByCode 授权码登录。
func (a *AuthAPI) LoginByCode(ctx context.Context, code string) (*models.CredResult, error) {
	log.Info().Str("action", "token_login").Msg("开始授权码登录流程")

	var response struct {
		Data struct {
			Content string `json:"content"`
		} `json:"data"`
	}

	if err := json.Unmarshal([]byte(code), &response); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	if response.Data.Content == "" {
		return nil, fmt.Errorf("无效的数据结构，缺少content字段")
	}

	if _, err := config.SaveToken(response.Data.Content); err != nil {
		return nil, fmt.Errorf("保存Token失败: %w", err)
	}

	return a.GetCredByToken(ctx, response.Data.Content)
}

func (a *AuthAPI) GetCredByToken(ctx context.Context, token string) (*models.CredResult, error) {
	grantCode, err := a.getGrantCode(ctx, token)
	if err != nil {
		return nil, err
	}

	reqBody := map[string]any{
		"code": grantCode,
		"kind": 1,
	}

	credResult, err := client.ExecuteRequest[models.CredResult](ctx, a.client, http.MethodPost, GenerateCredURL, reqBody, client.WithoutSign())
	if err != nil {
		return nil, err
	}

	if credResult.GetCode() != 0 {
		return nil, fmt.Errorf("%s", credResult.GetMessage())
	}
	return credResult, nil
}

func (a *AuthAPI) getGrantCode(ctx context.Context, token string) (string, error) {
	result, err := client.ExecuteRequest[models.GrantResult](ctx, a.client, http.MethodPost, GrantURL, map[string]any{"appCode": AppCode, "token": token, "type": 0}, client.WithoutSign())
	if err != nil {
		return "", err
	}
	return result.Data.Code, nil
}
