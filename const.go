package httpc

const (
	charsetUTF8 = "charset=UTF-8"

	// HeaderContentType 内容类型
	HeaderContentType = "Content-Type"
	// HeaderContentDisposition 内容说明
	HeaderContentDisposition = "Content-Disposition"

	// MIMEApplicationJSON json 类型
	MIMEApplicationJSON = "application/json"
	// MIMETextHTML html 类型
	MIMETextHTML = "text/html"
	// MIMETextHTMLUTF8 html 类型附带 utf-8 声明
	MIMETextHTMLCharsetUTF8 = MIMETextHTML + "; " + charsetUTF8
	// MIMEApplicationXML xml 类型（普通用户不可读）
	MIMEApplicationXML = "application/xml"
	// MIMETextXML xml 类型 （普通用户可读）
	MIMETextXML = "text/xml"
	// MIMETextPlain 文本类型
	MIMETextPlain = "text/plain"
	// MIMEXWWWFormURLEncoded 简单表单
	MIMEXWWWFormURLEncoded = "x-www-form-urlencoded"
)
