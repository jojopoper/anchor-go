package thread

import (
	_DB "github.com/jojopoper/freeAnchor/models/db"
	_F "github.com/jojopoper/freeAnchor/models/file"
	_L "github.com/jojopoper/freeAnchor/models/log"
	_ck "github.com/jojopoper/go-models/checker"
)

const (
	UpdateTomlFileName string = "UpdateTomlFile"
)

// UpdateTomlFile 更新Toml文件内容
type UpdateTomlFile struct {
	_ck.CheckBase
	fCtl *_F.TomlFileController
}

// Init 初始化
func (ths *UpdateTomlFile) Init(interval int, f *_F.TomlFileController) _ck.ICheckInterface {
	ths.CheckBase.Init(interval)
	ths.BeginStop()
	ths.SetContinue(false)
	ths.SetName(UpdateTomlFileName)
	ths.SetExeFunc(ths.exe)
	ths.fCtl = f
	return ths
}

func (ths *UpdateTomlFile) exe() {
	_L.LoggerInstance.InfoPrint("Updating toml file [%s]...\n", ths.fCtl.FilePath)
	datas := make([]*_DB.TAnchorHistory, 0)
	err := _DB.DatabaseInstance.GetRecords(_DB.DbAnchorHistoryOperation, "status=1", "id", -1, false, &datas)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("There is not any records in database: error is \n%+v\n", err)
		ths.SetContinue(true)
		return
	}
	for _, itm := range datas {
		ths.fCtl.AddCurrenc(itm.AssetName, itm.AssetAddr, itm.AnchorName)
	}
	err = ths.fCtl.FlashFile()
	if err != nil {
		_L.LoggerInstance.ErrorPrint("Update TOML file has error: \n%+v\n", err)
		ths.SetContinue(true)
		return
	}
	ths.SetContinue(false)
	_L.LoggerInstance.InfoPrint("Update toml file complete\n")
}
