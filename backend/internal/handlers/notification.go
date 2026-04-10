package handlers

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/pkg/broadcaster"
)

// GetEvents 建立SSE连接
func GetEvents(c *gin.Context) {
	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	client := &broadcaster.SSEClient{
		Channel: make(chan string, 10),
	}
	broadcaster.AddClient(client)
	defer func() {
		broadcaster.RemoveClient(client)
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
