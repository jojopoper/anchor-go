package db

import (
	"fmt"

	"github.com/go-xorm/xorm"
)

// AnchorHistoryOperation 锚点历史数据库控制定义
type AnchorHistoryOperation struct {
	OperationBase
	updateCmd string
}

// Init 初始化
func (ths *AnchorHistoryOperation) Init(e *xorm.Engine) (err error) {
	err = ths.OperationBase.Init(e)
	ths.OperationBase.updateFunc = ths.update
	ths.updateCmd = "UPDATE `t_anchor_history` SET `account_id` = ?, " +
		"`create_time` = ?, `close_timeout` = ?, `anchor_name` = ?, `asset_addr` = ?, " +
		"`asset_name` = ?, `trans_hash` = ?, `paging_token` = ?, `status` = ?, " +
		"`left_balance` = ? WHERE `id`=?"
	return
}

// GetKey 获取Key值
func (ths *AnchorHistoryOperation) GetKey() string {
	return DbAnchorHistoryOperation
}

// Query 操作
func (ths *AnchorHistoryOperation) Query(qtype int, v ...interface{}) error {
	if v == nil || len(v) < 1 {
		return fmt.Errorf("[AnchorHistoryOperation:Query] Query get anchor history , must input TAnchorHistory struct pointer as parameter")
	}
	// switch qtype {
	// case QtGetLastRecord:
	// 	if len(v) != 2 {
	// 		return fmt.Errorf("Input parameter has to 2, Parameter cnt,typeof struct")
	// 	}
	// 	return ths.GetLastRecord(v[0].(int), v[1])
	// case QtQuaryRecords:
	// 	if len(v) != 4 {
	// 		return fmt.Errorf("Input parameter has to 3, Parameter conditions,cnt,isdesc,typeof struct")
	// 	}
	// 	return ths.GetRecords(v[0].(string), v[1].(int), v[2].(bool), v[3])
	// }
	return ths.OperationBase.Query(qtype, v...)
}

func (ths *AnchorHistoryOperation) update(v interface{}) error {
	val := v.(*TAnchorHistory)
	_, err := ths.engine.Exec(ths.updateCmd, val.AccountID, val.CreateTime, val.CloseTimeout, val.AnchorName, val.AssetAddr,
		val.AssetName, val.TransHash, val.PagingToken, val.Status, val.LeftBalance, val.ID)
	// _, err := ths.engine.Where("`id`=?", val.ID).Update(v)
	return err
}
