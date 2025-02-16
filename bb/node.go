package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// 节点结构体
type Node struct {
	address      string
	upgrader     websocket.Upgrader
	connections  map[*websocket.Conn]bool
	eventManager *EventManager
}

// 创建新节点
func NewNode(address string, eventManager *EventManager) *Node {
	return &Node{
		address: address,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		connections:  make(map[*websocket.Conn]bool),
		eventManager: eventManager,
	}
}

// 处理连接
func (n *Node) handleConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("尝试升级 WebSocket 连接...")
	conn, err := n.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 升级失败:", err)
		return
	}
	log.Printf("成功连接到节点: %s", conn.RemoteAddr().String())
	n.connections[conn] = true
	n.eventManager.Trigger(Event{Type: EventConnected, Payload: conn.RemoteAddr().String()})
	go n.handleMessages(conn)
}

// 处理消息
func (n *Node) handleMessages(conn *websocket.Conn) {
	defer func() {
		log.Printf("关闭连接: %s", conn.RemoteAddr().String())
		delete(n.connections, conn)
		conn.Close()
		n.eventManager.Trigger(Event{Type: EventDisconnected, Payload: conn.RemoteAddr().String()})
	}()
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取消息出错: %v", err)
			} else {
				log.Printf("消息格式错误: %v", err)
				n.eventManager.Trigger(Event{Type: EventMessageFormatError, Payload: err})
			}
			break
		}
		log.Printf("收到消息: %s", msg.Content)
		n.eventManager.Trigger(Event{Type: EventMessageReceived, Payload: msg})
	}
}

// 发送消息到所有连接
func (n *Node) broadcastMessage(content string) {
	log.Printf("广播消息: %s", content)
	msg := Message{Content: content}
	for conn := range n.connections {
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("发送消息到连接 %v 失败: %v", conn.RemoteAddr(), err)
			n.eventManager.Trigger(Event{Type: EventMessageSendFailed, Payload: conn.RemoteAddr().String()})
			delete(n.connections, conn)
			conn.Close()
		} else {
			log.Printf("消息已发送到连接: %v", conn.RemoteAddr())
			n.eventManager.Trigger(Event{Type: EventMessageSent, Payload: conn.RemoteAddr().String()})
		}
	}
}

// 连接到其他节点
func (n *Node) connectToNode(addr string) {
	log.Printf("尝试连接到节点: %s", addr)
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
	}
	// 添加路径 /ws
	url := "ws://" + addr + "/ws" // <- 修改此处
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		if ne, ok := err.(net.Error); ok && ne.Timeout() {
			log.Printf("连接节点 %s 超时", addr)
			n.eventManager.Trigger(Event{Type: EventConnectionTimeout, Payload: addr})
		}
		log.Printf("连接到节点 %s 失败: %v", addr, err)
		return
	}
	log.Printf("成功连接到节点: %s", conn.RemoteAddr().String())
	n.connections[conn] = true
	n.eventManager.Trigger(Event{Type: EventConnected, Payload: conn.RemoteAddr().String()})
	go n.handleMessages(conn)
}
