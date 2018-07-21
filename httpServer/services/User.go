package services

import (
	violet "github.com/XMatrixStudio/Violet.SDK.Go"
	"github.com/ZhenlyChen/BugServer/httpServer/models"
	"github.com/globalsign/mgo"
)

// UserService 用户服务
type UserService interface {
	InitViolet(c violet.Config)
	// 登陆部分API
	Login(name, password string) (valid bool, email string, err error)
	GetUserFromViolet(code string) (ID, name string, err error)
	Register(name, email, password string) (err error)
	GetEmailCode(email string) error
	ValidEmail(email, vCode string) error
	// 用户信息 API
	GetUserBaseInfo(id string) (user UserBaseInfo)
	GetUserInfo(id string) (user models.Users, err error)
	SetUserInfo(id string, info models.UserInfo) error
}

type userService struct {
	Model    *models.UserModel
	Violet   violet.Violet
	UserInfo map[string]UserBaseInfo
	Service  *Service
}

// UserBaseInfo 用户个性信息
type UserBaseInfo struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Gender int    `json:"gender"`
}

// Login ...
func (s *userService) Login(name, password string) (valid bool, data string, err error) {
	res, err := s.Violet.Login(name, password)
	if err != nil {
		return
	}
	valid = res.Valid
	// 未激活邮箱
	if !valid {
		data = res.Email
		return
	}
	// 登陆成功
	data = res.Code
	return
}

// GetUserFromViolet ...
func (s *userService) GetUserFromViolet(code string) (ID, name string, err error) {
	// 获取用户Token
	tokenRes, err := s.Violet.GetToken(code)
	if err != nil {
		return
	}
	// 保存数据并获取用户信息
	if user, tErr := s.Model.GetUserByVID(tokenRes.UserID); tErr == nil { // 数据库已存在该用户
		ID = user.ID.Hex()
		name = user.Info.NikeName
		s.Model.SetUserToken(user.ID.Hex(), tokenRes.Token)
	} else if tErr == mgo.ErrNotFound { // 数据库不存在此用户
		userInfoRes, tErr := s.Violet.GetUserBaseInfo(tokenRes.UserID, tokenRes.Token)
		if err != nil {
			return "", "", tErr
		}
		bsonID, tErr := s.Model.AddUser(tokenRes.UserID, userInfoRes.Name, userInfoRes.Email)
		err = tErr
		ID = bsonID.Hex()
		name = "new_user"
	} else { // 其他错误
		err = tErr
	}
	return
}

func (s *userService) Register(name, email, password string) error {
	return s.Violet.Register(name, email, password)
}

func (s *userService) GetEmailCode(email string) error {
	return s.Violet.GetEmailCode(email)
}

func (s *userService) ValidEmail(email, vCode string) error {
	return s.Violet.ValidEmail(email, vCode)
}

func (s *userService) GetUserInfo(id string) (user models.Users, err error) {
	return s.Model.GetUserByID(id)
}

func (s *userService) SetUserInfo(id string, info models.UserInfo) error {
	if info.NikeName == "new_user" {
		return ErrNotAllow
	}
	users, err := s.Model.GetUsers()
	if err != nil {
		return ErrNotAllow
	}
	for _, user := range users {
		if user.Info.NikeName == info.NikeName {
			return ErrNotAllow
		}
	}
	s.GetUserBaseInfo(id)
	s.UserInfo[id] = UserBaseInfo{
		Avatar: info.Avatar,
		Name: info.NikeName,
		Gender: info.Gender,
	}
	return s.Model.SetUserInfo(id, info)
}

func (s *userService) InitViolet(c violet.Config) {
	s.Violet = violet.NewViolet(c)
}

// GetUserBaseInfo 从缓存中读取用户基本信息，如果不存在则从数据库中读取
func (s *userService) GetUserBaseInfo(id string) (user UserBaseInfo) {
	user, ok := s.UserInfo[id]
	if !ok {
		userInfo, err := s.GetUserInfo(id)
		if err != nil {
			return UserBaseInfo{
				Name:   "匿名用户",
				Avatar: "default",
				Gender: 0,
			}
		}
		user = UserBaseInfo{
			Name:   userInfo.Info.NikeName,
			Avatar: userInfo.Info.Avatar,
			Gender: userInfo.Info.Gender,
		}
		s.UserInfo[id] = user
	}
	return
}
