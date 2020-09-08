package models

type RoomManager struct {
	Rooms      map[string]*Room
	Register   chan *Room
	Unregister chan *Room
}

func (rm *RoomManager) Start() {
	for {
		select {
		case room := <-rm.Register:
			// 新房间注册
			rm.Rooms[room.Rid] = room
		case room := <-rm.Unregister:
			// 房间注销
			// 注销所有的client
			for _, client := range room.Clients {
				room.Unregister <- client
			}

			// 关闭通道
			close(room.Register)
			close(room.Unregister)
			close(room.Broadcast)

			// 删除room
			delete(rm.Rooms, room.Rid)
		}
	}
}
