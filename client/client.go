package client

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"maps"

	"github.com/ErikaCarnea/Skland/models"
	"github.com/rs/zerolog/log"
)

const (
	DefaultTimeout = 30 * time.Second
	UserAgent      = "Skland/1.21.0 (com.hypergryph.skland; build:102100065; iOS 17.6.0; ) Alamofire/5.7.1" //Skland/1.50.0 (com.hypergryph.skland; build:105000018; Android 28; ) Okhttp/4.11.0
)

type httpClient struct {
	client    *http.Client
	headers   map[string]string
	signToken string
}

type signHeader struct {
	Platform  string `json:"platform"`
	Timestamp string `json:"timestamp"`
	DId       string `json:"dId"`
	VName     string `json:"vName"`
}

func NewClient() SklandHTTPClient {
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
		},
	}
}

func (c *httpClient) SetCred(cred string) {
	c.headers["cred"] = cred
}

func (c *httpClient) SetSignToken(token string) {
	c.signToken = token
}

// func (c *httpClient) doRequest(method, url string, body any, headers map[string]string, timeout time.Duration) (*http.Response, error) {
// 	var reqBody io.Reader
// 	if body != nil {
// 		jsonBody, err := json.Marshal(body)
// 		if err != nil {
// 			return nil, err
// 		}
// 		reqBody = bytes.NewBuffer(jsonBody)
// 	}

// 	req, err := http.NewRequest(method, url, reqBody)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 设置headers
// 	for k, v := range headers {
// 		req.Header.Set(k, v)
// 	}

// 	// 使用自定义超时（如果提供）
// 	client := c.client
// 	if timeout > 0 {
// 		clientCopy := *c.client
// 		clientCopy.Timeout = timeout
// 		client = &clientCopy
// 	}

// 	return client.Do(req)
// }

func (c *httpClient) GetSignHeaders(urlStr, method string, bodyBytes []byte) (map[string]string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var bodyOrQuery string
	if method == http.MethodGet {
		bodyOrQuery = u.RawQuery
	} else if len(bodyBytes) > 0 {
		bodyOrQuery = string(bodyBytes)
	}

	sign, headerCa, err := c.generateSignature(u.Path, bodyOrQuery)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"sign": sign,
	}
	maps.Copy(headers, headerCa)

	return headers, nil
}

func (c *httpClient) ExecuteRequest(method, urlStr string, reqBody any, respTarget models.APIResponse, opts ...RequestOptions) error {
	options := RequestOptions{
		NeedSign: true,
		Headers:  make(map[string]string),
	}
	if len(opts) > 0 {
		options = opts[0]
	}

	headers := make(map[string]string)
	maps.Copy(headers, c.headers)
	maps.Copy(headers, options.Headers)

	// 预序列化请求体（如果存在）
	var bodyBytes []byte
	var err error
	if reqBody != nil {
		bodyBytes, err = json.Marshal(reqBody)
		if err != nil {
			return fmt.Errorf("序列化请求体失败: %w", err)
		}
		headers["Content-Type"] = "application/json"
	}

	// 生成签名 (如果需要)
	if options.NeedSign {
		signHeaders, err := c.GetSignHeaders(urlStr, method, bodyBytes)
		if err != nil {
			return fmt.Errorf("生成签名头失败: %w", err)
		}
		maps.Copy(headers, signHeaders)
	}

	// 构造请求
	var reqBodyReader io.Reader
	if len(bodyBytes) > 0 {
		reqBodyReader = bytes.NewReader(bodyBytes)
	}
	req, err := http.NewRequest(method, urlStr, reqBodyReader)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := c.client
	if options.Timeout > 0 {
		ctx, cancel := context.WithTimeout(req.Context(), options.Timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求发送失败: %w", err)
	}
	defer CloseResponse(resp) // 统一关闭

	return c.parseResponse(resp, respTarget)
}

func (c *httpClient) parseResponse(resp *http.Response, target models.APIResponse) error {
	body, err := ReadResponseBody(resp)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	// if code := target.GetCode(); code != 0 {
	// 	return fmt.Errorf("API 返回错误: %s (code: %d)", target.GetMessage(), code)
	// }

	return nil
}

func CloseResponse(resp *http.Response) {
	if resp != nil && resp.Body != nil {
		if err := resp.Body.Close(); err != nil {
			log.Warn().Err(err).Msg("关闭响应体失败")
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
