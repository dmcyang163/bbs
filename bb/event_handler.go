package main

import "log"

// 事件处理接口
type EventHandlerInterface interface {
	Handle(event Event)
}

// 事件处理器结构体
type EventHandler struct{}

func (h *EventHandler) Handle(event Event) {
	log.Printf("处理事件: %s", event.Type)
	switch event.Type {
	case EventConnected:
		log.Printf("已连接到节点: %s\n", event.Payload.(string))
	case EventDisconnected:
		log.Printf("与节点 %s 断开连接\n", event.Payload.(string))
	case EventMessageSent:
		log.Printf("消息已发送到节点: %s\n", event.Payload.(string))
	case EventMessageReceived:
		msg := event.Payload.(Message)
		log.Printf("收到消息: %s\n", msg.Content)
	case EventConnectionTimeout:
		log.Printf("连接节点 %s 超时\n", event.Payload.(string))
	case EventMessageSendFailed:
		log.Printf("消息发送到节点 %s 失败\n", event.Payload.(string))
	case EventMessageFormatError:
		log.Printf("接收到的消息格式错误: %v\n", event.Payload)
	}
}
