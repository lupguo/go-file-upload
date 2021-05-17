package app

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/lupguo/go-file-upload/config"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type IApp interface {
	UploadImage(files []os.File) (url []string)
	Download(urls []string) (url []string)
}

type App struct {
	ctx context.Context
}

func NewApp(ctx context.Context) *App {
	return &App{
		ctx: ctx,
	}
}

// UploadImage 上传图片
func (app *App) UploadImage(files []os.File) (url []string) {
	panic("implement me")
}

// Download 下载网络URL图片
func (app *App) Download(urls []string) (url []string) {
	panic("implement me")
}

// GetCtx 返回app的ctx
func (app *App) GetCtx() context.Context {
	return app.ctx
}

// GetSavePath 获取图片存储路径
func GetSavePath() string {
	root := config.GetString("upload.root")
	subdir := time.Now().Format("2006/01")
	return fmt.Sprintf("%s/%s", root, subdir)
}

// GetRandFilename 在存储路径下，返回一个ext后缀的文件名
func GetRandFilename(path string, ext string) string {
	return fmt.Sprintf("%s/img_%10d%s", path, rand.Int63(), ext)
}

// ParseFilesToURLs 将上传图片本地地址，替换成远程URL地址返回
func ParseFilesToURLs(uploads []string) []string {
	root := config.GetString("upload.root")
	domain := config.GetString("upload.domain")
	var urls []string
	for _, filename := range uploads {
		url := strings.Replace(filename, root, domain, 1)
		urls = append(urls, url)
	}
	return urls
}
