package httpc

import (
	"bytes"
	"io/ioutil"
	"net/url"
	"reflect"
	"testing"
)

func TestOneWayCall(t *testing.T) { // 单向测试（所有变量的作用域应该只向后进行泄露）
	client := New()

	client.SetBaseURL("https://example.com")

	if client.BaseURL != "" {
		t.Error("BaseURL 设置泄露")
	}

	client.SetURLQuery("name", "niconiconi")
	if client.URLQuery != nil {
		t.Error("URLQuery 设置泄露")
	}

	client.SetHeader(ContentType, MIMEJson)
	if client.Headers != nil {
		t.Error("Header 设置泄露")
	}

	client.SetBody("niconiconi")
	if client.Body != nil {
		t.Error("Body 设置泄露")
	}
}

func TestSetBaseURL(t *testing.T) {
	exampleURL := "https://example.com"
	if SetBaseURL(exampleURL).BaseURL != exampleURL {
		t.Error("内部保存 url 与原始值不一致")
	}
}

func TestAddURLQueryS(t *testing.T) {
	if !reflect.DeepEqual(AddURLQueryS("name=niconiconi").AddURLQuery("name", "elissa").URLQuery, url.Values{"name": {"niconiconi", "elissa"}}) {
		t.Error("内部保存参数与添加不一致")
	}
}

func TestSetURLQueryS(t *testing.T) {
	query := url.Values{"name": []string{"niconiconi"}}
	if !reflect.DeepEqual(SetURLQueryS(query.Encode()).URLQuery, query) {
		t.Error("内部保存查询参数与原始值不一致")
	}
}

func TestClient_DeleteURLQuery(t *testing.T) {
	if SetURLQuery("name=niconiconi").DeleteURLQuery().URLQuery != nil {
		t.Error("删除 URL 查询参数失败")
	}
}

func TestURLQuery(t *testing.T) {
	if !reflect.DeepEqual(AddURLQueryS("name=niconiconi").AddURLQuery("name", "elissa").SetURLQueryS("name=foobar").URLQuery, url.Values{"name": {"foobar"}}) {
		t.Error("覆写查询参数失败")
	}
	if !reflect.DeepEqual(AddURLQuery("name", "niconiconi").AddURLQuery("name", "elissa").SetURLQueryS("name=foobar").URLQuery, url.Values{"name": {"foobar"}}) {
		t.Error("覆写查询参数失败")
	}
	if !reflect.DeepEqual(SetURLQueryS("name=foobar").AddURLQueryS("name=niconiconi").AddURLQueryS("name=elissa").URLQuery, url.Values{"name": {"foobar", "niconiconi", "elissa"}}) {
		t.Error("查询参数添加失败")
	}
	if !reflect.DeepEqual(SetURLQuery("name", "foobar").AddURLQueryS("name=niconiconi").AddURLQueryS("name=elissa").URLQuery, url.Values{"name": {"foobar", "niconiconi", "elissa"}}) {
		t.Error("查询参数添加失败")
	}
}

func TestSetHeader(t *testing.T) {
	if SetContentType(MIMEJson).Headers[ContentType] != MIMEJson {
		t.Error("设置 Content-Type 失败")
	}
}

func TestSetBody(t *testing.T) {
	me := []byte("niconiconi")

	b, err := ioutil.ReadAll(SetBody(bytes.NewReader(me)).Body)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(me, b) {
		t.Error("内部保存 Body 与原始值不一致")
	}
}
