package db

import "time"

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
