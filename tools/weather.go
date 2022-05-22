package tools

// 获取天气请求相关
import (
	"encoding/json"
	"github.com/tencent-connect/botgo/dto"
	"io/ioutil"
	"log"
	"net/http"
)

//WeatherResp 定义了返回天气数据的结构
type WeatherResp struct {
	Success    string `json:"success"` //标识请求是否成功，0表示成功，1表示失败
	ResultData Result `json:"result"`  //请求成功时，获取的数据
	Msg        string `json:"msg"`     //请求失败时，失败的原因
}

//Result 定义了具体天气数据结构
type Result struct {
	Days            string `json:"days"`             //日期，例如2022-03-01
	Week            string `json:"week"`             //星期几
	CityNm          string `json:"citynm"`           //城市名
	Temperature     string `json:"temperature"`      //当日温度区间
	TemperatureCurr string `json:"temperature_curr"` //当前温度
	Humidity        string `json:"humidity"`         //湿度
	Weather         string `json:"weather"`          //天气情况
	Wind            string `json:"wind"`             //风向
	Winp            string `json:"winp"`             //风力
	TempHigh        string `json:"temp_high"`        //最高温度
	TempLow         string `json:"temp_low"`         //最低温度
	WeatherIcon     string `json:"weather_icon"`     //气象图标
}

//GetWeatherByCity 获取对应城市的天气数据
func GetWeatherByCity(cityName string) *WeatherResp {
	url := "http://api.k780.com/?app=weather.today&cityNm=" + cityName + "&appkey=65826&sign=169e5e8f3531633d6085daba6993ee2c&format=json"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("天气预报接口请求异常, err = ", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("天气预报接口数据异常, err = ", err)
		return nil
	}
	var weatherData WeatherResp
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		log.Fatalln("解析数据异常 err = ", err, body)
		return nil
	}
	if weatherData.Success != "1" {
		log.Fatalln("返回数据问题 err = ", weatherData.Msg)
		return nil
	}
	return &weatherData
}

func (weather *WeatherResp) CreateArkObjArray() []*dto.ArkObj {
	objectArray := []*dto.ArkObj{
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: weather.ResultData.CityNm + " " + weather.ResultData.Weather + " " + weather.ResultData.Days + " " + weather.ResultData.Week,
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "当日温度区间：" + weather.ResultData.Temperature,
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "当前温度：" + weather.ResultData.TemperatureCurr,
				},
			},
		},
		{
			[]*dto.ArkObjKV{
				{
					Key:   "desc",
					Value: "当前湿度：" + weather.ResultData.Humidity,
				},
			},
		},
	}
	return objectArray
}

func (weather *WeatherResp) CreateEmbed() *dto.Embed {
	return &dto.Embed{
		Title: weather.ResultData.CityNm + " " + weather.ResultData.Weather,
		Thumbnail: dto.MessageEmbedThumbnail{
			URL: weather.ResultData.WeatherIcon,
		},
		Fields: []*dto.EmbedField{
			{
				Name: weather.ResultData.Days + " " + weather.ResultData.Week,
			},
			{
				Name: "当日温度区间：" + weather.ResultData.Temperature,
			},
			{
				Name: "当前温度：" + weather.ResultData.TemperatureCurr,
			},
			{
				Name: "最高温度：" + weather.ResultData.TempHigh,
			},
			{
				Name: "最低温度：" + weather.ResultData.TempLow,
			},
			{
				Name: "当前湿度：" + weather.ResultData.Humidity,
			},
		},
	}
}
