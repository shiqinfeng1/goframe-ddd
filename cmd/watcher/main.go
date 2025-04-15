package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/shiqinfeng1/goframe-ddd/internal/server/watcher"
)

func main() {
	r := gin.Default()

	r.POST("/image/upgrade/:version", watcher.UpgradeImage)

	if err := r.Run(":31083"); err != nil {
		log.Fatalf("启动服务器时出错: %v", err)
	}
}
