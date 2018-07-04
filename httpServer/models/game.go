package models

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// GameModel 游戏数据库
type GameModel struct {
	DB *mgo.Collection
}

// Game ...
type Game struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	Version    int           `bson:"version" json:"version"`    // 版本
	Title      string        `bson:"title" json:"title"`      // 版本标题
	VersionStr string        `bson:"versionStr" json:"versionStr"` // 版本号
}

// GetNewestVersion 获取最新版本信息
func (m *GameModel) GetNewestVersion() (res Game) {
	m.DB.Find(nil).Sort("version").One(&res)
	return
}

// SetNewVersion 设置最新版本
func (m *GameModel) SetNewVersion(data Game) error {
	data.ID = bson.NewObjectId()
	return m.DB.Insert(&data)
}
