package dal

import (
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"time"
)

type IUser interface {
	Get(id int64) (*model.UserInfo, error)
	GetByWeId(weid string) (*model.UserInfo, error)
	GetListPage(index, size int, query map[string]interface{}) ([]model.UserInfo, int64, error)
	GetListByIntegral(from, to int64, index, size int) ([]model.UserInfo, int64, error)
	Add(u *model.UserInfo) error
	Update(u *model.UserInfo) error
	SetAccount(id int64, account string) error
	SetIntegral(id, integral int64) error
	GetUserAllCoin(id int64) ([]model.UserCoin, error)
	GetUserCoin(id int64, coin string) (*model.UserCoin, error)
	SetUserCoin(id int64, coin string, v float64) error
}
type IBot interface {
	GetListDone(done model.TaskDone) ([]model.BotTask, error)
	GetListPage(index, size int, query map[string]interface{}) ([]model.BotTask, int64, error)
	Get(id int64) (*model.BotTask, error)
	Add(v *model.BotTask) error
	Update(v *model.BotTask) error
	SetDone(id int64, done model.TaskDone) error
	Delete(id int64) error
}
type IRoom interface {
	GetList() ([]model.RoomInfo, error)
	Get(id int64) (*model.RoomInfo, error)
	GetWeId(weid string) (*model.RoomInfo, error)
	Add(v *model.RoomInfo) error
	Update(v *model.RoomInfo) error
	Delete(id int64) error
}
type IActivity interface {
	GetListDone(done model.TaskDone) ([]model.ActivityInfo, error)
	GetListPage(index, size int, query map[string]interface{}) ([]model.ActivityInfo, int64, error)
	Get(id int64) (*model.ActivityInfo, error)
	Add(v *model.ActivityInfo) error
	Update(v *model.ActivityInfo) error
	SetDone(id int64, done model.TaskDone) error
	Delete(id int64) error
	GetActivityLog(actID int64) ([]model.ActivityLog, error)
	GetActivityDayLog(actID int64, t time.Time) ([]model.ActivityLog, error)
	SetActivityLog(v *model.ActivityLog) error
}
