# 简介
链式调用 http 客户端  
理论上只兼容最新稳点版 go 编译器，如果你在老版本上可用，那只是巧合，本库不做任何保证
本库保留随时修改的权力，不做任何稳定承诺

# 理念
单我们调用 http 请求时，最后一步操作就是发起请求（GET POST PUT PATCH DELETE 等）。
在发起请求之前可能会需要进行各种操作加工我们需要发起的请求，所以我我将 http 调用改成了链式调用。在强求发起之前，可选的执行各种操作


## 基本使用

### 发起请求

    res, err := httpc.get("https://example.com")
    if err != nil{
        panic(err)
    }
    defer res.Close()

    fmt.Println(res.Text()) // 以 text 显示响应内容

除了了最基础的 get 支持其它 http 方法

### 设置 Body

    m := struct{
        Name string
    }{
        Name:"niconiconi"
    }

    // 支持 struct map io.Reader []byte string slice
    // 其中 struct map slice 会在 header 里自动添加 Content-Type 为 "application/json"
    res,err := httpc.SetBody(m).post(https://example.com)
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
