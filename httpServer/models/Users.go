package models

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// UserModel 用户数据库
type UserModel struct {
	DB *mgo.Collection
}

// Users 用户
type Users struct {
	ID       bson.ObjectId `bson:"_id"`   // 用户ID
	VioletID bson.ObjectId `bson:"vid"`   // VioletID
	Name     string        `bson:"name"`  // 用户唯一名字
	Email    string        `bson:"email"` // 邮箱
	Info     UserInfo      `bson:"info"`  // 用户个性信息
	Token    string        `bson:"token"` // Violet 访问令牌
	Level    int64         `bson:"level"` // 用户等级
}

// 性别
const (
	GenderMan int = iota
	GenderWoman
	GenderUnknown
)

// UserInfo 用户个性信息
type UserInfo struct {
	Avatar   string `bson:"avatar"`   // 头像URL
	Gender   int    `bson:"gender"`   // 性别
	NikeName string `bson:"nikeName"` // 昵称
}

// AddUser 添加用户
func (m *UserModel) AddUser(vID, name, email, token, avatar string, gender int) (bson.ObjectId, error) {
	newUser := bson.NewObjectId()
	err := m.DB.Insert(&Users{
		ID:       newUser,
		VioletID: bson.ObjectIdHex(vID),
		Name:     name,
		Email:    email,
		Info: UserInfo{
			NikeName: "user_" + newUser.Hex(),
			Gender:   gender,
			Avatar:   avatar,
		},
		Token: token,
	})
	if err != nil {
		return "", err
	}
	return newUser, nil
}

// SetUserInfo 设置用户信息
func (m *UserModel) SetUserInfo(id string, info UserInfo) (err error) {
	_, err = m.DB.UpsertId(bson.ObjectIdHex(id), bson.M{"$set": info})
	return
}

// SetUserName 设置用户名
func (m *UserModel) SetUserName(id, name string) (err error) {
	_, err = m.DB.UpsertId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"name": name}})
	return
}

// SetUserToken 设置Token
func (m *UserModel) SetUserToken(id, token string) (err error) {
	_, err = m.DB.UpsertId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"token": token}})
	return
}

// GetUserByID 根据ID查询用户
func (m *UserModel) GetUserByID(id string) (user Users, err error) {
	err = m.DB.FindId(bson.ObjectIdHex(id)).One(&user)
	return
}

// GetUserByVID 根据VioletID查询用户
func (m *UserModel) GetUserByVID(id string) (user Users, err error) {
	err = m.DB.Find(bson.M{"vid": bson.ObjectIdHex(id)}).One(&user)
	return
}

// GetUserByName 根据名字获取用户
func (m *UserModel) GetUserByName(name string) (user Users, err error) {
	err = m.DB.Find(bson.M{"name": name}).One(&user)
	return
}

// GetUserByEmail 根据邮箱获取用户
func (m *UserModel) GetUserByEmail(email string) (user Users, err error) {
	err = m.DB.Find(bson.M{"email": email}).One(&user)
	return
}
