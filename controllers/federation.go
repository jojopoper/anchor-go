package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
	_F "github.com/jojopoper/freeAnchor/models/federation"
	_L "github.com/jojopoper/freeAnchor/models/log"
)

//FederationController 联邦服务器
type FederationController struct {
	beego.Controller
}

// Get get function
func (ths *FederationController) Get() {
	SetCxtRespHeader(ths.Ctx)
	msg := &_F.ResponseMsg{}
	msg.DecodeParam(ths.Input())
	if msg.Status == 200 {
		msg.Execute()
	}

	if msg.Status == 200 {
		ths.Data["json"] = msg.Msg
		_L.LoggerInstance.InfoPrint("[%s] [%+v] \r\nget message success!\r\n", ths.Ctx.Input.IP(), ths.Input())
	} else {
		ths.Data["json"] = map[string]interface{}{
			"detail": msg.ErrMsg,
		}
		http.Error(ths.Ctx.ResponseWriter, "", msg.Status)
		_L.LoggerInstance.InfoPrint("[%s] [%+v] \r\nget message failure! [errmsg : %s]\r\n", ths.Ctx.Input.IP(), ths.Input(), msg.ErrMsg)
	}
	ths.ServeJSON(false)
}
