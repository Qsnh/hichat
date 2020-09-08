package models

var RM = RoomManager{
	Rooms:      make(map[string]*Room),
	Register:   make(chan *Room),
	Unregister: make(chan *Room),
}
