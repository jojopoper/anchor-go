package thread

import (
	"time"

	_D "github.com/jojopoper/freeAnchor/models/db"
	_L "github.com/jojopoper/freeAnchor/models/log"
	_ck "github.com/jojopoper/go-models/checker"
)

const (
	CheckInvalidDataName string = "CheckInvalidData"
)

// CheckInvalidData 循环检查数据库中的数据是否有超时的
type CheckInvalidData struct {
	_ck.CheckBase
}

// Init 初始化
func (ths *CheckInvalidData) Init(interval int) _ck.ICheckInterface {
	ths.CheckBase.Init(interval)
	ths.SetName(CheckInvalidDataName)
	ths.SetExeFunc(ths.exe)
	return ths
}

func (ths *CheckInvalidData) exe() {
	_L.LoggerInstance.InfoPrint("Checking database ...\n")
	tmpdatas := make([]*_D.TAnchorHistory, 0)
	condi := "status=1"
	err := _D.DatabaseInstance.GetRecords(_D.DbAnchorHistoryOperation, condi, "id", -1, true, &tmpdatas)
	var msg *_ck.CheckMessage
	if err == nil {
		for _, itm := range tmpdatas {
			if itm.CloseTimeout.Unix() <= time.Now().Unix() {
				if msg == nil {
					msg = new(_ck.CheckMessage)
				}
				itm.Status = 0
				err = _D.DatabaseInstance.Update(_D.DbAnchorHistoryOperation, itm)
				if err != nil {
					_L.LoggerInstance.ErrorPrint("[CheckInvalidData:exe()] Update data has error: \n%+v\n", err)
				}
				time.Sleep(500 * time.Millisecond)
			}
		}
		if msg != nil {
			ths.Report(ths, msg.SetMessage(UpdatedMessage))
		}
	} else {
		_L.LoggerInstance.ErrorPrint("[CheckInvalidData:exe()] Checking database status has error: \n%+v\n", err)
	}
	_L.LoggerInstance.InfoPrint("Check database complete\n")
}
