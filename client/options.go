package client

var (
	// SignedRequest 需要签名的请求选项
	SignedRequest = RequestOptions{
		NeedSign: true,
		Headers:  make(map[string]string),
	}

	// UnsignedRequest 不需要签名的请求选项
	UnsignedRequest = RequestOptions{
		NeedSign: false,
		Headers:  make(map[string]string),
	}
)
