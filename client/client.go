package client

import (
	"Skland/models"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	DefaultTimeout = 30 * time.Second
	UserAgent      = "Skland/1.0.1 (com.hypergryph.skland; build:100001014; Android 31; ) Okhttp/4.11.0"
)

type httpClient struct {
	client    *http.Client
	headers   map[string]string
	signToken string
	mu        sync.Mutex
}

type signHeader struct {
	Platform  string `json:"platform"`
	Timestamp string `json:"timestamp"`
	DId       string `json:"dId"`
	VName     string `json:"vName"`
}

func NewClient() SklandHttpClient {
	return &httpClient{
		client: &http.Client{
			Timeout: DefaultTimeout,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 100,              // 每个主机最大空闲连接数
				IdleConnTimeout:     90 * time.Second, // 空闲连接超时时间
				DisableCompression:  false,
			},
		},
		headers: map[string]string{
			"User-Agent":      UserAgent,
			"Accept-Encoding": "gzip",
			"Connection":      "keep-alive",
		},
	}
}

func (c *httpClient) SetCred(cred string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.headers["cred"] = cred
}

func (c *httpClient) SetSignToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.signToken = token
}

func (c *httpClient) DoRequest(method, url string, body any, headers map[string]string) (*http.Response, error) {
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

func (c *httpClient) GetSignHeaders(urlStr, method string, body any) (map[string]string, error) {
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

func (c *httpClient) ExecuteRequest(method, urlStr string, reqBody any, respTarget models.ApiResponse) error {
	// 1. 生成签名头
	headers, err := c.GetSignHeaders(urlStr, method, reqBody)
	if err != nil {
		return fmt.Errorf("生成签名头失败: %w", err)
	}
	headers["Content-Type"] = "application/json"

	// 2. 发送请求
	resp, err := c.DoRequest(method, urlStr, reqBody, headers)
	if err != nil {
		return fmt.Errorf("请求发送失败: %w", err)
	}
	defer CloseResponse(resp) // 统一关闭

	return c.parseResponse(resp, respTarget)
}

func (c *httpClient) parseResponse(resp *http.Response, target models.ApiResponse) error {
	body, err := ReadResponseBody(resp)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if code := target.GetCode(); code != 0 {
		return fmt.Errorf("API 返回错误: %s (code: %d)", target.GetMessage(), code)
	}

	return nil
}

func CloseResponse(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("关闭响应体失败")
		}
	}
}

func ReadResponseBody(resp *http.Response) ([]byte, error) {
	var reader io.Reader
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		gzReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("创建gzip读取器失败: %w", err)
		}
		defer gzReader.Close()
		reader = gzReader
	case "deflate":
		flateReader := flate.NewReader(resp.Body)
		defer flateReader.Close()
		reader = flateReader
	default:
		reader = resp.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %w", err)
	}
	return body, nil
}

func (c *httpClient) generateSignature(path, bodyOrQuery string) (string, map[string]string, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix()-2)
	header := signHeader{
		Timestamp: timestamp,
	}

	headerJson, err := json.Marshal(header)
	if err != nil {
		return "", nil, err
	}
	sStr := fmt.Sprintf("%s%s%s%s", path, bodyOrQuery, timestamp, string(headerJson))

	h := hmac.New(sha256.New, []byte(c.signToken))
	h.Write([]byte(sStr))
	sha := hex.EncodeToString(h.Sum(nil))

	md5Hash := md5.Sum([]byte(sha))
	sign := hex.EncodeToString(md5Hash[:])

	headerCa := map[string]string{
		"timestamp": header.Timestamp,
	}

	return sign, headerCa, nil
}
