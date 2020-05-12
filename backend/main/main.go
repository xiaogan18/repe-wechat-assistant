package main

import (
	"flag"
	"github.com/xiaogan18/repe-wechat-assistant/backend/bll"
	"github.com/xiaogan18/repe-wechat-assistant/backend/dal"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/log"
	"github.com/xiaogan18/repe-wechat-assistant/backend/web"
	"math/rand"
	"time"
)

var args = SetupArgs{
	Log:    flag.Int("log", 3, "log level"),
	Config: flag.String("config", "repe.yaml", "yaml config file path"),
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	conf, err := initConfig(*args.Config)
	if err != nil {
		log.Fatal("init config error", "err", err)
		return
	}
	log.SetLevel(log.Level(*args.Log))
	db, err := dal.NewDbContext(dal.DbConfig{
		DriverName: conf.Database.DriverName,
		DataSource: conf.Database.DataSource,
	})
	if err != nil {
		log.Fatal("init db context error", "err", err)
	}
	botSvr := bll.NewBotTaskServer(db.Bot, db.Room)
	cmdConf := &conf.CommandConfig
	actvSvr := bll.NewActivityServer(botSvr, db, cmdConf)
	cmdSvr := bll.NewCmdServer(db, actvSvr, cmdConf)
	if err = botSvr.Start(); err != nil {
		log.Fatal("start bot server error", "err", err)
	}
	if err = actvSvr.Start(); err != nil {
		log.Fatal("start activity server error", "err", err)
	}
	if err = cmdSvr.Start(); err != nil {
		log.Fatal("start command server error", "err", err)
	}
	weber := web.NewWebEngine(db, cmdSvr, actvSvr, botSvr)
	go func() {
		if err = weber.ListenBot(conf.System.BotListen); err != nil {
			log.Fatal("listen bot error", "err", err)
		}
	}()
	if err = weber.Listen(conf.System.WebListen); err != nil {
		log.Fatal("listen web error", "err", err)
	}
}
