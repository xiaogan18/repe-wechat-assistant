package bll

import (
	"github.com/xiaogan18/repe-wechat-assistant/backend/dal"
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/cache"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/log"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	cacheKeyBotTask = "bot_underway"
)

func NewBotTaskServer(botDb dal.IBot, roomDb dal.IRoom) *BotTaskServer {
	return &BotTaskServer{
		ch:      make(chan *TaskEvent, 1024),
		taskHot: cache.NewHotMap(64),
		db:      botDb,
		roomDB:  roomDb,
		cacher:  cache.NewCacher(),
	}
}

type TaskEvent struct {
	TaskID    int64
	To        *model.UserInfo
	Room      *model.RoomInfo
	Content   string
	KeepAlive bool
}
type BotTaskServer struct {
	ch      chan *TaskEvent
	taskHot *cache.HotMap
	db      dal.IBot
	roomDB  dal.IRoom
	cacher  cache.ICache
	quit    chan struct{}
	mux     sync.Mutex
}

func (t *BotTaskServer) Start() error {
	t.quit = make(chan struct{})
	if err := t.readTaskFromData(); err != nil {
		return err
	}
	go t.loopUnderwayTask()
	log.Info("start BotTaskServer")
	return nil
}
func (t *BotTaskServer) Stop() {
	close(t.quit)
	close(t.ch)
	t.cacher.Clear()
	log.Info("stop BotTaskServer")
}
func (t *BotTaskServer) NewTask(task *model.BotTask) error {
	if err := t.db.Add(task); err != nil {
		return err
	}
	if task.TaskType == uint8(model.BotTaskImmediate) {
		//即时任务直接处理
		t.putTask(task)
		return nil
	}
	// timer任务
	return t.cacheTask(task)
}
func (t *BotTaskServer) Update(task *model.BotTask) error {
	if err := t.db.Update(task); err != nil {
		return err
	}
	return t.cacheTask(task)
}
func (t *BotTaskServer) Put(ev *TaskEvent) {
	log.Trace("put bot task", "id", ev.TaskID, "keepalive", ev.KeepAlive, "content", ev.Content)
	t.ch <- ev
}
func (t *BotTaskServer) Next() *TaskEvent {
	select {
	case ev, ok := <-t.ch:
		if !ok {
			return nil
		}
		if ev.TaskID > 0 && !ev.KeepAlive {
			t.DoneTask(ev.TaskID)
		}
		log.Trace("pop bot task", "id", ev.TaskID, "keepalive", ev.KeepAlive, "content", ev.Content)
		return ev
	default:
		return nil
	}
}
func (t *BotTaskServer) DoneTask(id int64) error {
	if err := t.db.SetDone(id, model.TaskDoneAlready); err != nil {
		return err
	}
	data, err := t.cacher.Get(cacheKeyBotTask)
	if err == nil {
		ls := data.(map[int64]*model.BotTask)
		delete(ls, id)
	}
	return nil
}
func (t *BotTaskServer) cacheTask(task *model.BotTask) error {
	t.mux.Lock()
	defer t.mux.Unlock()
	var cac map[int64]*model.BotTask
	data, err := t.cacher.Get(cacheKeyBotTask)
	if err == cache.ErrorNotFound {
		cac = make(map[int64]*model.BotTask)
		cac[task.Id] = task
		return t.cacher.Set(cacheKeyBotTask, cac)
	} else if err != nil {
		return err
	}
	cac = data.(map[int64]*model.BotTask)
	cac[task.Id] = task
	return nil
}
func (t *BotTaskServer) readTaskFromData() error {
	ls, err := t.db.GetListDone(model.TaskDoneNot)
	if err != nil {
		return err
	}
	cac := make(map[int64]*model.BotTask)
	for i := range ls {
		// 处理掉即时任务
		if ls[i].TaskType == uint8(model.BotTaskImmediate) {
			t.putTask(&ls[i])
			continue
		}
		cac[ls[i].Id] = &ls[i]
	}
	t.cacher.Set(cacheKeyBotTask, cac)
	return nil
}
func (t *BotTaskServer) putTask(task *model.BotTask) {
	if task.Id > 0 {
		//去重
		if t.taskHot.Has(task.Id) {
			return
		}
		t.taskHot.Push(task.Id, struct{}{})
	}
	tk := &TaskEvent{
		TaskID:    task.Id,
		To:        nil,
		Room:      nil,
		Content:   task.Content,
		KeepAlive: task.TaskType == uint8(model.BotTaskEveryday),
	}
	if task.RoomId > 0 {
		room, err := t.roomDB.Get(task.RoomId)
		if err == nil {
			tk.Room = room
		}
	}
	t.Put(tk)
}
func (t *BotTaskServer) loopUnderwayTask() {
	tk := time.Tick(31 * time.Second)
	for {
		select {
		case <-tk:
			t.checkUnderWayTask()
		case <-t.quit:
			return
		}
	}
}
func (t *BotTaskServer) checkUnderWayTask() {
	t.mux.Lock()
	defer t.mux.Unlock()
	data, err := t.cacher.Get(cacheKeyBotTask)
	if err != nil {
		return
	}
	now := time.Now()
	ls := data.(map[int64]*model.BotTask)
	for _, v := range ls {
		if v.Done == uint8(model.TaskDoneAlready) {
			delete(ls, v.Id)
			continue
		}
		log.Trace("discover task plan", "id", v.Id, "name", v.TaskName, "time", v.TaskTime, "type", v.TaskType)
		planTime := strings.Split(v.TaskTime, ":")
		if len(planTime) < 2 {
			log.Error("invalid task plan time", "id", v.Id, "t", v.TaskTime)
			continue
		}
		hour, _ := strconv.ParseInt(planTime[0], 10, 32)
		minute, _ := strconv.ParseInt(planTime[1], 10, 32)
		if int(hour) == now.Hour() && int(minute) == now.Minute() {
			t.putTask(v)
		}
		// 让每日任务重新开启
		if v.TaskType == uint8(model.BotTaskEveryday) {
			if int(hour) == now.Hour() && int(minute)+1 == now.Minute() {
				t.taskHot.Remove(v.Id)
			}
		} else {
			// 删除过期任务
			if int(hour) < now.Hour() || (int(hour) == now.Hour() && int(minute) < now.Minute()) {
				delete(ls, v.Id)
			}
		}

	}

}
