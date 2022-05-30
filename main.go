package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/webview/webview"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.GET("/uploads/:path", UploadsController)
		router.GET("/api/v1/addresses", AddressesController)
		router.POST("/api/v1/texts", TextsController)
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

	//time.Sleep(time.Second * 3)
	//cmd := exec.Command(`open`, "http://localhost:8080/static/index.html")
	//err := cmd.Start()
	//if err != nil {
	//	return
	//}
	//
	//chSignal := make(chan os.Signal, 1)
	//signal.Notify(chSignal, os.Interrupt)
	//select {
	//case <-chSignal:
	//	err := cmd.Process.Kill()
	//	if err != nil {
	//		return
	//	}
	//}
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

	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	//w.SetTitle("Hello world!")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://localhost:8080/static/index.html")
	w.Run()
}

func TextsController(c *gin.Context) {
	var json struct {
		Raw string `json:"raw"`
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		exe, err := os.Executable() //获取当前执行文件的路径
		if err != nil {             //如果获取失败
			log.Fatal(err) //输出错误
		}
		dir := filepath.Dir(exe) //获取当前执行文件的目录
		if err != nil {
			log.Fatal(err)
		}
		filename := uuid.New().String()          //生成一个文件名
		uploads := filepath.Join(dir, "uploads") //获取uploads目录的绝对路径
		err = os.MkdirAll(uploads, os.ModePerm)  //创建uploads目录
		if err != nil {
			log.Fatal(err)
		}
		fullpath := path.Join("uploads", filename+".txt")                            //拼接文件的绝对路径(不含目录)
		err = ioutil.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644) //将json.Raw 写入文件
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath}) //返回文件的绝对路径（不含exe所在目录）
	}

}

func AddressesController(c *gin.Context) {
	addrs, _ := net.InterfaceAddrs()
	var result []string
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"addresses": result})
}

func GetUploadsDir() (uploads string) {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe)
	uploads = filepath.Join(dir, "uploads")
	return
}

func UploadsController(c *gin.Context) {
	if path := c.Param("path"); path != "" {
		target := filepath.Join(GetUploadsDir(), path)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+path)
		c.Header("Content-Type", "application/octet-stream")
		c.File(target)
	} else {
		c.Status(http.StatusNotFound)
	}
}
