package httpc

import (
	"io"
	"net/http"
	"os"
)

// NewRequest 新建请求
func NewRequest(method string, url string, headers map[string]string, body io.Reader) (*http.Request, error) { // 带认证的 考虑一下怎么处理
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	userAgent := os.Getenv("UserAgent")
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	for k, v := range headers {
		if v == "" {
			continue
		}
		req.Header.Set(k, v)
	}

	return req, err
}
