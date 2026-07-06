package client

import "time"

// RequestOption 函数选项模式：每个选项是一个 func(*RequestOptions)。
type RequestOption func(*RequestOptions)

// RequestOptions 请求配置。
type RequestOptions struct {
	NeedSign bool
	Headers  map[string]string
	Timeout  time.Duration
}

// WithSign 请求需要签名。
// 注意：NeedSign 默认为 true，此选项为显式声明用途保留，实际不需要调用。
func WithSign() RequestOption {
	return func(o *RequestOptions) {
		o.NeedSign = true
	}
}

// WithoutSign 请求不需要签名。
func WithoutSign() RequestOption {
	return func(o *RequestOptions) {
		o.NeedSign = false
	}
}

// WithTimeout 设置请求超时。
func WithTimeout(d time.Duration) RequestOption {
	return func(o *RequestOptions) {
		o.Timeout = d
	}
}

// WithHeader 添加单个自定义请求头。
func WithHeader(key, value string) RequestOption {
	return func(o *RequestOptions) {
		if o.Headers == nil {
			o.Headers = make(map[string]string)
		}
		o.Headers[key] = value
	}
}
