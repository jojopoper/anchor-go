package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_L "github.com/jojopoper/freeAnchor/models/log"
	_ "github.com/lib/pq"
)

// DatabaseInstance 数据库访问唯一实例
var DatabaseInstance *CenterManager

// DatabaseInfo 数据库基本信息定义
type DatabaseInfo struct {
	DbType    DatabaseType
	AliasName string
	Host      string
	Port      string
	UserName  string
	Password  string
	IsDebug   bool
}

// CenterManager 数据库管理实例定义
type CenterManager struct {
	dbInfo     DatabaseInfo
	dbEngine   *xorm.Engine
	operations map[string]OperationInterface
}

// CreateDBInstance 创建实例
func CreateDBInstance(dbConfig DatabaseInfo) *CenterManager {
	ret := &CenterManager{}
	ret.Init(dbConfig)
	return ret
}

// Init 初始化数据库
func (ths *CenterManager) Init(dbConfig DatabaseInfo) {
	_L.LoggerInstance.InfoPrint("[CenterManager:InitDB] Init database begin\r\n")
	ths.dbInfo = dbConfig

	ths.initEngine()

	// 注册Orm的数据库表
	ths.ormRegModels(ths.dbEngine)

	ths.initOperation()
	_L.LoggerInstance.InfoPrint("[CenterManager:InitDB] Init database success\r\n")
}

func (ths *CenterManager) initEngine() {
	ths.dbEngine = nil
	switch ths.dbInfo.DbType {
	case MySqlDriver:
		ths.dbEngine = ths.getMySQLEngine()
	case PostgresDriver:
		ths.dbEngine = ths.getPostgresEngine()
	}
	if ths.dbEngine == nil {
		_L.LoggerInstance.ErrorPrint("[CenterManager:initEngine] Undefined db type = %s\r\n", ths.dbInfo.DbType)
		panic(1)
	}
	ths.dbEngine.ShowDebug = ths.dbInfo.IsDebug
	ths.dbEngine.ShowInfo = ths.dbInfo.IsDebug
	ths.dbEngine.ShowSQL = ths.dbInfo.IsDebug
	ths.dbEngine.ShowErr = true
	ths.dbEngine.ShowWarn = true
}

func (ths *CenterManager) getMySQLEngine() *xorm.Engine {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", //Asia%2FShanghai
		ths.dbInfo.UserName, ths.dbInfo.Password, ths.dbInfo.Host, ths.dbInfo.Port, ths.dbInfo.AliasName)
	ret, err := xorm.NewEngine(string(MySqlDriver), dataSourceName)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[CenterManager:getMySQLEngine] Create MySql has error! \r\n\t%v\r\n", err)
		return nil
	}
	err = ret.Ping()
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[CenterManager:getMySQLEngine] Create MySql Ping error! \r\n\t %v\r\n", err)
		return nil
	}
	return ret
}

func (ths *CenterManager) getPostgresEngine() *xorm.Engine {
	dataSourceName := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s sslmode=disable",
		ths.dbInfo.AliasName, ths.dbInfo.UserName, ths.dbInfo.Password, ths.dbInfo.Host, ths.dbInfo.Port)
	ret, err := xorm.NewEngine(string(PostgresDriver), dataSourceName)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("[Manager:getPostgresEngine] Create Postgres has error! \r\n\t%v\r\n", err)
		return nil
	}
	return ret
}

// ormRegModels 初始化数据库表
func (ths *CenterManager) ormRegModels(eng *xorm.Engine) {
	err := eng.Sync(new(TAnchorHistory))
	if err != nil {
		_L.LoggerInstance.InfoPrint("[CenterManager:ormRegModels] XORM Engine Sync is err %v\r\n", err)
		panic(1)
	}
}

func (ths *CenterManager) initOperation() {
	if ths.operations == nil {
		ths.operations = make(map[string]OperationInterface)
	}

	// 注册 LuckyTempOperation 操作
	aho := &AnchorHistoryOperation{}
	aho.Init(ths.dbEngine)
	ths.operations[aho.GetKey()] = aho
}

// GetOperation 得到对应的操作数据库控制器
func (ths *CenterManager) GetOperation(key string) OperationInterface {
	return ths.operations[key]
}

// Query 通用数据库执行接口
func (ths *CenterManager) Query(operationKey string, qtype int, val ...interface{}) error {
	return ths.operations[operationKey].Query(qtype, val...)
}

// Add 添加记录
func (ths *CenterManager) Add(operationKey string, val interface{}) error {
	return ths.operations[operationKey].Query(QtAddOneRecord, val)
}

// Get 获取记录
func (ths *CenterManager) Get(operationKey string, val interface{}) error {
	return ths.operations[operationKey].Query(QtGetRecord, val)
}

// Remove 删除记录
func (ths *CenterManager) Remove(operationKey string, val interface{}) error {
	return ths.operations[operationKey].Query(QtDeleteRecord, val)
}

// Update 更新记录
func (ths *CenterManager) Update(operationKey string, val interface{}) error {
	return ths.operations[operationKey].Query(QtUpdateRecord, val)
}

// GetLastRecord 获取最新cnt条记录
func (ths *CenterManager) GetLastRecord(operationKey string, val ...interface{}) error {
	return ths.operations[operationKey].Query(QtGetLastRecord, val...)
}

// GetRecords 获取最新cnt条记录
func (ths *CenterManager) GetRecords(operationKey string, val ...interface{}) error {
	return ths.operations[operationKey].Query(QtQuaryRecords, val...)
}
