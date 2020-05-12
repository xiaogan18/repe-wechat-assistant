package bll

import (
	"errors"
	"fmt"
	"github.com/xiaogan18/repe-wechat-assistant/backend/dal"
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/cache"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	cacheKeyActivityUnderway string = "actv_underway"
	cacheKeyActivityPrefix   string = "actvinfo_"
	cacheKeyActivityJoined   string = "actvjoined_"
)

func NewActivityServer(task *BotTaskServer, db *dal.DbContext, conf *ContentConfig) *ActivityServer {
	return &ActivityServer{
		taskServer: task,
		cacher:     cache.NewCacher(),
		db:         db.Activity,
		roomDB:     db.Room,
		userDB:     db.User,
		mux:        sync.Mutex{},
		conf:       conf,
	}
}

type ActivityServer struct {
	taskServer    *BotTaskServer
	cacher        cache.ICache
	db            dal.IActivity
	roomDB        dal.IRoom
	userDB        dal.IUser
	addCurrentCmd func(string, int)
	mux           sync.Mutex
	conf          *ContentConfig
	quit          chan struct{}
}

func (t *ActivityServer) Start() error {
	t.quit = make(chan struct{})
	if err := t.loadUnderwayFromData(); err != nil {
		return err
	}
	go t.autoCloseActivity()
	log.Info("start activity server")
	return nil
}
func (t *ActivityServer) Stop() {
	close(t.quit)
	t.cacher.Clear()
	log.Info("stop activity server")
}
func (t *ActivityServer) setCommandTrigger(f func(string, int)) {
	t.addCurrentCmd = f
}
func (t *ActivityServer) NewActivity(actv *model.ActivityInfo, room *model.RoomInfo) error {
	log.Trace("begin new activity", "room", room.WeId, "create", actv.CreateBy, "name", actv.Name, "coin", actv.CoinType)
	fmt.Println(actv)
	// 检查发起人余额
	if actv.CreateBy > 0 {
		if actv.CoinSum <= 0 || actv.Capacity <= 0 || actv.CoinType == model.CoinIntegral {
			return errors.New("invalid arguments")
		}
		var balance float64
		if u, err := t.userDB.GetUserCoin(actv.CreateBy, actv.CoinType); err != nil {
			return err
		} else if u != nil {
			balance = u.Balance
		}
		if actv.CoinSum > balance {
			return errors.New("余额不足")
		}
	}
	// 是否已存在口令
	if t.addCurrentCmd != nil {
		t.addCurrentCmd(actv.Command, 1)
	}
	if err := t.db.Add(actv); err != nil {
		return err
	}
	// 加入缓存
	if err := t.cacheActivity(actv); err != nil {
		return err
	}
	// 扣除
	if actv.CreateBy > 0 {
		if err := t.setUserCoin(actv.CreateBy, actv.CoinType, -actv.CoinSum); err != nil {
			return err
		}
	}
	// 通知
	cpStr := "*"
	if actv.Capacity > 0 {
		cpStr = fmt.Sprintf("%v", actv.Capacity)
	}
	t.notify(room, nil, fmt.Sprintf(t.conf.Activity.TextNew, actv.Content, actv.CoinSum, actv.CoinType, cpStr))
	return nil
}
func (t *ActivityServer) CloseActivity(actvID int64) error {
	actv, err := t.getActivityInfo(actvID)
	if err != nil {
		return err
	}
	if actv.Done != uint8(model.TaskDoneNot) {
		return nil
	}
	// 修改数据库数据
	actv.Done = uint8(model.TaskDoneAlready)
	if err = t.db.SetDone(actvID, model.TaskDoneAlready); err != nil {
		return err
	}
	// 删除缓存
	t.cacher.Remove(t.cacheKey(actvID))
	t.cacher.Remove(t.cacheKeyJoined(actvID))
	// 口令移除
	if t.addCurrentCmd != nil {
		t.addCurrentCmd(actv.Command, -1)
	}
	// 修改在线集合
	t.mux.Lock()
	underway := t.allUnderwayActivity()
	for i := range underway {
		if underway[i] == actv.Id {
			underway = append(underway[:i], underway[i+1:]...)
			break
		}
	}
	if err = t.cacher.Set(cacheKeyActivityUnderway, underway); err != nil {
		t.mux.Unlock()
		return err
	}
	t.mux.Unlock()
	// 活动结束
	room, _ := t.roomDB.Get(actv.RoomId)
	rewardText := ""
	// 长期类活动结束，不通知奖励情况
	if actv.ActivityType == uint8(model.ActivityTypeOnce) {
		alogs, err := t.db.GetActivityLog(actv.Id)
		if err != nil {
			return err
		}
		for _, v := range alogs {
			u, _ := t.userDB.Get(v.UserId)
			uname := fmt.Sprint(v.UserId)
			if u != nil {
				uname = u.WeName
			}
			rewardText += fmt.Sprintf("\n%v:%v %s", uname, v.Reward, v.RewardCoin)
		}
	}
	t.notify(room, nil, fmt.Sprintf(t.conf.Activity.TextClose, actv.Name, rewardText))
	// 剩余量返还
	if actv.CoinRest == 0 || actv.CreateBy == 0 {
		return nil
	}
	if err := t.setUserCoin(actv.CreateBy, actv.CoinType, actv.CoinRest); err != nil {
		log.Error("return rest coin error", "actv_id", actv.Id, "rest", actv.CoinRest, "to", actv.CreateBy)
		return nil
	}
	creater, _ := t.userDB.Get(actv.CreateBy)
	t.notify(room, creater, fmt.Sprintf(t.conf.Activity.TextCoinReturn, actv.CoinRest, actv.CoinType))
	return nil
}
func (t *ActivityServer) Join(u *model.UserInfo, actv *model.ActivityInfo) error {
	// 检查重复参加
	if !t.setJoined(u.Id, actv.Id) {
		return errors.New("already joined")
	}
	// 检查要求
	if actv.Done == uint8(model.TaskDoneAlready) {
		return errors.New("activity already done")
	}
	if actv.Capacity > 0 && actv.Capacity <= actv.Joined {
		return errors.New("joined full")
	}
	if actv.IntegralRequire > u.Integral {
		return errors.New("require integral not enough")
	}
	if actv.IntegralCost > u.Integral {
		return errors.New("integral insufficient")
	}
	var reward = t.calcReward(actv)
	joinLog := &model.ActivityLog{
		ActivityId: actv.Id,
		UserId:     u.Id,
		RewardCoin: actv.CoinType,
		Reward:     reward,
		JoinTime:   time.Now(),
	}
	if err := t.db.SetActivityLog(joinLog); err != nil {
		return err
	}
	// 更新剩余量
	if actv.Capacity > 0 {
		actv.Joined++
		actv.CoinRest -= reward
	}
	if err := t.db.Update(actv); err != nil {
		// 回滚
		t.db.Delete(joinLog.Id)
		return fmt.Errorf("deduct activity coin err:%v", err)
	}
	// 转账
	if err := t.setUserCoin(u.Id, actv.CoinType, reward); err != nil {
		log.Error("set activity reward error", "actv_id", actv.Id, "to", u.Id, "err", err)
		return err
	}
	// 反馈通知
	room, _ := t.roomDB.Get(actv.RoomId)
	t.notify(room, u, fmt.Sprintf(t.conf.Activity.TextJoin, actv.Name, reward, joinLog.RewardCoin))
	// 参与人数已满，结束活动
	if actv.Capacity > 0 && actv.Joined >= actv.Capacity {
		if err := t.CloseActivity(actv.Id); err != nil {
			log.Error("close activity error", "id", actv.Id, "err", err)
		}
	}
	return nil
}
func (t *ActivityServer) calcReward(actv *model.ActivityInfo) float64 {
	var reward float64
	// 不限制参与人数时，coinSum为单人奖励
	switch model.RewardType(actv.RewardType) {
	case model.RewardTypeRandom:
		// 最后一个直接拿到剩余量
		if actv.Joined+1 == actv.Capacity {
			reward = actv.CoinRest
		} else if actv.Capacity == 0 {
			reward = randomFloat64(actv.CoinSum)
		} else {
			reward = randomFloat64(actv.CoinRest)
		}
	case model.RewardTypeAvg:
		if actv.Capacity == 0 {
			reward = actv.CoinSum
		} else {
			reward = actv.CoinSum / float64(actv.Capacity)
		}
	}
	return reward
}
func (t *ActivityServer) setUserCoin(userID int64, coin string, sum float64) error {
	if coin == model.CoinIntegral {
		return t.userDB.SetIntegral(userID, int64(sum))
	}
	return t.userDB.SetUserCoin(userID, coin, sum)
}

