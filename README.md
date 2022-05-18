# QQ 频道签到机器人

## 使用说明

1. 进入频道聊天室
2. @moe-测试中 签到，非”签到”消息回复提示：“请输入‘签到’”
3. 如当天已签到，回复：”今日已签到，请勿重复操作“，如未签到，执行签到逻辑，回复签到成功，和当月签到天数

## 方案设计
### 技术栈
| Websocket | RESTful | Log    | DB    |
|-----------| ------- |--------|-------|
| Gorilla   | Resty   | Logrus | Redis |
### 工作流
1. bot 向 QQ开放平台gateway发起http请求.
2. 开发平台升级 http 到 websocket 并返回 websocket url
3. bot 拿到 url 结合botToken，intents进行鉴权，
4. 开始发送心跳
5. 监听websocket,根据类型（重连，心跳，事件分发）处理消息。
6. 事件分发主要分为鉴权成功（Ready)和其他业务事件，这里是@消息事件
7. 进入签到业务逻辑部分，使用Redis的BitMap类型做存储，userId,年，月拼接成key，0代表未签到，1代表已签到
8. 根据key和offset（多少号）查询bit位值，为1则回复已签到，为0则执行签到逻辑设为1，并 count当月已签到天数，回复用户签到成功






