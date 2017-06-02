package controllers

import (
	"github.com/astaxie/beego"
)

// MainController 主控制器
type MainController struct {
	beego.Controller
}

// Get Http get func
func (ths *MainController) Get() {
	SetCxtRespHeader(ths.Ctx)
	ths.Data["Website"] = "ebitgo.com"
	ths.Data["Email"] = "support@ebitgo.com"
	ths.TplName = "index.tpl"
}

// Post Http post func
func (ths *MainController) Post() {
	SetCxtRespHeader(ths.Ctx)
	ths.Data["Website"] = "ebitgo.com"
	ths.Data["Email"] = "support@ebitgo.com"
	ths.TplName = "index.tpl"
}
