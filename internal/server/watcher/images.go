package watcher

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/shiqinfeng1/goframe-ddd/pkg/dockerctl"
)

func listImages() (string, error) {
	output, err := gproc.ShellExec(gctx.New(), "./run images")
	if err != nil {
		return "", fmt.Errorf("执行 ./run images 时出错: %s: %w", string(output), err)
	}
	return output, nil
}
func listComposeImages() (string, error) {
	output, err := gproc.ShellExec(gctx.New(), "./run composeimages")
	if err != nil {
		return "", fmt.Errorf("执行 ./run images 时出错: %s: %w", string(output), err)
	}
	return output, nil
}
func parseDockerImages(output string) ([]string, error) {
	var imageTags []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	// 跳过标题行
	if scanner.Scan() {
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				repository := fields[0]
				tag := fields[1]
				// 排除 <none> 标签
				if tag != "<none>" {
					imageTag := fmt.Sprintf("%s:%s", repository, tag)
					imageTags = append(imageTags, imageTag)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("解析输出出错: %w", err)
	}

	return imageTags, nil
}

func parseComposeImages(output string) ([]string, error) {
	var imageTags []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	// 跳过标题行
	if scanner.Scan() {
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) >= 3 {
				repository := fields[1]
				tag := fields[2]
				// 排除 <none> 标签
				if tag != "<none>" {
					imageTag := fmt.Sprintf("%s:%s", repository, tag)
					imageTags = append(imageTags, imageTag)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("解析输出出错: %w", err)
	}

	return imageTags, nil
}

// upgradeImage 处理升级镜像版本的 API 请求
func Images(c *gin.Context) {
	result, err := listImages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dockerctl.HandlerResponse{Code: -1, Message: fmt.Sprintf("获取镜像列表出错: %v", err)})
		return
	}
	images, err := parseDockerImages(result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dockerctl.HandlerResponse{Code: -1, Message: fmt.Sprintf("获取镜像列表出错: %v", err)})
		return
	}
	c.JSON(http.StatusOK, dockerctl.HandlerResponse{Code: 0, Message: "查询镜像列表成功", Data: images})
}

// upgradeImage 处理升级镜像版本的 API 请求
func ComposeImages(c *gin.Context) {
	result, err := listComposeImages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dockerctl.HandlerResponse{Code: -1, Message: fmt.Sprintf("获取镜像列表出错: %v", err)})
		return
	}
	images, err := parseComposeImages(result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dockerctl.HandlerResponse{Code: -1, Message: fmt.Sprintf("获取镜像列表出错: %v", err)})
		return
	}
	c.JSON(http.StatusOK, dockerctl.HandlerResponse{Code: 0, Message: "查询镜像列表成功", Data: images})
}
