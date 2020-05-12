package bll

import (
	"fmt"
	"github.com/xiaogan18/repe-wechat-assistant/backend/dal"
	"github.com/xiaogan18/repe-wechat-assistant/backend/model"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/cache"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/log"
	"strings"
)

const (
	cacheKeyCommandPrefix string = "cmd_key_"
	commandSplitChar      string = " "
)

func NewCommand(cmd, from, room string, mention []string) *Command {
	args := strings.Split(cmd, commandSplitChar)
	if len(args) == 0 {
		return nil
	}
	c := &Command{
		From:    from,
		Room:    room,
		Cmd:     args[0],
		Mention: mention,
	}
	if len(args) > 1 {
		c.Arg = args[1:]
	}
	return c
}

type Command struct {
	Room    string
	From    string
	Cmd     string
	Mention []string
	Arg     []string
}

func NewCmdServer(ctx *dal.DbContext, actvServer *ActivityServer, conf *ContentConfig) *CommandServer {
	cmdSvr := &CommandServer{
		ch:         make(chan *Command, 1024),
		userDB:     ctx.User,
		roomDB:     ctx.Room,
		cacher:     cache.NewCacher(),
		actvServer: actvServer,
		conf:       conf,
	}
	actvServer.setCommandTrigger(cmdSvr.addCurrentCmd)
	return cmdSvr
}

type CommandServer struct {
	ch         chan *Command
	execSort   []string
	userDB     dal.IUser
	roomDB     dal.IRoom
	cacher     cache.ICache
	actvServer *ActivityServer
	conf       *ContentConfig
	quit       chan struct{}
}

func (t *CommandServer) Start() error {
	t.quit = make(chan struct{})
	go t.listenCmd()
	log.Info("start command server")
	return nil
}
func (t *CommandServer) Stop() {
	close(t.quit)
	close(t.ch)
	t.cacher.Clear()
	log.Info("stop command server ")
}
func (t *CommandServer) ExistCmd(cmd string) bool {
	// 检查缓存指令
	_, err := t.cacher.Get(t.cacheKey(cmd))
	if err != nil {
		return false
	}
	return true
}
func (t *CommandServer) Put(cmd *Command) {
	if len(cmd.From) == 0 {
		return
	}
	select {
	case t.ch <- cmd:
	default:
		log.Warn("command server chan full", "cmd", cmd)
	}
}
func (t *CommandServer) cacheKey(cmd string) string {
	return fmt.Sprintf("%s%s", cacheKeyCommandPrefix, cmd)
}
func (t *CommandServer) addCurrentCmd(cmd string, times int) {
	k := t.cacheKey(cmd)
	old, err := t.cacher.Get(k)
	if err == nil {
		times += old.(int)
	}
	log.Debug("add current cmd", "cmd", cmd, "count", times)
	if times <= 0 {
		t.cacher.Remove(k)
	} else {
		t.cacher.Set(k, times)
	}
}
func (t *CommandServer) listenCmd() {
	for {
		select {
		case cmd, ok := <-t.ch:
			if !ok {
				return
			}
			if t.ExistCmd(cmd.Cmd) {
				if err := t.handleCmd(cmd); err != nil {
					log.Debug("handle cmd err", "err", err, "cmd", cmd.Cmd, "from", cmd.From)
				}
			}
		case <-t.quit:
			return
		}
	}
}
func (t *CommandServer) getOrSetUser(weid string) (*model.UserInfo, error) {
	// 获取from user
	user, err := t.userDB.GetByWeId(weid)
	if err != nil {
		if err != dal.ErrNotFound {
			return nil, err
		}
		// set user
		user = &model.UserInfo{
			WeId: weid,
		}
		return user, t.userDB.Add(user)
	}
	return user, nil
}
func (t *CommandServer) getOrSetRoom(roomID string) (*model.RoomInfo, error) {
	room, err := t.roomDB.GetWeId(roomID)
	if err != nil {
		if err != dal.ErrNotFound {
			return nil, err
		}
		// set
		room = &model.RoomInfo{
			WeId: roomID,
		}
		return room, t.roomDB.Add(room)
	}
	return room, nil
}
func (t *CommandServer) handleCmd(cmd *Command) error {
	log.Trace("handle command", "cmd", cmd.Cmd, "args", cmd.Arg, "from", cmd.From, "room", cmd.Room)
	var user *model.UserInfo
	var room *model.RoomInfo
	var err error
	if len(cmd.From) > 0 {
		user, err = t.getOrSetUser(cmd.From)
		if err != nil {
			return err
		}
	}
	if len(cmd.Room) > 0 {
		room, err = t.getOrSetRoom(cmd.Room)
		if err != nil {
			return err
		}
	}
	if room == nil && user == nil {
		return nil
	}

	//按顺序模拟执行
	execFunc := func(actv *model.ActivityInfo) bool {
		if actv.RoomId > 0 {
			if room == nil || room.Id != actv.RoomId {
				return true
			}
		} else {
			if room != nil {
				return true
			}
		}
		if actv.Command != cmd.Cmd {
			return true
		}
		if err := t.actvServer.Join(user, actv); err != nil {
			log.Debug("simulate exec failed", "user", user.Id, "actv", actv.Id, "err", err)
			return true
		}
		return false
	}
	t.actvServer.foreachUnderway(execFunc)
	return nil
}