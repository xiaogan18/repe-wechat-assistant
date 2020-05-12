package bll

type ContentConfig struct {
	MessageTemplate     string `yaml:"MessageTemplate"`     //消息母版 1.message
	ErrorInvalidArgs    string `yaml:"ErrorInvalidArgs"`    //格式错误时提示 1.description
	ErrorPrivateCommand string `yaml:"ErrorPrivateCommand"` //提醒需要私聊的口令
	Activity            struct {
		EverydayReopenHour   int    `yaml:"EverydayReopenHour"`   //每日活动重置时间（小时）
		EverydayReopenMinute int    `yaml:"EverydayReopenMinute"` //每日活动重置时间（分钟）
		TextNew              string `yaml:"TextNew"`              // 新活动开展文本 1.content 2.reward 3.coinName
		TextClose            string `yaml:"TextClose"`            // 活动结束文本 1.actvName
		TextJoin             string `yaml:"TextJoin"`             // 加入活动文本 1.actvName 2.reward 3.coinName
		TextCoinReturn       string `yaml:"TextCoinReturn"`       // 活动剩余量返回文本 1.coinSum 2.coinName
	} `yaml:"Activity"`
	RedPacket struct {
		Command         string `yaml:"Command"`         //口令
		Description     string `yaml:"Description"`     //格式描述
		Reply           string `yaml:"Reply"`           //回复
		DefaultCommand  string `yaml:"DefaultCommand"`  //默认参加口令
		DefaultDeadTime int64  `yaml:"DefaultDeadTime"` //默认超时时间（秒）
		InsufficientErr string `yaml:"InsufficientErr"` //余额不足时提示
	} `yaml:"RedPacket"` //发红包
	CmdRedPacket CommandContent `yaml:"CmdRedPacket"`
	GiveReward   CommandContent `yaml:"GiveReward"`
	CoinPrice    struct {
		Command     string `yaml:"Command"`     //口令
		Description string `yaml:"Description"` //格式描述
		Reply       string `yaml:"Reply"`       //回复
		SearchErr   string `yaml:"SearchErr"`   //查询失败文本
	} `yaml:"CoinPrice"` //行情
	Topup    CommandContent `yaml:"Topup"`    //充值
	Withdraw CommandContent `yaml:"Withdraw"` //提现
	Account  CommandContent `yaml:"Account"`  //账户
	Balance  CommandContent `yaml:"Balance"`  //余额
	Help     CommandContent `yaml:"Help"`     //帮助
}

type CommandContent struct {
	Command     string `yaml:"Command"`     //口令
	Description string `yaml:"Description"` //格式描述
	Reply       string `yaml:"Reply"`       //回复
}
