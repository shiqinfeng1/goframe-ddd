package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shiqinfeng1/goframe-ddd/internal/watcher"
)

func main() {
	r := gin.Default()

	r.POST("/image/upgrade/:version", watcher.UpgradeImage)
	r.POST("/image/list", watcher.Images)
	r.POST("/image/runnings", watcher.ComposeImages)

	if err := r.Run(":31083"); err != nil {
		log.Fatalf("启动服务器时出错: %v", err)
	}
}

// package main

// import "gofr.dev/pkg/gofr"

// func main() {
// 	app := gofr.New()

// 	app.GET("/greet", func(ctx *gofr.Context) (any, error) {
// 		return "Hello World!", nil
// 	})

// 	app.Run() // listens and serves on localhost:8000
// }
