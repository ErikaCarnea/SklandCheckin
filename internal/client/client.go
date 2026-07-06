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

	"github.com/ErikaCarnea/Skland/internal/models"
	"github.com/rs/zerolog/log"
)

const (
	DefaultTimeout          = 30 * time.Second
	SignTimestampOffset     = -2 // 签名时间戳偏移量（秒）
	UserAgent               = "Skland/1.21.0 (com.hypergryph.skland; build:102100065; iOS 17.6.0; ) Alamofire/5.7.1"
)

type httpClient struct {
	client    *http.Client
	headers   map[string]string
	signToken string
}

// signHeader 签名请求头结构。
// 所有字段固定为空字符串，仅 timestamp 每次计算。
type signHeader struct {
	Platform  string `json:"platform"`
	Timestamp string `json:"timestamp"`
	DID       string `json:"dId"`
	VName     string `json:"vName"`
}

func NewClient() SklandHTTPClient {
	return &httpClient{
		client: &http.Client{
			Timeout: DefaultTimeout,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
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

func (c *httpClient) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *httpClient) DefaultHeaders() map[string]string {
	return c.headers
}

func (c *httpClient) GetSignHeaders(ctx context.Context, urlStr, method string, bodyBytes []byte) (map[string]string, error) {
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

	headers := map[string]string{"sign": sign}
	maps.Copy(headers, headerCa)

	return headers, nil
}

// ExecuteRequest 泛型 HTTP 请求执行器。
// T 必须实现 models.APIResponse（GetCode / GetMessage）。
// 返回值直接给出精确类型，无需调用方先声明零值再传递指针。
func ExecuteRequest[T models.APIResponse](ctx context.Context, c HTTPClient, method, urlStr string, reqBody any, opts ...RequestOption) (*T, error) {
	options := RequestOptions{
		NeedSign: true,
		Headers:  make(map[string]string),
	}
	for _, opt := range opts {
		opt(&options)
	}

	headers := make(map[string]string)
	maps.Copy(headers, c.DefaultHeaders())
	maps.Copy(headers, options.Headers)

	var bodyBytes []byte
	var err error
	if reqBody != nil {
		bodyBytes, err = json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("序列化请求体失败: %w", err)
		}
		headers["Content-Type"] = "application/json"
	}

	if options.NeedSign {
		signHeaders, err := c.GetSignHeaders(ctx, urlStr, method, bodyBytes)
		if err != nil {
			return nil, fmt.Errorf("生成签名头失败: %w", err)
		}
		maps.Copy(headers, signHeaders)
	}

	var reqBodyReader io.Reader
	if len(bodyBytes) > 0 {
		reqBodyReader = bytes.NewReader(bodyBytes)
	}
	req, err := http.NewRequestWithContext(ctx, method, urlStr, reqBodyReader)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if options.Timeout > 0 {
		timeoutCtx, cancel := context.WithTimeout(req.Context(), options.Timeout)
		defer cancel()
		req = req.WithContext(timeoutCtx)
	}

	resp, err := c.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("请求发送失败: %w", err)
	}
	defer CloseResponse(resp)

	body, err := ReadResponseBody(resp)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
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
	timestamp := fmt.Sprintf("%d", time.Now().Unix()+SignTimestampOffset)
	header := signHeader{
		Timestamp: timestamp,
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", nil, err
	}
	sStr := fmt.Sprintf("%s%s%s%s", path, bodyOrQuery, timestamp, string(headerJSON))

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
