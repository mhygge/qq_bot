package bot

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"qq_bot/redis"
	"strings"
	"time"
)

const (
	singKey   = "user:%s:%d:%d" //用户打卡信息 %s=发送方uid %d=年 %d=月 bitmap
	botAtInfo = "<@!9681314665609974950> "
)

func Process(payload WsPayload) string {
	content := strings.Replace(payload.Data.Content, botAtInfo, "", -1)
	if content != "签到" {
		return "请输入\"签到\""
	}
	ctx := context.Background()
	now := time.Now()
	key := fmt.Sprintf(singKey, payload.Data.Author.Id, now.Year(), now.Month())
	offset := int64(now.Day())
	bit, err := redis.GlobalRedis.GetBit(ctx, key, offset).Uint64()
	if err != nil {
		logrus.Errorf("get redis key error,%v", err)
		return "server error"
	}
	if bit == 1 {
		return "今日已签到，请勿重复操作"
	}
	_, err = redis.GlobalRedis.SetBit(ctx, key, offset, 1).Result()
	if err != nil {
		logrus.Errorf("SetBit fail uid=%s|key=%s|offset=%d|err=%s", payload.Data.Author.Id, key, offset, err)
		return "server error"
	}

	count, err := redis.GlobalRedis.BitCount(ctx, key, nil).Result()
	if err != nil {
		logrus.Errorf("BitCount fail key=%s|err=%s", key, err)
		return "server error"
	}
	return fmt.Sprintf("签到成功,本月共签到: %v次", count)
}
