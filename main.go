package main

import (
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func main() {
	// 监听端口
	listener, err := net.Listen("tcp", "0.0.0.0:server_port")
	if err != nil {
		log.Fatalf("监听端口出错: %v", err)
	}
	defer listener.Close()

	log.Println("等待客户端连接...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("接受连接出错: %v", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// 配置 SSH 客户端
	config := &ssh.ClientConfig{
		User: "your_username",
		Auth: []ssh.AuthMethod{
			ssh.Password("your_password"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 注意：生产环境需替换为安全的主机密钥验证
	}

	// 建立 SSH 连接
	sshConn, chans, reqs, err := ssh.NewClientConn(conn, conn.RemoteAddr().String(), config)
	if err != nil {
		log.Printf("建立 SSH 连接出错: %v", err)
		return
	}
	defer sshConn.Close()

	client := ssh.NewClient(sshConn, chans, reqs)
	session, err := client.NewSession()
	if err != nil {
		log.Printf("创建 SSH 会话出错: %v", err)
		return
	}
	defer session.Close()

	// 执行命令
	output, err := session.CombinedOutput("ls -l")
	if err != nil {
		log.Printf("执行命令出错: %v", err)
		return
	}
	log.Println(string(output))
}
