package model

import "time"

// 用户信息
type UserInfo struct {
	Id          int64  `json:"id"`
	WeId        string `json:"weId" xorm:"'we_id' varchar(50) notnull unique index"` // 微信id
	WeName      string `json:"weName" xorm:"'we_name'"`                              // 微信名
	BindAccount string `json:"bindAccount" xorm:"'bind_account' varchar(50)"`        // 绑定的交易所账户
	Integral    int64  `json:"integral" xorm:"'integral'"`                           // 积分
}

// 用户持币信息
type UserCoin struct {
	Id      int64   `json:"id"`
	UserId  int64   `json:"userId" xorm:"'user_id' notnull"`         //user id
	Coin    string  `json:"coin" xorm:"'coin' varchar(50) notnull"`  //币类型
	Balance float64 `json:"balance" xorm:"'balance' decimal(22,12)"` // 余额
}

// 机器人任务
type BotTask struct {
	Id       int64     `json:"id"`
	RoomId   int64     `json:"roomId" xorm:"'room_id'"`     // 群聊id
	TaskType uint8     `json:"taskType" xorm:"'task_type'"` // 任务类型 enum> BotTaskType
	TaskName string    `json:"taskName" xorm:"'task_name'"` // 任务命名
	TaskTime string    `json:"taskTime" xorm:"'task_time'"` // 定时任务执行时间 HH:mm:ss
	Content  string    `json:"content" xorm:"'content'"`    // 任务内容（机器人发送的文本）
	Done     uint8     `json:"done" xorm:"'done'"`          // 是否已完成 enum> TaskDone
	DoneTime time.Time `json:"doneTime" xorm:"'done_time'"` // 完成时间
}

// 群聊信息
type RoomInfo struct {
	Id     int64  `json:"id"`
	WeId   string `json:"weId" xorm:"'we_id' varchar(50) notnull"` // 微信标识
	WeName string `json:"weName" xorm:"'we_name' notnull"`         // 微信名
	Remark string `json:"remark" xorm:"'remark'"`                  // 备注
}

// 活动信息
type ActivityInfo struct {
	Id              int64     `json:"id"`
	Name            string    `json:"name" xorm:"'name'"`                         // 活动命名
	Content         string    `json:"content" xorm:"'content'"`                   // 活动内容文本
	ActivityType    uint8     `json:"activityType" xorm:"'activity_type'"`        // 活动类型 enum> ActivityType
	RoomId          int64     `json:"roomId" xorm:"'room_id'"`                    // 群聊id
	Capacity        int       `json:"capacity" xorm:"'capacity'"`                 // 可最大参与人数
	Joined          int       `json:"joined" xorm:"'joined'"`                     // 参与的人数
	CoinType        string    `json:"coinType" xorm:"'coin_type'"`                // 奖励币种
	CoinSum         float64   `json:"coinSum" xorm:"'coin_sum' decimal(22,12)"`   // 奖励总币数
	CoinRest        float64   `json:"coinRest" xorm:"'coin_rest' decimal(22,12)"` // 奖励剩余量
	RewardType      uint8     `json:"rewardType" xorm:"'reward_type'"`            // 奖励分配方式 enum> RewardType
	IntegralRequire int64     `json:"integralRequire" xorm:"'integral_require'"`  // 参与需要积分门槛
	IntegralCost    int64     `json:"integralCost" xorm:"'integral_cost'"`        // 参与消耗积分
	Command         string    `json:"command" xorm:"'command'"`                   // 参与口令关键字
	Deadtime        int64     `json:"deadtime" xorm:"'deadtime'"`                 // 自动结束时间 s
	Done            uint8     `json:"done" xorm:"'done'"`                         // 是否结束 enum> TaskDone
	DoneTime        time.Time `json:"doneTime" xorm:"'done_time'"`                // 结束时间
	CreateBy        int64     `json:"createBy" xorm:"'create_by'"`                // 创建人
	CreateTime      time.Time `json:"createTime" xorm:"'create_time'"`            //创建时间
}

// 活动参与记录
type ActivityLog struct {
	Id         int64     `json:"id"`
	ActivityId int64     `json:"activityId" xorm:"'activity_id'"`
	UserId     int64     `json:"userId" xorm:"'user_id'"`
	RewardCoin string    `json:"rewardCoin" xorm:"'reward_coin'"`       // 奖励币种
	Reward     float64   `json:"reward" xorm:"'reward' decimal(22,12)"` // 奖励币数
	JoinTime   time.Time `json:"joinTime" xorm:"'join_time'"`           // 参与时间
}