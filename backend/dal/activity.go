package dal

import (
	"fmt"
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/go-xorm/xorm"
	"time"
)

type activityInfo struct {
	db *xorm.Engine
}

func (t activityInfo) GetListDone(done model.TaskDone) ([]model.ActivityInfo, error) {
	var ls []model.ActivityInfo
	err := t.db.Where("done=?", uint8(done)).Desc("id").Find(&ls)
	return ls, err
}
func (t activityInfo) GetListPage(index, size int, query map[string]interface{}) ([]model.ActivityInfo, int64, error) {
	var r []model.ActivityInfo
	ses := t.db.Limit(size, index*size)
	if len(query) > 0 {
		if done, ext := query["done"]; ext {
			ses = ses.Where("done=?", done)
		}
		if rm, ext := query["room"]; ext {
			ses = ses.Where("room_id=?", rm)
		}
		if name, ext := query["name"]; ext {
			ses = ses.Where("name like concat('%',?,'%')", name)
		}
		if actvType, ext := query["activityType"]; ext {
			ses = ses.Where("activity_type=?", actvType)
		}
	}
	count, err := ses.FindAndCount(&r)
	return r, count, err
}

func (t activityInfo) Get(id int64) (*model.ActivityInfo, error) {
	v := new(model.ActivityInfo)
	b, err := t.db.ID(id).Get(v)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, nil
	}
	return v, nil
}

func (t activityInfo) Add(v *model.ActivityInfo) error {
	v.Id = 0
	v.CreateTime = time.Now()
	v.CoinRest = v.CoinSum
	_, err := t.db.Insert(v)
	return err
}
func (t activityInfo) Update(v *model.ActivityInfo) error {
	_, err := t.db.ID(v.Id).Update(v)
	return err
}
func (t activityInfo) SetDone(id int64, done model.TaskDone) error {
	v := &model.ActivityInfo{Id: id, Done: uint8(done)}
	var err error
	if done == model.TaskDoneAlready {
		v.DoneTime = time.Now()
		_, err = t.db.ID(id).Cols("done", "done_time").Update(v)
	} else if done == model.TaskDoneNot {
		v.Joined = 0
		v.CoinRest = v.CoinSum
		_, err = t.db.ID(id).Cols("done", "done_time").Update(v)
	}
	return err
}
func (t activityInfo) Delete(id int64) error {
	_, err := t.db.ID(id).Delete(new(model.ActivityInfo))
	return err
}

func (t activityInfo) GetActivityLog(actID int64) ([]model.ActivityLog, error) {
	var ls []model.ActivityLog
	err := t.db.Where("activity_id=?", actID).Find(&ls)
	return ls, err
}
func (t activityInfo) GetActivityDayLog(actID int64, tm time.Time) ([]model.ActivityLog, error) {
	var ls []model.ActivityLog
	from := fmt.Sprintf("%d-%d-%d 00:00:00.00000", tm.Year(), tm.Month(), tm.Day())
	to := fmt.Sprintf("%d-%d-%d 23:59:59.99999", tm.Year(), tm.Month(), tm.Day())
	err := t.db.Where("activity_id=?", actID).And("join_time between ? and ?", from, to).Find(&ls)
	return ls, err
}
func (t activityInfo) SetActivityLog(v *model.ActivityLog) error {
	_, err := t.db.Insert(v)
	return err
}
func (t activityInfo) DeleteActivityLog(id int64) error {
	_, err := t.db.ID(id).Delete(new(model.ActivityLog))
	return err
}
