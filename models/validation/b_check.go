package validation

import (
	"fmt"
)

// CheckParamValid 检查输入参数合法
type CheckParamValid struct {
	resultData *OperationResult
}

// ValidParam 参数检查
func (ths *CheckParamValid) ValidParam(val string, name string) error {
	if len(val) == 0 {
		return fmt.Errorf("%s is empty", name)
	}
	return nil
}
