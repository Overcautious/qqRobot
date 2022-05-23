package main

// 返回信息的格式

import (
	"QQRobot/tools"
	"github.com/tencent-connect/botgo/dto"
)

var defaultMsg []*dto.ArkObj

func init() {
	defaultMsg = []*dto.ArkObj{
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "小O暂时还不理解您的消息，可以试试其他消息噢，比如：",
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "/当前天气 北京",
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "/私信天气 深圳",
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "/打卡",
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "/打卡统计",
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "/疫情日报",
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "/私信疫情日报",
				},
			},
		},
	}
}

func getDefaultMsgArkForTemplate23() *dto.Ark {
	return &dto.Ark{
		TemplateID: 23,
		KV: []*dto.ArkKV{
			{
				Key:   "#DESC#",
				Value: "描述",
			},
			{
				Key:   "#PROMPT#",
				Value: "#PROMPT",
			},
			{
				Key: "#LIST#",
				Obj: defaultMsg,
			},
		},
	}
}

// 创建23号当Ark
func createArkForTemplate23(resp tools.ToolsResp) *dto.Ark {
	return &dto.Ark{
		TemplateID: 23,
		KV:         createArkKvArray(resp),
	}
}

// 创建Ark需要当ArkKV数组
func createArkKvArray(resp tools.ToolsResp) []*dto.ArkKV {
	arkArray := make([]*dto.ArkKV, 3)
	arkArray[0] = &dto.ArkKV{
		Key:   "#DESC#",
		Value: "描述",
	}
	arkArray[1] = &dto.ArkKV{
		Key:   "#PROMPT#",
		Value: "#PROMPT",
	}
	arkArray[2] = &dto.ArkKV{
		Key: "#LIST#",
		Obj: resp.CreateArkObjArray(),
	}
	return arkArray
}

// 私发消息
func createEmbed(resp tools.ToolsResp) *dto.Embed {
	return resp.CreateEmbed()
}
