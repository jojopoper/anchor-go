package thread

import (
	"sync"

	"github.com/astaxie/beego"
	_F "github.com/jojopoper/freeAnchor/models/file"
	_L "github.com/jojopoper/freeAnchor/models/log"
	_ck "github.com/jojopoper/go-models/checker"
	_nf "github.com/jojopoper/go-models/notify"
)

// CheckManagerInstance 全局唯一实例
var CheckManagerInstance *CheckManager

// NewCheckManager 获取实例
func NewCheckManager() *CheckManager {
	ret := new(CheckManager)
	return ret.Init()
}

// CheckManager 检查管理器
type CheckManager struct {
	checker        map[string]_ck.ICheckInterface
	tomlController *_F.TomlFileController
	lock           *sync.Mutex
}

// Init 初始化
func (ths *CheckManager) Init() *CheckManager {
	ths.lock = new(sync.Mutex)
	ths.tomlController = new(_F.TomlFileController)
	ths.tomlController.Init(beego.AppConfig.String("tomlFile"))
	ths.tomlController.FaderationURL = beego.AppConfig.String("fadr_serv")
	ths.regChecker()
	return ths
}

// Check 启动检测
func (ths *CheckManager) Check() {
	for _, c := range ths.checker {
		if !c.IsRunning() {
			if c.IsBeginStart() {
				go c.Start()
			}
		}
	}
}

func (ths *CheckManager) regChecker() {
	ths.checker = make(map[string]_ck.ICheckInterface)

	ct := new(CheckTransaction)
	ct.Init(30, beego.AppConfig.String("accountid")).RegistManager(ths.checker).AddReportFunc(ths.checkerReport)
	utf := new(UpdateTomlFile)
	utf.Init(5, ths.tomlController).RegistManager(ths.checker).AddReportFunc(ths.checkerReport)
	cid := new(CheckInvalidData)
	cid.Init(1800).RegistManager(ths.checker).AddReportFunc(ths.checkerReport)
}

// func (ths *CheckManager) checkerReport(sender _ck.ICheckInterface, msg *_nf.ReportMessage) {
func (ths *CheckManager) checkerReport(sender interface{}, msg *_nf.ReportMessage) {
	var senderFrom _ck.ICheckInterface
	if sender != nil {
		senderFrom = sender.(_ck.ICheckInterface)
	}
	ths.lock.Lock()
	defer ths.lock.Unlock()
	if msg.GetError() != nil && senderFrom != nil {
		_L.LoggerInstance.ErrorPrint("[%s] Has error msg :\n%+v\n", senderFrom.Name(), msg.GetError())
	}
	if msg.GetMessage() != "" && senderFrom != nil {
		_L.LoggerInstance.InfoPrint("[%s] -> %s\n", senderFrom.Name(), msg.GetMessage())
	}

	if msg.GetMessage() == UpdatedMessage {
		ths.checker[UpdateTomlFileName].Start()
	}
}
