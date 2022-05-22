package models

import (
	"QQRobot/pkg/setting"
	"github.com/go-redis/redis"
	"log"
	"os"
	"time"
)

var (
	rdb *redis.Client
)

func init() {
	err := initClient()
	if err != nil {
		log.Println("Redis 连接失败")
		os.Exit(1)
	}

}

// 初始化连接
func initClient() error {
	sec, err := setting.Cfg.GetSection("redis")
	addr := sec.Key("ADDR").String()
	dbNum := sec.Key("DB").MustInt(0)

	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       dbNum, // 默认db
	})

	_, err = rdb.Ping().Result()
	return err
}

//GetFromHset 从user表中获取用户对应的ID
func GetFromHset(key, filed string) int64 {
	if rdb.HExists(filed, key).Val() {
		// 存在该用户
		i, _ := rdb.HGet(filed, key).Int64()
		return i
	} else {
		// 不存在，插入
		length, _ := rdb.HLen(filed).Result()
		log.Println("length =", length)
		rdb.HSet(filed, key, length)
		return length
	}
}

func SetToHset(key, value, filed string) bool {
	_, err := rdb.HSet(filed, key, value).Result()
	if err != nil {
		log.Println("SetToHset 错误, err =", err)
		os.Exit(1)
	}
	return true
}

func GetFromHsetString(key, filed string) string {
	val, _ := rdb.HGet(filed, key).Result()
	return val
}

// GetBitmap bitmap是否设为1
func GetBitmap(key string, offset int64) bool {
	log.Println("bitmap name :", key)
	val, err := rdb.GetBit(key, offset).Result()
	if err != nil {
		log.Println("GetBitmap 错误, err =", err)
		os.Exit(1)
	}
	return val == 1
}

// SetBitmap 参数key： bitmap名， offset 偏移量
func SetBitmap(key string, offset int64) bool {

	if GetBitmap(key, offset) {
		// 已经签到过
		return false
	} else {
		// 进行签到
		rdb.SetBit(key, offset, 1)
		return true
	}
}

func GetBitmapCount(key string, field string) int64 {
	j, err1 := rdb.HLen(field).Result()
	i, err := rdb.BitCount(key, &redis.BitCount{Start: 0, End: j}).Result()
	if err != nil || err1 != nil {
		log.Println("GetBitmapCount 错误, err =", err)
		os.Exit(1)
	}
	return i
}

func GetLongCheckUser(field string) (int, []int64) {
	currentTime := time.Now()
	daysCount := 0
	bitmapNameLast := "checkin_temp1"
	bitmapNameCur := "checkin_temp2"

	bitmapName := "checkin_" + currentTime.Format("2006-01-02")

	rdb.BitOpAnd(bitmapNameCur, bitmapName)
	rdb.BitOpAnd(bitmapNameLast, bitmapName)

	for GetBitmapCount(bitmapNameCur, field) > 0 {
		rdb.BitOpAnd(bitmapNameLast, bitmapName).Result() // 保存上一个状态
		daysCount++

		currentTime = currentTime.AddDate(0, 0, -1)
		bitmapName = "checkin_" + currentTime.Format("2006-01-02")
		rdb.BitOpAnd(bitmapNameCur, bitmapName)
	}

	index, _ := rdb.BitPos(bitmapNameLast, 1).Result()
	indexTop := (index/8 + 1) * 8

	var userList []int64
	for index >= 0 {
		for index < indexTop {
			val, _ := rdb.GetBit(bitmapNameLast, index).Result()
			if val == 1 {
				userList = append(userList, index)
			}
			index++
		}
		index, _ = rdb.BitPos(bitmapNameLast, 1, index/8+1).Result()
	}
	rdb.Del(bitmapNameLast)
	rdb.Del(bitmapNameCur)

	return daysCount, userList
}

// RedisSAdd 向Redis的Set中插入数据
func RedisSAdd(key string, value interface{}) bool {
	result, err := rdb.SAdd(key, value).Result()

	if err != nil {
		log.Println("RedisSAdd 插入错误， err =", err)
	}
	return result == 1
}
