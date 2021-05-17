package handler

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/lupguo/go-file-upload/app"
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
	app *app.App
}

// NewHandler 创建一个Handler处理器，转发给App
func NewHandler(app *app.App) *Handler {
	return &Handler{
		ctx: app.GetCtx(),
		app: app,
	}
}

// IndexHandle 首页
func (h *Handler) IndexHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "首页，Request[%s], Host: %s, URL: %s\n", r.Host, r.Method, r.URL)
}

// Upload 图片上传
func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	// parse multipart/form-data
	maxSize := config.GetInt64("upload.maxSize")
	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		HTTPServerError(w, "r.ParseMultipartForm got err, %s", err)
		return
	}

	// init save storage path
	savePath := app.GetSavePath()
	if err := os.MkdirAll(savePath, 0755); err != nil {
		HTTPServerError(w, "os.MkdirAll(%s), err:%s", savePath, err)
		return
	}

	// parse multipart form
	mulForm := r.MultipartForm
	var uploads []string
	for k, _ := range mulForm.File {
		// k is post file key
		file, header, err := r.FormFile(k)
		if err != nil {
			log.Printf("r.FormFile() got err, %s", err)
			continue
		}

		// log print
		log.Printf("upload file headers, header=>%v, filename=%s, size=%d", header.Header, header.Filename, header.Size)

		// upload mime check
		mime := header.Header.Get("Content-Type")
		if !app.IsAllowedUpload(mime) {
			log.Printf("mime(%s) is not allowed mime type", mime)
			continue
		}

		// save file to storage path
		filename := app.GetRandFilename(savePath, app.ExtByMIME(mime))
		create, err := os.Create(filename)
		if err != nil {
			HTTPServerError(w, "fail to create file:%s, got err %s", filename, err)
			return
		}

		// copy
		_, err = io.Copy(create, file)
		if err != nil {
			log.Printf("copy file %s got err %s", header.Filename, err)
			return
		}

		// new files
		uploads = append(uploads, filename)
		log.Printf("storage local files, %v", uploads)
	}

	// turn to links
	urls := app.ParseFilesToURLs(uploads)
	log.Printf("image urls: %v", urls)

	// return json
	ret := map[string][]string{
		"uploads_urls": urls,
	}
	HTTPJson(w, ret)
}

// Download 图片下载
func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	// base check
	urls, _ := r.URL.Query()["url"]
	if len(urls) == 0 {
		HTTPServerError(w, "upload empty url!")
		return
	}

	// 下载图片，并存储到指定位置
	savePath := app.GetSavePath()
	var saveFiles []string
	for _, url := range urls {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("download url file got err:%s", err)
			continue
		}

		// create file
		ctype := resp.Header.Get("Content-Type")
		filename := app.GetRandFilename(savePath, app.ExtByMIME(ctype))
		save, err := os.Create(filename)
		if err != nil {
			log.Printf("os.Create() got err, %s", err)
			return
		}

		// write to file
		resp.Write(save)
		saveFiles = append(saveFiles, filename)
	}

	log.Printf( "download url success, save files %v", saveFiles)

	return
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