// 检查重复参加
func (t *ActivityServer) setJoined(uid, actvID int64) bool {
	key := t.cacheKeyJoined(actvID)
	data, err := t.cacher.Get(key)
	if err != nil {
		t.cacher.Set(key, []int64{uid})
		return true
	}
	joined := data.([]int64)
	for _, v := range joined {
		if v == uid {
			return false
		}
	}
	t.cacher.Set(key, append(joined, uid))
	return true
}
func (t *ActivityServer) cacheKeyJoined(actvID int64) string {
	return fmt.Sprintf("%s%d", cacheKeyActivityJoined, actvID)
}
func (t *ActivityServer) allUnderwayActivity() []int64 {
	data, err := t.cacher.Get(cacheKeyActivityUnderway)
	if err != nil {
		return nil
	}
	underway := data.([]int64)
	return underway
}
func (t *ActivityServer) foreachUnderway(f func(actv *model.ActivityInfo) bool) {
	t.mux.Lock()
	ids := t.allUnderwayActivity()
	t.mux.Unlock()
	for _, id := range ids {
		actv, err := t.getActivityInfo(id)
		if err == nil {
			if !f(actv) {
				return
			}
		}
	}
}

// 启动时检查数据库
func (t *ActivityServer) loadUnderwayFromData() error {
	ls, err := t.db.GetListDone(model.TaskDoneNot)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	for i, actv := range ls {
		log.Debug("discover undone activity", "id", actv.Id, "name", actv.Name, "cmd", actv.Command)
		// 检查过期时间
		if actv.Deadtime > 0 {
			if actv.CreateTime.Unix()+actv.Deadtime <= now {
				if err = t.db.SetDone(actv.Id, model.TaskDoneAlready); err != nil {
					return err
				}
				continue
			}
		}
		if t.addCurrentCmd != nil {
			t.addCurrentCmd(actv.Command, 1)
		}
		if err = t.cacheActivity(&ls[i]); err != nil {
			return err
		}
		var logs []model.ActivityLog
		switch actv.ActivityType {
		case uint8(model.ActivityTypeOnce):
			logs, err = t.db.GetActivityLog(actv.Id)
		case uint8(model.ActivityTypeEveryday):
			logs, err = t.db.GetActivityDayLog(actv.Id, time.Now())
		}
		if err != nil {
			return err
		}
		for _, lg := range logs {
			t.setJoined(lg.UserId, actv.Id)
		}

	}
	return nil
}
func (t *ActivityServer) autoCloseActivity() {
	tk := time.Tick(time.Second)
	checkSwitch := true
	for {
		select {
		case <-tk:
			t.mux.Lock()
			underway := t.allUnderwayActivity()
			t.mux.Unlock()
			now := time.Now().Unix()
			for _, actvID := range underway {
				actv, err := t.getActivityInfo(actvID)
				if err != nil {
					log.Error("get activity error", "id", actvID)
					continue
				}
				// 过期移除
				if actv.Deadtime > 0 && actv.CreateTime.Unix()+actv.Deadtime <= now {
					if err := t.CloseActivity(actvID); err != nil {
						log.Error("close activity error", "id", actvID, "err", err)
					}
				}
			}
			// 每日检查
			checkSwitch = t.checkReopenActivity(checkSwitch)
		case <-t.quit:
			return
		}
	}
}

