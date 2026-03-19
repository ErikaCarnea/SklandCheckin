package client

import "time"

type RequestOptions struct {
	NeedSign bool
	Headers  map[string]string
	Timeout  time.Duration
}

// SignedRequest 返回一个需要签名的请求选项副本
func SignedRequest() RequestOptions {
	return RequestOptions{
		NeedSign: true,
		Headers:  make(map[string]string),
	}
}

// UnsignedRequest 返回一个不需要签名的请求选项副本
func UnsignedRequest() RequestOptions {
	return RequestOptions{
		NeedSign: false,
		Headers:  make(map[string]string),
	}
}
