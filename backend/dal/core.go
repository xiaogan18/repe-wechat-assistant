package dal

import (
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/cache"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type DbContext struct {
	eg       *xorm.Engine
	User     IUser
	Bot      IBot
	Room     IRoom
	Activity IActivity
}

func NewDbContext(cf DbConfig) (*DbContext, error) {
	eg, err := xorm.NewEngine(cf.DriverName, cf.DataSource)
	if err != nil {
		return nil, err
	}
	err = eg.Sync2(
		new(model.UserInfo),
		new(model.UserCoin),
		new(model.BotTask),
		new(model.RoomInfo),
		new(model.ActivityInfo),
		new(model.ActivityLog),
	)
	if err != nil {
		return nil, err
	}
	ctx := &DbContext{eg: eg}
	ctx.registerSubDb()
	return ctx, err
}
func (t *DbContext) registerSubDb() {
	t.User = userInfo{db: t.eg, cach: cache.NewHotMap(100), indexs: make(map[string]int64)}
	t.Bot = botTask{db: t.eg}
	t.Room = roomInfo{db: t.eg, cach: cache.NewHotMap(10), indexs: make(map[string]int64)}
	t.Activity = activityInfo{db: t.eg}
}
