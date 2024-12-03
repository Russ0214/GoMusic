package websocket

import (
    "sync"
    "encoding/json"
)

type Notifier struct {
    connections map[string]*Connection
    mu          sync.RWMutex
}

type Connection struct {
    UserID string
    Conn   *websocket.Conn
}

type Message struct {
    Type    string `json:"type"`
    TaskID  string `json:"task_id"`
    Message string `json:"message"`
    ETA     int    `json:"eta"` // 预计剩余时间(秒)
}

func (n *Notifier) SendUpdate(userID string, msg Message) {
    n.mu.RLock()
    conn, exists := n.connections[userID]
    n.mu.RUnlock()
    
    if exists {
        data, _ := json.Marshal(msg)
        conn.Conn.WriteMessage(websocket.TextMessage, data)
    }
} 