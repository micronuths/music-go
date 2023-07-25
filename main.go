package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"music-go/config"
	"music-go/middleware"
	"music-go/model"
	"music-go/router"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("gin.Version: ", gin.Version)
	if config.Conf.GoConf.Env != model.DevelopmentMode {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()

	}
	logFile, err := os.OpenFile(config.Conf.GoConf.LogDir, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		fmt.Printf(err.Error())
		os.Exit(-1)
	}
	gin.DefaultWriter = io.MultiWriter(logFile)

	app := gin.New()
	maxSize := int64(32)
	app.MaxMultipartMemory = maxSize << 20 // 3 MiB

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	app.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	app.Use(gin.Recovery())
	app.Use(middleware.CORSMiddleware())
	router.Route(app)

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(fmt.Sprintf("http://127.0.0.1:%d", config.Conf.GoConf.Port)); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Println("The router has been deployed successfully.")
	}()

	log.Printf("Start to listening the incoming requests on http address: %d\n", config.Conf.GoConf.Port)
	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", config.Conf.GoConf.Port), app).Error())
	//
}

func initServer(address string, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:              address,
		Handler:           router,
		ReadHeaderTimeout: 20 * time.Second,
		WriteTimeout:      20 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}

// pingServer 确保 http server是工作的.
func pingServer(url string) error {
	for i := 0; i < 10; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == 200 {
			fmt.Printf("在浏览器中打开：%s\n", url)
			openWebsite(url)
			return nil
		}
		// 等待间隔
		fmt.Println(err, "Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}

	return errors.New("Cannot connect to the router.")
}

// 用默认浏览器打开网站
func openWebsite(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	default: // Linux 和其他Unix-like系统
		cmd = exec.Command("xdg-open", url)
	}

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
}
