package resp

//--------------backend
type PageSearch struct {
	Index int `json:"index" form:"index"`
	Size  int `json:"size" form:"size"`
}
type LoginUser struct {
	Code     string `json:"code"`
	Password string `json:"password"`
}
type UserSearch struct {
	Page         PageSearch `json:"page"`
	Name         string     `json:"name"`
	IntegralFrom int64      `json:"integralFrom"`
	IntegralTo   int64      `json:"integralTo"`
}
type ActivitySearch struct {
	Page         PageSearch `json:"page"`
	Name         string     `json:"name"`
	Done         int        `json:"done"`
	ActivityType int        `json:"activityType"`
}

//---------------bot
type BotMessage struct {
	User    string   `json:"user"`
	Room    string   `json:"room"`
	Content string   `json:"content"`
	Mention []string `json:"mention"`
}
type WeContact struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

//----------------wechat
type UserCoinRequest struct {
	Coin   string  `json:"coin"`
	Amount float64 `json:"amount"`
}
