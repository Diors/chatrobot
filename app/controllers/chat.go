package controllers

import (
	"encoding/json"

	"github.com/revel/revel"
	"xonesoft.cn/chatrobot/app/api"
)

type Chat struct {
	*revel.Controller
}

type event struct {
	Message message `json:"message"`
}

type message struct {
	ChatId    string `json:"chat_id"`
	ChatType  string `json:"chat_type"`
	Content   string `json:"content"`
	MessageId string `json:"message_id"`
}

type content struct {
	Text string `json:"text"`
}

func (c Chat) ChatTxt() revel.Result {
	var jsondata map[string]interface{}
	c.Params.BindJSON(&jsondata)
	revel.AppLog.Debug("Get Request Message:", jsondata)
	if jsondata["type"] == "url_verification" {
		data := make(map[string]interface{})
		data["challenge"] = jsondata["challenge"]
		return c.RenderJSON(data)
	} else {
		b, _ := json.Marshal(jsondata["event"])
		revel.AppLog.Debug("收到消息内容:" + string(b))
		var e event
		err1 := json.Unmarshal(b, &e)
		if err1 != nil {
			panic(err1)
		}
		var msgContent content
		err2 := json.Unmarshal([]byte(e.Message.Content), &msgContent)
		if err2 != nil {
			panic(err2)
		}
		revel.AppLog.Debug("收到消息内容:", e.Message)

		//异步回复消息
		defer api.Reply(e.Message.ChatId, e.Message.MessageId, msgContent.Text)

		revel.AppLog.Debug("接受消息，异步返回答复")
		return c.RenderText("OK")
	}
}

/**
curl --location --request POST 'https://open.feishu.cn/open-apis/im/v1/messages/om_ae552979c1973903333271be0ddcd41e/reply' \
--header 'Authorization: Bearer t-g1044s3CEEPAXCBVBCRIPPYV676XZFAHOXIYKNHH' \
--header 'Content-Type: application/json; charset=utf-8' \
-d '{
    "content": "{\"text\":\"<at user_id=\\\"ou_155184d1e73cbfb8973e5a9e698e74f2\\\">Tom</at> test content\"}",
    "msg_type": "text",
    "uuid": "a0d69e20-1dd1-458b-k525-dfeca4015204"
}'

curl --location --request POST 'https://open.feishu.cn/open-apis/im/v1/messages/om_952661c3d24a25dcdf3d4ea49016a5ec/reply' \
--header 'Authorization: Bearer t-g1045453FCHRAW4LQBX75KXJUZFSAOVAN7BTSEWJ' \
--header 'Content-Type: application/json; charset=utf-8' \
-d '{"content":"{\"text\":\"您好，我是AI助手，请问有什么需要帮助的吗？\"}","msg_type":"text","receive_id":"oc_51e9767c972a0e923c9fd1843e7f9a78","uuid":"3320c078-10ad-488c-a242-f1e6178fc1ec"} '

map[event:map[message:map[chat_id:oc_6b79c410b430535db2e0ac9eef934d20 chat_type:group content:{"text":"@_user_1 ccc"}
create_time:1682621194394 mentions:[map[id:map[open_id:ou_9899f6616df58f654ed88128f2301981
union_id:on_cd3b0106e3015f0f8022ebbf2999e4f4 user_id:] key:@_user_1 name:XAI tenant_key:13ddfdb6778e5740]]
message_id:om_ae552979c1973903333271be0ddcd41e message_type:text]
sender:map[sender_id:map[open_id:ou_8459ef0131ffdd0002a23ae3a729a236
union_id:on_3c6ef2d8b742ee95484357ea948e0ba2 user_id:f557d334] sender_type:user tenant_key:13ddfdb6778e5740]]
header:map[app_id:cli_a4a7d0486e38100e create_time:1682621194585 event_id:11cfd88846a4c77e3bda1148f34ebdad
event_type:im.message.receive_v1 tenant_key:13ddfdb6778e5740 token:yLcZlilWloN2GjILxHkFxYWFPayjxOO4] schema:2.0]

const APP_ID="cli_a4a7d0486e38100e"
const APP_SECRET="ETI7XjwPy0vztgwZ6KGh4cTRQtHmwp8r"

curl --location --request POST 'https://open.feishu.cn/open-apis/auth/v3/app_access_token/internal' \
--header 'Content-Type: application/json; charset=utf-8' \
-d '{"app_id": "cli_a4a7d0486e38100e","app_secret": "ETI7XjwPy0vztgwZ6KGh4cTRQtHmwp8r"}'

{"app_access_token":"t-g1044s3CEEPAXCBVBCRIPPYV676XZFAHOXIYKNHH","code":0,"expire":7200,"msg":"ok","tenant_access_token":"t-g1044s3CEEPAXCBVBCRIPPYV676XZFAHOXIYKNHH"}
**/
