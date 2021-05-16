package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/lupguo/go-file-upload/application"
	"github.com/lupguo/go-file-upload/config"
)

// BusinessHandle 业务处理
type BusinessHandle interface {
	IndexHandle(w http.ResponseWriter, r *http.Request)
	Upload(w http.ResponseWriter, r *http.Request)
	Download(w http.ResponseWriter, r *http.Request)
}

// Handler 业务处理
type Handler struct {
	ctx context.Context
	app *application.App
}

// NewHandler 创建一个Handler处理器，转发给App
func NewHandler(app *application.App) *Handler {
	return &Handler{
		ctx: app.GetCtx(),
		app: app,
	}
}

// IndexHandle 首页
func (h *Handler) IndexHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "首页，Request[%s], Host: %s, URL: %s\n", r.Host, r.Method, r.URL)
}

var extMIME = map[string]string{
	"image/jpeg":               ".jpg",
	"image/gif":                ".gif",
	"image/png":                ".png",
	"image/bmp":                ".bmp",
	"image/webp":               ".webp",
	"image/x-icon":             ".ico",
	"image/vnd.microsoft.icon": ".ico",
}

var mimeExt = map[string]string{
	".jpg":  "image/jpeg",
	".gif":  "image/gif",
	".png":  "image/png",
	".bmp":  "image/bmp",
	".webp": "image/webp",
	".ico":  "image/x-icon",
}

// GetImageExtByMIME 通过mime获取图片的后缀名，返回以.为前缀的扩展，如.jpg
func GetImageExtByMIME(mime string) string {
	if ext, ok := extMIME[mime]; ok {
		return ext
	}
	return ".omg"
}

// IsAllowedUpload 是否允许上传的MIMEType
func IsAllowedUpload(mime string) bool {
	if _, ok := extMIME[mime]; ok {
		return true
	}
	return false
}

// Upload 图片上传
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	// parse multipart/form-data
	maxSize := config.GetInt64("upload.maxSize")
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		log.Printf("ParseMultipartForm() got err, %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// init save storage path
	root := config.GetString("upload.root")
	subdir := time.Now().Format("2006/01")
	savePath := fmt.Sprintf("%s/%s", root, subdir)
	if err := os.MkdirAll(savePath, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// parse multipart form
	mulForm := r.MultipartForm
	var uploads []string
	for k, headers := range mulForm.File {
		_ = headers
		// k is post file key
		file, header, err := r.FormFile(k)
		if err != nil {
			panic(err)
		}

		// log print
		log.Printf("upload file headers, header=>%v, filename=%s, size=%d", header.Header, header.Filename, header.Size)

		// upload mime check
		mime := header.Header.Get("Content-Type")
		if !IsAllowedUpload(mime) {
			log.Printf("not allowed mime type, %s", mime)
			continue
		}

		// save file to storage path
		ext := GetImageExtByMIME(mime)
		filename := fmt.Sprintf("%s/img_%10d%s",savePath,rand.Int63(), ext)
		create, err := os.Create(filename)
		if err != nil {
			return
		}
		_, err = io.Copy(create, file)
		if err != nil {
			panic(err)
			return
		}

		// new files
		uploads = append(uploads, filename)
		log.Printf("storage local files, %v", uploads)
	}

	// transform http image link
	domain := config.GetString("upload.domain")
	var urls []string
	for _, upname := range uploads {
		url := strings.Replace(upname, root, domain, 1)
		urls = append(urls, url)
	}
	log.Printf("image urls: %v", urls)

	// fmt.Fprintf(w, "接收POST 图片文件，生成文件名称，存储到指定位置, Request[%s], Host: %s, URL: %s\n", r.Host, r.Method, r.URL)
	ret := map[string][]string{
		"uploads_urls": urls,
	}
	retJson, _ := json.Marshal(ret)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", retJson)
}

// Download 图片下载
func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "接收POST URL，下载POST URL图片，存储到指定目录")
}

// ShowImage 图片显示
func (h *Handler) ShowImage(w http.ResponseWriter, r *http.Request) {
	root := config.GetString("upload.root")
	fname := strings.Replace(r.URL.Path, "/upload_imgs", root, 1)
	log.Printf("fname=>%s", fname)

	// read data content, set content type and write to response
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Printf("read data got err, %v", err)
		return
	}

	// got mime type
	contentType := http.DetectContentType(data)

	w.Header().Set("Content-Type", contentType)
	// w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	_, err = w.Write(data)
	if err != nil {
		log.Printf("w.Write() got err, %s", err)
		return
	}
}
