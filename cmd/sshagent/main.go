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
		log.Println("SSH agent cannot run as root when single-user mode is enabled.")
		log.Println("To disable single-user mode unset SHELLHUB_SINGLE_USER_PASSWORD env.")
		os.Exit(1)
	}
	// 非root用户
	if os.Geteuid() != 0 && cfg.SingleUserPassword == "" {
		log.Println("When running as non-root user you need to set password for single-user mode by SHELLHUB_SINGLE_USER_PASSWORD environment variable.")
		log.Println("You can use openssl passwd utility to generate password hash. The following algorithms are supported: bsd1, apr1, sha256, sha512.")
		log.Println("Example: SHELLHUB_SINGLE_USER_PASSWORD=$(openssl passwd -6)")
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
