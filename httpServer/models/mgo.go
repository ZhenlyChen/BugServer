package models

import (
	"log"

	"github.com/globalsign/mgo"
)

// Mongo 数据库配置
type Mongo struct {
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	Name     string `yaml:"Name"`
}

// Model ...
type Model struct {
	Config Mongo
	DB     *mgo.Database
	User   UserModel
	Game   GameModel
}

// InitMongo 初始化数据库
func (m *Model) InitMongo(conf Mongo) error {
	m.Config = conf
	if m.DB != nil {
		m.DB.Session.Close()
	}
	session, err := mgo.Dial(
		"mongodb://" +
			conf.User +
			":" + conf.Password +
			"@" + conf.Host +
			":" + conf.Port +
			"/" + conf.Name)
	if err != nil {
		return err
	}
	m.DB = session.DB(conf.Name)
	m.User.DB = m.DB.C("users")
	m.Game.DB = m.DB.C("game")
	log.Printf("MongoDB Connect Success!")
	return nil
}

// NewModel ...
func NewModel(c Mongo) (*Model, error) {
	model := new(Model)
	err := model.InitMongo(c)
	return model, err
}
