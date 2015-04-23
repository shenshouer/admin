package controllers

import (
	r "github.com/revel/revel"

	"github.com/shenshouer/admin/app/models"
	"github.com/shenshouer/admin/app/services"
)

type (
	Base struct {
		*r.Controller
		services.Service
	}
)

func (this *Base) Before() r.Result {
	logger.Infof("Before UserId[%s] Path[%s]", this.Session.Id(), this.Request.URL.Path)
	this.UserId = this.Session.Id()

	var err error
	this.MongoSession, err = models.CopyMonotonicSession(this.UserId)
	if err != nil {
		logger.Errorln(err)
		return this.RenderError(err)
	}
	return nil
}

func (this *Base) After() r.Result {
	defer func() {
		if this.MongoSession != nil {
			models.CloseSession(this.UserId, this.MongoSession)
			this.MongoSession = nil
		}
	}()

	logger.Infof("After UserId[%s] Path[%s]", this.UserId, this.Request.URL.Path)
	return nil
}

func (this *Base) Panic() r.Result {
	defer func() {
		models.CloseSession(this.UserId, this.MongoSession)
		this.MongoSession = nil
	}()

	logger.Infof("Panic UserId[%s] Path[%s]", this.UserId, this.Request.URL.Path)
	return nil
}

func (this *Base) Base() *Base {
	return this
}

func (this *Base) Services() *services.Service {
	return &this.Service
}
