package api

import (
	"Skland/client"
	"Skland/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const AppCode = "4ca99fa6b56cc2ba"

type AuthAPI struct {
	client *client.HttpClient
}

func NewAuthAPI(c *client.HttpClient) *AuthAPI {
	return &AuthAPI{client: c}
}

func (a *AuthAPI) Login(account models.AccountInfo) (string, error) {
	reqBody := map[string]string{
		"phone":    account.Phone,
		"password": account.Password,
	}

	resp, err := a.client.DoRequest(
		http.MethodPost,
		"https://as.hypergryph.com/user/auth/v1/token_by_phone_password",
		reqBody,
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	var result models.LoginResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Status != 0 {
		return "", fmt.Errorf(result.Message)
	}

	return result.Data.Token, nil
}

func (a *AuthAPI) GetCredByToken(token string) (*models.CredResult, error) {
	grantCode, err := a.getGrantCode(token)
	if err != nil {
		return nil, err
	}

	reqBody := map[string]interface{}{
		"code": grantCode,
		"kind": 1,
	}

	resp, err := a.client.DoRequest(
		http.MethodPost,
		"https://zonai.skland.com/api/v1/user/auth/generate_cred_by_code",
		reqBody,
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	var credResult models.CredResult
	if err := json.NewDecoder(resp.Body).Decode(&credResult); err != nil {
		return nil, err
	}

	if credResult.Code != 0 {
		return nil, fmt.Errorf(credResult.Message)
	}

	return &credResult, nil
}

func (a *AuthAPI) getGrantCode(token string) (string, error) {
	reqBody := map[string]interface{}{
		"appCode": AppCode,
		"token":   token,
		"type":    0,
	}

	resp, err := a.client.DoRequest(
		http.MethodPost,
		"https://as.hypergryph.com/user/oauth2/v2/grant",
		reqBody,
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	var result models.GrantResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Status != 0 {
		return "", fmt.Errorf(result.Message)
	}
	return result.Data.Code, nil
}
