package client

import (
	"Skland/models"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultTimeout = 30 * time.Second
	UserAgent      = "Skland/1.0.1 (com.hypergryph.skland; build:100001014; Android 31; ) Okhttp/4.11.0"
)

type HttpClient struct {
	client    *http.Client
	headers   map[string]string
	signToken string
}

func NewClient() *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: DefaultTimeout,
			Transport: &http.Transport{
				DisableCompression: false,
			},
		},
		headers: map[string]string{
			"User-Agent":      UserAgent,
			"Accept-Encoding": "gzip",
			"Connection":      "close",
		},
	}
}

func (c *HttpClient) SetCred(cred string) {
	c.headers["cred"] = cred
}

func (c *HttpClient) SetSignToken(token string) {
	c.signToken = token
}

func (c *HttpClient) DoRequest(method, url string, body any, headers map[string]string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	// 设置headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return c.client.Do(req)
}

func ReadResponseBody(resp *http.Response) ([]byte, error) {
	var reader io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer gzReader.Close()
		reader = gzReader
	case "deflate":
		reader = flate.NewReader(resp.Body)
		defer reader.(io.ReadCloser).Close()
	default:
		reader = resp.Body
	}
	return io.ReadAll(reader)
}

func (c *HttpClient) GetSignHeaders(urlStr, method string, body any) (map[string]string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var bodyOrQuery string
	if method == http.MethodGet {
		bodyOrQuery = u.RawQuery
	} else if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyOrQuery = string(jsonBody)
	}

	sign, headerCa, err := c.generateSignature(u.Path, bodyOrQuery)
	if err != nil {
		return nil, err
	}
	headers := map[string]string{
		"sign": sign,
	}
	for k, v := range headerCa {
		headers[k] = v
	}

	return headers, nil
}

func (c *HttpClient) PrintAllPlayersInfo(bindings []models.Binding) error {
	for _, binding := range bindings {
		urlStr := fmt.Sprintf("https://zonai.skland.com/api/v1/game/player/info?uid=%s", binding.Uid)

		headers, err := c.GetSignHeaders(urlStr, http.MethodGet, nil)
		if err != nil {
			return fmt.Errorf("获取签名头失败(UID:%s): %w", binding.Uid, err)
		}
		headers["Content-Type"] = "application/json"

		resp, err := c.DoRequest(http.MethodGet, urlStr, nil, headers)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		body, err := ReadResponseBody(resp)
		if err != nil {
			return err
		}

		content := string(body)

		fmt.Printf("=== 玩家 %s (%s) ===\n", binding.NickName, binding.Uid)
		fmt.Println(content)
		fmt.Println("========================")
	}
	return nil
}
