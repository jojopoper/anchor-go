package routers

import (
	"github.com/astaxie/beego"
	"github.com/jojopoper/freeAnchor/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/federation", &controllers.FederationController{})
	beego.Router("/query", &controllers.QueryController{})
}
