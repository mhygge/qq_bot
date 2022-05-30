package main

import (
	"context"
	"github.com/robfig/cron/v3"
	"qq_bot/bot"
	"qq_bot/conf"
	"qq_bot/redis"
)

func main() {
	conf.ConfigInit()
	redis.Init(conf.GlobalConfig)
	bot.Start()
	ctx := context.Background()
	c := cron.New(cron.WithSeconds())
	dailySpec := "@daily"
	c.AddFunc(dailySpec, func() {
		redis.GlobalRedis.Del(ctx, "today")
	})
	weeklySpec := "@weekly"
	c.AddFunc(weeklySpec, func() {
		redis.GlobalRedis.Del(ctx, "week")
	})

}
