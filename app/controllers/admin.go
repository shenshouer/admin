package controllers

import (
	r "github.com/revel/revel"
)

type Admin struct {
	*r.Controller
}

func (this Admin) Index() r.Result {
	return this.Render()
}
