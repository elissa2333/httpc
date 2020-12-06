# 简介
链式调用 http 客户端  
理论上只兼容最新稳点版 go 编译器，如果你在老版本上可用，那只是巧合，本库不做任何保证
本库保留随时修改的权力，不做任何稳定承诺

# 理念
单我们调用 http 请求时，最后一步操作就是发起请求（GET POST PUT PATCH DELETE 等）。
在发起请求之前可能会需要进行各种操作加工我们需要发起的请求，所以我我将 http 调用改成了链式调用。在发起之前，可选的执行各种操作。
每次 请求加工以 http 方法结束，如果你需要多次使用某一种方法，请接收链式调用某个步骤，然后再进行操作 具体请参考代理设置 和 基础 URL 设置


## 基本使用

### 发起请求

    res, err := httpc.Get("https://example.com")
    if err != nil{
        panic(err)
    }
    defer res.Body.Close()

    fmt.Println(res.Text()) // 以 text 显示响应内容

除了 Get 还支持 Options Head Post Put Patch Delete。  
如果你想自己构建请求方法，可以使用 Do 方法，注意改方法会越过大量的链式调用。  
Call 方法于基础方法相同，只是请求方式是使用字符串进行指定的。

### 设置 Body

    m := struct{
        Name string
    }{
        Name:"niconiconi",
    }

    // 支持 struct map io.Reader []byte string slice
    // 其中 struct map slice 会在 header 里自动添加 Content-Type 为 "application/json"
	res,err := httpc.SetBody(m).Post("https://postman-echo.com/post")
	if err != nil{
		panic(err)
	}
	

	fmt.Println(res.Text())

### FormData
文件上传

	file, err := os.Open("./foo.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var rows []httpc.FromDataRow
	rows = append(rows, httpc.FromDataRow{
		Key:   "file", // 对端指定的键名
		Value: file.Name(),
		Data:  file, // 如果不是数据则为 nil
	})

	res, err := httpc.SetFromData(rows...).Post("https://example.com")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	io.Copy(ioutil.Discard, res.Body) // 读取响应内容


### 设置 Header
    res,err := httpc.SetHeader("Token","123456789").Get("https://postman-echo.com/get")
	// res,err := httpc.SetHeaders(map[string]string{"foo":"bar"}) // 使用 map 设置 headers
	if err != nil{
		panic(err)
	}
	defer res.Body.Close()

	fmt.Println(res.Text())

除了直接指定 Header 本库还封装了常用的 Herader SetUserAgent 和 SetContentType

### 设置代理
因为设置理念是 调用作用域只向右进行，所以设置代理后请需要接收后才可以多调用

    //client := httpc.SetProxy("socks5://127.0.0.1:1080")
	client := httpc.SetProxy("http://127.0.0.1:8118")

    // httpc.SetProxy("http://127.0.0.1:8118").Get("https://example.com") // 本次生效代理

	res,err := client.Get("https://example.com")
	if err != nil{
		panic(err)
	}
	defer res.Body.Close()

	fmt.Println(res.Text())

### 设置超时时间

	res, err := httpc.SetTimeout(time.Second * 30).Get("https://postman-echo.com/get")


### 基础 URL 设置
再某些时候我们会对某个网址进行重复调用，可变部分仅后面的 URL 这时我们就可以使用  SetBaseURL 设置一个值，然后就只用填写，可变部分了。
我个人是使用在 对某个网址的 API 调用上。

	client := httpc.SetBaseURL("https://postman-echo.com")

	res,err := client.Get("/get")
	if err != nil{
		panic(err)
	}

	// 某个访问不是相同网址时可以暂时移除设置的 URL （实用于封装过后）
	exampleRes,err := client.DeleteBaseURL().Get("https://example.com")
	if err != nil{
		panic(err)
	}

	fmt.Println(exampleRes.Text())
	
	fmt.Println(res.Text())
	postRes,posrErr := client.Post("/post")
	if posrErr != nil{
		panic(err)
	}

	fmt.Println(postRes.Text())


### 设置 URL 查询参数

	res,err := httpc.SetURLQuery("name","niconiconi","foobar").Get("https://postman-echo.com/get")
	//res,err := httpc.SetURLQueryS("name=niconiconi&token=123456").Get("https://postman-echo.com/get")
	if err != nil{
		panic(err)
	}

	fmt.Println(res.Text())

除了 SetURLQuery 还有 AddURLQuery，Set 如果存在相同的键信息会直接进行覆盖。添加如果存在相同的键会对值进行合并处理

### 读取响应
除了上面的 `Text()` 还用 `ToJSON()` 可以使用，在我们调用 API 时大部分响应是以 json 格式进行返回的我们可以将响应读取到 struct 和 map 如果 你需要注意传入的类型必须为指针。

res, err := httpc.Get("https://postman-echo.com/get")
	if err != nil {
		panic(err)
	}

	m := map[string]interface{}{}

	if err := res.ToJSON(&m); err != nil {
		panic(err)
	}

	fmt.Println(m)

### 复合使用
比如说将我们需要的操作串在一起

	client := httpc.SetBaseURL("https://postman-echo.com")

	m := struct{
		Age int `json:"age"`
	}{
	Age: 8,
	}

	res,err := client.SetHeader("token","123456789").SetURLQuery("name","niconiconi").SetBody(m).Post("/post")
	if err != nil{
		panic(err)
	}

	fmt.Println(res.Text())