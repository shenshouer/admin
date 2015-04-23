package models

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/revel/revel"
	"github.com/shenshouer/admin/app/base"
	"github.com/shenshouer/logging"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const (
	MASTER_SESSION    = "master"    // 强连接session 即主连接
	MONOTONIC_SESSION = "monotonic" // 单调的session
)

var (
	_This  *mongoManager                     // 单例的引用
	logger logging.Logger = base.NewLogger() // 日志记录器
)

type (
	// mongoManager 所需要的 包含了拨号(dial)与session信息
	mongoSession struct {
		mongoDBDialInfo *mgo.DialInfo
		mongoSession    *mgo.Session
	}

	// mongoManager 管理一个session的map，类似jdbc的数据库连接池
	mongoManager struct {
		sessions map[string]*mongoSession
	}

	// mongocall定义了一个类型的函数，可用于在MongoDB执行代码
	MongoCall func(*mgo.Collection) error
)

// 启动，携带管理器到运行时状态
func Startup(sessionId string) (err error) {
	logger.Infof("Startup sessionId[%s]", sessionId)

	// 创建 Mongo Manager
	_This = &mongoManager{
		sessions: map[string]*mongoSession{},
	}

	// 日志输出mongodb的链接信息
	logger.Infof("MongoDB : Addr[%s]", revel.Config.StringDefault("mgo.host", ""))
	logger.Infof("MongoDB : Database[%s]", revel.Config.StringDefault("mgo.database", ""))
	logger.Infof("MongoDB : Username[%s]", revel.Config.StringDefault("mgo.username", ""))

	hosts := strings.Split(revel.Config.StringDefault("mgo.host", ""), ",")
	fmt.Print("==============>> ", hosts)
	// 创建强连接与不变连接session
	err = CreateSession(sessionId,
		"strong",
		MASTER_SESSION,
		hosts,
		revel.Config.StringDefault("mgo.database", ""),
		revel.Config.StringDefault("mgo.username", ""),
		revel.Config.StringDefault("mgo.password", ""))

	if err != nil {
		logger.Errorf("Startup Monogodb err[%v]", err)
	}

	err = CreateSession(sessionId,
		"monotonic",
		MONOTONIC_SESSION,
		hosts,
		revel.Config.StringDefault("mgo.database", ""),
		revel.Config.StringDefault("mgo.username", ""),
		revel.Config.StringDefault("mgo.password", ""))

	logger.Infoln("Startup", sessionId, "Completed")

	if err != nil {
		logger.Errorf("Startup Monogodb err[%v]", err)
	}
	return err
}

// 关闭系统时优雅地关闭管理器
func Shutdown(sessionId string) (err error) {
	logger.Infoln(sessionId, "Shutdown")

	// 关闭数据库连接
	for _, session := range _This.sessions {
		CloseSession(sessionId, session.mongoSession)
	}

	return err
}

// CreateSession 创建一个可供使用的连接池
func CreateSession(sessionId string, mode string, sessionName string, hosts []string, databaseName string, username string, password string) (err error) {
	logger.Infof("CreateSession Mode[%s] SessionName[%s] sessionId[%s] Hosts[%s] DatabaseName[%s] Username[%s]", mode, sessionName, sessionId, hosts, databaseName, username)

	// 创建数据库连接对象
	mongoSession := &mongoSession{
		mongoDBDialInfo: &mgo.DialInfo{
			Addrs:    hosts,
			Timeout:  60 * time.Second,
			Database: databaseName,
			Username: username,
			Password: password,
		},
	}
	fmt.Println("===========>> mongoSession ", mongoSession.mongoDBDialInfo)
	// 创建主session
	mongoSession.mongoSession, err = mgo.DialWithInfo(mongoSession.mongoDBDialInfo)
	if err != nil {
		logger.Errorf("create session named [%s] err[%v]", sessionId, err)
		return err
	}

	switch mode {
	case "strong":
		// Reads and writes will always be made to the master server using a
		// unique connection so that reads and writes are fully consistent,
		// ordered, and observing the most up-to-date data.
		// http://godoc.org/labix.org/v2/mgo#Session.SetMode
		mongoSession.mongoSession.SetMode(mgo.Strong, true)
		break

	case "monotonic":
		// Reads may not be entirely up-to-date, but they will always see the
		// history of changes moving forward, the data read will be consistent
		// across sequential queries in the same session, and modifications made
		// within the session will be observed in following queries (read-your-writes).
		// http://godoc.org/labix.org/v2/mgo#Session.SetMode
		mongoSession.mongoSession.SetMode(mgo.Monotonic, true)
	}

	// Have the session check for errors
	// http://godoc.org/labix.org/v2/mgo#Session.SetSafe
	mongoSession.mongoSession.SetSafe(&mgo.Safe{})

	// 将连接信息添加到map
	_This.sessions[sessionName] = mongoSession

	return err
}

