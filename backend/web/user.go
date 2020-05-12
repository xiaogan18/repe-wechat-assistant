package web

import (
	"errors"
	"github.com/xiaogan18/repe-wechat-assistant/backend/dal"
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/log"
	"github.com/xiaogan18/repe-wechat-assistant/backend/web/resp"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type userApi struct {
	db *dal.DbContext
}

func (t userApi) Register(engine *gin.RouterGroup) {
	engine.POST("/topup", t.topUp)
	engine.POST("/withdraw", t.withdraw)
	engine.GET("/:id", t.getUser)
	engine.PUT("/", t.putUser)
}
func (t userApi) getUser(ctx *gin.Context) {
	weid := ctx.Param("id")
	if len(weid) == 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	user, err := t.db.User.GetByWeId(weid)
	if !assertError(ctx, err) {
		return
	}
	coins, err := t.db.User.GetUserAllCoin(user.Id)
	if !assertError(ctx, err) {
		return
	}
	result := &resp.UserResult{
		UserInfo: *user,
		Title:    "",
		Coins:    coins,
	}
	ctx.JSON(http.StatusOK, result)
}
func (t userApi) putUser(ctx *gin.Context) {
	f1 := ctx.PostForm("id")
	id, err := strconv.ParseInt(f1, 10, 64)
	if !assertError(ctx, err) {
		return
	}
	newAcc := ctx.PostForm("bind_account")
	if len(newAcc) == 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	err = t.db.User.SetAccount(id, newAcc)
	if !assertError(ctx, err) {
		return
	}
	success(ctx)
}
func (t userApi) topUp(ctx *gin.Context) {
	err := errors.New("not implement yet")
	assertError(ctx, err)
}
func (t userApi) withdraw(ctx *gin.Context) {
	err := errors.New("not implement yet")
	assertError(ctx, err)
}

type userManagerApi struct {
	db *dal.DbContext
}

func (t userManagerApi) Register(engine *gin.RouterGroup) {
	engine.GET("/", t.getList)
	engine.GET("/:id", t.get)
	engine.PUT("/:id", t.update)
}
func (t userManagerApi) getList(ctx *gin.Context) {
	query := make(map[string]interface{})
	if name, err := queryString(ctx, "name"); err == nil {
		query["weName"] = name
	}
	if v, err := queryInt64(ctx, "integral_from"); err == nil {
		query["integralFrom"] = v
	}
	if v, err := queryInt64(ctx, "integral_to"); err == nil {
		query["integralTo"] = v
	}
	page := queryPage(ctx)
	log.Trace("search user list", "name", query["weName"], "from", query["integralFrom"], "to", query["integralTo"], "page", page)
	ls, count, err := t.db.User.GetListPage(page.Index, page.Size, query)
	if !assertError(ctx, err) {
		return
	}
	successData(ctx, ls, count)
}
func (t userManagerApi) get(ctx *gin.Context) {
	p1 := ctx.Param("id")
	id, err := strconv.ParseInt(p1, 10, 64)
	if err != nil || id <= 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	user, err := t.db.User.Get(id)
	if !assertError(ctx, err) {
		return
	}
	coins, err := t.db.User.GetUserAllCoin(user.Id)
	if !assertError(ctx, err) {
		return
	}
	result := &resp.UserResult{
		UserInfo: *user,
		Title:    "",
		Coins:    coins,
	}
	successData(ctx, result, 0)
}
func (t userManagerApi) update(ctx *gin.Context) {
	user := new(model.UserInfo)
	err := ctx.BindJSON(user)
	if !assertError(ctx, err) {
		return
	}
	if user.Id <= 0 {
		assertError(ctx, errInvalidParams)
		return
	}
	oldUser, err := t.db.User.Get(user.Id)
	if !assertError(ctx, err) {
		return
	}
	//只允许修改的字段
	oldUser.BindAccount = user.BindAccount
	oldUser.Integral = user.Integral
	err = t.db.User.Update(oldUser)
	if !assertError(ctx, err) {
		return
	}
	successData(ctx, nil, 0)
}
