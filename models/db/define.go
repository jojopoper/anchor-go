package db

import (
	"time"

	"github.com/go-xorm/xorm"
)

// DatabaseType 数据库类型
type DatabaseType string

const (
	DbAnchorHistoryOperation = "AnchorHistoryOperation"

	QtGetRecord = iota + 1
	QtCheckRecord
	QtGetLastRecord
	QtQuaryRecords
	QtQuaryAllRecord
	QtAddOneRecord
	QtUpdateRecord
	QtDeleteRecord
	QtSeachRecord

	MySqlDriver    DatabaseType = "mysql"
	sqliteDriver   DatabaseType = "sqlite3"
	PostgresDriver DatabaseType = "postgres"

	LanguageChinese = "cn"
	LanguageEnglish = "en"
)

// OperationInterface 数据库接口定义
type OperationInterface interface {
	Init(e *xorm.Engine) error
	GetKey() string
	Query(qtype int, v ...interface{}) error
	GetEngine() *xorm.Engine
}

// TAnchorHistory anchor 历史数据库表定义
type TAnchorHistory struct {
	ID           uint64    `xorm:"pk autoincr 'id'"`
	AccountID    string    `xorm:"account_id"` // 监控的地址
	CreateTime   time.Time `xorm:"DateTime"`   // 创建时间
	CloseTimeout time.Time `xorm:"DateTime"`   // 关闭时间
	AnchorName   string    `xorm:"varchar(128)"`
	AssetAddr    string    `xorm:"varchar(64) notnull"`
	AssetName    string    `xorm:"varchar(12) notnull"`
	TransHash    string    `xorm:"varchar(128)"`
	PagingToken  string    `xorm:"varchar(64) notnull"`
	Status       int       `xorm:"default(0)"` // 0:close 1:open
	LeftBalance  float64
}
