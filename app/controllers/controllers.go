package controllers

import (
	"fmt"
	"github.com/shenshouer/logging"
	"regexp"
	"strconv"
	"strings"

	r "github.com/revel/revel"

	"github.com/shenshouer/admin/app/base"
)

const (
	CSessionRole = "CSessionRole"
)

var (
	// 日志记录器
	logger    logging.Logger = base.NewLogger()
	staticReg *regexp.Regexp = regexp.MustCompile(`static|admin.tologin|admin.getsecuritycode.*`)
)

// 初始化
func init() {
	r.InterceptFunc(checkAuthentication, r.BEFORE, r.ALL_CONTROLLERS)
}

func checkAuthentication(c *r.Controller) r.Result {
	action := strings.ToLower(c.Action)

	fmt.Println("拦截器", action, staticReg.MatchString(action))
	/*
		if !staticReg.MatchString(action) {
			sessionUserId, isExists := c.Session[CSessionRole]
			if !isExists {
				return c.Redirect(Admin.ToLogin)
			}
			fmt.Println(sessionUserId)
		}
	*/
	return nil
}

// 转换intStr 为uint类型 如出错，则使用默认_default
func parseUintOrDefault(intStr string, _default uint64) uint64 {
	if value, err := strconv.ParseUint(intStr, 0, 64); err != nil {
		return _default
	} else {
		return value
	}
}

// 转换intStr 为int类型 如出错，则使用默认_default
func parseIntOrDefault(intStr string, _default int64) int64 {
	if value, err := strconv.ParseInt(intStr, 0, 64); err != nil {
		return _default
	} else {
		return value
	}
}

// 转换int64 为string类型，如果出错，则使用默认_default
func parseStringOrDefault(intValue int64) string {
	return strconv.FormatInt(intValue, 10)
}

type loginResult struct {
	Code string `json:"code"`
	Data string `json:"data"`
}
