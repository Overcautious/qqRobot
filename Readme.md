
# 简介
该项目为QQ频道机器人，使用前需要在 /conf/app.ini 中修改QQ机器人提供的`appid`和`token`

目前提供的接口有：
```
CmdDirectChatMsg   = "/私信天气"
CmdNowWeather      = "/当前天气"
CmdClockIn         = "/打卡"
CmdCheckInStatist  = "/打卡统计"
CmdCovidInfo       = "/疫情日报"
CmdDirectCovidInfo = "/私信疫情日报"
```

## Todo:
- [ ] 提供百度检索接口
- [ ] 提供图片搜索接口
- [ ] 提供点歌接口（分享歌曲链接至群）
- [ ] 歌词识曲

以及一些互动功能
- [ ] 棋盘室（提供象棋、五子棋，供特定频道使用，人人对战，人机对战）
- [ ] 古诗词背诵（机器人提问，用户选择选项）


# 实现功能：

## 1. 查询日期功能
#### 使用方法：
```
@腾讯机器人 /当前天气 西安
@腾讯机器人 /私信天气 西安
```
支持共屏返回天气结果和私信返回天气结果，根据城市不同，返回不同的城市的当日天气

#### 实现设计
利用[NowApi](https://www.nowapi.com/api)调用查询接口，将查询返回值作为响应结果返回

## 2. 打卡功能
#### 使用方法
```
@腾讯机器人 /打卡
@腾讯机器人 /打卡统计
```
- 若当天打卡，则会返回 "签到成功" 字样。并且返回今日签到人数和用户累计签到天数
- 若当天已打卡，则会返回 "今日已签到" 字样。并且返回今日签到人数和用户累计签到天数
- 打卡统计，会返回今天打卡人数，已经历史中打卡连续次数最多的用户
#### 实现设计
利用bitmap记录打卡状态，1为打卡成功，并且每天都会有一个新的bitmap来记录所有用户的打卡状态。
将多天的打卡记录bitmap进行与操作，操作后依旧为1的位置，表示该用户在这几天都是打卡状态

## 3. 疫情信息查询功能
#### 使用方法
```
@腾讯机器人 /疫情日报
@腾讯机器人 /私信疫情日报
```
返回本土今日新增确诊、新增无症状、以及累计确诊病例
#### 实现设计
使用api接口，请求疫情信息，并返回