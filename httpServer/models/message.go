package models

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// MessageModel 消息据库
type MessageModel struct {
	DB *mgo.Collection
}

// UserMessage 消息系统
type UserMessage struct {
	ID       bson.ObjectId `bson:"_id"`
	UserID   bson.ObjectId `bson:"userId"`  // 用户ID
	Messages []Message     `bson:"message"` // 消息
}

// Message 消息
type Message struct {
	Title  string        `bson:"title" json:"title"`   // 标题
	Type   string        `bson:"type" json:"type"`     // 类型
	Msg    string        `bson:"msg" json:"msg"`       // 信息内容
	UserID bson.ObjectId `bson:"userId" json:"userId"` // 信息主体
}

const (
	TypeMessage_System = "system" // 系统信息
	TypeMessage_Friend = "friend" // 好友信息
)

// GetNewestVersion 获取最新版本信息
func (m *MessageModel) GetNewestVersion() (res Game) {
	m.DB.Find(nil).Sort("version").One(&res)
	return
}
