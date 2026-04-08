package handlers

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// SSEClient SSE客户端
type SSEClient struct {
	Channel chan string
}

// 全局SSE客户端管理器
var clients = make(map[*SSEClient]bool)

// Broadcast 向所有客户端广播消息
func Broadcast(event string, data interface{}) {
	message := fmt.Sprintf("event: %s\ndata: %s\n\n", event, data)
	for client := range clients {
		select {
		case client.Channel <- message:
		default:
			// 客户端缓冲区满，跳过
		}
	}
}

// GetEvents 建立SSE连接
func GetEvents(c *gin.Context) {
	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	client := &SSEClient{
		Channel: make(chan string, 10),
	}
	clients[client] = true
	defer func() {
		delete(clients, client)
		close(client.Channel)
	}()

	// 心跳 ticker
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// 监听客户端断开和心跳
	for {
		select {
		case message, ok := <-client.Channel:
			if !ok {
				return
			}
			_, err := io.WriteString(c.Writer, message)
			if err != nil {
				return
			}
			c.Writer.Flush()
		case <-ticker.C:
			// 发送心跳
			_, err := io.WriteString(c.Writer, ": heartbeat\n\n")
			if err != nil {
				return
			}
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}
