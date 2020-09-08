package api

import (
	"fmt"
	"github.com/Qsnh/hichat/config"
	"github.com/Qsnh/hichat/helper"
	"github.com/Qsnh/hichat/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func RoomResume(c echo.Context) error {
	rid := c.Request().FormValue("rid")
	sign := c.Request().FormValue("sign")
	if rid == "" || sign == "" {
		return fmt.Errorf("参数不能为空")
	}
	if helper.Md5(rid+config.APP_SECRET) != sign {
		return fmt.Errorf("请求非法")
	}
	if models.RM.Rooms[rid] == nil {
		return fmt.Errorf("房间不存在")
	}
	models.RM.Rooms[rid].Forbidden = false
	return c.String(http.StatusOK, "success")
}

func RoomPause(c echo.Context) error {
	rid := c.Request().FormValue("rid")
	sign := c.Request().FormValue("sign")
	if rid == "" || sign == "" {
		return fmt.Errorf("参数不能为空")
	}
	if helper.Md5(rid+config.APP_SECRET) != sign {
		return fmt.Errorf("请求非法")
	}
	if models.RM.Rooms[rid] == nil {
		return fmt.Errorf("房间不存在")
	}
	models.RM.Rooms[rid].Forbidden = true
	return c.String(http.StatusOK, "success")
}
