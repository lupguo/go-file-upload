## tips
- form表单post传递时候，表单参数放在HTTP的body内, HTTP的请求header字段: 
    - `content-type: applicaon/x-www-form-urlencoded`
    - 可以通过`enctype`属性改变Form传输数据的内容编码类型，支持:
        - `application/x-www-urlencoded`: 默认值
        - `text/plain`
        - `multipart/form-data`: 支持分段传给服务端
- go语言通过`ParseMultipartForm`对`multipart/form-data`类型进行解析，分三步:
    1. 执行`r.ParseMultipartform()`，得到`form := r.MultiPartForm`
    2. 迭代处理`form.File和form.Value`
    3. 文件处理，重点处理`fh := form.File`，内容在`fh.Header`、`fh.Filename`、`fh.Size`
    - multipart.Form
    - FileHeader
- go语言中`r.ParseMultiparForm(maxMemory)`，存储大小maxMemory字节+10M字节大小内容
    - `filename=""`的内容，如果内容大于(maxMemory+10M)则报错
    - `filename!=""`的文件，如果内容大于maxMemory，则通过`os.CreateTemp`创建临时文件，存储进磁盘
-     

```go

// -- multipart表单，包含文件上传，必须在ParseMultipartForm()后才可用
// $GOROOT/src/net/http/request.go

type Request struct {
... ...
// MultipartForm is the parsed multipart form, including file uploads.
// This field is only available after ParseMultipartForm is called.
// The HTTP client ignores MultipartForm and uses Body instead.
MultipartForm *multipart.Form
... ...
}

// -- multipart表单，Form包含值和文件
// $GOROOT/src/mime/multipart/formdata.go

// Form is a parsed multipart form.
// Its File parts are stored either in memory or on disk,
// and are accessible via the *FileHeader's Open method.
// Its Value parts are stored as strings.
// Both are keyed by field name.
type Form struct {
    Value map[string][]string
    File  map[string][]*FileHeader
}

// -- FileHeader包含上传文件的基本信息
// $GOROOT/src/mime/multipart/formdata.go
type FileHeader struct {
    Filename string // 上传文件原始名称
    Header   textproto.MIMEHeader // MIME TYPE
    Size     int64  // 上传文件字节大小

    content []byte  // 上传文件内容(部分/全部)
    tmpfile string  // 多出maxMemory部分，剩余内容放零食文件
}

// -- MIMEHeader支持
// $GOROOT/src/net/textproto/header.go

// A MIMEHeader represents a MIME-style header mapping
// keys to sets of values.
type MIMEHeader map[string][]string
``` 

## curl执行post文件命令
```shell
curl --location --request POST ':8080/upload' \
--form 'name="hey"' \
--form 'age="23"' \
--form 'file1=@"/data/file.log"' \
--form 'file3=@"/data/console.json"'
```

## 图片渲染问题
- https://www.sanarias.com/blog/1214PlayingwithimagesinHTTPresponseingolang