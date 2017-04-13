package validation

import (
	"fmt"

	"github.com/astaxie/beego"
)

const (
	NoError            = 0
	CommomError        = 1
	UnknownFormatError = 2
	DataBaseError      = 3
	ParamInvalid       = 100
)

// OperationResult 反馈结果结构定义
type OperationResult struct {
	ErrorMsg   string      `json:"error"`
	CodeID     int         `json:"codeid"`
	ResultData interface{} `json:"data"`
	Language   string      `json:"language"`
}

// SetError 设置错误信息
func (ths *OperationResult) SetError(id int, msg string) {
	ths.CodeID = id
	ths.ErrorMsg = msg
}

// IOperation 操作接口定义
type IOperation interface {
	GetResultData() *OperationResult
	QueryExecute() *OperationResult
	DecodeContext(ctl *beego.Controller)
}

// OperationModel 统一操作定义
type OperationModel struct {
	operation IOperation
}

// GetResultData 按照接口要求的获取结果
func (ths *OperationModel) GetResultData() *OperationResult {
	if ths.operation != nil {
		return ths.operation.GetResultData()
	}
	return &OperationResult{
		ErrorMsg: "Undefined parameter format[2]",
		CodeID:   UnknownFormatError,
	}
}

// QueryExecute 按照接口要求执行过程
func (ths *OperationModel) QueryExecute() *OperationResult {
	if ths.operation != nil {
		return ths.operation.QueryExecute()
	}
	return &OperationResult{
		ErrorMsg: "Can not find out result object",
		CodeID:   CommomError,
	}
}

// DecodeContext 按照接口要求解码获取参数
func (ths *OperationModel) DecodeContext(ctl *beego.Controller) {
	paramType := ctl.Input().Get("t")

	switch paramType {
	case "get":
		ths.operation = &QueryOperation{}
	case "ask":
		ths.operation = &AskManageAccount{}
	default:
		fmt.Println(paramType)
		return
	}
	ths.operation.DecodeContext(ctl)
}
