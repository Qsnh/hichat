package api

import (
	"fmt"
	"github.com/Qsnh/hichat/config"
	"github.com/Qsnh/hichat/helper"
	"github.com/Qsnh/hichat/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

func WsHandler(c echo.Context) error {
	res := c.Response()
	req := c.Request()
	//将http协议升级成websocket协议
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if err != nil {
		http.NotFound(res, req)
		return fmt.Errorf("无法切换websocket协议")
	}

	rid := req.FormValue("room")
	id := req.FormValue("id")
	nickname := req.FormValue("nickname")
	avatar := req.FormValue("avatar")
	sign := req.FormValue("sign")

	// 链接有效期校验
	timestamp := req.FormValue("timestamp")
	t, _ := strconv.ParseInt(timestamp, 10, 64)
	nowT := time.Now().Unix()
	if t < nowT || t-10 > nowT {
		// 误差控制在10s之内
		conn.Close()
		return fmt.Errorf("链接已过期")
	}

	if rid == "" || id == "" || nickname == "" || avatar == "" || sign == "" {
		conn.Close()
		return fmt.Errorf("参数不完整")
	}

	// sign校验
	if sign != helper.Md5(rid+id+nickname+avatar+timestamp+config.APP_SECRET) {
		conn.Close()
		return fmt.Errorf("sign校验失败")
	}

	// 判断房间是否注册
	room := models.RM.Rooms[rid]
	if room == nil {
		room = &models.Room{
			Rid:        rid,
			Clients:    make(map[int]*models.Client),
			Register:   make(chan *models.Client),
			Unregister: make(chan *models.Client),
			Broadcast:  make(chan models.Message),
		}
		models.RM.Register <- room
		go room.Start()
	}

	// 客户端注册
	cid, _ := strconv.Atoi(id)
	client := &models.Client{
		Id:       cid,
		Nickname: nickname,
		Avatar:   avatar,
		Rid:      room.Rid,
		Conn:     conn,
		Receive:  make(chan models.Message),
	}

	go client.Send()
	go client.Read()

	room.Register <- client

	return nil
}
