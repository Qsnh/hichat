package main

import (
	"github.com/Qsnh/hichat/api"
	"github.com/Qsnh/hichat/config"
	"github.com/Qsnh/hichat/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	go models.RM.Start()

	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/", api.Index)
	e.GET("/room/resume", api.RoomResume)
	e.GET("/room/pause", api.RoomPause)
	e.GET("/ws", api.WsHandler)

	e.Logger.Fatal(e.Start(config.APP_PORT))
}
