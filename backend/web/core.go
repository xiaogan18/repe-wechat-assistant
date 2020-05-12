package web

import (
	"errors"
	"fmt"
	"github.com/xiaogan18/repe-wechat-assistant/backend/bll"
	"github.com/xiaogan18/repe-wechat-assistant/backend/dal"
	"github.com/xiaogan18/repe-wechat-assistant/backend/tools/log"
	"github.com/xiaogan18/repe-wechat-assistant/backend/web/resp"
	"github.com/astaxie/beego"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	errInvalidParams = errors.New("invalid params")
)

type WebEngine struct {
	eng       *gin.Engine
	db        *dal.DbContext
	cmdSvr    *bll.CommandServer
	actvSvr   *bll.ActivityServer
	taskSvr   *bll.BotTaskServer
}

func NewWebEngine(db *dal.DbContext, cmdSvr *bll.CommandServer, actvSvr *bll.ActivityServer, taskSvr *bll.BotTaskServer) *WebEngine {
	return &WebEngine{
		db:        db,
		cmdSvr:    cmdSvr,
		actvSvr:   actvSvr,
		taskSvr:   taskSvr,
	}
}
func (t *WebEngine) Listen(url string) error {
	eng := gin.Default()
	t.eng = eng
	// 部署vue
	execPath := getCurrentDirectory()
	if _, err := os.Stat(execPath + "/dist"); err == nil {
		eng.LoadHTMLGlob(execPath + "/dist/*.html")           // 添加入口index.html
		eng.Static("/repe/static", execPath+"/dist/static")   // 添加资源路径
		eng.StaticFile("/repe/", execPath+"/dist/index.html") //前端接口
	}
	eng.LoadHTMLGlob(execPath + "/tmpl/*")

	register(eng,"/example",exampleActv{})
	brouter := eng.Group("/b")
	register(brouter, "/user", userManagerApi{db: t.db})
	register(brouter, "/actv", actvApi{db: t.db, actvServer: t.actvSvr})
	register(brouter, "/room", roomApi{db: t.db})
	register(brouter, "/task", taskApi{db: t.db, taskSvr: t.taskSvr})
	return eng.Run(url)
}
func (t *WebEngine) ListenBot(url string) error {
	log.Debug("listen bot server", "url", url)
	eng := gin.New()
	register(eng, "/bot", botApi{db: t.db, cmdSvr: t.cmdSvr, taskSvr: t.taskSvr})
	// FIXME should be delete
	eng.LoadHTMLFiles("tmpl/example.tmpl")
	register(eng, "/example", exampleApi{db: t.db})
	return eng.Run(url)
}

func register(eng gin.IRouter, path string, api IRouter) {
	gp := eng.Group(path)
	api.Register(gp)
}
func assertError(ctx *gin.Context, err error) bool {
	if err != nil {
		rd := ctx.Request.Body
		body, _ := ioutil.ReadAll(rd)
		fmt.Println(string(body))
		log.Debug("handle request error", "err", err, "url", ctx.Request.URL)
		ctx.String(http.StatusBadRequest, err.Error())
		return false
	}
	return true
}
func success(ctx *gin.Context) {
	ctx.String(http.StatusOK, "success")
}
func successData(ctx *gin.Context, data interface{}, total int64) {
	result := &resp.ResponseResult{
		Code:    1,
		Message: "",
		Data:    data,
		Total:   total,
	}
	ctx.JSON(http.StatusOK, result)
}
func errorData(ctx *gin.Context, data interface{}, message string) {
	result := &resp.ResponseResult{
		Code:    -1,
		Message: message,
		Data:    data,
		Total:   0,
	}
	ctx.JSON(http.StatusOK, result)
}

type IRouter interface {
	Register(engine *gin.RouterGroup)
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		beego.Debug(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
