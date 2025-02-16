package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("请提供节点地址，例如: go run main.go :8080")
	}
	address := os.Args[1]

	// 创建事件管理器，此时会自动注册默认事件处理器
	eventManager := NewEventManager()

	// 创建节点
	node := NewNode(address, eventManager)

	// 启动服务端
	http.HandleFunc("/ws", node.handleConnection)
	go func() {
		log.Printf("节点 %s 开始监听...", address)
		if err := http.ListenAndServe(address, nil); err != nil {
			log.Fatalf("监听 %s 失败: %v", address, err)
		}
	}()

	// 修改后的 main.go 输入处理部分
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("请输入要发送的消息（输入 'q' 退出，输入 'c <地址>' 连接到其他节点）: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input == "q" {
				// 退出逻辑
			} else if len(input) > 2 && strings.HasPrefix(input, "c ") {
				addr := strings.TrimSpace(input[2:])
				node.connectToNode(addr)
			} else {
				node.broadcastMessage(input)
			}
		}
	}()

	// 处理系统信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("收到退出信号，关闭节点...")
	for conn := range node.connections {
		conn.Close()
	}
}
