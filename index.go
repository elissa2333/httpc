package httpc

// 封装 http 客户端 理论上本应用的所有 http请求因使用本包提供的 http 客户端进行发起

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	stdURL "net/url"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

// LinkError 错误链
type LinkError struct {
	Value error

	Next *LinkError
}

// NewLinkError 新建错误链
func NewLinkError(err error) *LinkError {
	if err == nil {
		panic("input error value is nil")
	}
	return &LinkError{
		Value: err,
	}
}

// Shift 在头部添加数据
func (e *LinkError) Shift(err error) *LinkError {
	if err != nil {
		return &LinkError{Value: err, Next: e}
	}

	return e
}

// Push 在尾部添加数据
func (e *LinkError) Push(err error) *LinkError {
	if err != nil {
		for node := e; node != nil; node = node.Next {
			if node.Next == nil {
				node.Next = &LinkError{Value: err}
				break
			}
		}
	}

	return e
}

func (e *LinkError) Error() string {
	var out []string
	for node := e; node != nil; node = node.Next {
		out = append(out, node.Value.Error())
	}

	return strings.Join(out, " -> ")
}

func (e *LinkError) Unwrap() error {
	return e.Next
}

// Client 客户端
type Client struct {
	BaseURL  string        // 基础 URL
	URLQuery stdURL.Values // URL 查询参数

	Headers map[string]string // 头部信息
	Body    io.Reader         // 内容
	Client  http.Client       // 客户端

	Error *LinkError // 错误
}

// New 新建客户端
func New() *Client {
	return UseClient(http.Client{})
}

// UseClient 使用指定客户端发送请求
func UseClient(httpClient http.Client) *Client {
	return &Client{
		Client: httpClient, // 最后使用的时候进行了空检测
	}
}

// UseDNS 使用指定 dns 解析域名
func (c Client) UseDNS(dns ...string) *Client {
	tr := &http.Transport{}
	if c.Client.Transport != nil {
		current, ok := c.Client.Transport.(*http.Transport)
		if ok {
			tr = current
		}
	}

	dialer := net.Dialer{}
	tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		addrs, err := dnsResolver(host, dns...)
		if err != nil {
			return nil, err
		}
		if len(addrs) != 0 {
			addr = net.JoinHostPort(addrs[0], port)
		}
		return dialer.DialContext(ctx, network, addr)
	}

	c.Client.Transport = tr

	return &c
}

var systemHosts map[string]string

func useHostsPsrseIP(host string) string {
	if systemHosts == nil {
		hostsPATH := "/etc/hosts"
		switch runtime.GOOS {
		case "windows":
			hostsPATH = "/windows/system32/drivers/etc/hosts"
		}

		b, err := os.ReadFile(hostsPATH)
		if err != nil {
			return ""
		}

		systemHosts = map[string]string{}

		for _, v := range strings.Split(string(b), "\n") {
			if v != "" && !strings.ContainsRune(v, '#') {
				sp := strings.Fields(v)
				if len(sp) < 2 {
					continue
				}

				for _, v := range sp[1:] {
					systemHosts[v] = sp[0]
				}
			}
		}
	}

	return systemHosts[host]
}

// SetProxy 设置代理
func (c Client) SetProxy(proxy string) *Client {
	tr := &http.Transport{}

	if c.Client.Transport != nil {
		current, ok := c.Client.Transport.(*http.Transport)
		if ok {
			tr = current
		}
	}

	p, err := stdURL.Parse(proxy)
	if err != nil {
		return c.handleError(err)
	}

	tr.Proxy = http.ProxyURL(p)

	c.Client.Transport = tr

	return &c
}

// DeleteProxy 删除代理
func (c Client) DeleteProxy() *Client {
	tr, ok := c.Client.Transport.(*http.Transport)
	if !ok {
		c.handleError(errors.New("client Transport interface type is not a *http.Transport"))
	}
	if tr != nil {
		tr.Proxy = nil
	}

	return &c
}

// SetTimeout 设置超时时间
func (c Client) SetTimeout(duration time.Duration) *Client {
	c.Client.Timeout = duration
	return &c
}

// SetBaseURL 设置基础 URL 后续访问将默认拼接本URL
func (c Client) SetBaseURL(url string) *Client {
	c.BaseURL = url
	return &c
}

// DeleteBaseURL 删除基础 URL （如果有设置的话）
func (c Client) DeleteBaseURL() *Client {
	c.BaseURL = ""
	return &c
}

// AddURLQuery 添加 URL 查询参数
func (c Client) AddURLQuery(key string, values ...string) *Client {
	v, ok := c.URLQuery[key]
	if ok {
		c.URLQuery[key] = append(v, values...)
		return &c
	}

	return c.SetURLQuery(key, values...)
}

