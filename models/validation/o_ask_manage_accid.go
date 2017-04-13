package validation

import "github.com/astaxie/beego"

// AskManageAccount 请求中心打款账户定义
type AskManageAccount struct {
	CheckParamValid
}

// GetResultData 返回处理结果
func (ths *AskManageAccount) GetResultData() *OperationResult {
	if ths.resultData.CodeID == NoError {
		ths.resultData.ResultData = beego.AppConfig.String("accountid")
	}
	return ths.resultData
}

// QueryExecute 处理过程
func (ths *AskManageAccount) QueryExecute() *OperationResult {
	return ths.resultData
}

// DecodeContext 解码获取参数
func (ths *AskManageAccount) DecodeContext(ctl *beego.Controller) {
	ths.resultData = new(OperationResult)
}
