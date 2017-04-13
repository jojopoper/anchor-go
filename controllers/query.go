package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego"
	_L "github.com/jojopoper/freeAnchor/models/log"
	_V "github.com/jojopoper/freeAnchor/models/validation"
)

//QueryController 查询服务器
type QueryController struct {
	beego.Controller
}

// Get get function
func (ths *QueryController) Get() {
	SetCxtRespHeader(ths.Ctx)
	ths.Ctx.Output.Header("content-type", "application/json;charset=utf8")

	_L.LoggerInstance.DebugPrint("\r\n[%s] \r\nParams : \r\n%+v\r\n", ths.Ctx.Input.IP(), ths.Input())
	op := &_V.OperationModel{}
	op.DecodeContext(&ths.Controller)
	ret := op.QueryExecute()
	if ret.CodeID == 0 {
		ret = op.GetResultData()
	}

	if ret.CodeID > 0 {
		_L.LoggerInstance.ErrorPrint("[%s] Validation Error : \r\n\t%s\r\n", ths.Ctx.Input.IP(), ret.ErrorMsg)
	}
	data, _ := json.Marshal(ret)
	ths.Data["json"] = string(data)
	ths.ServeJSON(false)
}
