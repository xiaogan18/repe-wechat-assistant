package dal

import (
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/cache"
	"github.com/go-xorm/xorm"
	"sync"
)

type userInfo struct {
	db     *xorm.Engine
	cach   *cache.HotMap
	indexs map[string]int64
	sync.RWMutex
}

func (t userInfo) getCache(id int64) *model.UserInfo {
	t.RLock()
	defer t.RUnlock()
	v := t.cach.Get(id)
	if v == nil {
		return nil
	}
	return v.(*model.UserInfo)
}
func (t userInfo) getCacheByWeid(weid string) *model.UserInfo {
	t.RLock()
	defer t.RUnlock()
	id, ok := t.indexs[weid]
	if !ok {
		return nil
	}
	v := t.cach.Get(id)
	return v.(*model.UserInfo)
}
func (t userInfo) setCache(u *model.UserInfo) {
	t.Lock()
	t.indexs[u.WeId] = u.Id
	t.cach.Push(u.Id, u)
	t.Unlock()
}
func (t userInfo) Get(id int64) (*model.UserInfo, error) {
	if u := t.getCache(id); u != nil {
		return u, nil
	}
	r := &model.UserInfo{Id: id}
	if b, err := t.db.Get(r); err != nil {
		return nil, err
	} else if !b {
		return nil, ErrNotFound
	}
	t.setCache(r)
	return r, nil
}
func (t userInfo) GetByWeId(weid string) (*model.UserInfo, error) {
	if u := t.getCacheByWeid(weid); u != nil {
		return u, nil
	}
	r := &model.UserInfo{WeId: weid}
	if b, err := t.db.Where("we_id=?", weid).Get(r); err != nil {
		return nil, err
	} else if !b {
		return nil, ErrNotFound
	}
	t.setCache(r)
	return r, nil
}
func (t userInfo) GetListPage(index, size int, query map[string]interface{}) ([]model.UserInfo, int64, error) {
	var r []model.UserInfo
	ses := t.db.Limit(size, index*size)
	if len(query) > 0 {
		if name, ext := query["weName"]; ext {
			ses = ses.Where("we_name like concat('%',?,'%')", name)
		}
		if from, ext := query["integralFrom"]; ext {
			ses = ses.Where("integral>=?", from)
		}
		if from, ext := query["integralTo"]; ext {
			ses = ses.Where("integral<?", from)
		}
	}
	count, err := ses.FindAndCount(&r)
	return r, count, err
}
func (t userInfo) GetListByIntegral(from, to int64, index, size int) ([]model.UserInfo, int64, error) {
	cond := t.db.Where("integral between ? and ?", from, to)
	var r []model.UserInfo
	count, err := cond.Limit(size, index*size).Asc("integral").FindAndCount(&r)
	return r, count, err
}
func (t userInfo) Add(u *model.UserInfo) error {
	_, err := t.db.Insert(u)
	return err
}
func (t userInfo) Update(u *model.UserInfo) error {
	_, err := t.db.ID(u.Id).Update(u)
	t.setCache(u)
	return err
}
func (t userInfo) SetAccount(id int64, account string) error {
	u := &model.UserInfo{Id: id, BindAccount: account}
	_, err := t.db.ID(id).Cols("bind_account").Update(u)
	if user := t.getCache(id); user != nil {
		user.BindAccount = account
	}
	return err
}
func (t userInfo) SetIntegral(id, integral int64) error {
	u, err := t.Get(id)
	if err != nil {
		return err
	}
	u.Integral = u.Integral + integral
	_, err = t.db.ID(id).Cols("integral").Update(u)
	return err
}
func (t userInfo) GetUserAllCoin(id int64) ([]model.UserCoin, error) {
	var ls []model.UserCoin
	err := t.db.Where("user_id=?", id).Find(&ls)
	return ls, err
}
func (t userInfo) GetUserCoin(id int64, coin string) (*model.UserCoin, error) {
	uc := new(model.UserCoin)
	b, err := t.db.Where("user_id=? and coin=?", id, coin).Get(uc)
	if err != nil {
		return nil, err
	}
	if !b {
		return nil, nil
	}
	return uc, nil
}
func (t userInfo) SetUserCoin(id int64, coin string, v float64) error {
	uc, err := t.GetUserCoin(id, coin)
	if err != nil {
		return err
	}
	if uc != nil {
		uc.Balance += v
		_, err = t.db.ID(uc.Id).Cols("balance").Update(uc)
	} else {
		uc = &model.UserCoin{UserId: id, Coin: coin, Balance: v}
		_, err = t.db.Insert(uc)
	}
	return err
}
