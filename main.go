package main

import (
	"QQRobot/pkg/setting"
	"QQRobot/tools"
	"context"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	"log"
	"os"
	"strconv"
	"time"
)

//Config 定义了配置文件的结构
type Config struct {
	AppID uint64 `yaml:"appid"` //机器人的appid
	Token string `yaml:"token"` //机器人的token
}

var config Config
var api openapi.OpenAPI
var ctx context.Context
var channelId = "" //保存子频道的id

//定义常量
const (
	CmdDirectChatMsg   = "/私信天气"
	CmdNowWeather      = "/当前天气"
	CmdClockIn         = "/打卡"
	CmdCheckInStatist  = "/打卡统计"
	CmdCovidInfo       = "/疫情日报"
	CmdDirectCovidInfo = "/私信疫情日报"
)

//atMessageEventHandler 处理 @机器人 的消息
func atMessageEventHandler(event *dto.WSPayload, data *dto.WSATMessageData) error {
	channelId = data.ChannelID                //当@机器人时，保存ChannelId，主动消息需要 channelId 才能发送出去
	res := message.ParseCommand(data.Content) //去掉@结构和清除前后空格
	log.Println("cmd = " + res.Cmd + " content = " + res.Content)
	cmd := res.Cmd         ///对于像 /私信天气 城市名 指令，cmd 为 私信天气
	content := res.Content //content 为 城市名
	switch cmd {
	case CmdNowWeather: //获取当前天气 指令是 /天气 城市名
		if len(content) == 0 {
			api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{
				MsgID:   data.ID,
				Content: "请输入查询天气城市",
			})
		}
		webData := tools.GetWeatherByCity(content)
		if webData != nil {
			//MsgID 表示这条消息的触发来源，如果为空字符串表示主动消息
			//Ark 传入数据时表示发送的消息是Ark
			api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Ark: createArkForTemplate23(webData)})
		}
	case CmdDirectChatMsg: //私信天气消息到用户
		webData := tools.GetWeatherByCity(content)
		if webData != nil {
			//创建私信会话
			directMsg, err := api.CreateDirectMessage(ctx, &dto.DirectMessageToCreate{
				SourceGuildID: data.GuildID,
				RecipientID:   data.Author.ID,
			})
			if err != nil {
				log.Println("私信创建出错了，err = ", err)
			}
			//发送私信消息
			//Embed 传入数据时表示发送的是 Embed
			api.PostDirectMessage(ctx, directMsg, &dto.MessageToCreate{Embed: createEmbed(webData)})
		}
	case CmdClockIn:
		userId := data.Author.ID // 获取用户ID，string类型
		userName := data.Author.Username
		//log.Println("userID =", userId)
		checkData := tools.GetCheckByUserIdUserName(userId, userName)
		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Ark: createArkForTemplate23(checkData)})
	case CmdCovidInfo:
		covidData := tools.GetCovidInfo()
		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Ark: createArkForTemplate23(covidData)})

	case CmdDirectCovidInfo:
		covidData := tools.GetCovidInfo()
		if covidData != nil {
			//创建私信会话
			directMsg, err := api.CreateDirectMessage(ctx, &dto.DirectMessageToCreate{
				SourceGuildID: data.GuildID,
				RecipientID:   data.Author.ID,
			})
			if err != nil {
				log.Println("私信创建出错了，err = ", err)
			}
			//发送私信消息
			//Embed 传入数据时表示发送的是 Embed
			api.PostDirectMessage(ctx, directMsg, &dto.MessageToCreate{Embed: createEmbed(covidData)})
		}

	case CmdCheckInStatist:
		checkNums := tools.GetTodayCheckNums() // 今日签到人数
		days, userNameList := tools.GetTopCheckUser()
		var userList string
		//log.Println("userNameList.length =", len(userNameList))
		for i, v := range userNameList {
			if i < len(userNameList)-1 {
				userList += v + "，"
			} else {
				userList += v
			}
		}
		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{
			MsgID:   data.ID,
			Content: "今日签到人数：" + strconv.FormatInt(checkNums, 10) + " \n" + "最长连续打卡次数：" + strconv.FormatInt(int64(days), 10) + " \n" + "打卡冠军为：" + userList,
		})

	default:
		api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{MsgID: data.ID, Ark: getDefaultMsgArkForTemplate23()})

	}
	return nil
}

//第一步： 获取机器人的配置信息，即机器人的appid和token
func init() {
	sec, err := setting.Cfg.GetSection("robot")
	if err != nil {
		log.Println("读取配置文件出错， err = ", err)
		os.Exit(1)
	}
	config.AppID = uint64(sec.Key("appid").MustUint())
	config.Token = sec.Key("token").String()
	log.Println(config)
}

func main() {

	//第二步：生成token，用于校验机器人的身份信息
	token := token.BotToken(config.AppID, config.Token)
	//第三步：获取操作机器人的API对象
	api = botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
	//获取context
	ctx = context.Background()
	//第四步：获取websocket
	ws, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln("websocket错误， err = ", err)
	}

	var atMessage event.ATMessageEventHandler = atMessageEventHandler

	intent := websocket.RegisterHandlers(atMessage)     // 注册处理函数
	botgo.NewSessionManager().Start(ws, token, &intent) // 启动socket监听
}
