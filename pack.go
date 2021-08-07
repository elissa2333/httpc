package httpc

import (
	"context"
	"net/http"
	"time"
)

// DefaultClient 默认客户端
var DefaultClient = New()

// UseDNS 使用指定 dns 解析域名
func UseDNS(dns ...string) *Client {
	return DefaultClient.UseDNS(dns...)
}

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

// CallWithContext 使用指定 http 方法访问 url
func CallWithContext(ctx context.Context, method string, url string) (*Response, error) {
	return DefaultClient.CallWithContext(ctx, method, url)
}

// Options 获取目的资源所支持的通信选项
func Options(url string) (*Response, error) {
	return DefaultClient.Options(url)
}

// OptionsWithContext 获取目的资源所支持的通信选项（可取消）
func OptionsWithContext(ctx context.Context, url string) (*Response, error) {
	return DefaultClient.OptionsWithContext(ctx, url)
}

// Head 获取响应头
func Head(url string) (*Response, error) {
	return DefaultClient.Head(url)
}

//  HeadWithContext 获取响应头（可取消）
func HeadWithContext(ctx context.Context, url string) (*Response, error) {
	return DefaultClient.HeadWithContext(ctx, url)
}

// Get 发送 GET 请求
func Get(url string) (*Response, error) {
	return DefaultClient.Get(url)
}

// GetWithContext 发送 GET 请求（可取消）
func GetWithContext(ctx context.Context, url string) (*Response, error) {
	return DefaultClient.GetWithContext(ctx, url)
}

// Post 发送 POST 请求
func Post(url string) (*Response, error) {
	return DefaultClient.Post(url)
}

// PostWithContext 发送 POST 请求（可取消）
func PostWithContext(ctx context.Context, url string) (*Response, error) {
	return DefaultClient.PostWithContext(ctx, url)
}

// Put 发送 PUT 请求
func Put(url string) (*Response, error) {
	return DefaultClient.Put(url)
}

// PutWithContext 发送 PUT 请求（可取消）
func PutWithContext(ctx context.Context, url string) (*Response, error) {
	return DefaultClient.PutWithContext(ctx, url)
}

// Patch 发送 PATCH 请求）
func Patch(url string) (*Response, error) {
	return DefaultClient.Patch(url)
}

// PatchWithContext 发送 PATCH 请求（可取消）
func PatchWithContext(ctx context.Context, url string) (*Response, error) {
	return DefaultClient.PatchWithContext(ctx, url)
}

// Delete 发送 DELETE 请求
func Delete(url string) (*Response, error) {
	return DefaultClient.Delete(url)
}

// DeleteWithContext 发送 DELETE 请求（可取消）
func DeleteWithContext(ctx context.Context, url string) (*Response, error) {
	return DefaultClient.DeleteWithContext(ctx, url)
}