// AddURLQueryS 以 string 方式添加 URL 查询参数
// name=niconiconi
func (c Client) AddURLQueryS(query string) *Client {
	if c.URLQuery != nil {
		if query != "" {
			query = c.URLQuery.Encode() + "&" + query
		}
	}

	return c.SetURLQueryS(query)
}

// SetURLQuery 设置 URL 查询参数
func (c Client) SetURLQuery(key string, values ...string) *Client {
	if c.URLQuery == nil {
		c.URLQuery = stdURL.Values{key: values}
	} else {
		c.URLQuery[key] = values
	}

	return &c
}

// SetURLQueryS 以 string 的方式设置查询参数
func (c Client) SetURLQueryS(query string) *Client {
	values, err := stdURL.ParseQuery(query)
	if err != nil {
		return c.handleError(err)
	}
	c.URLQuery = values

	return &c
}

// DeleteURLQuery 删除 URL 查询参数
func (c Client) DeleteURLQuery() *Client {
	c.URLQuery = nil
	return &c
}

// SetHeader 设置 header 字段
func (c Client) SetHeader(key, value string) *Client {
	if c.Headers == nil {
		c.Headers = map[string]string{}
	}

	c.Headers[key] = value
	return &c
}

// SetHeaders 设置 headers
func (c Client) SetHeaders(headers map[string]string) *Client {
	this := &c
	for k, v := range headers {
		this = this.SetHeader(k, v)
	}

	return this
}

// DeleteHeaders 删除指定 headers
func (c Client) DeleteHeaders(keys ...string) *Client {
	for _, v := range keys {
		delete(c.Headers, v)
	}

	return &c
}

// SetUserAgent 设置浏览器标识
func (c Client) SetUserAgent(UA string) *Client {
	return c.SetHeader("User-Agent", UA)
}

// SetContentType 设置内容类型
func (c Client) SetContentType(contentType string) *Client {
	return c.SetHeader(HeaderContentType, contentType)
}

func (c Client) setBody(body io.Reader) *Client {
	c.Body = body
	return &c
}

// SetBody 设置内容
func (c Client) SetBody(body interface{}) *Client {
	switch value := body.(type) {
	case io.Reader:
		return c.setBody(value)
	case []byte:
		return c.setBody(bytes.NewReader(value))
	case string:
		return c.setBody(strings.NewReader(value))
	case nil:
		return &c
	}

	refType := reflect.TypeOf(body)
	if refType.Kind() == reflect.Ptr {
		refType = refType.Elem()
	}
	switch refType.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		b, err := json.Marshal(body)
		if err != nil {
			return c.handleError(err)
		}

		this := &c
		if c.Headers[HeaderContentType] == "" {
			this = this.SetContentType(MIMEApplicationJSON)
		}
		return this.setBody(bytes.NewReader(b))
	}

	return c.handleError(errors.New("no content to set"))
}

// FromDataRow formData 单条数据
type FromDataRow struct {
	Key   string    // 键
	Value string    // 值 （如果上传的是文件则这里是文件名）
	Data  io.Reader // 数据 （如果上传的是文件则这里是文件不然则为空）
}

// SetFromData 设置表单数据
func (c Client) SetFromData(rows ...FromDataRow) *Client {
	if rows == nil {
		return &c
	}

	buf := &bytes.Buffer{}
	mp := multipart.NewWriter(buf)
	for _, v := range rows {
		if v.Data != nil {
			fw, err := mp.CreateFormFile(v.Key, v.Value)
			if err != nil {
				return c.handleError(err)
			}
			if _, err := io.Copy(fw, v.Data); err != nil {
				return c.handleError(err)
			}
		} else {
			fw, err := mp.CreateFormField(v.Key)
			if err != nil {
				return c.handleError(err)
			}
			if _, err := fw.Write([]byte(v.Value)); err != nil {
				return c.handleError(err)
			}
		}
	}

	// 终止写入不然对方会报部分上传
	if err := mp.Close(); err != nil {
		return c.handleError(err)
	}

	return c.SetContentType(mp.FormDataContentType()).SetBody(buf)
}

// Do 发送自定义请求（无法共享链式调用的数据）
func (c Client) Do(req *http.Request) (*Response, error) {
	if c.Error != nil {
		return nil, c.Error
	}

	resp, err := c.Client.Do(req)
	return &Response{Response: resp}, err
}

