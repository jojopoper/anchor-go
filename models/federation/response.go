package federation

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/astaxie/beego"
	_D "github.com/jojopoper/freeAnchor/models/db"
	_api "github.com/jojopoper/stellarApi"
)

// TransactionDef transaction defined
type TransactionDef struct {
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	MemoType string `json:"memo_type"`
	Memo     string `json:"memo"`
}

// ResponseBaseMsg 返回正确消息定义
type ResponseBaseMsg struct {
	Address   string `json:"stellar_address"`
	Accountid string `json:"account_id"`
	MemoType  string `json:"memo_type"`
	Memo      string `json:"memo"`
}

// ResponseMsg 返回消息定义
type ResponseMsg struct {
	Status    int
	Msg       *ResponseBaseMsg
	ErrMsg    string `json:"detail"`
	QParam    string
	TypeParam string
}

// DecodeParam 获取输入参数
func (ths *ResponseMsg) DecodeParam(vals url.Values) *ResponseMsg {
	ths.Status = 200
	ths.QParam = vals.Get("q")
	err := ths.checkValue(ths.QParam, "q")
	if err != nil {
		ths.Status = 501
		ths.ErrMsg = err.Error()
		return ths
	}
	ths.TypeParam = vals.Get("type")
	err = ths.checkValue(ths.TypeParam, "type")
	if err != nil {
		ths.Status = 501
		ths.ErrMsg = err.Error()
	}
	return ths
}

// Execute 执行
func (ths *ResponseMsg) Execute() *ResponseMsg {
	ths.TypeParam = strings.ToLower(ths.TypeParam)
	switch ths.TypeParam {
	case "name":
		ths.nameRequest()
	case "id":
		ths.getIDRequest()
	case "txid":
		ths.txidRequest()
	default:
		ths.Status = 501
		ths.ErrMsg = fmt.Sprintf("Undefined type [%s] string!", ths.TypeParam)
	}
	return ths

}

func (ths *ResponseMsg) checkValue(val, flag string) (err error) {
	if len(val) == 0 {
		err = fmt.Errorf("Input parameter [%s] is not invalid or nil.", flag)
	}
	return
}

func (ths *ResponseMsg) nameRequest() {
	name := ths.getNickName(ths.QParam)
	if len(name) > 0 {
		ths.getFromAnchor(name)
	}
}

func (ths *ResponseMsg) getNickName(param string) string {
	name := strings.ToLower(param)
	suffixIndex := strings.LastIndex(name, "*anchor.ebitgo.com")
	if suffixIndex == -1 || suffixIndex == 0 {
		ths.Status = 501
		ths.ErrMsg = fmt.Sprintf("type=name has error [ The value of 'q' format is error ]")
		return ""
	}
	return strings.Trim(string([]byte(name)[0:suffixIndex]), " ")
}

func (ths *ResponseMsg) getFromAnchor(n string) {
	tmp := &_D.TAnchorHistory{
		AnchorName: n,
		Status:     1,
	}
	err := _D.DatabaseInstance.Get(_D.DbAnchorHistoryOperation, tmp)
	if err != nil {
		ths.Status = 404
		ths.ErrMsg = err.Error()
	} else {
		ths.Msg = new(ResponseBaseMsg)
		ths.Status = 200
		ths.Msg.Accountid = tmp.AssetAddr
		ths.Msg.Address = n + "*anchor.ebitgo.com"
	}
}

func (ths *ResponseMsg) getIDRequest() {
	id := ths.QParam
	tmp := &_D.TAnchorHistory{
		AssetAddr: id,
		Status:    1,
	}
	err := _D.DatabaseInstance.Get(_D.DbAnchorHistoryOperation, tmp)
	if err != nil {
		ths.Status = 404
		ths.ErrMsg = err.Error()
	} else {
		ths.Msg = new(ResponseBaseMsg)
		ths.Status = 200
		ths.Msg.Accountid = id
		ths.Msg.Address = tmp.AnchorName + "*anchor.ebitgo.com"
	}
}

func (ths *ResponseMsg) txidRequest() {
	txid := strings.ToLower(ths.QParam)
	trans := _api.NewTransactionRequest(txid)
	param := &_api.QueryParameters{}
	param.UseTestNetwork, _ = beego.AppConfig.Bool("stellar_test_network")
	param.HttpType = _api.ClientHttp
	err := trans.GetTransInfo(nil, param)
	if err != nil {
		ths.Status = 404
		ths.ErrMsg = fmt.Sprintf("Get transaction [%s] info has error : \n%v\n", txid, err)
		return
	}
	if trans.Status == 0 {
		ths.Status = 200
		ths.Msg = new(ResponseBaseMsg)
		mt, _ := trans.GetMemo()
		ths.Msg.Accountid = trans.GetSourceAccount()
		tmp := &_D.TAnchorHistory{
			AssetAddr: ths.Msg.Accountid,
			Status:    1,
		}
		err = _D.DatabaseInstance.Get(_D.DbAnchorHistoryOperation, tmp)
		if err == nil {
			ths.Msg.Address = tmp.AnchorName + "*anchor.ebitgo.com"
		}
		ths.Msg.Memo = trans.GetMemoString()
		ths.Msg.MemoType = mt.String()
	} else {
		ths.Status = trans.Status
		ths.ErrMsg = trans.Detail
	}
}
