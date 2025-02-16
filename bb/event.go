package main

import "log"

// 事件类型定义
const (
	EventConnected          = "connected"
	EventDisconnected       = "disconnected"
	EventMessageSent        = "message_sent"
	EventMessageReceived    = "message_received"
	EventConnectionTimeout  = "connection_timeout"
	EventMessageSendFailed  = "message_send_failed"
	EventMessageFormatError = "message_format_error"
)

// 消息结构体
type Message struct {
	Content string `json:"content"`
}

// 事件结构体
type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// 事件管理器结构体
type EventManager struct {
	eventHandlers map[string][]EventHandlerInterface
}

// 创建新的事件管理器，并注册默认事件处理器
func NewEventManager() *EventManager {
	log.Println("初始化事件管理器...")
	em := &EventManager{
		eventHandlers: make(map[string][]EventHandlerInterface),
	}
	handler := &EventHandler{}
	em.RegisterDefaultHandlers(handler)
	log.Println("默认事件处理器已注册")
	return em
}

// 注册事件处理器
func (em *EventManager) On(eventType string, handler EventHandlerInterface) {
	em.eventHandlers[eventType] = append(em.eventHandlers[eventType], handler)
}

// 触发事件
func (em *EventManager) Trigger(event Event) {
	if handlers, ok := em.eventHandlers[event.Type]; ok {
		for _, handler := range handlers {
			go handler.Handle(event) // 异步处理事件
		}
	}
}

// 批量注册事件处理器
func (em *EventManager) RegisterDefaultHandlers(handler EventHandlerInterface) {
	events := []string{
		EventConnected,
		EventDisconnected,
		EventMessageSent,
		EventMessageReceived,
		EventConnectionTimeout,
		EventMessageSendFailed,
		EventMessageFormatError,
	}
	for _, event := range events {
		em.On(event, handler)
	}
}
