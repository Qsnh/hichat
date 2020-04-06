package main

type RoomManager struct {
	rooms       map[string]*Room
	register    chan *Room
	unregisgter chan *Room
}

type Room struct {
	rid        string `json:rid`
	clients    map[int]*Client
	count      int `json:count`
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
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
					// client已注册，那就先删除掉
					// unregister通道是同步的，所以这里会阻塞直到完成
					r.unregister <- r.clients[client.id]
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
			go func() {
				for _, c := range r.clients {
					c.receive <- m
				}
			}()
		}
	}
}
