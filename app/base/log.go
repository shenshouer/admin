package base

import (
	"logging"
)

// 使用外包logging项目中的日志系统
func NewLogger() logging.Logger {
	return logging.NewSimpleLogger()
}
