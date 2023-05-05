package handler

import (
	"encoding/json"

	"github.com/revel/revel"
	"xonesoft.cn/chatrobot/app/api"
)

type MessageReceiveEvent struct {
	Sender  MessageSender `json:"sender"`
	Message MessageDetail `json:"message"`
}

type MessageSender struct {
	SendId    MessageSenderID `json:"sender_id"`
	SendType  string          `json:"sender_type"`
	TenantKey string          `json:"tenant_key"`
}

type MessageSenderID struct {
	UnionId string `json:"union_id"`
	UserId  string `json:"user_id"`
	OpenId  string `json:"open_id"`
}

type MessageDetail struct {
	MessageId   string            `json:"message_id"`
	RootId      string            `json:"root_id"`
	ParentId    string            `json:"parent_id"`
	CreateTime  string            `json:"create_time"`
	ChatId      string            `json:"chat_id"`
	ChatType    string            `json:"chat_type"`
	MessageType string            `json:"message_type"`
	Content     string            `json:"content"`
	Mentions    []MessageMentions `json:"mentions"`
}

type MessageMentions struct {
	Key       string           `json:"key"`
	Id        MessageMentionId `json:"id"`
	Name      string           `json:"name"`
	TenantKey string           `json:"tenant_key"`
}

type MessageMentionId struct {
	UnionId string `json:"union_id"`
	UserId  string `json:"user_id"`
	OpenId  string `json:"open_id"`
}

func MessageReceivedEventHandler(event interface{}) {
	revel.AppLog.Debug("开始处理订阅的消息接收事件.")
	b, errMarshal := json.Marshal(event)
	if errMarshal != nil {
		panic(errMarshal)
	}
	revel.AppLog.Debug("消息事件内容:" + string(b))
	var msgevent MessageReceiveEvent
	errUnmarshal := json.Unmarshal(b, &msgevent)
	if errUnmarshal != nil {
		panic(errUnmarshal)
	}
	receiveid := msgevent.Message.ChatId
	messageid := msgevent.Message.MessageId
	msgcontent := msgevent.Message.Content
	revel.AppLog.Debug(receiveid + " " + messageid + " " + msgcontent)
	api.Reply(receiveid, messageid, msgcontent)
}
