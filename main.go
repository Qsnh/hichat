package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

const APPSECRET = "hello"

var RM = RoomManager{
	rooms:       make(map[string]*Room),
	register:    make(chan *Room),
	unregisgter: make(chan *Room),
}

func main() {
	go RM.start()

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(res, "hello hichat.")
	})
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("listening 0.0.0.0:8911")

	err := http.ListenAndServe(":8911", nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}
}

func wsHandler(res http.ResponseWriter, req *http.Request) {
	//将http协议升级成websocket协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		fmt.Println("升级websocket协议失败")
		http.NotFound(res, req)
		return

	}

	rid := req.FormValue("room")
	id := req.FormValue("id")
	nickname := req.FormValue("nickname")
	avatar := req.FormValue("avatar")
	sign := req.FormValue("sign")

	if rid == "" || id == "" || nickname == "" || avatar == "" || sign == "" {
		fmt.Println("请填写完整的参数")
		conn.Close()
		return
	}

	// sign校验
	if sign != Md5(rid+id+nickname+avatar+APPSECRET) {
		fmt.Println("sign校验失败")
		conn.Close()
		return
	}

	// 判断房间是否注册
	room := RM.rooms[rid]
	if room == nil {
		room = &Room{
			rid:        rid,
			clients:    make(map[int]*Client),
			count:      0,
			register:   make(chan *Client),
			unregister: make(chan *Client),
			broadcast:  make(chan Message),
		}
		RM.register <- room
		go room.start()
	}

	// 客户端注册
	cid, _ := strconv.Atoi(id)
	client := &Client{
		id:       cid,
		nickname: nickname,
		avatar:   avatar,
		rid:      room.rid,
		conn:     conn,
		receive:  make(chan Message),
	}

	go client.send()
	go client.read()

	room.register <- client
}
