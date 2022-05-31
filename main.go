//go:build linux || darwin || windows

// linux darwin windows：表示linux darwin windows都可以编译

package main

import (
	"embed"
	"github.com/harryruan/File_Transter/server"
	"github.com/webview/webview"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	port := "27149"
	// 创建gin引擎
	go func() {
		server.Run()
	}()

	// 启动webview
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	//w.SetTitle("Hello world!")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://localhost:" + port + "/static/index.html")
	w.Run()
}
