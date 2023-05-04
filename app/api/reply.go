package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/revel/revel"
)

type ResponseContent struct {
	Text string `json:"text"`
}

type ResponseData struct {
	Content     string `json:"content"`
	MessageType string `json:"msg_type"`
	ReceiveId   string `json:"receive_id"`
	UUID        string `json:"uuid"`
}

func buildUrl(messageId string) string {
	baseUrl := "https://open.feishu.cn/open-apis/im/v1/messages/"
	return baseUrl + messageId + "/reply"
}

func Reply(receiveId string, messageId string, message string) {
	if checkDuplicateMsgId(messageId) {
		revel.AppLog.Info("重复接受到相同messageId[" + messageId + "]的消息,不再重复处理。")
		return
	}

	var respData ResponseData
	var respContent ResponseContent
	respContent.Text = generateResp(message)

	respData.UUID = uuid.NewString()
	respData.MessageType = "text"
	respData.ReceiveId = receiveId
	bcontent, _ := json.Marshal(respContent)
	respData.Content = string(bcontent)

	posturl := buildUrl(messageId)
	revel.AppLog.Debug("回复消息URL:" + posturl)

	body, _ := json.Marshal(&respData)
	revel.AppLog.Debug("回复消息内容:" + string(body))

	r, err := http.NewRequest("POST", posturl, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	token := GetToken()
	revel.AppLog.Debug("Token is " + token)
	r.Header.Add("Authorization", "Bearer "+token)
	r.Header.Add("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		panic(res.Status)
	}
}

func generateResp(message string) string {
	return CompletionRequest(message)
}

func checkDuplicateMsgId(messageId string) bool {
	_, err := LocalCache.Get([]byte(messageId))
	if err != nil {
		if err.Error() == "Entry not found" {
			//记录messageid，避免重复处理。保留时长1小时
			LocalCache.Set([]byte(messageId), []byte("1"), 60*60)
			return false
		} else {
			panic(err)
		}
	} else {

		return true
	}
}
