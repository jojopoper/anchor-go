package thread

import (
	_DB "github.com/jojopoper/freeAnchor/models/db"
	_F "github.com/jojopoper/freeAnchor/models/file"
	_L "github.com/jojopoper/freeAnchor/models/log"
)

const (
	UpdateTomlFileName string = "UpdateTomlFile"
)

// UpdateTomlFile 更新Toml文件内容
type UpdateTomlFile struct {
	CheckBase
	fCtl *_F.TomlFileController
}

// Init 初始化
func (ths *UpdateTomlFile) Init(interval int, f *_F.TomlFileController) ICheckInterface {
	ths.CheckBase.Init(interval)
	ths.beginStart = false
	ths.keyName = UpdateTomlFileName
	ths.fCtl = f
	ths.CheckBase.exe = ths.exe
	return ths
}

func (ths *UpdateTomlFile) exe() {
	_L.LoggerInstance.InfoPrint("Updating toml file [%s]...\n", ths.fCtl.FilePath)
	datas := make([]*_DB.TAnchorHistory, 0)
	err := _DB.DatabaseInstance.GetRecords(_DB.DbAnchorHistoryOperation, "status=1", "id", -1, false, &datas)
	if err != nil {
		_L.LoggerInstance.ErrorPrint("There is not any records in database: error is \n%+v\n", err)
		ths.continueRun = true
		return
	}
	for _, itm := range datas {
		ths.fCtl.AddCurrenc(itm.AssetName, itm.AssetAddr, itm.AnchorName)
	}
	err = ths.fCtl.FlashFile()
	if err != nil {
		_L.LoggerInstance.ErrorPrint("Update TOML file has error: \n%+v\n", err)
		ths.continueRun = true
		return
	}
	ths.continueRun = false
	_L.LoggerInstance.InfoPrint("Update toml file complete\n")
}
