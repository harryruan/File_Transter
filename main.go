package main

import (
	"github.com/gin-gonic/gin"
	"github.com/webview/webview"
	"net/http"
)

func main() {
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		router.GET("/", func(c *gin.Context) {
			c.String(http.StatusOK, "Hello World!")
		})
		router.Run(":8080")
	}()
	//time.Sleep(time.Second * 3)
	//exec.Command(`open`, "http://localhost:8080").Start()
	//select {}

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
	w.SetTitle("Hello world!")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://localhost:8080")
	w.Run()
}
