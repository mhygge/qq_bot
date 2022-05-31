package bot

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	Constant "qq_bot/const"
	"time"
)

const configFile = "conf/config.yaml"

const (
	BOT  = "Bot"
	BEAR = "Bear"
)

type Author struct {
	Avatar   string `json:"avatar"`
	Bot      bool   `json:"bot"`
	Id       string `json:"id"`
	Username string `json:"username"`
}

type Data struct {
	Id                string          `json:"id"`
	Author            Author          `json:"author"`
	ChannelId         string          `json:"channel_id"`
	Content           string          `json:"content"`
	GuildId           string          `json:"guild_id"`
	Seq               int64           `json:"seq"`
	HeartbeatInterval int64           `json:"heartbeat_interval"`
	SessionId         string          `json:"session_id"`
	Token             string          `json:"token"`
	Intents           Constant.Intent `json:"intents"`
}
type WsPayload struct {
	Opcode int64  `json:"op"`
	Seq    int64  `json:"s"`
	Type   string `json:"t"`
	Data   Data   `json:"d"`
}

var heartbeatInterval int64
var seq int64
var sessionId string
var conn *websocket.Conn
var send chan []byte

// 消息处理器，持有 openapi 对象
func Start() {
	ConnectWs()
	Identify()
	go Heartbeat()
	go Listening()
	select {}
}

func Listening() {
	for true {
		var payload WsPayload
		if err := conn.ReadJSON(&payload); err != nil {
			logrus.Errorf("listen error. %v", err)
			if websocket.IsUnexpectedCloseError(err, 4009) {
				sessionId = ""
				seq = 0
				ConnectWs()
				Identify()
			} else {
				//4009允许重连
				ConnectWs()
				Resume()
			}
		}
		msg, _ := json.Marshal(payload)
		logrus.Infof("event Received: %v.\n", string(msg))
		opSelect(payload)
	}
}

func opSelect(payload WsPayload) {
	switch payload.Opcode {
	case Constant.WSHello:
		go Heartbeat()
		break
	case Constant.WSDispatchEvent:
		// 记录消息序列号，心跳用
		seq = payload.Seq
		eventDispatch(payload)
		break
	case Constant.WSReconnect:
		logrus.Info("重新连接")
		ConnectWs()
		Resume()
		break
	case Constant.WSHeartbeatAck:
		logrus.Info("接收到心跳响应")
		break
	default:
		break
	}
}

func eventDispatch(payload WsPayload) {
	switch payload.Type {
	case Constant.EventReady:
		// 鉴权成功
		sessionId = payload.Data.SessionId
		break
	case Constant.EventGuildCreate:
		break
	case Constant.EventAtMessageCreate:
		SendMessage(payload.Data, Process(payload))
		break
	case "RESUMED":
		logrus.Info("received RESUMED message")
	default:
		break
	}
}

func SendMessage(data Data, content string) {
	body := make(map[string]string)
	body["content"] = fmt.Sprintf("<@!%v>\n", data.Author.Id) + content
	body["msg_id"] = data.Id
	client := resty.New()
	resp, err := client.R().SetAuthToken(GetToken()).
		SetAuthScheme("Bot").
		SetPathParam("channel_id", data.ChannelId).
		SetBody(body).Post(GetURL(messagesURI))
	if err != nil {
		logrus.Errorf("send message error,%v", err)
	}
	logrus.Infof("send message success,%v", resp)
}

func Resume() {
	payload := WsPayload{
		Data: Data{
			SessionId: sessionId,
			Seq:       seq,
			Token:     getBotToken(),
			Intents:   Constant.IntentGuilds | Constant.IntentGuildMembers | Constant.IntentGuildAtMessage,
		},
		Opcode: Constant.WSResume,
	}
	if err := conn.WriteJSON(&payload); err != nil {
		logrus.Errorf("Resume error. %v", err)
		return
	}
	logrus.Infoln("resume success")
}

func Identify() {
	payload := WsPayload{
		Data: Data{
			SessionId: sessionId,
			Seq:       seq,
			Token:     getBotToken(),
			Intents:   Constant.IntentGuilds | Constant.IntentGuildMembers | Constant.IntentGuildAtMessage,
		},
		Opcode: Constant.WSIdentity,
	}
	if err := conn.WriteJSON(&payload); err != nil {
		logrus.Errorf("auth send error. %v", err)
	}
}

func Heartbeat() {

	duration := time.Duration(heartbeatInterval) * time.Millisecond
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	data := make(map[string]int64)
	data["op"] = 1
	for range ticker.C {
		data["d"] = seq
		logrus.Info("ticker ticker ticker ... send Heartbeat\n")
		if err := conn.WriteJSON(data); err != nil {
			logrus.Errorf("Heartbeat send error. %v", err)

		}
	}

}

func ConnectWs() {
	url := getWsUrl()
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	conn = c
	if err != nil {
		logrus.Errorf("const connect error,%v", err)
		os.Exit(-1)
	}
	var payload WsPayload
	err = conn.ReadJSON(&payload)
	if err != nil {
		logrus.Errorf("read json error,%v", err)
	}
	msg, _ := json.Marshal(payload)
	logrus.Infof("const opcode:%v", string(msg))
	heartbeatInterval = payload.Data.HeartbeatInterval
}

//func CloseWs() {
// conn.
//}

type WebsocketAP struct {
	URL string `json:"url"`
}

func getWsUrl() string {
	client := resty.New()
	resp, err := client.R().SetAuthToken(GetToken()).SetAuthScheme(BOT).SetResult(WebsocketAP{}).Get(GetURL(gatewayURI))
	if err != nil {
		logrus.Errorf("request gateway failed,err: %v", err)
		os.Exit(-1)
	}
	return resp.Result().(*WebsocketAP).URL
}

func GetToken() string {
	var conf struct {
		AppID string `yaml:"appid"`
		Token string `yaml:"token"`
	}
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		logrus.Errorf("read token from file failed, err: %v", err)
		os.Exit(-1)
	}
	if err = yaml.Unmarshal(content, &conf); err != nil {
		logrus.Errorf("parse config failed, err: %v", err)
		os.Exit(-1)
	}

	return conf.AppID + "." + conf.Token
}
func getBotToken() string {
	return BOT + " " + GetToken()
}
