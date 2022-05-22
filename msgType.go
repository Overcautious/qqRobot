package main

// 返回信息的格式

import (
	"QQRobot/tools"
	"github.com/tencent-connect/botgo/dto"
)

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
