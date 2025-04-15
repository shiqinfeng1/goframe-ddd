package watcher

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl"
)

// updateEnvFile 根据传入的版本号更新.env文件中的镜像版本
func updateYmlFile(version string) error {
	// 打开模板文件
	templateFile, err := os.Open("docker-compose.yml.tmpl")
	if err != nil {
		return err
	}
	defer templateFile.Close()

	// 创建一个新的文件用于写入替换后的内容
	outputFile, err := os.Create("docker-compose.yml")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// 创建一个扫描器来逐行读取模板文件
	scanner := bufio.NewScanner(templateFile)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()
		// 替换 ${IMAGE_VERSION} 为实际的版本号
		newLine := strings.ReplaceAll(line, "${IMAGE_VERSION}", version)
		// 将替换后的行写入输出文件
		_, err := writer.WriteString(newLine + "\n")
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// 刷新写入缓冲区
	if err := writer.Flush(); err != nil {
		return err
	}
	return nil
}

// restartContainers 执行 docker compose up 命令重启容器
func restartContainers() error {
	r, err := gproc.ShellExec(gctx.New(), "./run upgrade")
	if err != nil {
		return fmt.Errorf("执行 ./run upgrade 时出错: %s: %w", string(r), err)
	}
	return nil
}

// upgradeImage 处理升级镜像版本的 API 请求
func UpgradeImage(c *gin.Context) {
	version := c.Param("version")
	if version == "" {
		c.JSON(http.StatusBadRequest, dockerctl.HandlerResponse{Code: -1, Message: "请提供有效的镜像版本号"})
		return
	}

	if err := updateYmlFile(version); err != nil {
		c.JSON(http.StatusInternalServerError, dockerctl.HandlerResponse{Code: -1, Message: fmt.Sprintf("更新yml文件时出错: %v", err)})
		return
	}

	if err := restartContainers(); err != nil {
		c.JSON(http.StatusInternalServerError, dockerctl.HandlerResponse{Code: -1, Message: fmt.Sprintf("重启容器时出错: %v", err)})
		return
	}

	c.JSON(http.StatusOK, dockerctl.HandlerResponse{Code: 0, Message: "容器镜像版本升级成功"})
}
