package tools

import "github.com/tencent-connect/botgo/dto"

type ToolsResp interface {
	CreateArkObjArray() []*dto.ArkObj
	CreateEmbed() *dto.Embed
}
