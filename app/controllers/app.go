package controllers

import (
	r "github.com/revel/revel"

	"github.com/shenshouer/admin/app/base"

	"bytes"
	"fmt"
	"strconv"
	"time"
)

type App struct {
	*r.Controller
}

func (this App) Index() r.Result {
	return this.Render()
}

// 跳转至登录页面
func (this App) ToLogin() r.Result {
	return this.RenderTemplate("App/login.html")
}

// 获取验证码
func (this App) GetSecurityCode(timestamp int64) r.Result {
	// 时间戳参数，第一次加载为1，后续加载为当前的时间戳，可以用来验证客户端刷新频率
	// 如：刷新频率过高，直接限制当前客户端等
	//fmt.Println("GetSecurityCode", timestamp)

	d := make([]byte, 4)
	s := base.NewLen(4)
	ss := ""
	d = []byte(s)

	for v := range d {
		d[v] %= 10
		ss += strconv.FormatInt(int64(d[v]), 32)
	}

	this.Session["securityCode"] = ss
	fmt.Println(this.Session["securityCode"])

	buff := &bytes.Buffer{}
	base.NewImage(d, 100, 40).WriteTo(buff)

	return this.RenderBinary(buff, "img", r.Attachment, time.Now())
}

// 登陆验证
func (this App) Login() r.Result {
	// 从请求参数中获取当前的请求表单值
	username := this.Params.Get("username")
	password := this.Params.Get("password")
	securityCode := this.Params.Get("securityCode")

	fmt.Println("para ", securityCode, "sess", this.Session["securityCode"], username, password)
	if securityCode != this.Session["securityCode"] {
		return this.RenderJson(&loginResult{Code: "err", Data: "验证码错误"})
	} else {
		if username != "shenshouer@yahoo.com.cn" || password != "123456" {
			return this.RenderJson(&loginResult{Code: "err", Data: "账号或密码错误"})
		} else {
			//TODO 将验证完成后的用户标识设置到session
			return this.RenderJson(&loginResult{Code: "ok", Data: "/admin/index"})
		}
	}

	return nil
}

func (this App) Logout() r.Result {
	// 清除当前用户的信息
	//strId := this.Session[CSessionRole]

	for k := range this.Session {
		delete(this.Session, k)
	}

	return this.Redirect(App.Index)
}
