package main

import (
	"qq_bot/bot"
	"qq_bot/conf"
	"qq_bot/redis"
)

func main() {
	conf.ConfigInit()
	redis.Init(conf.GlobalConfig)
	bot.Start()
}
