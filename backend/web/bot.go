package web

import (
	"github.com/xiaogan18/repe-wechat-assistant/backend/bll"
	"github.com/xiaogan18/repe-wechat-assistant/backend/dal"
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/log"
	"github.com/xiaogan18/repe-wechat-assistant/backend/web/resp"
	"github.com/gin-gonic/gin"
	"net/http"
)

type botApi struct {
	cmdSvr  *bll.CommandServer
	taskSvr *bll.BotTaskServer
	db      *dal.DbContext
}

func (t botApi) Register(engine *gin.RouterGroup) {
	engine.GET("/task", t.getTask)
	engine.POST("/cmd", t.postCommand)
	engine.POST("/user", t.postUser)
	engine.POST("/room", t.postRoom)
	engine.POST("/sync/room", t.syncRooms)
	engine.POST("/sync/user", t.syncUsers)
}
func (t botApi) getTask(ctx *gin.Context) {
	nextTask := t.taskSvr.Next()
	if nextTask == nil {
		return
	}
	msg := &resp.BotMessage{
		Content: nextTask.Content,
	}
	if nextTask.To != nil {
		msg.User = nextTask.To.WeId
	}
	if nextTask.Room != nil {
		msg.Room = nextTask.Room.WeId
	}
	ctx.JSON(http.StatusOK, msg)
}
func (t botApi) postCommand(ctx *gin.Context) {
	msg := &resp.BotMessage{}
	if !assertError(ctx, ctx.BindJSON(msg)) {
		return
	}
	log.Trace("receive command", "user", msg.User, "room", msg.Room, "cmd", msg.Content)
	if len(msg.User) == 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	command := bll.NewCommand(msg.Content, msg.User, msg.Room, msg.Mention)
	t.cmdSvr.Put(command)
	success(ctx)
}
func (t botApi) postUser(ctx *gin.Context) {
	var postUser resp.WeContact
	err := ctx.BindJSON(&postUser)
	if err != nil || len(postUser.Id) == 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	u, err := t.db.User.GetByWeId(postUser.Id)
	if err == dal.ErrNotFound {
		err = t.db.User.Add(&model.UserInfo{
			WeId:   postUser.Id,
			WeName: postUser.Name,
		})
	} else if err == nil && u.WeName != postUser.Name {
		u.WeName = postUser.Name
		err = t.db.User.Update(u)
	}
	if !assertError(ctx, err) {
		return
	}
	success(ctx)
}
func (t botApi) postRoom(ctx *gin.Context) {
	var postRoom resp.WeContact
	err := ctx.BindJSON(&postRoom)
	if err != nil || len(postRoom.Id) == 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	u, err := t.db.Room.GetWeId(postRoom.Id)
	if err == dal.ErrNotFound {
		err = t.db.Room.Add(&model.RoomInfo{
			WeId:   postRoom.Id,
			WeName: postRoom.Name,
		})
	} else if err == nil && u.WeName != postRoom.Name {
		u.WeName = postRoom.Name
		err = t.db.Room.Update(u)
	}
	if !assertError(ctx, err) {
		return
	}
	success(ctx)
}
func (t botApi) syncUsers(ctx *gin.Context) {
	var users []resp.WeContact
	err := ctx.BindJSON(&users)
	if !assertError(ctx, err) {
		return
	}
	for _, v := range users {
		u, err := t.db.User.GetByWeId(v.Id)
		if err == dal.ErrNotFound {
			err = t.db.User.Add(&model.UserInfo{
				WeId:   v.Id,
				WeName: v.Name,
			})
		} else if err == nil && v.Name != u.WeName {
			u.WeName = v.Name
			err = t.db.User.Update(u)
		}
		if !assertError(ctx, err) {
			log.Error("set user error", "id", v.Id, "name", v.Name, "err", err)
			return
		}
	}
	success(ctx)
}
func (t botApi) syncRooms(ctx *gin.Context) {
	var rooms []resp.WeContact
	err := ctx.BindJSON(&rooms)
	if !assertError(ctx, err) {
		return
	}
	for _, v := range rooms {
		r, err := t.db.Room.GetWeId(v.Id)
		if err == dal.ErrNotFound {
			err = t.db.Room.Add(&model.RoomInfo{
				WeId:   v.Id,
				WeName: v.Name,
			})
		} else if err == nil && v.Name != r.WeName {
			r.WeName = v.Name
			err = t.db.Room.Update(r)
		}
		if !assertError(ctx, err) {
			return
		}
	}
	success(ctx)
}
