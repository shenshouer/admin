package base

import (
	"github.com/shenshouer/logging"
)

// 使用外包logging项目中的日志系统
func NewLogger() logging.Logger {
	return logging.NewSimpleLogger()
}
