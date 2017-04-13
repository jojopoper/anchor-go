package thread

// ReportFunction 结果报告函数指针定义
type ReportFunction func(sender ICheckInterface, msg *CheckMessage)

// executeFunction 执行函数的定义
type executeFunction func()

// RunFlag 运行标识
type RunFlag int

const (
	ForceStop RunFlag = 1
)

// ICheckInterface 线程运行接口定义
type ICheckInterface interface {
	Name() string
	Start()
	Stop()
	IsRunning() bool
	IsBeginStart() bool
	RegistManager(m map[string]ICheckInterface) ICheckInterface
	AddReportFunc(f ReportFunction)
}
