package dal

import (
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/cache"
	"github.com/go-xorm/xorm"
	"sync"
)

type roomInfo struct {
	db     *xorm.Engine
	cach   *cache.HotMap
	indexs map[string]int64
	sync.RWMutex
}

func (t roomInfo) getCache(id int64) *model.RoomInfo {
	t.RLock()
	defer t.RUnlock()
	v := t.cach.Get(id)
	if v == nil {
		return nil
	}
	return v.(*model.RoomInfo)
}
func (t roomInfo) getCacheByWeid(weid string) *model.RoomInfo {
	t.RLock()
	defer t.RUnlock()
	id, ok := t.indexs[weid]
	if !ok {
		return nil
	}
	v := t.cach.Get(id)
	return v.(*model.RoomInfo)
}
func (t roomInfo) setCache(u *model.RoomInfo) {
	t.Lock()
	t.indexs[u.WeId] = u.Id
	t.cach.Push(u.Id, u)
	t.Unlock()
}

func (t roomInfo) GetList() ([]model.RoomInfo, error) {
	var ls []model.RoomInfo
	err := t.db.Find(&ls)
	return ls, err
}
func (t roomInfo) Get(id int64) (*model.RoomInfo, error) {
	if rm := t.getCache(id); rm != nil {
		return rm, nil
	}
	v := new(model.RoomInfo)
	b, err := t.db.ID(id).Get(v)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, ErrNotFound
	}
	t.setCache(v)
	return v, nil
}
func (t roomInfo) GetWeId(weid string) (*model.RoomInfo, error) {
	if rm := t.getCacheByWeid(weid); rm != nil {
		return rm, nil
	}
	v := new(model.RoomInfo)
	b, err := t.db.Where("we_id=?", weid).Get(v)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, ErrNotFound
	}
	t.setCache(v)
	return v, nil
}
func (t roomInfo) Add(v *model.RoomInfo) error {
	_, err := t.db.Insert(v)
	return err
}
func (t roomInfo) Update(v *model.RoomInfo) error {
	_, err := t.db.ID(v.Id).Update(v)
	t.setCache(v)
	return err
}
func (t roomInfo) Delete(id int64) error {
	_, err := t.db.ID(id).Delete(new(model.RoomInfo))
	t.Lock()
	t.cach.Remove(id)
	t.Unlock()
	return err
}
