package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket用のアップグレーダー
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// クライアントの管理を行う構造体
type Client struct {
	Conn *websocket.Conn
	Pool *Pool
}

// WebSocketのコネクションを管理するプール
type Pool struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// 新しいプールを作成
func NewPool() *Pool {
	return &Pool{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// プールの管理
func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			p.Clients[client] = true
			fmt.Println("Client registered")
		case client := <-p.Unregister:
			delete(p.Clients, client)
			client.Conn.Close()
			fmt.Println("Client unregistered")
		case message := <-p.Broadcast:
			for client := range p.Clients {
				if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
					fmt.Println("Error sending message:", err)
					client.Conn.Close()
					delete(p.Clients, client)
				}
			}
		}
	}
}

// WebSocket接続を開始するためのハンドラー
func WebSocketHandler(pool *Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
			return
		}

		client := &Client{
			Conn: conn,
			Pool: pool,
		}

		pool.Register <- client

		go func() {
			defer func() {
				pool.Unregister <- client
			}()
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("Error reading message:", err)
					return
				}
				fmt.Printf("Message received: %s\n", message)
				pool.Broadcast <- message
			}
		}()
	}
}
