package client

import (
	"Skland/models"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
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
	mu        sync.Mutex
}

func NewClient() *HttpClient {
	return &HttpClient{
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

func (c *HttpClient) SetCred(cred string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.headers["cred"] = cred
}

func (c *HttpClient) SetSignToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
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
			return nil, fmt.Errorf("创建gzip读取器失败: %w", err)
		}
		defer func(gzReader *gzip.Reader) {
			if err := gzReader.Close(); err != nil {
				log.Printf("关闭gzip读取器失败: %v", err)
			}
		}(gzReader)
		reader = gzReader
	case "deflate":
		flateReader := flate.NewReader(resp.Body)
		defer func(flateReader io.ReadCloser) {
			if err := flateReader.Close(); err != nil {
				log.Printf("关闭deflate读取器失败: %v", err)
			}
		}(flateReader)
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

func (c *HttpClient) ExecuteRequest(
	method, urlStr string,
	reqBody any,
	respTarget models.APIResponse) error {
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
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("关闭响应体失败: %v", closeErr)
		}
	}()

	// 3. 读取并解析响应
	body, err := ReadResponseBody(resp)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	if err := json.Unmarshal(body, respTarget); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	// 4. 检查业务状态码
	if code := respTarget.GetCode(); code != 0 {
		return fmt.Errorf("API 返回错误: %s (code: %d)", respTarget.GetMessage(), code)
	}

	return nil
}
