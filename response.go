package httpc

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

// Response 包装 http.Response
type Response struct {
	*http.Response
}

// Bytes 以 byetes 的方式显示响应
func (r Response) Bytes() ([]byte, error) {
	defer r.Body.Close()

	if r.StatusCode < http.StatusOK && r.StatusCode > http.StatusIMUsed {
		return nil, errors.New("http status code unsuccessful response")
	}

	const maxSize = 1024 * 1024 * 20 // 20M
	body, err := readOnlySpecifiedSize(r.Body, maxSize)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// String 将响应以解析为字符串输出
func (r Response) String() string {
	body, _ := r.Bytes()

	return string(body)
}

// ToJSON 以 JSON 的方式解析响应 接受一个 *map 和 *struct
func (r Response) ToJSON(ptr interface{}) error {
	contentType := r.Header.Get(ContentType)

	// 不进行转化
	if !strings.Contains(contentType, "json") && !strings.Contains(contentType, "javascript") && !strings.Contains(contentType, "xml") { // json jsonp xml -> application/json application/javascript text/xml application/xml
		return errors.New("cannot be converted to json")
	}

	body, err := r.Bytes()
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

	purify := make([]byte, n) // 消除多余空数据
	copy(purify, cache)

	return purify, nil
}
