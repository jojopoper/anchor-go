package db

import (
	"fmt"
	"sync"

	"github.com/go-xorm/xorm"
)

// GernalFunc 通用方法定义
type GernalFunc func(v interface{}) error

// OperationBase 基础Operation定义
type OperationBase struct {
	engine     *xorm.Engine
	locker     *sync.Mutex
	addFunc    GernalFunc
	getFunc    GernalFunc
	removeFunc GernalFunc
	updateFunc GernalFunc
}

// Init 定义
func (ths *OperationBase) Init(e *xorm.Engine) error {
	if e == nil {
		return fmt.Errorf("[OperationBase:Init] Can not get database engine")
	}
	ths.locker = &sync.Mutex{}
	ths.engine = e
	ths.addFunc = ths.add
	ths.removeFunc = ths.remove
	ths.getFunc = ths.get
	return nil
}

// GetEngine 获取数据库引擎，如果需要操作必须注意使用线程锁
func (ths *OperationBase) GetEngine() *xorm.Engine {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	return ths.engine
}

// Query 操作
func (ths *OperationBase) Query(qtype int, v ...interface{}) error {
	if v == nil || len(v) < 1 {
		return fmt.Errorf("[OperationBase:Query] Query function , must input typeof struct pointer as parameter")
	}

	ths.locker.Lock()
	defer ths.locker.Unlock()
	switch qtype {
	case QtAddOneRecord:
		return ths.addFunc(v[0])
	case QtGetRecord:
		return ths.getFunc(v[0])
	case QtDeleteRecord:
		return ths.removeFunc(v[0])
	case QtUpdateRecord:
		return ths.updateFunc(v[0])
	case QtGetLastRecord:
		if len(v) != 3 {
			return fmt.Errorf("Input parameter has to 3, Parameter cnt,orderkey,typeof struct")
		}
		return ths.getLastRecord(v[0].(int), v[1].(string), v[2])
	case QtQuaryRecords:
		if len(v) != 5 {
			return fmt.Errorf("Input parameter has to 5, Parameter conditions,orderkey,cnt,isdesc,typeof struct")
		}
		return ths.getRecords(v[0].(string), v[1].(string), v[2].(int), v[3].(bool), v[4])
	}

	return fmt.Errorf("[OperationBase:Query] Query type is not defined (%d)", qtype)
}

func (ths *OperationBase) get(v interface{}) error {
	b, err := ths.engine.Get(v)
	if err != nil {
		return fmt.Errorf("[OperationBase:get] %v", err)
	}
	if !b {
		return fmt.Errorf("[OperationBase:get] Get is not exist")
	}
	return nil
}

func (ths *OperationBase) add(v interface{}) error {
	_, err := ths.engine.InsertOne(v)
	return err
}

func (ths *OperationBase) remove(v interface{}) error {
	_, err := ths.engine.Delete(v)
	return err
}

// GetLastRecord 得到最新的cnt条记录
func (ths *OperationBase) getLastRecord(cnt int, orderKey string, v interface{}) error {
	session := ths.engine.NewSession()
	defer session.Close()
	if cnt > 0 {
		session = session.Limit(cnt)
	}
	return session.Desc(orderKey).Find(v)
}

// GetRecords 得到一定条件的数据
func (ths *OperationBase) getRecords(conditions, orderKey string, cnt int, isdesc bool, v interface{}) error {
	session := ths.engine.NewSession()
	defer session.Close()
	if len(conditions) > 0 {
		session = session.Where(conditions)
	}
	if cnt > 0 {
		session = session.Limit(cnt)
	}
	if isdesc {
		return session.Desc(orderKey).Find(v)
	}
	return session.Asc(orderKey).Find(v)
}
