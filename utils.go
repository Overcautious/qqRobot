package main

import (
	"github.com/robfig/cron"
	"github.com/tencent-connect/botgo/dto"
)

//注册定时器
func registerMsgPush() {
	var activeMsgPush = func() {
		if channelId != "" {
			// MsgID为空，表示主动消息
			api.PostMessage(ctx, channelId, &dto.MessageToCreate{MsgID: "", Content: "当前天气是：晴天"})
		}
	}
	timer := cron.New()
	//cron表达式由6部分组成，从左到右分别表示 秒 分 时 日 月 星期
	//*表示任意值  ？表示不确定值，只能用于星期和日
	//这里表示每天15:53分发送消息
	timer.AddFunc("0 53 15 * * ?", activeMsgPush)
	timer.Start()
}
