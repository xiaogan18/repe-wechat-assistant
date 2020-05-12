package web

import (
	"github.com/xiaogan18/repe-wechat-assistant/backend/dal"
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

type roomApi struct {
	db *dal.DbContext
}

func (t roomApi) Register(engine *gin.RouterGroup) {
	engine.GET("/", t.getList)
	engine.GET("/:id", t.get)
	engine.PUT("/:id", t.update)
}
func (t roomApi) getList(ctx *gin.Context) {
	ls, err := t.db.Room.GetList()
	if !assertError(ctx, err) {
		return
	}
	successData(ctx, ls, int64(len(ls)))
}
func (t roomApi) get(ctx *gin.Context) {
	p1 := ctx.Param("id")
	id, err := strconv.ParseInt(p1, 10, 64)
	if err != nil || id <= 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	room, err := t.db.Room.Get(id)
	if !assertError(ctx, err) {
		return
	}
	successData(ctx, room, 0)
}
func (t roomApi) update(ctx *gin.Context) {
	p1 := ctx.Param("id")
	id, err := strconv.ParseInt(p1, 10, 64)
	if err != nil || id <= 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	oldRm, err := t.db.Room.Get(id)
	if !assertError(ctx, err) {
		return
	}
	rm := new(model.RoomInfo)
	err = ctx.BindJSON(rm)
	if err != nil {
		assertError(ctx, errInvalidParams)
	}
	//只允许修改字段
	oldRm.Remark = rm.Remark
	err = t.db.Room.Update(oldRm)
	if !assertError(ctx, err) {
		return
	}
	successData(ctx, nil, 0)

}
