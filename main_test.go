package main

import (
	"qq_bot/bot"
	"testing"
)

func TestBot(t *testing.T) {
	t.Run("getToken", func(t *testing.T) {
		token := bot.GetToken()
		t.Log(token)
	})
	t.Run("connectWs", func(t *testing.T) {
		bot.ConnectWs()
	})
	t.Run("identify", func(t *testing.T) {
		bot.ConnectWs()
		bot.Identify()
	})
	t.Run("resume", func(t *testing.T) {
		bot.ConnectWs()
		//心跳为41秒
		//time.Sleep(1 * time.Minute)
		bot.Resume()
		bot.Listening()
	})
	t.Run("heartBeat", func(t *testing.T) {
		bot.ConnectWs()
		bot.Heartbeat()
	})
	t.Run("sendMessage", func(t *testing.T) {
		bot.ConnectWs()
		//bot.SendMessage( )
	})
}
