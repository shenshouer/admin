package models

import (
	"labix.org/v2/mgo/bson"
)

type (
	// 地址
	Address struct {
		Street   string `bson:"street" json:"street"`                // 街道
		District string `bson:"district, omitempty" json:"district"` // 行政区，如北京的海淀区
		County   string `bson:"county, omitempty" json:"county"`     // 县
		City     string `bson:"city" json:"city"`                    // 市
		Province string `bson:"province, omitempty" json:"province"` // 省份
		State    string `bson:"state" json:"state"`                  // 国家
	}
	// 联系方式
	Contact struct {
		Telephone int    `bson:"telephone" json:"telephone"` // 电话
		QQ        int    `bson:"qq" json:"qq"`               // qq
		Email     string `bson:"email" json:"email"`         // email
	}
	// 用户档案
	Profile struct {
		WorkAddress *Address `bson:"workAddress, omitempty" json:"workAddress"` // 工作单位地址
		HomeAddress *Address `bson:"homeAddress, omitempty" json:"homeAddress"` // 家庭住址
		Contact     *Contact `bson:"contact, omitempty" json:"contact"`         // 联系方式
	}
	// 用户
	User struct {
		Id       bson.ObjectId `bson:"_id,omitempty" json:"id"`  // 使用Mongodb 的ID，忽略空（omit empty）
		Username string        `bson:"username" json:"username"` // 用户名
		Password string        `bson:"password" json:"password"` // 密码
	}
)
