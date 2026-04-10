package broadcaster

import (
	"encoding/json"
	"fmt"
	"sync"
)

// SSEClient SSE客户端
type SSEClient struct {
	Channel chan string
}

// 全局SSE客户端管理器
var (
	clients = make(map[*SSEClient]bool)
	mu      sync.RWMutex
)

// Broadcast 向所有客户端广播消息
func Broadcast(event string, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	message := fmt.Sprintf("event: %s\ndata: %s\n\n", event, jsonData)

	mu.RLock()
	defer mu.RUnlock()

	for client := range clients {
		select {
		case client.Channel <- message:
		default:
			// 客户端缓冲区满，跳过
		}
	}
}

// AddClient 添加新客户端
func AddClient(client *SSEClient) {
	mu.Lock()
	defer mu.Unlock()
	clients[client] = true
}

// RemoveClient 移除客户端
func RemoveClient(client *SSEClient) {
	mu.Lock()
	defer mu.Unlock()
	delete(clients, client)
}

// GetClientCount 获取客户端数量
func GetClientCount() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(clients)
}