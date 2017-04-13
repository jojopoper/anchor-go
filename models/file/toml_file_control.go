package file

import (
	"fmt"
	"io"
	"os"
	"path"
	"sync"
)

// TomlFileController toml文件控制
type TomlFileController struct {
	FaderationURL string
	currencs      []*TomlCurrenciesInfo
	locker        *sync.Mutex
	FilePath      string
}

// Init 初始化
func (ths *TomlFileController) Init(fpath string) {
	ths.locker = new(sync.Mutex)
	ths.currencs = make([]*TomlCurrenciesInfo, 0)
	ths.FilePath = fpath
}

// AddCurrencies 添加资产
func (ths *TomlFileController) AddCurrencies(c *TomlCurrenciesInfo) {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	ths.currencs = append(ths.currencs, c)
}

// AddCurrenc 添加资产
func (ths *TomlFileController) AddCurrenc(code, issuer, nickname string) {
	ths.AddCurrencies(&TomlCurrenciesInfo{
		Code:     code,
		Issuer:   issuer,
		Nickname: nickname,
	})
}

// ToString 返回文件Currencies字符串
func (ths *TomlFileController) ToString() string {
	ret := fmt.Sprintf("FEDERATION_SERVER=\"%s\"\n", ths.FaderationURL)
	for _, itm := range ths.currencs {
		ret += fmt.Sprintf("\n%s", itm.ToString())
	}
	return ret
}

// FlashFile 写数据到文件中，并清除currencs
func (ths *TomlFileController) FlashFile() error {
	ths.locker.Lock()
	defer ths.locker.Unlock()
	dir := path.Dir(ths.FilePath)
	err := os.MkdirAll(dir, os.FileMode(0777))
	if err != nil {
		return err
	}
	f, err := os.Create(ths.FilePath)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = io.WriteString(f, ths.ToString())
	ths.currencs = make([]*TomlCurrenciesInfo, 0)
	return err
}
