package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Verfiy() revel.Result {
	var jsondata map[string]interface{}
	c.Params.BindJSON(&jsondata)
	revel.AppLog.Debug("Feishu Verfisy: " + string(c.Params.JSON))
	data := make(map[string]interface{})
	data["challenge"] = jsondata["challenge"]
	revel.AppLog.Debug("Reply Challenge: ", data)
	return c.RenderJSON(data)
}
