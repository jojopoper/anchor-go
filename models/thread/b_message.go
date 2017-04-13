package thread

// CheckMessage 检测消息载体定义
type CheckMessage struct {
	err  error
	msg  string
	data interface{}
}

// SetCheckMessage 设置消息
func (ths *CheckMessage) SetCheckMessage(e error, m string, d interface{}) *CheckMessage {
	ths.err = e
	ths.msg = m
	ths.data = d
	return ths
}

// SetError 设置错误
func (ths *CheckMessage) SetError(e error) *CheckMessage {
	return ths.SetCheckMessage(e, "", nil)
}

// GetError 获取错误
func (ths *CheckMessage) GetError() error {
	return ths.err
}

// SetMessage 设置消息
func (ths *CheckMessage) SetMessage(m string) *CheckMessage {
	return ths.SetCheckMessage(nil, m, nil)
}

// GetMessage 获取消息
func (ths *CheckMessage) GetMessage() string {
	return ths.msg
}

// SetData 设置数据
func (ths *CheckMessage) SetData(d interface{}) *CheckMessage {
	return ths.SetCheckMessage(nil, "", d)
}

// GetData 获取数据
func (ths *CheckMessage) GetData() interface{} {
	return ths.data
}
