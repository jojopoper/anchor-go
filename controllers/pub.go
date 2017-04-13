package controllers

import "github.com/astaxie/beego/context"

// SetCxtRespHeader 设置Context的ResponseHeader，为跨域访问
func SetCxtRespHeader(cxt *context.Context) {
	cxt.Output.Header("Access-Control-Allow-Headers", "Origin, No-Cache, X-Requested-With, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With, Accept")
	cxt.Output.Header("Access-Control-Allow-Origin", "*")
	cxt.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
}