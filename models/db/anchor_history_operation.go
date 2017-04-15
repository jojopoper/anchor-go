package db

import (
	"github.com/go-xorm/xorm"
	_db "github.com/jojopoper/go-models/db"
)

const (
	DbAnchorHistoryOperation string = "AnchorHistoryOperation"
)

// AnchorHistoryOperation 锚点历史数据库控制定义
type AnchorHistoryOperation struct {
	_db.OperationBase
}

// Init 初始化
func (ths *AnchorHistoryOperation) Init(e *xorm.Engine) (err error) {
	err = ths.OperationBase.Init(e)
	ths.OperationKey = DbAnchorHistoryOperation
	ths.SetGernalFun(_db.UpdateFuncType, ths.update)
	ths.UpdateSqlCmd = "UPDATE `t_anchor_history` SET `account_id` = ?, " +
		"`create_time` = ?, `close_timeout` = ?, `anchor_name` = ?, `asset_addr` = ?, " +
		"`asset_name` = ?, `trans_hash` = ?, `paging_token` = ?, `status` = ?, " +
		"`left_balance` = ? WHERE `id`=?"
	return
}

func (ths *AnchorHistoryOperation) update(v interface{}) error {
	val := v.(*TAnchorHistory)
	_, err := ths.GetEngine().Exec(ths.UpdateSqlCmd, val.AccountID, val.CreateTime, val.CloseTimeout, val.AnchorName, val.AssetAddr,
		val.AssetName, val.TransHash, val.PagingToken, val.Status, val.LeftBalance, val.ID)
	// _, err := ths.engine.Where("`id`=?", val.ID).Update(v)
	return err
}
