package file

import (
	"fmt"
)

// TomlCurrenciesInfo 资产定义
type TomlCurrenciesInfo struct {
	Code string
	Issuer string
	Nickname string
}

// ToString 返回文件Currencies字符串
func (ths *TomlCurrenciesInfo) ToString() string {
	ret := fmt.Sprintf("[[CURRENCIES]]\ncode=\"%s\"\nissuer=\"%s\"\n",ths.Code, ths.Issuer)
	if len(ths.Nickname) != 0 {
		ret += fmt.Sprintf("# nickname=%s\n", ths.Nickname)
	}
	return ret
}