package httpc

// 封装 http 客户端 理论上本应用的所有 http请求因使用本包提供的 http 客户端进行发起

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	stdURL "net/url"
	"reflect"
	"strings"
	"time"
)

// Client 客户端
type Client struct {
	BaseURL  string        // 基础 URL
	URLQuery stdURL.Values // URL 查询参数

	Headers map[string]string // 头部信息
	Body    io.Reader         // 内容
	Client  http.Client       // 客户端

	Error error // 错误
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

// SetProxy 设置代理
func (c Client) SetProxy(proxy string) *Client {
	p, err := stdURL.Parse(proxy)
	if err != nil {
		return c.handleError(err)
	}
	c.Client.Transport = &http.Transport{Proxy: http.ProxyURL(p)}

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
	return c.SetHeader(ContentType, contentType)
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
			c.Error = err
			return &c
		}

		this := &c
		if c.Headers[ContentType] == "" {
			this = this.SetContentType(MIMEJson)
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

// Call 使用指定 http 方法访问 url
func (c Client) Call(method string, url string) (*Response, error) {
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

	req, err := NewRequest(method, url, c.Headers, c.Body)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
}

// Options 获取目的资源所支持的通信选项
func (c Client) Options(url string) (*Response, error) {
	return c.Call(http.MethodOptions, url)
}

// Head 获取响应头
func (c Client) Head(url string) (*Response, error) {
	return c.Call(http.MethodHead, url)
}

// Get 发送 GET 请求（获取数据）
func (c Client) Get(url string) (*Response, error) {
	return c.Call(http.MethodGet, url)
}

// Post 发送 POST 请求（添加数据）
func (c Client) Post(url string) (*Response, error) {
	return c.Call(http.MethodPost, url)
}

//// PostForm http.PostForm
//func (c Client) PostForm(url string, values url.Values) (*Response, error) {
//	return c.SetContentType(MIMEXWWWFormURLEncoded).SetBody(values.Encode()).Post(url)
//}

// Put 发送 PUT 请求（覆盖更新）
func (c Client) Put(url string) (*Response, error) {
	return c.Call(http.MethodPut, url)
}

// Patch 发送 PATCH 请求（部分更新）
func (c Client) Patch(url string) (*Response, error) {
	return c.Call(http.MethodPatch, url)
}

// Delete 发送 DELETE 请求（删除资源）
func (c Client) Delete(url string) (*Response, error) {
	return c.Call(http.MethodDelete, url)
}

// handleError 处理错误
func (c Client) handleError(err error) *Client {
	if c.Error == nil {
		c.Error = err
	} else if err != nil {
		c.Error = fmt.Errorf("%w -> %s", c.Error, err)
	}

	return &c
}
