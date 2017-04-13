package thread

import (
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	_DB "github.com/jojopoper/freeAnchor/models/db"
	_L "github.com/jojopoper/freeAnchor/models/log"
	_api "github.com/jojopoper/stellarApi"
	_x "github.com/stellar/go/xdr"
)

const (
	CheckTransName string = "CheckTransaction"
	UpdatedMessage string = "UpdatedDatabaseMsg"
)

// CheckTransaction 检查账户transaction
type CheckTransaction struct {
	CheckBase
	cursor   string
	accTrans *_api.AccountTransactions
	accid    string // 需要检测的地址
}

// Init 初始化
func (ths *CheckTransaction) Init(interval int, accid string) ICheckInterface {
	ths.CheckBase.Init(interval)
	ths.keyName = CheckTransName
	ths.accid = accid
	ths.CheckBase.exe = ths.exe
	ths.accTrans = _api.NewAccountTransactions(ths.accid)
	return ths
}

func (ths *CheckTransaction) exe() {
	_L.LoggerInstance.InfoPrint("Checking base account transaction ...\n")
	err := ths.initCursor()
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[CheckTransaction : do] Init cursor has error: \n%+v\n", err)
		return
	}
	for {
		cnt, err := ths.updateCursor()
		if err != nil {
			_L.LoggerInstance.ErrorPrint("[CheckTransaction : do] Update cursor has error: \n%+v\n", err)
			return
		}
		if cnt < 10 {
			_L.LoggerInstance.InfoPrint("Checking complete!\n")
			return
		}
	}
	_L.LoggerInstance.InfoPrint("Check base account transaction complete\n")
}

func (ths *CheckTransaction) initCursor() error {
	if len(ths.cursor) > 0 {
		return nil
	}
	_L.LoggerInstance.InfoPrint("Init cursor ...\n")
	// 如果当前没有cursor就从数据库中读取一个，如果数据库里面也没有，就使用配置文件中的值
	ths.getStartCursor(nil)
	if len(ths.cursor) == 0 {
		ths.cursor = beego.AppConfig.String("initCorsor")
		_L.LoggerInstance.InfoPrint("Get cursor [%s] from config...\n", ths.cursor)
	}
	return nil
}

func (ths *CheckTransaction) getStartCursor(wt *sync.WaitGroup) {
	if wt != nil {
		defer wt.Done()
	}
	_L.LoggerInstance.InfoPrint("Read cursor from database ...\n")
	datas := make([]*_DB.TAnchorHistory, 0)
	err := _DB.DatabaseInstance.GetLastRecord(_DB.DbAnchorHistoryOperation, 1, "id", &datas)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("Get last cursor from db has error :\n%+v\n", err)
		return
	}
	if 0 < len(datas) {
		ths.cursor = datas[0].PagingToken
	} else {
		_L.LoggerInstance.InfoPrint(" Init cursor read from database is empty!\n")
	}
}

// updateCursor 更新cursor
func (ths *CheckTransaction) updateCursor() (int, error) {
	if len(ths.cursor) == 0 {
		return -1, fmt.Errorf("[CheckTransaction:updateCursor] Can not get current cursor")
	}

	_L.LoggerInstance.InfoPrint("Begin read 10 paging token from stellar network...\n")
	queryParam := &_api.QueryParameters{
		Size:      10,
		Cursor:    ths.cursor,
		OrderType: _api.AscOrderType,
	}
	queryParam.HttpType = _api.ClientHttp
	queryParam.UseTestNetwork, _ = beego.AppConfig.Bool("stellar_test_network")
	err := ths.accTrans.GetTransInfo(nil, queryParam)
	if err != nil {
		return -1, fmt.Errorf("Get prev transaction from horizon has error : \n account id : %s\n error : %+v", ths.accid, err)
	}
	if ths.accTrans.TransactionSize == 0 {
		return 0, nil
	}
	var msg *CheckMessage
	for idx := ths.accTrans.TransactionSize - 1; idx >= 0; idx-- {
		ret := ths.getDataFromOrig(ths.accTrans.GetAccTransUnit(idx))
		if ret != nil {
			ths.cursor = ret.PagingToken
			if ret.Status == 1 { // 添加或更新
				msg, err = ths.status1Func(ret)
			} else if ret.Status == -1 { // 删除
				msg, err = ths.status2Func(ret)
			} else if ret.Status == -100 {
				msg, err = ths.status3Func(ret)
			}
			if err != nil {
				return -1, err
			}
		}
	}
	if msg != nil {
		ths.report(ths, msg.SetMessage(UpdatedMessage))
	}
	return ths.accTrans.TransactionSize, nil
}

