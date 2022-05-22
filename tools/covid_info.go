package tools

import (
	"encoding/json"
	"github.com/tencent-connect/botgo/dto"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type CovidInfoResp struct {
	Info string `json:"info"`
	Data struct {
		Diseaseh5 struct {
			ChinaTotal struct {
				Mtime             string `json:"mtime"`             // 更新时间
				LocalConfirmAdd   int    `json:"localConfirmAdd"`   // 新增本土确诊
				LocalWzzAdd       int    `json:"localWzzAdd"`       // 新增本土无症状
				LocalConfirmH5    int    `json:"localConfirmH5"`    // 现有本土确诊
				NowLocalWzz       int    `json:"nowLocalWzz"`       // 现有本土无症状
				HighRiskAreaNum   int    `json:"highRiskAreaNum"`   // 高风险地区
				MediumRiskAreaNum int    `json:"mediumRiskAreaNum"` // 中风险地区

			} `json:"chinaTotal"`
		} `json:"diseaseh5Shelf"`
	} `json:"data"`
}

func GetCovidInfo() *CovidInfoResp {
	var covid_url = "https://api.inews.qq.com/newsqa/v1/query/inner/publish/modules/list?modules=statisGradeCityDetail,diseaseh5Shelf"
	resp, err := http.Get(covid_url)
	if err != nil {
		log.Fatalln("疫情接口请求异常, err = ", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("天气预报接口数据异常, err = ", err)
		return nil
	}
	var covidResp CovidInfoResp
	err = json.Unmarshal(body, &covidResp)
	if err != nil {
		log.Fatalln("解析数据异常 err = ", err, body)
		return nil
	}
	return &covidResp
}

func (covidInfoResp *CovidInfoResp) CreateArkObjArray() []*dto.ArkObj {
	objectArray := []*dto.ArkObj{
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "今日新增本土确诊：" + strconv.FormatInt(int64(covidInfoResp.Data.Diseaseh5.ChinaTotal.LocalConfirmAdd), 10),
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "今日新增本土无症状：" + strconv.FormatInt(int64(covidInfoResp.Data.Diseaseh5.ChinaTotal.LocalWzzAdd), 10),
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "现有本土确诊：" + strconv.FormatInt(int64(covidInfoResp.Data.Diseaseh5.ChinaTotal.LocalConfirmH5), 10),
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "现有本土无症状：" + strconv.FormatInt(int64(covidInfoResp.Data.Diseaseh5.ChinaTotal.NowLocalWzz), 10),
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "中风险地区数量：" + strconv.FormatInt(int64(covidInfoResp.Data.Diseaseh5.ChinaTotal.MediumRiskAreaNum), 10),
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "高风险地区数量：" + strconv.FormatInt(int64(covidInfoResp.Data.Diseaseh5.ChinaTotal.HighRiskAreaNum), 10),
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "更新时间：" + covidInfoResp.Data.Diseaseh5.ChinaTotal.Mtime,
				},
			},
		},
	}
	return objectArray
}

func (covidInfoResp *CovidInfoResp) CreateEmbed() *dto.Embed {
	return &dto.Embed{
		Title: "疫情日报",
		Fields: []*dto.EmbedField{
			{
				Name: "今日新增本土确诊：" + strconv.FormatInt(int64(covidInfoResp.Data.Diseaseh5.ChinaTotal.LocalConfirmAdd), 10),
			},
			{
				Name: "今日新增本土无症状：" + strconv.FormatInt(int64(covidInfoResp.Data.Diseaseh5.ChinaTotal.LocalWzzAdd), 10),
			},
			{
				Name: "现有本土确诊：" + strconv.FormatInt(int64(covidInfoResp.Data.Diseaseh5.ChinaTotal.LocalConfirmH5), 10),
			},
			{
				Name: "现有本土无症状：" + strconv.FormatInt(int64(covidInfoResp.Data.Diseaseh5.ChinaTotal.NowLocalWzz), 10),
			},
			{
				Name: "中风险地区数量：" + strconv.FormatInt(int64(covidInfoResp.Data.Diseaseh5.ChinaTotal.MediumRiskAreaNum), 10),
			},
			{
				Name: "高风险地区数量：" + strconv.FormatInt(int64(covidInfoResp.Data.Diseaseh5.ChinaTotal.HighRiskAreaNum), 10),
			},
			{
				Name: "更新时间：" + covidInfoResp.Data.Diseaseh5.ChinaTotal.Mtime,
			},
		},
	}
}