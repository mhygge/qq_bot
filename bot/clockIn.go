package bot

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	client "qq_bot/redis"
	"strconv"
	"strings"
	"time"
)

const (
	singKey   = "user:%s:%d:%d" //用户打卡信息 %s=发送方uid %d=年 %d=月 bitmap
	botAtInfo = "<@!9681314665609974950> "
)

//解耦
var nowFn = time.Now

func Process(payload WsPayload) string {
	content := strings.Replace(payload.Data.Content, botAtInfo, "", -1)
	//if content != "签到" {
	//	return "请输入\"签到\""
	//}
	ctx := context.Background()
	now := nowFn()
	key := fmt.Sprintf(singKey, payload.Data.Author.Id, now.Year(), now.Month())
	if content == "签到" {

		offset := int64(now.Day())
		bit, err := client.GlobalRedis.GetBit(ctx, key, offset).Uint64()
		if err != nil {
			logrus.Errorf("get redis key error,%v", err)
			return "server error"
		}
		if bit == 1 {
			return "今日已签到，请勿重复操作"
		}
		_, err = client.GlobalRedis.SetBit(ctx, key, offset, 1).Result()
		if err != nil {
			logrus.Errorf("SetBit fail uid=%s|key=%s|offset=%d|err=%s", payload.Data.Author.Id, key, offset, err)
			return "server error"
		}

		count, err := client.GlobalRedis.BitCount(ctx, key, nil).Result()
		if err != nil {
			logrus.Errorf("BitCount fail key=%s|err=%s", key, err)
			return "server error"
		}
		client.GlobalRedis.ZAddNX(ctx, "today", &redis.Z{
			Score:  float64(time.Now().Unix()),
			Member: payload.Data.Author.Username,
		})
		//位次
		rank := client.GlobalRedis.ZRank(ctx, "today", payload.Data.Author.Username).Val()
		client.GlobalRedis.ZIncrBy(ctx, "week", 1-float64(rank)/100000, payload.Data.Author.Username)
		return fmt.Sprintf("签到成功,本月共签到: %v次", count)
	} else if content == "本周统计" {
		monday := now.Day() - int(now.Weekday())

		field := client.GlobalRedis.BitField(ctx, key, "GET", "u7", monday)
		formatInt := strconv.FormatInt(field.Val()[0], 2)
		//补0到7位
		if len(formatInt) < 7 {
			for i := 0; i < 7-len(formatInt); i++ {
				formatInt = "0" + formatInt
			}
		}
		var signedDays []string
		var unSignedDays []string
		for i := 0; i < len(formatInt); i++ {
			if formatInt[i] == '1' {
				signedDays = append(signedDays, strconv.Itoa(i))
			} else {
				unSignedDays = append(unSignedDays, strconv.Itoa(i))
			}
		}
		return fmt.Sprintf("本周第%s天打卡，第%s天未打卡", strings.Join(signedDays, ","), strings.Join(unSignedDays, ","))
	} else if content == "日排行" {
		result, err := client.GlobalRedis.ZRange(ctx, "today", 0, 9).Result()
		if err != nil {
			return ""
		}
		return strings.Join(result, ",")
	} else if content == "周排行" {
		result, err := client.GlobalRedis.ZRevRange(ctx, "week", 0, 9).Result()
		if err != nil {
			return ""
		}
		return strings.Join(result, ",")
	} else {
		return "请输入\"签到\"或\"本周统计\"或\"日排行\"或\"周排行\""
	}

}
