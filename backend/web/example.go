package web

import (
	"github.com/xiaogan18/repe-wechat-assistant/backend/dal"
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type exampleActv struct {
}

func (t exampleActv) Register(engine *gin.RouterGroup) {
	engine.GET("/actv", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "activity.tmpl", struct{}{})
	})
}

type exampleApi struct {
	db *dal.DbContext
}

func (t exampleApi) Register(engine *gin.RouterGroup) {
	engine.GET("/:id/me", func(ctx *gin.Context) {
		room := getRoomFromQuery(ctx, t.db)
		roomId := ""
		if room != nil {
			roomId = room.WeId
		}
		ip := ctx.ClientIP()
		result := struct {
			User string `json:"user"`
			Room string `json:"room"`
		}{
			User: ip,
			Room: roomId,
		}

		ctx.JSON(http.StatusOK, result)
	})
	engine.GET("/:id", func(ctx *gin.Context) {
		room := getRoomFromQuery(ctx, t.db)
		if room == nil {
			room = &model.RoomInfo{
				Id:     0,
				WeId:   "",
				WeName: "私聊",
				Remark: "",
			}
		}
		ctx.HTML(http.StatusOK, "example.tmpl", struct {
			Room *model.RoomInfo
		}{
			Room: room,
		})
	})
}

func getRoomFromQuery(ctx *gin.Context, db *dal.DbContext) *model.RoomInfo {
	id := ctx.Param("id")
	if id == "private" {
		return nil
	}
	roomID, err := strconv.ParseInt(id, 10, 64)
	if !assertError(ctx, err) {
		return nil
	}
	room, err := db.Room.Get(roomID)
	if !assertError(ctx, err) {
		return nil
	}
	return room
}
