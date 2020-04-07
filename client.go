package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	id       int    `json:"id"`
	rid      string `json:"rid"`
	nickname string `json:"nickname"`
	avatar   string `json:"avatar"`
	conn     *websocket.Conn
	receive  chan Message
}

func (c *Client) send() {
	defer func() {
		RM.rooms[c.rid].unregister <- c
	}()

	for {
		select {
		case m, ok := <-c.receive:
			if ok == false {
				break
			}
			jsonMessage, _ := json.Marshal(m)
			c.conn.WriteMessage(websocket.TextMessage, jsonMessage)
		}
	}
}

func (c *Client) read() {
	defer func() {
		RM.rooms[c.rid].unregister <- c
	}()

	for {
		_, m, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		message := Message{
			User: MessageUser{
				Id:       c.id,
				Nickname: c.nickname,
				Avatar:   c.avatar,
			},
			T: "message",
			C: string(m),
		}
		// 广播消息到room
		RM.rooms[c.rid].broadcast <- message
	}
}