// CallWithContext 使用指定 http 方法访问 url
func (c Client) CallWithContext(ctx context.Context, method string, url string) (*Response, error) {
	if c.BaseURL != "" {
		url = c.BaseURL + url
	}

	if c.URLQuery != nil {
		parseURL, err := stdURL.Parse(url)
		if err != nil {
			return nil, err
		}
		withQuery := parseURL.Query()
		for k, v := range c.URLQuery {
			for _, j := range v {
				withQuery.Add(k, j)
			}
		}
		parseURL.RawQuery = withQuery.Encode()

		url = parseURL.String()
	}

	req, err := NewRequestWithContext(ctx, method, url, c.Headers, c.Body)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Options 获取目的资源所支持的通信选项
func (c Client) Options(url string) (*Response, error) {
	return c.OptionsWithContext(context.Background(), url)
}

// OptionsWithContext 获取目的资源所支持的通信选项（可取消）
func (c Client) OptionsWithContext(ctx context.Context, url string) (*Response, error) {
	return c.CallWithContext(ctx, http.MethodOptions, url)
}

// Head 获取响应头
func (c Client) Head(url string) (*Response, error) {
	return c.HeadWithContext(context.Background(), url)
}

// HeadWithContext 获取响应头（可取消）
func (c Client) HeadWithContext(ctx context.Context, url string) (*Response, error) {
	return c.CallWithContext(ctx, http.MethodHead, url)
}

// Get 发送 GET 请求
func (c Client) Get(url string) (*Response, error) {
	return c.GetWithContext(context.Background(), url)
}

// GetWithContext 发送 GET 请求（可取消）
func (c Client) GetWithContext(ctx context.Context, url string) (*Response, error) {
	return c.CallWithContext(ctx, http.MethodGet, url)
}

// Post 发送 POST 请求
func (c Client) Post(url string) (*Response, error) {
	return c.PostWithContext(context.Background(), url)
}

// PostWithContext 发送 POST 请求（可取消）
func (c Client) PostWithContext(ctx context.Context, url string) (*Response, error) {
	return c.CallWithContext(ctx, http.MethodPost, url)
}

//// PostForm http.PostForm
//func (c Client) PostForm(url string, values url.Values) (*Response, error) {
//	return c.SetContentType(MIMEXWWWFormURLEncoded).SetBody(values.Encode()).Post(url)
//}

// Put 发送 PUT 请求
func (c Client) Put(url string) (*Response, error) {
	return c.PutWithContext(context.Background(), url)
}

// CallWithContext 发送 PUT 请求（可取消）
func (c Client) PutWithContext(ctx context.Context, url string) (*Response, error) {
	return c.CallWithContext(ctx, http.MethodPut, url)
}

// Patch 发送 PATCH 请求
func (c Client) Patch(url string) (*Response, error) {
	return c.PatchWithContext(context.Background(), url)
}

// PatchWithContext 发送 PATCH 请求（可取消）
func (c Client) PatchWithContext(ctx context.Context, url string) (*Response, error) {
	return c.CallWithContext(ctx, http.MethodPatch, url)
}

// Delete 发送 DELETE 请求
func (c Client) Delete(url string) (*Response, error) {
	return c.DeleteWithContext(context.Background(), url)
}

// DeleteWithContext 发送 DELETE 请求（可取消）
func (c Client) DeleteWithContext(ctx context.Context, url string) (*Response, error) {
	return c.CallWithContext(ctx, http.MethodDelete, url)
}

// handleError 处理错误
func (c Client) handleError(err error) *Client {
	if err != nil {
		if c.Error == nil {
			c.Error = NewLinkError(err)
		} else {
			c.Error = c.Error.Shift(err)
		}
	}

	return &c
}

func dnsResolver(domain string, dnsList ...string) (addrs []string, err error) {
	addr := net.ParseIP(domain)
	if len(addr) != 0 {
		return []string{domain}, nil
	}
	if host := useHostsPsrseIP(domain); host != "" {
		return []string{host}, nil
	}

	var dnsErr error
	var appendErr = func(err error) {
		if dnsErr == nil {
			dnsErr = err
		} else {
			dnsErr = fmt.Errorf("%w -> %s", dnsErr, err)
		}
	}
	var purifyDNS []string
	for k, v := range dnsList {
		if v == "" {
			continue
		}

		host, port, err := net.SplitHostPort(v)
		if err != nil {
			addrErr, ok := err.(*net.AddrError)
			if !ok {
				appendErr(fmt.Errorf("%d: %w", k, err))
				continue
			}

			// 端口不存在自动加一个
			if addrErr.Err == "missing port in address" {
				host = v
				port = "53"
			} else {
				appendErr(fmt.Errorf("%d: %w", k, err))
				continue
			}
		}

		if host == "" {
			appendErr(fmt.Errorf("%d: %w", k, errors.New("dns host is empty")))
			continue
		}

		purifyDNS = append(purifyDNS, net.JoinHostPort(host, port))
	}
	if len(purifyDNS) == 0 {
		return net.LookupHost(domain)
	}

	for _, v := range purifyDNS {
		r := net.DefaultResolver
		r.PreferGo = true
		r.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
			var d net.Dialer // 官方库就是这样用的
			return d.DialContext(ctx, network, v)
		}

		addrs, err := r.LookupHost(context.Background(), domain)
		if err == nil {
			return addrs, nil
		}

		appendErr(fmt.Errorf("%s: %w", v, err))
	}

	return nil, dnsErr
}
