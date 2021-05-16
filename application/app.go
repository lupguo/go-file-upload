package application

import (
	"context"
	"math/rand"
	"os"
	"time"
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

