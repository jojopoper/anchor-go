package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/astaxie/beego"
	_D "github.com/jojopoper/freeAnchor/models/db"
	_L "github.com/jojopoper/freeAnchor/models/log"
	_T "github.com/jojopoper/freeAnchor/models/thread"
	_ "github.com/jojopoper/freeAnchor/routers"
	_db "github.com/jojopoper/go-models/db"
	_api "github.com/jojopoper/stellarApi"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	_L.LoggerInstance = _L.NewLoggerInstance(fmt.Sprintf("log/anchor.%s", time.Now().Format("2006-01-02_15.04.05.000")))
	_L.LoggerInstance.OpenDebug = true
	_L.LoggerInstance.SetLogFunCallDepth(4)

	dbInfo := _db.DatabaseInfo{
		DbType:    _db.DatabaseType(beego.AppConfig.String("dbtype")),
		AliasName: beego.AppConfig.String("aliasname"),
		Host:      "127.0.0.1",
		Port:      "3306",
		UserName:  "root",
		Password:  "1234",
	}
	dbInfo.IsDebug = beego.AppConfig.String("runmode") == "dev"
	_D.DatabaseInstance = _D.CreateDBInstance(dbInfo)
	_api.SetHorizonBand(_api.FlyHorizon)
	_T.CheckManagerInstance = _T.NewCheckManager()
	_T.CheckManagerInstance.Check()
	beego.Run()
}
