package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.StaticFS("/static", http.FS(staticFiles))
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/static/") {
				reader, err := staticFiles.Open("index.html")
				if err != nil {
					log.Fatal(err)
				}
				defer func(reader fs.File) {
					err := reader.Close()
					if err != nil {
						log.Fatal(err)
					}
				}(reader)
				stat, err := reader.Stat()
				if err != nil {
					log.Fatal(err)
				}
				c.DataFromReader(http.StatusOK, stat.Size(), "text/html;characters=UTF-8", reader, nil) //如果页面出现乱码，请检查页面编码是否为UTF-8
			} else {
				c.Status(http.StatusNotFound)
			}
		})
		err := router.Run(":8080")
		if err != nil {
			return
		}
	}()

	time.Sleep(time.Second * 3)
	cmd := exec.Command(`open`, "http://localhost:8080/static/index.html")
	err := cmd.Start()
	if err != nil {
		return
	}

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	select {
	case <-chSignal:
		err := cmd.Process.Kill()
		if err != nil {
			return
		}
	}
	//exec.Command(`open`, `https://baidu.com`).Start()
	//exec.Command(`open`, `https://www.baidu.com`).Start()
	//var ui lorca.UI
	//currentDir, _ := os.Getwd()
	//dir := filepath.Join(currentDir, ".cache")
	//fmt.Printf("%v\n", ui)
	// Create UI with basic HTML passed via data URI
	//ui, err := lorca.New("https://baidu.com", "", 800, 600,
	//	"--disable-sync", "--disable-transtate")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//chSignal := make(chan os.Signal, 1)
	//signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	//select {
	////case <-ui.Done():
	//case <-chSignal:
	//}
	////ui.Close()

	//debug := true
	//w := webview.New(debug)
	//defer w.Destroy()
	//w.SetTitle("Hello world!")
	//w.SetSize(800, 600, webview.HintNone)
	//w.Navigate("http://localhost:8080")
	//w.Run()
}
