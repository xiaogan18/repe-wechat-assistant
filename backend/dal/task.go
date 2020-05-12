package dal

import (
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/go-xorm/xorm"
	"time"
)

type botTask struct {
	db *xorm.Engine
}

func (t botTask) GetListDone(done model.TaskDone) ([]model.BotTask, error) {
	var ls []model.BotTask
	err := t.db.Where("done=?", uint8(done)).Desc("id").Find(&ls)
	return ls, err
}
func (t botTask) GetListPage(index, size int, query map[string]interface{}) ([]model.BotTask,int64,error){
	var r []model.BotTask
	ses := t.db.Limit(size, index*size)
	if len(query) > 0 {
		if done, ext := query["done"]; ext {
			ses = ses.Where("done=?", done)
		}
		if rm, ext := query["room"]; ext {
			ses = ses.Where("room_id=?", rm)
		}
		if name, ext := query["name"]; ext {
			ses = ses.Where("task_name like concat('%',?,'%')", name)
		}
		if actvType, ext := query["type"]; ext {
			ses = ses.Where("task_type=?", actvType)
		}
	}
	count, err := ses.FindAndCount(&r)
	return r, count, err
}
func (t botTask) Get(id int64) (*model.BotTask, error) {
	v := new(model.BotTask)
	b, err := t.db.ID(id).Get(v)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, nil
	}
	return v, nil
}
func (t botTask) Add(v *model.BotTask) error {
	v.Id=0
	v.Done=0
	_, err := t.db.Insert(v)
	return err
}
func (t botTask) Update(v *model.BotTask) error {
	_, err := t.db.ID(v.Id).Update(v)
	return err
}
func (t botTask) SetDone(id int64, done model.TaskDone) error {
	v := &model.BotTask{Done: uint8(done)}
	if done == model.TaskDoneAlready {
		v.DoneTime = time.Now()
	}
	_, err := t.db.ID(id).Cols("done", "done_time").Update(v)
	return err
}
func (t botTask) Delete(id int64) error {
	v := &model.BotTask{Id: id}
	_, err := t.db.ID(id).Delete(v)
	return err
}
