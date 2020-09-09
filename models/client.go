package models

import (
	"encoding/json"
	"github.com/Qsnh/hichat/config"
	"github.com/Qsnh/hichat/helper"
	"github.com/gorilla/websocket"
	"html"
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

		// 解析消息
		content := ReceiveMessage{}
		if err = json.Unmarshal(m, &content); err != nil {
			helper.SL.Errorf("消息解析失败,错误信息:%s,消息内容:%s", err, m)
			continue
		}
		if content.Type == config.HEARTBEAT_MESSAGE_TYPE {
			continue
		}

		// 过滤xss
		newMessage := ReceiveMessage{
			Type:  content.Type,
			Value: html.EscapeString(content.Value),
		}

		// 构建Message
		newMessageStr, _ := json.Marshal(newMessage)
		message := Message{
			User: MessageUser{
				Id:       c.Id,
				Nickname: c.Nickname,
				Avatar:   c.Avatar,
			},
			C: string(newMessageStr),
		}

		// 广播消息到room
		RM.Rooms[c.Rid].Broadcast <- message
	}
}
