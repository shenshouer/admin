package controllers

import (
	r "github.com/revel/revel"

	"github.com/shenshouer/admin/app/services"
)

type User struct {
	Base
}

func init() {
	r.InterceptMethod((*User).Before, r.BEFORE)
	r.InterceptMethod((*User).After, r.AFTER)
	r.InterceptMethod((*User).Panic, r.PANIC)
}

func (this *User) User(userId string) r.Result {
	user, err := services.QueryUser(this.Services(), userId)

	if err != nil {
		return this.RenderText(err.Error())
	}
	return this.RenderJson(user)
}
