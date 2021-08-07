package httpc

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// Response 包装 http.Response
type Response struct {
	*http.Response
}

// IsSuccessFul 响应成功
func (r Response) IsSuccessful() bool {
	return r.StatusCode > 200 && r.StatusCode < 300
}

// Bytes 以 byetes 的方式显示响应
func (r *Response) Bytes(maxSize ...int) ([]byte, error) {
	defer r.Body.Close()

	if r.StatusCode < http.StatusOK && r.StatusCode > http.StatusIMUsed {
		return nil, errors.New("http status code unsuccessful response")
	}

	size := 64 * 1024 // 本值越大越费内存
	if len(maxSize) != 0 && maxSize[0] > 0 {
		size = maxSize[0]
	}
	body, err := readOnlySpecifiedSize(r.Body, size)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Text 将响应以解析为字符串输出
func (r *Response) Text(maxSize ...int) string {
	body, _ := r.Bytes(maxSize...)

	return string(body)
}

// ToJSON 以 JSON 的方式解析响应 接受一个 *map 和 *struct
func (r *Response) ToJSON(ptr interface{}, maxSize ...int) error {
	body, err := r.Bytes(maxSize...)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, ptr)
}

// readOnlySpecifiedSize 仅允许读取的尺寸（防止读巨型文件时爆内存）
func readOnlySpecifiedSize(src io.Reader, maxSize int) ([]byte, error) {
	cache := make([]byte, maxSize+1)

	n, err := io.ReadFull(src, cache)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return nil, err
		}
	}

	if n > maxSize {
		return nil, errors.New("there more data")
	}

	return cache[:n], nil
}