// 每日零点重开每日活动
func (t *ActivityServer) checkReopenActivity(swc bool) bool {
	now := time.Now()
	if now.Hour() != t.conf.Activity.EverydayReopenHour {
		return true
	}
	switch now.Minute() {
	case t.conf.Activity.EverydayReopenMinute:
		if !swc {
			return false
		}
		// 重开每日活动
		underway := t.allUnderwayActivity()
		for _, id := range underway {
			actv, err := t.getActivityInfo(id)
			if err != nil {
				log.Error("get activity error when reopen", "id", id, "err", err)
				continue
			}
			// 删除参加记录缓存
			if actv.ActivityType == uint8(model.ActivityTypeEveryday) {
				t.cacher.Remove(t.cacheKeyJoined(actv.Id))
			}
			// 修改状态
			if err = t.db.SetDone(id, model.TaskDoneNot); err != nil {
				log.Error("set activity status error", "id", id, "err", err)
			}
		}
		return false
	default:
		return true
	}
}

// 发送主动通知
func (t *ActivityServer) notify(room *model.RoomInfo, user *model.UserInfo, content string) {
	ev := &TaskEvent{
		Content: fmt.Sprintf(t.conf.MessageTemplate, content),
	}
	if room != nil {
		ev.Room = room
	}
	if user != nil {
		ev.To = user
	}
	t.taskServer.Put(ev)
}
func (t *ActivityServer) getActivityInfo(actvID int64) (*model.ActivityInfo, error) {
	// 从缓存取
	var actv *model.ActivityInfo
	v, err := t.cacher.Get(t.cacheKey(actvID))
	if err == cache.ErrorNotFound {
		// 从数据库取
		if actv, err = t.db.Get(actvID); err != nil {
			return nil, err
		}
		// 进行中的活动加入缓存
		if actv.Done == uint8(model.TaskDoneNot) {
			if err = t.cacheActivity(actv); err != nil {
				log.Error("cache activity error", "id", actv)
			}
		}
		return actv, nil
	} else if err != nil {
		return nil, err
	}
	actv = v.(*model.ActivityInfo)
	return actv, nil
}
func (t *ActivityServer) cacheActivity(actv *model.ActivityInfo) error {
	// 缓存在线集合
	t.mux.Lock()
	defer t.mux.Unlock()
	newUnderway := t.allUnderwayActivity()
	if len(newUnderway) == 0 {
		newUnderway = []int64{actv.Id}
	} else {
		newUnderway = append(newUnderway, actv.Id)
	}
	if err := t.cacher.Set(cacheKeyActivityUnderway, newUnderway); err != nil {
		return err
	}
	// 缓存信息
	return t.cacher.Set(t.cacheKey(actv.Id), actv)
}
func (t *ActivityServer) cacheKey(id int64) string {
	return fmt.Sprintf("%s%d", cacheKeyActivityPrefix, id)
}

// 随机一个浮点数
func randomFloat64(max float64) float64 {
	str := fmt.Sprintf("%v", max)
	nums := strings.Split(str, ".")
	// 随机整数位和小数位
	var left, right int64
	v, err := strconv.ParseInt(nums[0], 10, 64)
	if err == nil && v > 0 {
		left = rand.Int63n(v)
	}
	// 整数位大于一时，直接随机一个小数
	if v > 0 {
		return float64(left) + rand.Float64()
	}
	if len(nums) > 1 {
		v, err = strconv.ParseInt(nums[1], 10, 64)
		if err == nil && v > 0 {
			right = rand.Int63n(v)
		}
	}
	if right == 0 {
		return float64(left)
	} else {
		fstr := fmt.Sprintf("%d", right)
		for len(fstr) != len(nums[1]) {
			fstr = fmt.Sprintf("0%s", fstr)
		}
		r, err := strconv.ParseFloat(fmt.Sprintf("0.%s", fstr), 64)
		if err != nil {
			panic(fstr)
		}
		fmt.Printf("left=%v right=%v fstr=%v len=%v r=%v\n", left, right, fstr, len(nums[1]), r)
		return r
	}
}
