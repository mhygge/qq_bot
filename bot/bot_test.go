package bot

import (
	"bou.ke/monkey"
	"github.com/sirupsen/logrus"
	Constant "qq_bot/const"
	"testing"
	"time"
)

var payloadTests = []WsPayload{
	{
		Opcode: Constant.WSHello,
	},
	{
		Opcode: Constant.WSReconnect,
		Seq:    0,
		Type:   "",
		Data: Data{
			Id: "",
			Author: Author{
				Avatar:   "",
				Bot:      false,
				Id:       "",
				Username: "",
			},
			ChannelId:         "",
			Content:           "",
			GuildId:           "",
			Seq:               0,
			HeartbeatInterval: 0,
			SessionId:         "111",
			Token:             "",
			Intents:           0,
		},
	},
	{
		Opcode: Constant.WSDispatchEvent,
		Seq:    0,
		Type:   "RESUMED",
		Data: Data{
			Id: "",
			Author: Author{
				Avatar:   "",
				Bot:      false,
				Id:       "",
				Username: "",
			},
			ChannelId:         "",
			Content:           "",
			GuildId:           "",
			Seq:               0,
			HeartbeatInterval: 0,
			SessionId:         "",
			Token:             "",
			Intents:           0,
		},
	},
	{
		Opcode: Constant.WSDispatchEvent,
		Seq:    0,
		Type:   Constant.EventAtMessageCreate,
		Data: Data{
			Id: "",
			Author: Author{
				Avatar:   "",
				Bot:      false,
				Id:       "",
				Username: "",
			},
			ChannelId:         "",
			Content:           "",
			GuildId:           "",
			Seq:               0,
			HeartbeatInterval: 0,
			SessionId:         "",
			Token:             "",
			Intents:           0,
		},
	},
	//{
	//	Opcode: 0,
	//	Seq:    0,
	//	Type:   "",
	//	Data:   Data{
	//		Id:                "",
	//		Author:            Author{
	//			Avatar:   "",
	//			Bot:      false,
	//			Id:       "",
	//			Username: "",
	//		},
	//		ChannelId:         "",
	//		Content:           "",
	//		GuildId:           "",
	//		Seq:               0,
	//		HeartbeatInterval: 0,
	//		SessionId:         "",
	//		Token:             "",
	//		Intents:           0,
	//	},
	//},
	//{
	//	Opcode: 0,
	//	Seq:    0,
	//	Type:   "",
	//	Data:   Data{
	//		Id:                "",
	//		Author:            Author{
	//			Avatar:   "",
	//			Bot:      false,
	//			Id:       "",
	//			Username: "",
	//		},
	//		ChannelId:         "",
	//		Content:           "",
	//		GuildId:           "",
	//		Seq:               0,
	//		HeartbeatInterval: 0,
	//		SessionId:         "",
	//		Token:             "",
	//		Intents:           0,
	//	},
	//},
}

//函数太简单没必要单测
func TestOpSelect(t *testing.T) {
	for _, payload := range payloadTests {
		monkey.Patch(Heartbeat, func() {
			duration := time.Duration(1000) * time.Millisecond
			ticker := time.NewTicker(time.Duration(duration))
			defer ticker.Stop()
			for range ticker.C {
				logrus.Info("ticker ticker ticker ... send Heartbeat")
			}
		})
		monkey.Patch(ConnectWs, func() {
			logrus.Infof("const opcode: 10")
		})
		monkey.Patch(Resume, func() {
			logrus.Infoln("resume success")
		})

		OpSelect(payload)
	}

}