func (ths *CheckTransaction) status1Func(r *_DB.TAnchorHistory) (msg *CheckMessage, err error) {
	if len(r.AnchorName) > 0 {
		tmp := &_DB.TAnchorHistory{
			AnchorName: r.AnchorName,
		}
		err = _DB.DatabaseInstance.Get(_DB.DbAnchorHistoryOperation, tmp)
		if err == nil && tmp.AssetAddr != r.AssetAddr {
			_L.LoggerInstance.ErrorPrint("Nick name is repeated asset_addr[%s] has use the nicknanme[%s], Accountid [%s] repeat ask!\n",
				tmp.AssetAddr, tmp.AnchorName, r.AssetAddr)
			return nil, nil
		}
	}
	tmp := &_DB.TAnchorHistory{
		AssetAddr: r.AssetAddr,
		// AssetName: r.AssetName,
	}
	_L.LoggerInstance.DebugPrint("Findout data before :\n%+v\n", tmp)
	err = _DB.DatabaseInstance.Get(_DB.DbAnchorHistoryOperation, tmp)
	_L.LoggerInstance.DebugPrint("Findout data after :\n%+v\nerror:\n%+v\n", tmp, err)
	if err != nil {
		leftDay := int64(math.Floor(r.LeftBalance))
		r.CloseTimeout = r.CreateTime.Add(time.Duration(leftDay) * time.Hour * 24)
		err = _DB.DatabaseInstance.Add(_DB.DbAnchorHistoryOperation, r)
		msg = new(CheckMessage)
	} else {
		if tmp.AssetName == r.AssetName {
			if tmp.CreateTime.Unix() < r.CreateTime.Unix() {
				leftDay := int64(math.Floor(r.LeftBalance))
				if tmp.Status == 1 {
					r.LeftBalance += tmp.LeftBalance
					r.CloseTimeout = tmp.CloseTimeout.Add(time.Duration(leftDay) * time.Hour * 24)
				} else {
					r.CloseTimeout = r.CreateTime.Add(time.Duration(leftDay) * time.Hour * 24)
				}
				r.ID = tmp.ID
				if len(r.AssetName) == 0 {
					r.AssetName = tmp.AssetName
				}
				if len(r.AnchorName) == 0 {
					r.AnchorName = tmp.AnchorName
				}
				err = _DB.DatabaseInstance.Update(_DB.DbAnchorHistoryOperation, r)
				msg = new(CheckMessage)
			}
		} else {
			leftDay := int64(math.Floor(r.LeftBalance))
			r.CloseTimeout = r.CreateTime.Add(time.Duration(leftDay) * time.Hour * 24)
			err = _DB.DatabaseInstance.Add(_DB.DbAnchorHistoryOperation, r)
			msg = new(CheckMessage)
		}
	}
	return
}

func (ths *CheckTransaction) status2Func(r *_DB.TAnchorHistory) (msg *CheckMessage, err error) {
	tmp := &_DB.TAnchorHistory{
		AssetAddr:  r.AssetAddr,
		AssetName:  r.AssetName,
		AnchorName: r.AnchorName,
	}
	err = _DB.DatabaseInstance.Remove(_DB.DbAnchorHistoryOperation, tmp)
	if err == nil {
		msg = new(CheckMessage)
	}
	return
}

func (ths *CheckTransaction) status3Func(r *_DB.TAnchorHistory) (msg *CheckMessage, err error) {
	if r.LeftBalance <= 0 {
		return
	}
	if len(r.AssetName) > 0 {
		r.Status = 1
		return ths.status1Func(r)
	}
	return
}

func (ths *CheckTransaction) getDataFromOrig(tranInfo *_api.AccountTransUnit) *_DB.TAnchorHistory {
	if tranInfo == nil {
		return nil
	}
	ret := &_DB.TAnchorHistory{
		AccountID:   ths.accid,
		CreateTime:  tranInfo.LedgerCloseTime,
		AssetAddr:   tranInfo.GetSourceAccount(),
		TransHash:   tranInfo.Hash,
		PagingToken: tranInfo.PagingToken,
		Status:      -100,
	}
	t, m := tranInfo.GetMemo()
	if t == _x.MemoTypeMemoText {
		ths.SplitMemoText(m.(string), ret)
		// ret.AssetName = m.(string)
		if ret.Status == -1 { // 删除
			return ret
		}
	}

	var tp, cd, iu string
	if tranInfo.OperationCount > 0 {
		transEnv, _ := tranInfo.GetEnvelope()
		for _, opItm := range transEnv.Tx.Operations {
			pay, ok := opItm.Body.GetPaymentOp()
			if ok {
				if strings.Compare(pay.Destination.Address(), ths.accid) == 0 &&
					strings.Compare(ret.AssetAddr, ths.accid) != 0 {
					pay.Asset.MustExtract(&tp, &cd, &iu)
					if tp == "native" {
						ret.LeftBalance = float64(pay.Amount) / float64(10000000.0)
						if len(ret.AssetName) > 0 {
							ret.Status = 1
						}
					}
					break
				}
			}
		}
	}
	return ret
}

// 分割Memo字段获取内容
// [cmd]:[AssetCode]:[NickName]
func (ths *CheckTransaction) SplitMemoText(m string, r *_DB.TAnchorHistory) {
	if len(m) == 0 {
		return
	}
	strs := strings.Split(m, ":")
	if len(strs) != 3 {
		return
	}
	if strings.Trim(strings.ToLower(strs[0]), " ") == "rm" {
		r.Status = -1
	} else {
		r.Status = 0
	}
	r.AssetName = strings.Trim(strs[1], " ")
	r.AnchorName = strings.ToLower(strings.Trim(strs[2], " "))
}
