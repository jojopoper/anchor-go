package validation

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	_D "github.com/jojopoper/freeAnchor/models/db"
)

// QueryResultItem 正确结果子项返回结构定义
type QueryResultItem struct {
	Code        string    `json:"code"`
	NickName    string    `json:"nick_name"`
	Start       time.Time `json:"start_time"`
	End         time.Time `json:"close_time"`
	ConfirmHash string    `json:"hash"`
}

// QueryResult 正确结果返回结构定义
type QueryResult struct {
	AssetAddr string             `json:"asset_addr"`
	Items     []*QueryResultItem `json:"assets"`
}

// QueryOperation 查询接口定义
type QueryOperation struct {
	CheckParamValid
	assetAddr string
	nickname  string
	dbResults []*_D.TAnchorHistory
}

// GetResultData 返回处理结果
func (ths *QueryOperation) GetResultData() *OperationResult {
	if ths.resultData.CodeID == NoError {
		ret := &QueryResult{
			AssetAddr: ths.assetAddr,
			Items:     make([]*QueryResultItem, 0),
		}
		for _, itm := range ths.dbResults {
			if len(ret.AssetAddr) == 0 {
				ret.AssetAddr = itm.AssetAddr
			}
			subItm := &QueryResultItem{
				Code:        itm.AssetName,
				NickName:    itm.AnchorName,
				Start:       itm.CreateTime,
				End:         itm.CloseTimeout,
				ConfirmHash: itm.TransHash,
			}
			ret.Items = append(ret.Items, subItm)
		}

		ths.resultData.ResultData = ret
	}
	return ths.resultData
}

// QueryExecute 处理过程
func (ths *QueryOperation) QueryExecute() *OperationResult {
	if ths.resultData.CodeID == NoError {
		ths.dbResults = make([]*_D.TAnchorHistory, 0)
		condi := ""
		if len(ths.assetAddr) > 0 {
			condi = fmt.Sprintf("asset_addr='%s'", ths.assetAddr)
		}
		if len(ths.nickname) > 0 {
			if len(condi) > 0 {
				condi += " and "
			}
			condi += fmt.Sprintf("anchor_name='%s'", ths.nickname)
		}
		err := _D.DatabaseInstance.GetRecords(_D.DbAnchorHistoryOperation, condi, "id", -1, false, &ths.dbResults)
		if err != nil {
			ths.resultData.SetError(DataBaseError, err.Error())
		}
	}
	return ths.resultData
}

// DecodeContext 解码获取参数
func (ths *QueryOperation) DecodeContext(ctl *beego.Controller) {
	ths.resultData = new(OperationResult)
	ths.assetAddr = ctl.Input().Get("ac")
	ths.nickname = ctl.Input().Get("nn")
	if len(ths.assetAddr) == 0 && len(ths.nickname) == 0 {
		ths.resultData.SetError(CommomError, fmt.Sprintf("Can not recieve query parameter, at least query asset_name or nickname."))
	}
}
