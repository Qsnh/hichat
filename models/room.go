package models

import "strconv"

type Room struct {
	Rid        string `json:"rid"`
	Clients    map[int]*Client
	Count      int `json:"count"`
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
	Forbidden  bool
}

func (r *Room) Start() {
	defer func() {
		RM.Unregister <- r
	}()

	for {
		select {
		case client := <-r.Register:
			go func() {
				if r.Clients[client.Id] != nil {
					// client已注册，那就先释放原先的资源
					r.Clients[client.Id].Conn.Close()
				}
				r.Clients[client.Id] = client
				// 房间人数
				r.Count++
				// 发送用户加入的消息
				message := Message{
					User: MessageUser{
						Id:       client.Id,
						Nickname: client.Nickname,
						Avatar:   client.Avatar,
					},
					C: "{\"t\":\"connect\",\"count\":" + strconv.Itoa(r.Count) + "}",
				}
				r.Broadcast <- message
			}()
		case client := <-r.Unregister:
			go func() {
				// 关闭socket
				client.Conn.Close()
				// 房间人数
				r.Count--
				delete(r.Clients, client.Id)
				// 发送用户离开的消息
				message := Message{
					User: MessageUser{
						Id:       client.Id,
						Nickname: client.Nickname,
						Avatar:   client.Avatar,
					},
					C: "{\"t\":\"disconnect\",\"count\":" + strconv.Itoa(r.Count) + "}",
				}
				r.Broadcast <- message
			}()
		case m := <-r.Broadcast:
			// 未禁言
			if r.Forbidden == false {
				go func() {
					for _, c := range r.Clients {
						c.Receive <- m
					}
				}()
			}
		}
	}
}
