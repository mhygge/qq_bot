package bot

import (
	"fmt"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
	Constant "qq_bot/const"
	client "qq_bot/redis"
	"testing"
	"time"
)

var signPayloads = []WsPayload{
	{
		Opcode: Constant.WSDispatchEvent,
		Type:   Constant.EventAtMessageCreate,
		Data: Data{
			Author: Author{
				Id:       "1",
				Username: "user1",
			},
			Content: "签到",
		},
	},
	{
		Opcode: Constant.WSDispatchEvent,
		Type:   Constant.EventAtMessageCreate,
		Data: Data{
			Author: Author{
				Id:       "2",
				Username: "user2",
			},
			Content: "签到",
		},
	},
}
var user1Payload = signPayloads[0]
var user2Payload = signPayloads[1]

var rankPayloads = []WsPayload{
	{
		Opcode: Constant.WSDispatchEvent,
		Type:   Constant.EventAtMessageCreate,
		Data: Data{
			Author: Author{
				Id:       "1",
				Username: "user1",
			},
			ChannelId: "",
			Content:   "周排行",
		},
	},
	{
		Opcode: Constant.WSDispatchEvent,
		Type:   Constant.EventAtMessageCreate,
		Data: Data{
			Author: Author{
				Id:       "1",
				Username: "user1",
			},
			ChannelId: "",
			Content:   "日排行",
		},
	},
}

var weekRankPayload = rankPayloads[0]
var dailyRankPayload = rankPayloads[1]

func TestProcess(t *testing.T) {
	t.Run("签到", func(t *testing.T) {
		s, err := miniredis.Run()
		if err != nil {
			panic(err)
		}
		defer s.Close()

		mockClient := redis.NewClient(&redis.Options{
			Addr: s.Addr(),
		})
		stub := gostub.Stub(&client.GlobalRedis, mockClient)
		defer stub.Reset()
		reply := Process(signPayloads[0])
		assert.Equal(t, "签到成功,本月共签到: 1次", reply)
		//再次签到
		reply2 := Process(signPayloads[0])
		assert.Equal(t, "今日已签到，请勿重复操作", reply2)
	})
	t.Run("排行", func(t *testing.T) {

		s, err := miniredis.Run()
		if err != nil {
			panic(err)
		}
		defer s.Close()

		mockClient := redis.NewClient(&redis.Options{
			Addr: s.Addr(),
		})
		stub := gostub.Stub(&client.GlobalRedis, mockClient)
		defer stub.Reset()
		//user1 今天减2,3,6天打卡，user2今天减1,2,4,6天打卡
		signDays := [][]int{{2, 3, 6}, {1, 2, 4, 6}}
		for i, days := range signDays {
			for j, day := range days {
				nowFn = func() time.Time {
					now := time.Now()
					//今天：2022-06-11
					return time.Date(2022, 6, 11, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), time.Local).AddDate(0, 0, -1*day)
				}
				reply := Process(signPayloads[i])
				//此断言不跨月所以选6-11
				assert.Equal(t, fmt.Sprintf("签到成功,本月共签到: %v次", j+1), reply)
			}
		}
		//现在是11-6=5 即2022-06-05
		weeklyRankReply := Process(weekRankPayload)
		assert.Equal(t, "user2,user1", weeklyRankReply)
		dailyRankReply := Process(dailyRankPayload)
		assert.Equal(t, "user1,user2", dailyRankReply)
	})
	//本周统计中 BitMap操作 miniRedis不支持
}
