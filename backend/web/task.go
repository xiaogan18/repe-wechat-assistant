package web

import (
	"github.com/xiaogan18/repe-wechat-assistant/backend/bll"
	"github.com/xiaogan18/repe-wechat-assistant/backend/dal"
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/xiaogan18/repe-wechat-assistant/backend/web/resp"
	"github.com/gin-gonic/gin"
	"strconv"
)

type taskApi struct {
	db      *dal.DbContext
	taskSvr *bll.BotTaskServer
}

func (t taskApi) Register(engine *gin.RouterGroup) {
	engine.GET("/", t.getList)
	engine.GET("/:id", t.get)
	engine.POST("/", t.add)
	engine.PUT("/:id", t.update)
	engine.PUT("/:id/done", t.doneTask)
	engine.DELETE("/:id", t.delete)
}
func (t taskApi) getList(ctx *gin.Context) {
	query := make(map[string]interface{})
	done, _ := queryInt64(ctx, "done")
	query["done"] = done
	if rm, err := queryInt64(ctx, "room"); err == nil {
		query["room"] = rm
	}
	if name, err := queryString(ctx, "name"); err == nil {
		query["name"] = name
	}
	if actvType, err := queryInt64(ctx, "type"); err == nil {
		query["type"] = actvType
	}
	page := queryPage(ctx)
	ls, count, err := t.db.Bot.GetListPage(page.Index, page.Size, query)
	if !assertError(ctx, err) {
		return
	}
	result := make([]resp.BotTask, len(ls))
	for i := range ls {
		var room string
		if r, err := t.db.Room.Get(ls[i].RoomId); err == nil {
			room = r.WeName
		}
		result[i] = resp.BotTask{ls[i], room}
	}
	successData(ctx, result, count)
}
func (t taskApi) get(ctx *gin.Context) {
	p1 := ctx.Param("id")
	id, err := strconv.ParseInt(p1, 10, 64)
	if err != nil || id < 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	v, err := t.db.Bot.Get(id)
	if !assertError(ctx, err) {
		return
	}
	result := resp.BotTask{BotTask: *v}
	if r, err := t.db.Room.Get(v.RoomId); err == nil {
		result.Room = r.WeName
	}
	successData(ctx, result, 1)
}
func (t taskApi) add(ctx *gin.Context) {
	task := new(model.BotTask)
	err := ctx.BindJSON(task)
	if !assertError(ctx, err) {
		return
	}
	task.Done = uint8(model.TaskDoneNot)
	if err = t.taskSvr.NewTask(task); err != nil {
		assertError(ctx, err)
		return
	}
	successData(ctx, nil, 0)
}
func (t taskApi) doneTask(ctx *gin.Context) {
	p1 := ctx.Param("id")
	id, err := strconv.ParseInt(p1, 10, 64)
	if err != nil || id < 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	if err := t.taskSvr.DoneTask(id); err != nil {
		assertError(ctx, err)
	}
	successData(ctx, nil, 0)
}
func (t taskApi) update(ctx *gin.Context) {
	p1 := ctx.Param("id")
	id, err := strconv.ParseInt(p1, 10, 64)
	if err != nil || id < 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	task := new(model.BotTask)
	err = ctx.BindJSON(task)
	if !assertError(ctx, err) {
		return
	}
	oldTask, err := t.db.Bot.Get(id)
	if !assertError(ctx, err) {
		return
	}
	task.Id = id
	task.Done = oldTask.Done
	task.DoneTime = oldTask.DoneTime
	if err = t.taskSvr.Update(task); err != nil {
		assertError(ctx, err)
		return
	}
	successData(ctx, nil, 0)
}
func (t taskApi) delete(ctx *gin.Context) {
	p1 := ctx.Param("id")
	id, err := strconv.ParseInt(p1, 10, 32)
	if err != nil || id < 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	t.taskSvr.DoneTask(id)
	err = t.db.Bot.Delete(id)
	if !assertError(ctx, err) {
		return
	}
	successData(ctx, nil, 0)
}
