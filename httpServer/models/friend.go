package models

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// GameModel 好友数据库
type FriendModel struct {
	DB *mgo.Collection
}

// Friend 好友系统
type Friend struct {
	ID       bson.ObjectId   `bson:"_id"`
	UserID   bson.ObjectId   `bson:"userId"`   // 用户ID
	Friends  []bson.ObjectId `bson:"friends"`  // 朋友ID
	ToFriend []ToFriend      `bson:"toFriend"` // 好友申请
}

// 好友申请
type ToFriend struct {
	UserID bson.ObjectId `bson:"userId"`
	Msg    string        `bson:"msg"` // 留言
}

// GetByID 获取好友信息
func (m *FriendModel) GetByID(id string) (res Friend) {
	m.DB.FindId(bson.ObjectIdHex(id)).One(&res);
	return
}
