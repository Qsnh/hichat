package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const APPSECRET = "hello"

var RM = RoomManager{
	rooms:       make(map[string]*Room),
	register:    make(chan *Room),
	unregisgter: make(chan *Room),
}

func main() {
	// 启动房间管理
	go RM.start()

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, HiChat!")
	})
	e.GET("/ws", wsHandler)
	e.Logger.Fatal(e.Start(":8911"))
}

func wsHandler(c echo.Context) error {
	res := c.Response()
	req := c.Request()
	//将http协议升级成websocket协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return fmt.Errorf("升级websocket协议失败")
	}

	rid := req.FormValue("room")
	id := req.FormValue("id")
	nickname := req.FormValue("nickname")
	avatar := req.FormValue("avatar")
	sign := req.FormValue("sign")

	if rid == "" || id == "" || nickname == "" || avatar == "" || sign == "" {
		conn.Close()
		return fmt.Errorf("参数不完整")
	}

	// sign校验
	if sign != Md5(rid+id+nickname+avatar+APPSECRET) {
		conn.Close()
		return fmt.Errorf("sign校验失败")
	}

	// 判断房间是否注册
	room := RM.rooms[rid]
	if room == nil {
		room = &Room{
			rid:        rid,
			clients:    make(map[int]*Client),
			register:   make(chan *Client),
			unregister: make(chan *Client),
			broadcast:  make(chan Message),
		}
		RM.register <- room
		go room.start()
		go room.heartbeat()
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

	return nil
}
