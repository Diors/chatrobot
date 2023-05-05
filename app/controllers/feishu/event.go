package controllers

import (
	"github.com/revel/revel"
	"xonesoft.cn/chatrobot/app/controllers/feishu/handler"
)

type FeishuEventSubscribe struct {
	*revel.Controller
}

type FeishuEventV2 struct {
	Challenge string             `json:"challenge"`
	Token     string             `json:"token"`
	Type      string             `json:"type"`
	Schema    string             `json:"schema"`
	Header    FeishEventV2Header `json:"header"`
	Event     interface{}        `json:"event"`
}

type FeishEventV2Header struct {
	EventId    string `json:"event_id"`
	EventType  string `json:"event_type"`
	CreateTime string `json:"create_time"`
	Token      string `json:"token"`
	AppId      string `json:"app_id"`
	TenantKey  string `json:"tenant_key"`
}

func (c FeishuEventSubscribe) EventSubscribe() revel.Result {
	var event FeishuEventV2
	c.Params.BindJSON(&event)
	revel.AppLog.Debug("收到消息：", event)
	// 飞书配置请求地址验证消息
	if event.Type == "url_verification" {
		data := make(map[string]interface{})
		data["challenge"] = event.Challenge
		return c.RenderJSON(data)
	}

	// 验证订阅事件token
	if VerificationToken(event.Header.Token) == false {
		revel.AppLog.Error("验证事件订阅Token失败,请确认是否有其他应用调用该接口.")
		c.Response.Status = 404
		return c.Render()
	}

	//接口验证通过,处理订阅事件
	go handEvent(event.Header, event.Event)

	//直接返回成功
	return c.Render()
}

// 未加密token验证
func VerificationToken(token string) bool {
	definedToken, isfound := revel.Config.String("feishu.event.token")
	if isfound {
		return token == definedToken
	}
	return false
}

func handEvent(header FeishEventV2Header, event interface{}) {
	eventType := header.EventType
	if eventType == "im.message.receive_v1" {
		handler.MessageReceivedEventHandler(event)
	}
}
