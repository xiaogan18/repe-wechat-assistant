package model

// 任务类型
type BotTaskType uint8

const (
	BotTaskImmediate BotTaskType = iota // 即时任务
	BotTaskTiming                       // 定时任务
	BotTaskEveryday                     // 每日任务
)

// 完成状态
type TaskDone uint8

const (
	TaskDoneNot     TaskDone = iota //未完成
	TaskDoneAlready                 // 已完成
)

// 活动类型
type ActivityType uint8

const (
	ActivityTypeOnce     ActivityType = iota // 一次性活动
	ActivityTypeEveryday                     // 每日活动
	ActivityTypeLong                         // 长期活动
)

// 奖励方式
type RewardType uint8

const (
	RewardTypeRandom RewardType = iota // 随机奖励
	RewardTypeAvg                      // 平均奖励
)

// 管理员等级
type SystemUserLevel uint8

const (
	SystemUserLevelAdmin  SystemUserLevel = iota // 管理员
	SystemUserLevelOption                        // 操作员
)

const (
	CoinIntegral = "积分"
	ModuleCoin   = "coins"
)

type TransferStatus uint8

const (
	TransferStatusUnderway TransferStatus = iota
	TransferStatusSuccess
	TransferStatusFailed
)

type TransferType uint8

const (
	TransferTypeTopup TransferType = iota
	TransferTypeWithdraw
)
