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
	c.AddFunc("@daily", func() {
		redis.GlobalRedis.Del(ctx, "today")
	})
	c.AddFunc("@weekly", func() {
		redis.GlobalRedis.Del(ctx, "week")
	})

}
