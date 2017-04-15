package db

import (
	"github.com/go-xorm/xorm"
	_L "github.com/jojopoper/freeAnchor/models/log"
	_db "github.com/jojopoper/go-models/db"
)

// DatabaseInstance 数据库访问唯一实例
var DatabaseInstance *CenterManager

// CenterManager 数据库管理实例定义
type CenterManager struct {
	_db.ManageBase
}

// CreateDBInstance 创建实例
func CreateDBInstance(dbConfig _db.DatabaseInfo) *CenterManager {
	ret := &CenterManager{}
	ret.Init(dbConfig)
	return ret
}

// Init 初始化数据库
func (ths *CenterManager) Init(dbConfig _db.DatabaseInfo) {
	_L.LoggerInstance.InfoPrint("[CenterManager:Init] Init database begin\n")
	ths.SetRegModelsFunc(ths.ormRegModels)
	ths.SetInitOperaFunc(ths.initOperation)
	ths.ManageBase.Init(dbConfig)
	_L.LoggerInstance.InfoPrint("[CenterManager:Init] Init database success\n")
}

// ormRegModels 初始化数据库表
func (ths *CenterManager) ormRegModels(eng *xorm.Engine) {
	err := eng.Sync(new(TAnchorHistory))
	if err != nil {
		_L.LoggerInstance.InfoPrint("[CenterManager:ormRegModels] XORM Engine Sync is err %v\n", err)
		panic(1)
	}
}

func (ths *CenterManager) initOperation() {
	// 注册 AnchorHistoryOperation 操作
	_L.LoggerInstance.InfoPrint("Regist 'AnchorHistoryOperation' operation\n")
	ths.AppendOperation(new(AnchorHistoryOperation))
}