// CopyMasterSession 复制一个主session的副本给客户端使用
func CopyMasterSession(sessionId string) (*mgo.Session, error) {
	return CopySession(sessionId, MASTER_SESSION)
}

// CopyMonotonicSession 复制一个monotonic模式的session的副本给客户端使用
func CopyMonotonicSession(sessionId string) (*mgo.Session, error) {
	return CopySession(sessionId, MONOTONIC_SESSION)
}

// CopySession 复制一个具体的session给客户端使用
func CopySession(sessionId string, useSession string) (mongoSession *mgo.Session, err error) {
	logger.Infof("CopySession UseSession[%s] sessionId[%s]", useSession, sessionId)

	session := _This.sessions[useSession]

	if session == nil {
		err = fmt.Errorf("Unable To Locate Session %s", useSession)
		logger.Errorln(err, sessionId, "CopySession")
		return mongoSession, err
	}

	// 复制主session
	mongoSession = session.mongoSession.Copy()

	return mongoSession, err
}

// CloneMasterSession 克隆主session
func CloneMasterSession(sessionId string) (*mgo.Session, error) {
	return CloneSession(sessionId, MASTER_SESSION)
}

// CloneMonotonicSession 克隆monotonic session
func CloneMonotonicSession(sessionId string) (*mgo.Session, error) {
	return CloneSession(sessionId, MONOTONIC_SESSION)
}

// CloneSession 克隆
func CloneSession(sessionId string, useSession string) (mongoSession *mgo.Session, err error) {
	logger.Infof("CloneSession UseSession[%s] sessionId[%s]", useSession, sessionId)

	// 查找session对象
	session := _This.sessions[useSession]

	if session == nil {
		err = fmt.Errorf("Unable To Locate Session %s", useSession)
		logger.Errorf("CloneSession sessionId[%s], err[%v]", sessionId, err)
		return mongoSession, err
	}

	mongoSession = session.mongoSession.Clone()

	return mongoSession, err
}

// CloseSession 将连接归还到连接池
func CloseSession(sessionId string, mongoSession *mgo.Session) {
	logger.Infof("CloseSession sessionId[%s]", sessionId)

	mongoSession.Close()
}

// GetCollection 根据具体的数据库session与collection名称返回一个collection的引用
func GetCollection(mongoSession *mgo.Session, useDatabase string, useCollection string) (*mgo.Collection, error) {
	return mongoSession.DB(useDatabase).C(useCollection), nil
}

// CollectionExists 判断collection名称在具体的数据库中是否存在
func CollectionExists(sessionId string, mongoSession *mgo.Session, useDatabase string, useCollection string) bool {
	database := mongoSession.DB(useDatabase)
	collections, err := database.CollectionNames()

	if err != nil {
		return false
	}

	for _, collection := range collections {
		if collection == useCollection {
			return true
		}
	}

	return false
}

// ToString 将查询map转换成string
func ToString(queryMap bson.M) string {
	json, err := json.Marshal(queryMap)
	if err != nil {
		return ""
	}
	return string(json)
}

// Execute MongoDB 的 literal 功能
func Execute(sessionId string, mongoSession *mgo.Session, databaseName string, collectionName string, mongoCall MongoCall) (err error) {
	logger.Infof("Execute Database[%s] Collection[%s] sessionId[%s]", databaseName, collectionName, sessionId)

	// 获取一个具体的collection
	collection, err := GetCollection(mongoSession, databaseName, collectionName)
	if err != nil {
		logger.Errorln(err, sessionId, "Execute")
		return err
	}

	// 执行MongoCall接口/回调函数
	err = mongoCall(collection)
	if err != nil {
		logger.Errorln(err, sessionId, "Execute")
		return err
	}

	return err
}
