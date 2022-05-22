package tools

import (
	"QQRobot/models"
	"QQRobot/pkg/setting"
	"github.com/tencent-connect/botgo/dto"
	"log"
	"strconv"
	"time"
)

var (
	tableUseridToId   string
	tableIdToUsername string
)

type CheckResult struct {
	State     string `json:"state"`     // 签到是否成功
	CheckRank string `json:"checkRank"` // 签到排名
	CheckNums string `json:"checkNums"` // 累计签到天数
}

func GetCheckByUserIdUserName(userid, username string) *CheckResult {
	var state string
	var checkRank string
	var checkNums string
	if CheckIn(userid, username) {
		state = "1"
	} else {
		state = "0"
	}

	checkRank = strconv.FormatInt(GetTodayCheckNums(), 10)
	checkNums = strconv.FormatInt(GetUserCheckNums(userid), 10)

	checkResult := CheckResult{
		State:     state,
		CheckRank: checkRank,
		CheckNums: checkNums,
	}

	return &checkResult
}

func init() {
	sec, _ := setting.Cfg.GetSection("redis")
	tableUseridToId = sec.Key("TABLE_USERID_TO_ID").String()
	tableIdToUsername = sec.Key("TABLE_ID_TO_USERNAME").String()
}

// CheckIn 签到
func CheckIn(userid string, username string) bool {
	id := models.GetFromHset(userid, tableUseridToId)
	models.SetToHset(strconv.FormatInt(id, 10), username, tableIdToUsername) // 插入到映射表中

	log.Println("id =", id)
	bitmapName := "checkin_" + time.Now().Format("2006-01-02")
	return models.SetBitmap(bitmapName, id)
}

// GetTodayCheckNums 今日签到人数
func GetTodayCheckNums() int64 {
	bitmapName := "checkin_" + time.Now().Format("2006-01-02")
	return models.GetBitmapCount(bitmapName, tableUseridToId)
}

// GetUserCheckNums 统计该用户连续签到天数
func GetUserCheckNums(userid string) int64 {
	var count int64 = 0
	id := models.GetFromHset(userid, tableUseridToId)
	currentTime := time.Now()
	bitmapName := "checkin_" + currentTime.Format("2006-01-02")

	for models.GetBitmap(bitmapName, id) {
		count++
		currentTime = currentTime.AddDate(0, 0, -1)
		bitmapName = "checkin_" + currentTime.Format("2006-01-02")
	}

	return count
}

// GetTopCheckUser 统计累计连续签到次数最多的用户（一直到今天）
func GetTopCheckUser() (int, []string) {
	days, userList := models.GetLongCheckUser(tableUseridToId)
	var userNameList []string
	log.Println("vVVVVV")
	if days > 0 {
		for _, v := range userList {
			log.Println(v)
			userNameList = append(userNameList, models.GetFromHsetString(strconv.FormatInt(v, 10), tableIdToUsername))
		}
	}
	return days, userNameList
}

func (result *CheckResult) CreateArkObjArray() []*dto.ArkObj {

	if result.State == "1" {
		objectArray := []*dto.ArkObj{
			{
				[]*dto.ArkObjKV{
					{
						Key:   "desc",
						Value: "签到成功",
					},
				},
			},
			{
				[]*dto.ArkObjKV{
					{
						Key:   "desc",
						Value: "今日签到人数：" + result.CheckRank,
					},
				},
			},
			{
				[]*dto.ArkObjKV{
					{
						Key:   "desc",
						Value: "您已累计签到：" + result.CheckNums,
					},
				},
			},
		}
		return objectArray
	} else {
		objectArray := []*dto.ArkObj{
			{
				[]*dto.ArkObjKV{
					{
						Key:   "desc",
						Value: "您今日已签到",
					},
				},
			},
			{
				[]*dto.ArkObjKV{
					{
						Key:   "desc",
						Value: "您已累计签到：" + result.CheckNums,
					},
				},
			},
		}
		return objectArray
	}
}

func (result *CheckResult) CreateEmbed() *dto.Embed {
	var state string
	if result.State == "1" {
		state = "签到成功"
		return &dto.Embed{

			Title: "打卡消息：" + state,

			Fields: []*dto.EmbedField{
				{
					Name: "今日签到人数：" + result.CheckRank,
				},
				{
					Name: "您已累计签到天数：" + result.CheckNums,
				},
			},
		}
	} else {
		state = "您今日已签到"
		return &dto.Embed{
			Title: "打卡消息：" + state,
			Fields: []*dto.EmbedField{
				{
					Name: "您已累计签到天数：" + result.CheckNums,
				},
			},
		}
	}
}
