package httpc

import (
	"net/http"
	"time"
)

// DefaultClient 默认客户端
var DefaultClient = New()

// SetProxy 设置代理（支持 http 以及 socks5 代理）
func SetProxy(proxy string) *Client {
	return DefaultClient.SetProxy(proxy)
}

// SetTimeout 设置超时时间
func SetTimeout(duration time.Duration) *Client {
	return DefaultClient.SetTimeout(duration)
}

// SetBaseURL 设置基础 URL 后续访问将默认拼接本URL
func SetBaseURL(url string) *Client {
	return DefaultClient.SetBaseURL(url)
}

// AddURLQuery 添加 URL 查询参数
func AddURLQuery(key string, values ...string) *Client {
	return DefaultClient.AddURLQuery(key, values...)
}

// AddURLQueryS 以 string 的方式添加查询参数
func AddURLQueryS(query string) *Client {
	return DefaultClient.AddURLQueryS(query)
}

// SetURLQuery 设置 URL 查询参数
func SetURLQuery(key string, values ...string) *Client {
	return DefaultClient.SetURLQuery(key, values...)
}

// SetURLQueryS 以 string 的方式设置查询参数
func SetURLQueryS(query string) *Client {
	return DefaultClient.SetURLQueryS(query)
}

// SetHeader 设置 header
func SetHeader(key, value string) *Client {
	return DefaultClient.SetHeader(key, value)
}

// SetHeaders 设置 headers
func SetHeaders(headers map[string]string) *Client {
	return DefaultClient.SetHeaders(headers)
}

// SetUserAgent 设置浏览器标识
func SetUserAgent(ua string) *Client {
	return DefaultClient.SetUserAgent(ua)
}

// SetContentType 设置内容类型
func SetContentType(contentType string) *Client {
	return DefaultClient.SetContentType(contentType)
}

// SetBody 设置 body 内容
func SetBody(in interface{}) *Client {
	return DefaultClient.SetBody(in)
}

// SetFromData 设置表单数据
func SetFromData(rows ...FromDataRow) *Client {
	return DefaultClient.SetFromData(rows...)
}

// Do 发送自定义请求
func Do(req *http.Request) (*Response, error) {
	return DefaultClient.Do(req)
}

// Call 使用指定 http 方法访问 url
func Call(method string, url string) (*Response, error) {
	return DefaultClient.Call(method, url)
}

// Options 获取目的资源所支持的通信选项
func Options(url string) (*Response, error) {
	return DefaultClient.Options(url)
}

// Head 获取响应头
func Head(url string) (*Response, error) {
	return DefaultClient.Head(url)
}

// Get 发送 GET 请求（获取数据）
func Get(url string) (*Response, error) {
	return DefaultClient.Get(url)
}

// Post 发送 POST 请求（添加数据）
func Post(url string) (*Response, error) {
	return DefaultClient.Post(url)
}

// Put 发送 PUT 请求（覆盖更新）
func Put(url string) (*Response, error) {
	return DefaultClient.Put(url)
}

// Patch 发送 PATCH 请求（部分更新）
func Patch(url string) (*Response, error) {
	return DefaultClient.Patch(url)
}

// Delete 发送 DELETE 请求（删除资源）
func Delete(url string) (*Response, error) {
	return DefaultClient.Delete(url)
}
