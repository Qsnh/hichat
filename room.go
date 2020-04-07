package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type RoomManager struct {
	rooms       map[string]*Room
	register    chan *Room
	unregisgter chan *Room
}

type Room struct {
	rid        string `json:"rid"`
	clients    map[int]*Client
	count      int `json:"count"`
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
	forbidden  bool
}

func (rm *RoomManager) start() {
	for {
		select {
		case room := <-rm.register:
			rm.rooms[room.rid] = room
		case room := <-rm.unregisgter:
			// 注销所有的client
			for _, client := range room.clients {
				room.unregister <- client
			}

			// 关闭通道
			close(room.register)
			close(room.unregister)
			close(room.broadcast)

			// 删除room
			delete(rm.rooms, room.rid)
		}
	}
}

func (r *Room) start() {
	defer func() {
		RM.unregisgter <- r
	}()

	for {
		select {
		case client := <-r.register:
			go func() {
				if r.clients[client.id] != nil {
					// client已注册，那就先释放原先的资源
					r.clients[client.id].conn.Close()
				}
				r.clients[client.id] = client
				// 房间人数
				r.count++
				// 发送用户加入的消息
				message := Message{
					User: MessageUser{
						Id:       client.id,
						Nickname: client.nickname,
						Avatar:   client.avatar,
					},
					T: "connect",
					C: "",
				}
				r.broadcast <- message
			}()
		case client := <-r.unregister:
			go func() {
				// 关闭socket
				client.conn.Close()
				// 房间人数
				r.count--
				delete(r.clients, client.id)
				// 发送用户离开的消息

				message := Message{
					User: MessageUser{
						Id:       client.id,
						Nickname: client.nickname,
						Avatar:   client.avatar,
					},
					T: "disconnect",
					C: "",
				}
				r.broadcast <- message
			}()
		case m := <-r.broadcast:
			// 未禁言
			if r.forbidden == false {
				go func() {
					for _, c := range r.clients {
						c.receive <- m
					}
				}()
			}
		}
	}
}

func (r *Room) heartbeat() {
	timer := time.Tick(time.Second * 15)
	for {
		if RM.rooms[r.rid] == nil {
			break
		}
		// 发送心跳消息
		SL.Infof("room:%s,count:%d", r.rid, r.count)
		for _, client := range r.clients {
			c := client
			go func() {
				err := c.conn.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					r.unregister <- c
				}
			}()
		}
		<-timer
	}
}
