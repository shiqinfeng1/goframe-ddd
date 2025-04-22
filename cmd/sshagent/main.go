package main

import (
	"context"
	"log"
	"os"

	"github.com/shellhub-io/shellhub/pkg/agent"
)

func main() {
	ctx := context.Background()
	cfg := &agent.Config{
		ServerAddress:             "http://localhost:80",
		TenantID:                  "f6509edb-7ec6-4c51-a92a-0504ee6f6490",
		PrivateKey:                "./shellhub.key",
		SingleUserPassword:        "$6$.Y2sYclyb9kbOjD2$1FE8zFGcw2hhQK0eyBur/YeULj3dSOsxkqRsg8qNGXWlwIzGempiqukp3VRuV5XE2Ncb1lf.V4mRzCLcbrQjB1",
		KeepAliveInterval:         30,
		MaxRetryConnectionTimeout: 60,
	}
	ag, err := agent.NewAgentWithConfig(cfg, new(agent.HostMode))
	if err != nil {
		log.Fatal("new ssh agent fail:", err)
	}
	// root用户
	if os.Geteuid() == 0 && cfg.SingleUserPassword != "" {
		log.Println("当前系统只有root用户, 不能使用单用户模式.")
		log.Println("要取消单用户模式, 请不要配置SingleUserPassword.")
		os.Exit(1)
	}
	// 非root用户
	if os.Geteuid() != 0 && cfg.SingleUserPassword == "" {
		log.Println("非root用户必须设置密码(hash).")
		log.Println("需使用openssl密码工具生成密码hash. agent支持如下hash算法: bsd1, apr1, sha256, sha512.")
		log.Println("举例: openssl passwd -6")
		log.Println("See man openssl-passwd for more information.")
		os.Exit(1)
	}
	if err := ag.Initialize(); err != nil {
		log.Fatal("Failed to initialize agent:", err)
	}
	if err := ag.Listen(ctx); err != nil {
		log.Fatal("Failed to listen for connections:", err)
	}
	log.Println("exit ssh agent")
}
