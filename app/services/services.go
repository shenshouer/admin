package services

import (
	"labix.org/v2/mgo"

	"github.com/shenshouer/admin/app/base"
	"github.com/shenshouer/admin/app/models"
	"github.com/shenshouer/logging"
)

var (
	logger logging.Logger = base.NewLogger() // 日志记录器
)

type (
	// 包含所有service的共同的字段与行为
	Service struct {
		MongoSession *mgo.Session
		UserId       string
	}
)

// DBAction 对MongoDB执行查询和命令
func (this *Service) DBAction(collectionName string, mongoCall models.MongoCall) (err error) {
	return models.Execute(this.UserId,
		this.MongoSession,
		base.MONGO_DATABASE,
		collectionName,
		mongoCall)
}
