package web

import (
	"fmt"
	"github.com/xiaogan18/repe-wechat-assistant/backend/bll"
	"github.com/xiaogan18/repe-wechat-assistant/backend/dal"
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/log"
	"github.com/xiaogan18/repe-wechat-assistant/backend/web/resp"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type actvApi struct {
	db         *dal.DbContext
	actvServer *bll.ActivityServer
}

func (t actvApi) Register(engine *gin.RouterGroup) {
	engine.GET("/", t.getList)
	engine.GET("/:id", t.get)
	engine.GET("/:id/log", t.getLogList)
	engine.POST("/", t.add)
	engine.PUT("/:id/done", t.setDone)
}
func (t actvApi) get(ctx *gin.Context) {
	p1 := ctx.Param("id")
	id, err := strconv.ParseInt(p1, 10, 64)
	if err != nil {
		assertError(ctx, errInvalidParams)
		return
	}
	actv, err := t.db.Activity.Get(id)
	if !assertError(ctx, err) {
		return
	}
	ctx.JSON(http.StatusOK, actv)
}
func (t actvApi) getList(ctx *gin.Context) {
	query := make(map[string]interface{})
	done, _ := queryInt64(ctx, "done")
	query["done"] = done
	if rm, err := queryInt64(ctx, "room"); err == nil {
		query["room"] = rm
	}
	if name, err := queryString(ctx, "name"); err == nil {
		query["name"] = name
	}
	if actvType, err := queryInt64(ctx, "activityType"); err == nil {
		query["activityType"] = actvType
	}
	page := queryPage(ctx)
	ls, count, err := t.db.Activity.GetListPage(page.Index, page.Size, query)
	if !assertError(ctx, err) {
		return
	}
	log.Trace("get activity list", "done", done, "room", query["room"], "name", query["name"], "type", query["activityType"], "count", count)
	var result []resp.Activity
	for i := range ls {
		createBy := "-"
		if ls[i].CreateBy > 0 {
			if u, _ := t.db.User.Get(ls[i].CreateBy); u != nil {
				createBy = u.WeName
			} else {
				createBy = fmt.Sprintf("%d", ls[i].CreateBy)
			}
		}
		room := ""
		if r, _ := t.db.Room.Get(ls[i].RoomId); r != nil {
			room = r.WeName
		}
		result = append(result, resp.Activity{ls[i], createBy, room})
	}
	successData(ctx, result, count)
}
func (t actvApi) getLogList(ctx *gin.Context) {
	p1 := ctx.Param("id")
	id, err := strconv.ParseInt(p1, 10, 64)
	if err != nil {
		assertError(ctx, errInvalidParams)
		return
	}
	ls, err := t.db.Activity.GetActivityLog(id)
	if !assertError(ctx, err) {
		return
	}
	log.Trace("get activity log", "actv", id, "count", len(ls))
	result := make([]resp.ActivityLog, len(ls))
	for i, v := range ls {
		u := ""
		if juser, err := t.db.User.Get(v.UserId); err == nil {
			u = juser.WeName
		}
		result[i] = resp.ActivityLog{
			ActivityId: v.Id,
			User:       u,
			Reward:     v.Reward,
			RewardCoin: v.RewardCoin,
			JoinTime:   v.JoinTime.String(),
		}
	}
	successData(ctx, result, int64(len(result)))
}
func (t actvApi) add(ctx *gin.Context) {
	actv := new(model.ActivityInfo)
	err := ctx.BindJSON(actv)
	if !assertError(ctx, err) {
		return
	}
	if len(actv.Name) == 0 || len(actv.Command) == 0 || len(actv.Content) == 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	room, err := t.db.Room.Get(actv.RoomId)
	if !assertError(ctx, err) {
		return
	}
	actv.Done = uint8(model.TaskDoneNot)
	actv.CreateBy = 0
	actv.Joined = 0
	actv.CoinRest = actv.CoinSum
	// 每日活动无超时
	if actv.ActivityType == uint8(model.ActivityTypeEveryday) {
		actv.Deadtime = 0
	}
	err = t.actvServer.NewActivity(actv, room)
	if !assertError(ctx, err) {
		return
	}
	successData(ctx, nil, 0)
}
func (t actvApi) setDone(ctx *gin.Context) {
	p1 := ctx.Param("id")
	id, err := strconv.ParseInt(p1, 10, 64)
	if err != nil || id <= 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	if err = t.actvServer.CloseActivity(id); err != nil {
		assertError(ctx, err)
		return
	}
	successData(ctx, nil, 0)
}
