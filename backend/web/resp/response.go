package resp

import "github.com/xiaogan18/repe-wechat-assistant/backend/model"

type ResponseResult struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
}
type UserResult struct {
	model.UserInfo
	Title string           `json:"title"`
	Coins []model.UserCoin `json:"coins"`
}
type Activity struct {
	model.ActivityInfo
	CreatBy string `json:"createBy"`
	Room    string `json:"room"`
}
type ActivityLog struct {
	ActivityId int64   `json:"activityId"`
	User       string  `json:"user"`
	Reward     float64 `json:"reward"`
	RewardCoin string  `json:"rewardCoin"`
	JoinTime   string  `json:"joinTime"`
}
type BotTask struct {
	model.BotTask
	Room string `json:"room"`
}
