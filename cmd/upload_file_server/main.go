package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/lupguo/go-file-upload/application"
	"github.com/lupguo/go-file-upload/config"
	"github.com/lupguo/go-file-upload/handler"
)


func main() {
	listen := config.GetString("service.listen")
	fmt.Printf("go http file upload start(%s)...\n", listen)

	// 两个ctx
	baseCtx := func(net.Listener) context.Context {
		return context.Background()
	}
	connCtx := func(ctx context.Context, c net.Conn) context.Context {
		return context.Background()
	}

	// 创建一个http服务
	s := &http.Server{
		Addr:           listen,
		Handler:        nil,
		TLSConfig:      nil,
		ReadTimeout:    config.GetDuration("service.read_timeout"),
		WriteTimeout:   config.GetDuration("service.write_timeout"),
		MaxHeaderBytes: config.GetInt("service.max_header_bytes"),
		TLSNextProto:   nil,
		ErrorLog:       nil, // log
		BaseContext:    baseCtx,
		ConnContext:    connCtx,
	}

	// 基础服务器处理程序
	app := application.NewApp(context.Background())
	h := handler.NewHandler(app)
	http.HandleFunc("/", h.IndexHandle)
	http.HandleFunc("/upload", h.Upload)
	http.HandleFunc("/download", h.Download)
	http.HandleFunc("/upload_imgs/", h.ShowImage)

	// 服务启动和监听
	log.Fatal(s.ListenAndServe())
}
