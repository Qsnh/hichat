package models

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	Id       int    `json:"id"`
	Rid      string `json:"rid"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Conn     *websocket.Conn
	Receive  chan Message
}

func (c *Client) Send() {
	defer func() {
		RM.Rooms[c.Rid].Unregister <- c
	}()

	for {
		select {
		case m, ok := <-c.Receive:
			if ok == false {
				break
			}
			jsonMessage, _ := json.Marshal(m)
			c.Conn.WriteMessage(websocket.TextMessage, jsonMessage)
		}
	}
}

func (c *Client) Read() {
	defer func() {
		RM.Rooms[c.Rid].Unregister <- c
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		message := Message{
			User: MessageUser{
				Id:       c.Id,
				Nickname: c.Nickname,
				Avatar:   c.Avatar,
			},
			T: "message",
			C: string(m),
		}
		// 广播消息到room
		RM.Rooms[c.Rid].Broadcast <- message
	}
}
