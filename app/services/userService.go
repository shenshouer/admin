package services

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/shenshouer/admin/app/models"
)

func QueryUser(service *Service, userId string) (user *models.User, err error) {
	logger.Infoln("QueryUser", userId)

	// 查找具体的user id
	queryMap := bson.M{"_id": userId}
	// 执行查询
	user = &models.User{}
	err = service.DBAction("users", func(collection *mgo.Collection) error {
		return collection.Find(queryMap).One(user)
	})

	if err != nil {
		logger.Errorln("QueryUser() err", err)
		return nil, err
	}

	return user, nil
}

func CreateUser(service *Service, user *models.User) (err error) {
	err = service.DBAction("users", func(collection *mgo.Collection) error {
		return collection.Insert(&user)
	})
	return
}
